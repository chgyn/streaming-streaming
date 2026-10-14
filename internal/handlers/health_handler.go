package handlers

import (
	"net/http"
	"runtime"
	"time"

	"local-stream-hub/internal/services"

	"github.com/gin-gonic/gin"
)

var startTime = time.Now()

// HealthHandler fornece métricas operacionais, consumo de RAM e status do serviço.
type HealthHandler struct {
	tmdbService services.TMDBService
}

// NewHealthHandler instancia o manipulador de health check.
func NewHealthHandler(tmdbService services.TMDBService) *HealthHandler {
	return &HealthHandler{tmdbService: tmdbService}
}

// HealthCheck responde com o status do sistema, uptime e footprint de memória.
// GET /api/health
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	cacheStats := h.tmdbService.GetCacheStats()

	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "local-stream-hub",
		"uptime":  time.Since(startTime).String(),
		"memory": gin.H{
			"alloc_mb":       float64(m.Alloc) / 1024 / 1024,
			"total_alloc_mb": float64(m.TotalAlloc) / 1024 / 1024,
			"sys_mb":         float64(m.Sys) / 1024 / 1024,
			"num_gc":         m.NumGC,
		},
		"cache": cacheStats,
	})
}
