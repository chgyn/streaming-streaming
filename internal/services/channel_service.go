package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"local-stream-hub/internal/models"

	"github.com/patrickmn/go-cache"
)

const (
	allChannelsCacheKey   = "channels:all"
	categoriesCacheKey    = "channels:categories"
	defaultChannelTimeout = 12 * time.Second
)

// ChannelService gerencia o catálogo e busca de canais ao vivo com cache in-memory.
type ChannelService interface {
	GetChannels(ctx context.Context, category, query string) ([]models.Channel, error)
	GetCategories(ctx context.Context) ([]string, error)
}

type channelService struct {
	baseURL string
	client  *http.Client
	cache   *cache.Cache
	ttl     time.Duration
	mu      sync.Mutex // Evita múltiplos fetches simultâneos à API externa
}

// NewChannelService inicializa o serviço com cache in-memory e HTTP client otimizado.
func NewChannelService(baseURL string, cacheTTL time.Duration) ChannelService {
	transport := &http.Transport{
		MaxIdleConns:        50,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     60 * time.Second,
		TLSHandshakeTimeout: 5 * time.Second,
	}

	client := &http.Client{
		Timeout:   defaultChannelTimeout,
		Transport: transport,
	}

	c := cache.New(cacheTTL, cacheTTL*2)

	return &channelService{
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  client,
		cache:   c,
		ttl:     cacheTTL,
	}
}

// GetChannels retorna a lista filtrada de canais ao vivo.
func (s *channelService) GetChannels(ctx context.Context, category, query string) ([]models.Channel, error) {
	channels, err := s.getAllChannels(ctx)
	if err != nil {
		return nil, fmt.Errorf("falha ao carregar catálogo de canais: %w", err)
	}

	filtered := make([]models.Channel, 0, len(channels))
	targetCat := strings.TrimSpace(strings.ToLower(category))
	targetQuery := strings.TrimSpace(strings.ToLower(query))

	for _, ch := range channels {
		// Filtro de categoria
		if targetCat != "" && targetCat != "todos" && targetCat != "all" {
			if strings.ToLower(ch.Category) != targetCat {
				continue
			}
		}

		// Filtro de texto de busca
		if targetQuery != "" {
			nameMatch := strings.Contains(strings.ToLower(ch.Name), targetQuery)
			catMatch := strings.Contains(strings.ToLower(ch.Category), targetQuery)
			idMatch := strings.Contains(strings.ToLower(ch.ID), targetQuery)
			if !nameMatch && !catMatch && !idMatch {
				continue
			}
		}

		filtered = append(filtered, ch)
	}

	return filtered, nil
}

// GetCategories retorna todas as categorias únicas de canais ordenadas alfabeticamente.
func (s *channelService) GetCategories(ctx context.Context) ([]string, error) {
	if val, found := s.cache.Get(categoriesCacheKey); found {
		if cats, ok := val.([]string); ok {
			return cats, nil
		}
	}

	channels, err := s.getAllChannels(ctx)
	if err != nil {
		return nil, fmt.Errorf("falha ao obter categorias: %w", err)
	}

	catMap := make(map[string]struct{})
	for _, ch := range channels {
		cat := strings.TrimSpace(ch.Category)
		if cat != "" {
			catMap[cat] = struct{}{}
		}
	}

	categories := make([]string, 0, len(catMap))
	for cat := range catMap {
		categories = append(categories, cat)
	}
	sort.Strings(categories)

	s.cache.Set(categoriesCacheKey, categories, s.ttl)
	return categories, nil
}

// getAllChannels recupera a lista completa com cache e fallback resiliente.
func (s *channelService) getAllChannels(ctx context.Context) ([]models.Channel, error) {
	if val, found := s.cache.Get(allChannelsCacheKey); found {
		if channels, ok := val.([]models.Channel); ok {
			return channels, nil
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Double-check após adquirir o lock
	if val, found := s.cache.Get(allChannelsCacheKey); found {
		if channels, ok := val.([]models.Channel); ok {
			return channels, nil
		}
	}

	channels, err := s.fetchRemoteChannels(ctx)
	if err != nil {
		// Fallback: se houver dados expirados no cache, reutiliza para não derrubar o usuário
		log.Printf("[WARN] Falha ao atualizar canais da SuperFlix: %v", err)
		if val, found := s.cache.Get(allChannelsCacheKey); found {
			if cachedChannels, ok := val.([]models.Channel); ok {
				return cachedChannels, nil
			}
		}
		return nil, err
	}

	s.cache.Set(allChannelsCacheKey, channels, s.ttl)
	return channels, nil
}

// fetchRemoteChannels executa a requisição HTTP para a SuperFlixAPI.
func (s *channelService) fetchRemoteChannels(ctx context.Context) ([]models.Channel, error) {
	reqURL := fmt.Sprintf("%s/lista?category=canais&format=json", s.baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição de canais: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "LocalStreamHub/1.0")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro na conexão com SuperFlixAPI: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("status inesperado da SuperFlixAPI (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta de canais: %w", err)
	}

	// Verifica se a API retornou objeto de rate limit
	var rateLimitErr struct {
		Error      bool   `json:"error"`
		Message    string `json:"message"`
		RetryAfter int    `json:"retry_after"`
	}
	if err := json.Unmarshal(body, &rateLimitErr); err == nil && rateLimitErr.Error {
		return nil, fmt.Errorf("rate limit atingido na SuperFlixAPI: %s (retry_after: %ds)", rateLimitErr.Message, rateLimitErr.RetryAfter)
	}

	var response models.SuperflixChannelResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("erro ao decodificar JSON de canais: %w", err)
	}

	// Normaliza as URLs do embed caso venham com escaping ou relativas
	for i := range response.Data {
		if response.Data[i].EmbedURL == "" && response.Data[i].ID != "" {
			response.Data[i].EmbedURL = fmt.Sprintf("%s/canal/%s", s.baseURL, response.Data[i].ID)
		}
	}

	return response.Data, nil
}
