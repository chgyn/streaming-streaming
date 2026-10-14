package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"local-stream-hub/internal/models"

	"github.com/patrickmn/go-cache"
)

var (
	ErrNotFound      = errors.New("recurso não encontrado no TMDB")
	ErrTMDBAPIFailed = errors.New("falha na comunicação com a API do TMDB")
	ErrMissingAPIKey = errors.New("chave TMDB_API_KEY não configurada no servidor")
)

// TMDBService define o contrato para operações de metadados de streaming.
type TMDBService interface {
	GetTrending(ctx context.Context, mediaType, timeWindow string, page int) (*models.TMDBPaginatedResponse, error)
	Search(ctx context.Context, query, mediaType string, page int) (*models.TMDBPaginatedResponse, error)
	GetDetails(ctx context.Context, mediaType string, id int64) (*models.TMDBDetailResponse, error)
	GetCacheStats() map[string]interface{}
}

type tmdbService struct {
	client  *http.Client
	cache   *cache.Cache
	apiKey  string
	baseURL string
}

// NewTMDBService instancia um novo serviço com client HTTP reutilizável e cache in-memory.
func NewTMDBService(baseURL, apiKey string, defaultTTL, cleanupInterval time.Duration) TMDBService {
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 20,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 5 * time.Second,
	}

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
	}

	c := cache.New(defaultTTL, cleanupInterval)

	return &tmdbService{
		client:  client,
		cache:   c,
		apiKey:  strings.TrimSpace(apiKey),
		baseURL: strings.TrimRight(baseURL, "/"),
	}
}

// GetTrending busca títulos em alta com cache TTL de 60 minutos.
func (s *tmdbService) GetTrending(ctx context.Context, mediaType, timeWindow string, page int) (*models.TMDBPaginatedResponse, error) {
	if mediaType == "" {
		mediaType = "movie"
	}
	if timeWindow == "" {
		timeWindow = "week"
	}
	if page < 1 {
		page = 1
	}

	cacheKey := fmt.Sprintf("trending:%s:%s:%d", mediaType, timeWindow, page)
	if cached, found := s.cache.Get(cacheKey); found {
		if res, ok := cached.(*models.TMDBPaginatedResponse); ok {
			return res, nil
		}
	}

	endpoint := fmt.Sprintf("%s/trending/%s/%s", s.baseURL, url.PathEscape(mediaType), url.PathEscape(timeWindow))
	queryParams := url.Values{}
	queryParams.Set("page", strconv.Itoa(page))
	queryParams.Set("language", "pt-BR")

	reqURL := fmt.Sprintf("%s?%s", endpoint, queryParams.Encode())

	var response models.TMDBPaginatedResponse
	if err := s.doGet(ctx, reqURL, &response); err != nil {
		return nil, fmt.Errorf("falha ao buscar títulos em alta: %w", err)
	}

	s.cache.Set(cacheKey, &response, 60*time.Minute)
	return &response, nil
}

// Search realiza busca por termo no TMDB com cache TTL de 30 minutos.
func (s *tmdbService) Search(ctx context.Context, query, mediaType string, page int) (*models.TMDBPaginatedResponse, error) {
	query = strings.TrimSpace(query)
	if len(query) < 2 {
		return nil, fmt.Errorf("termo de busca deve ter no mínimo 2 caracteres")
	}

	if mediaType == "" {
		mediaType = "multi"
	}
	if page < 1 {
		page = 1
	}

	cacheKey := fmt.Sprintf("search:%s:%s:%d", mediaType, query, page)
	if cached, found := s.cache.Get(cacheKey); found {
		if res, ok := cached.(*models.TMDBPaginatedResponse); ok {
			return res, nil
		}
	}

	var pathSegment string
	switch mediaType {
	case "movie":
		pathSegment = "search/movie"
	case "tv":
		pathSegment = "search/tv"
	default:
		pathSegment = "search/multi"
	}

	endpoint := fmt.Sprintf("%s/%s", s.baseURL, pathSegment)
	queryParams := url.Values{}
	queryParams.Set("query", query)
	queryParams.Set("page", strconv.Itoa(page))
	queryParams.Set("language", "pt-BR")
	queryParams.Set("include_adult", "false")

	reqURL := fmt.Sprintf("%s?%s", endpoint, queryParams.Encode())

	var response models.TMDBPaginatedResponse
	if err := s.doGet(ctx, reqURL, &response); err != nil {
		return nil, fmt.Errorf("falha ao realizar busca no TMDB: %w", err)
	}

	s.cache.Set(cacheKey, &response, 30*time.Minute)
	return &response, nil
}

// GetDetails busca detalhes completos de um filme ou série com cache TTL de 120 minutos.
func (s *tmdbService) GetDetails(ctx context.Context, mediaType string, id int64) (*models.TMDBDetailResponse, error) {
	if mediaType != "movie" && mediaType != "tv" {
		return nil, fmt.Errorf("media_type inválido, deve ser 'movie' ou 'tv'")
	}
	if id <= 0 {
		return nil, fmt.Errorf("id inválido para consulta TMDB")
	}

	cacheKey := fmt.Sprintf("details:%s:%d", mediaType, id)
	if cached, found := s.cache.Get(cacheKey); found {
		if res, ok := cached.(*models.TMDBDetailResponse); ok {
			return res, nil
		}
	}

	endpoint := fmt.Sprintf("%s/%s/%d", s.baseURL, mediaType, id)
	queryParams := url.Values{}
	queryParams.Set("language", "pt-BR")

	reqURL := fmt.Sprintf("%s?%s", endpoint, queryParams.Encode())

	var response models.TMDBDetailResponse
	if err := s.doGet(ctx, reqURL, &response); err != nil {
		return nil, fmt.Errorf("falha ao obter detalhes do título %d: %w", id, err)
	}

	if response.MediaType == "" {
		response.MediaType = mediaType
	}

	s.cache.Set(cacheKey, &response, 120*time.Minute)
	return &response, nil
}

// GetCacheStats retorna informações de monitoramento do cache in-memory.
func (s *tmdbService) GetCacheStats() map[string]interface{} {
	return map[string]interface{}{
		"item_count": s.cache.ItemCount(),
	}
}

// doGet executa a requisição HTTP com timeout de contexto, headers de auth e decodificação JSON.
func (s *tmdbService) doGet(ctx context.Context, rawURL string, target interface{}) error {
	if s.apiKey == "" {
		return ErrMissingAPIKey
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("url inválida: %w", err)
	}

	isBearerToken := strings.HasPrefix(s.apiKey, "ey") || len(s.apiKey) > 40
	if !isBearerToken && s.apiKey != "" {
		q := parsedURL.Query()
		q.Set("api_key", s.apiKey)
		parsedURL.RawQuery = q.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedURL.String(), nil)
	if err != nil {
		return fmt.Errorf("erro ao criar requisição HTTP: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	if isBearerToken && s.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+s.apiKey)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao conectar à API externa: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return ErrNotFound
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return fmt.Errorf("%w: status %d - %s", ErrTMDBAPIFailed, resp.StatusCode, string(bodyBytes))
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("erro ao decodificar JSON do TMDB: %w", err)
	}

	return nil
}
