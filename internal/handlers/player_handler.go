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

// Resolve gera a URL segura de incorporação para filmes, episódios de séries ou canais ao vivo.
// GET /api/player/resolve?type={movie|tv|channel}&tmdb_id={id}&channel_id={id}&season={s}&episode={e}
func (h *PlayerHandler) Resolve(c *gin.Context) {
	var req models.PlayerResolveRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parâmetros inválidos para resolução do player", "details": err.Error()})
		return
	}

	var embedURL string
	switch req.Type {
	case "movie":
		if req.TMDBID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "para filmes ('movie'), 'tmdb_id' deve ser maior que zero"})
			return
		}
		embedURL = fmt.Sprintf("%s/filme/%d", h.superflixBaseURL, req.TMDBID)
	case "tv":
		if req.TMDBID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "para séries ('tv'), 'tmdb_id' deve ser maior que zero"})
			return
		}
		if req.Season <= 0 || req.Episode <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "para séries ('tv'), 'season' e 'episode' devem ser inteiros maiores que zero"})
			return
		}
		embedURL = fmt.Sprintf("%s/serie/%d/%d/%d", h.superflixBaseURL, req.TMDBID, req.Season, req.Episode)
	case "channel":
		if strings.TrimSpace(req.ChannelID) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "para canais ('channel'), 'channel_id' é obrigatório"})
			return
		}
		embedURL = fmt.Sprintf("%s/canal/%s", h.superflixBaseURL, strings.TrimSpace(req.ChannelID))
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "tipo inválido, utilize 'movie', 'tv' ou 'channel'"})
		return
	}

	c.JSON(http.StatusOK, models.PlayerResolveResponse{
		TMDBID:    req.TMDBID,
		ChannelID: req.ChannelID,
		Type:      req.Type,
		EmbedURL:  embedURL,
	})
}
