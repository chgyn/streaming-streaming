package handlers

import (
	"net/http"

	"local-stream-hub/internal/services"

	"github.com/gin-gonic/gin"
)

// ChannelHandler processa as requisições REST para catálogo e busca de canais ao vivo.
type ChannelHandler struct {
	service services.ChannelService
}

// NewChannelHandler instancia um novo ChannelHandler.
func NewChannelHandler(service services.ChannelService) *ChannelHandler {
	return &ChannelHandler{
		service: service,
	}
}

// GetChannels retorna a lista de canais ao vivo, com suporte a filtros por categoria e busca textual.
// GET /api/channels?category={slug}&q={termo}
func (h *ChannelHandler) GetChannels(c *gin.Context) {
	category := c.Query("category")
	query := c.Query("q")

	channels, err := h.service.GetChannels(c.Request.Context(), category, query)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error":   "falha ao carregar canais ao vivo",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": len(channels),
		"data":  channels,
	})
}

// GetCategories retorna todas as categorias únicas de canais para filtros na interface.
// GET /api/channels/categories
func (h *ChannelHandler) GetCategories(c *gin.Context) {
	categories, err := h.service.GetCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error":   "falha ao carregar categorias de canais",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": len(categories),
		"data":  categories,
	})
}
