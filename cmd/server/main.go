package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"local-stream-hub/internal/config"
	"local-stream-hub/internal/handlers"
	"local-stream-hub/internal/services"
	"local-stream-hub/web"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	if cfg.GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(corsMiddleware())

	// Inicialização de serviços e dependências
	cacheTTL := time.Duration(cfg.CacheTTLMinutes) * time.Minute
	cacheCleanup := time.Duration(cfg.CacheCleanupMinutes) * time.Minute
	tmdbService := services.NewTMDBService(cfg.TMDBBaseURL, cfg.TMDBAPIKey, cacheTTL, cacheCleanup)

	// Inicialização dos handlers
	tmdbHandler := handlers.NewTMDBHandler(tmdbService)
	playerHandler := handlers.NewPlayerHandler(cfg.SuperflixBaseURL)
	healthHandler := handlers.NewHealthHandler(tmdbService)

	// Rotas da API REST
	api := router.Group("/api")
	{
		api.GET("/trending", tmdbHandler.GetTrending)
		api.GET("/search", tmdbHandler.Search)
		api.GET("/details/:media_type/:id", tmdbHandler.GetDetails)
		api.GET("/player/resolve", playerHandler.Resolve)
		api.GET("/health", healthHandler.HealthCheck)
	}

	// Registro de SPA e fallback de arquivos embutidos
	handlers.RegisterSPARoutes(router, web.GetFileSystem())

	serverAddr := fmt.Sprintf(":%s", cfg.Port)
	srv := &http.Server{
		Addr:         serverAddr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Inicia o servidor HTTP em goroutine desacoplada
	go func() {
		log.Printf("[INFO] Local Stream Hub iniciado na porta %s (modo %s)", cfg.Port, cfg.GinMode)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("[FATAL] Erro ao iniciar servidor HTTP: %v", err)
		}
	}()

	// Graceful shutdown coordenado por sinais do sistema operacional
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[INFO] Encerrando servidor e drenando conexões ativas...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("[ERRO] Falha ao encerrar servidor de forma limpa: %v", err)
	}

	log.Println("[INFO] Servidor finalizado com sucesso.")
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
