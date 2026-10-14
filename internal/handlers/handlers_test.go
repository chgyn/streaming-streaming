package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"local-stream-hub/internal/handlers"
	"local-stream-hub/internal/models"
	"local-stream-hub/internal/services"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupTestRouter() (*gin.Engine, *handlers.PlayerHandler) {
	router := gin.New()

	playerHandler := handlers.NewPlayerHandler("https://superflixapi.quest")
	tmdbService := services.NewTMDBService("https://api.themoviedb.org/3", "dummy", 10*time.Minute, 1*time.Minute)
	tmdbHandler := handlers.NewTMDBHandler(tmdbService)
	healthHandler := handlers.NewHealthHandler(tmdbService)

	api := router.Group("/api")
	{
		api.GET("/player/resolve", playerHandler.Resolve)
		api.GET("/trending", tmdbHandler.GetTrending)
		api.GET("/health", healthHandler.HealthCheck)
	}

	mockFS := fstest.MapFS{
		"index.html": &fstest.MapFile{
			Data: []byte("<!DOCTYPE html><html><body>SPA Hub</body></html>"),
		},
		"assets/style.css": &fstest.MapFile{
			Data: []byte("body { background: #000; }"),
		},
	}

	handlers.RegisterSPARoutes(router, mockFS)

	return router, playerHandler
}

func TestPlayerResolve_Movie(t *testing.T) {
	router, _ := setupTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/player/resolve?type=movie&tmdb_id=550", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("esperava status 200, obteve %d. Body: %s", w.Code, w.Body.String())
	}

	var resp models.PlayerResolveResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("erro ao deserializar resposta: %v", err)
	}

	expectedURL := "https://superflixapi.quest/filme/550"
	if resp.EmbedURL != expectedURL {
		t.Fatalf("embed_url esperada: %s, obtida: %s", expectedURL, resp.EmbedURL)
	}
	if resp.Type != "movie" || resp.TMDBID != 550 {
		t.Fatalf("dados divergentes no response: %+v", resp)
	}
}

func TestPlayerResolve_TV(t *testing.T) {
	router, _ := setupTestRouter()

	// Caso de sucesso para série com temporada e episódio
	req, _ := http.NewRequest(http.MethodGet, "/api/player/resolve?type=tv&tmdb_id=1399&season=1&episode=3", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("esperava status 200, obteve %d. Body: %s", w.Code, w.Body.String())
	}

	var resp models.PlayerResolveResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("erro ao deserializar resposta: %v", err)
	}

	expectedURL := "https://superflixapi.quest/serie/1399/1/3"
	if resp.EmbedURL != expectedURL {
		t.Fatalf("embed_url esperada: %s, obtida: %s", expectedURL, resp.EmbedURL)
	}

	// Caso com falta de season e episode para série (deve falhar com 400 Bad Request)
	reqFail, _ := http.NewRequest(http.MethodGet, "/api/player/resolve?type=tv&tmdb_id=1399", nil)
	wFail := httptest.NewRecorder()
	router.ServeHTTP(wFail, reqFail)

	if wFail.Code != http.StatusBadRequest {
		t.Fatalf("esperava status 400 para série sem temporada/episódio, obteve %d", wFail.Code)
	}
}

func TestSPAFallback_ClientRouting(t *testing.T) {
	router, _ := setupTestRouter()

	// 1. Rota de cliente SPA aleatória deve retornar 200 com index.html
	req, _ := http.NewRequest(http.MethodGet, "/filme/550", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("esperava status 200 para rota SPA, obteve %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "SPA Hub") {
		t.Fatalf("esperava conteúdo de index.html no fallback, obteve: %s", w.Body.String())
	}

	// 2. Arquivo estático real deve ser servido
	reqAsset, _ := http.NewRequest(http.MethodGet, "/assets/style.css", nil)
	wAsset := httptest.NewRecorder()
	router.ServeHTTP(wAsset, reqAsset)

	if wAsset.Code != http.StatusOK {
		t.Fatalf("esperava status 200 para asset estático, obteve %d", wAsset.Code)
	}
	if !strings.Contains(wAsset.Body.String(), "background: #000") {
		t.Fatalf("conteúdo do asset css incorreto: %s", wAsset.Body.String())
	}
}

func TestAPIRouteNotFound_NeverReturnsHTML(t *testing.T) {
	router, _ := setupTestRouter()

	// Rota inexistente sob /api/* DEVE retornar status 404 em JSON, NUNCA o HTML da SPA
	req, _ := http.NewRequest(http.MethodGet, "/api/rota-que-nao-existe", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("esperava status 404 para rota /api inexistente, obteve %d", w.Code)
	}

	var jsonResp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &jsonResp); err != nil {
		t.Fatalf("resposta para /api 404 não é um JSON válido: %s", w.Body.String())
	}

	if jsonResp["error"] != "API route not found" {
		t.Fatalf("mensagem de erro esperada 'API route not found', obteve: %v", jsonResp["error"])
	}
}

func TestHealthCheck(t *testing.T) {
	router, _ := setupTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("esperava status 200 para health check, obteve %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("falha ao ler JSON de saúde: %v", err)
	}

	if resp["status"] != "healthy" {
		t.Fatalf("status de saúde inesperado: %v", resp["status"])
	}
}
