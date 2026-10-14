package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"local-stream-hub/internal/models"

	"github.com/gin-gonic/gin"
)

// PlayerHandler resolve as URLs normalizadas do iframe do player Superflix.
type PlayerHandler struct {
	superflixBaseURL string
}

// NewPlayerHandler instancia o handler do player.
func NewPlayerHandler(superflixBaseURL string) *PlayerHandler {
	return &PlayerHandler{
		superflixBaseURL: strings.TrimRight(superflixBaseURL, "/"),
	}
}

// Resolve gera a URL segura de incorporação para filmes ou episódios de séries.
// GET /api/player/resolve?type={movie|tv}&tmdb_id={id}&season={s}&episode={e}
func (h *PlayerHandler) Resolve(c *gin.Context) {
	var req models.PlayerResolveRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parâmetros inválidos para resolução do player", "details": err.Error()})
		return
	}

	var embedURL string
	switch req.Type {
	case "movie":
		embedURL = fmt.Sprintf("%s/filme/%d", h.superflixBaseURL, req.TMDBID)
	case "tv":
		if req.Season <= 0 || req.Episode <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "para séries ('tv'), 'season' e 'episode' devem ser inteiros maiores que zero"})
			return
		}
		embedURL = fmt.Sprintf("%s/serie/%d/%d/%d", h.superflixBaseURL, req.TMDBID, req.Season, req.Episode)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "tipo inválido, utilize 'movie' ou 'tv'"})
		return
	}

	c.JSON(http.StatusOK, models.PlayerResolveResponse{
		TMDBID:   req.TMDBID,
		Type:     req.Type,
		EmbedURL: embedURL,
	})
}
