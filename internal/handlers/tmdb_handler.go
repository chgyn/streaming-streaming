package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"local-stream-hub/internal/services"

	"github.com/gin-gonic/gin"
)

// TMDBHandler lida com as rotas de catálogo e metadados do TMDB.
type TMDBHandler struct {
	service services.TMDBService
}

// NewTMDBHandler cria uma nova instância do handler do TMDB.
func NewTMDBHandler(service services.TMDBService) *TMDBHandler {
	return &TMDBHandler{service: service}
}

// GetTrending retorna os títulos em alta (filmes ou séries).
// GET /api/trending?media_type={movie|tv}&time_window={day|week}&page={int}
func (h *TMDBHandler) GetTrending(c *gin.Context) {
	mediaType := c.DefaultQuery("media_type", "movie")
	timeWindow := c.DefaultQuery("time_window", "week")
	pageStr := c.DefaultQuery("page", "1")

	if mediaType != "movie" && mediaType != "tv" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parâmetro 'media_type' deve ser 'movie' ou 'tv'"})
		return
	}

	if timeWindow != "day" && timeWindow != "week" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parâmetro 'time_window' deve ser 'day' ou 'week'"})
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	res, err := h.service.GetTrending(c.Request.Context(), mediaType, timeWindow, page)
	if err != nil {
		h.handleTMDBError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

// Search pesquisa filmes e séries no TMDB por termo.
// GET /api/search?q={query}&type={multi|movie|tv}&page={int}
func (h *TMDBHandler) Search(c *gin.Context) {
	query := strings.TrimSpace(c.Query("q"))
	if len(query) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "o parâmetro de busca 'q' deve ter no mínimo 2 caracteres"})
		return
	}

	searchType := c.DefaultQuery("type", "multi")
	if searchType != "multi" && searchType != "movie" && searchType != "tv" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "o parâmetro 'type' deve ser 'multi', 'movie' ou 'tv'"})
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	res, err := h.service.Search(c.Request.Context(), query, searchType, page)
	if err != nil {
		h.handleTMDBError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

// GetDetails obtém detalhes completos de um filme ou série por ID.
// GET /api/details/:media_type/:id
func (h *TMDBHandler) GetDetails(c *gin.Context) {
	mediaType := c.Param("media_type")
	if mediaType != "movie" && mediaType != "tv" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parâmetro de rota 'media_type' deve ser 'movie' ou 'tv'"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido informado na rota"})
		return
	}

	res, err := h.service.GetDetails(c.Request.Context(), mediaType, id)
	if err != nil {
		h.handleTMDBError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *TMDBHandler) handleTMDBError(c *gin.Context, err error) {
	if errors.Is(err, services.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "título ou recurso não encontrado"})
		return
	}
	if errors.Is(err, services.ErrMissingAPIKey) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "TMDB_API_KEY não configurada",
			"details": "Defina a chave TMDB_API_KEY no arquivo .env para consultar o catálogo do TMDB.",
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao processar requisição do TMDB", "details": err.Error()})
}
