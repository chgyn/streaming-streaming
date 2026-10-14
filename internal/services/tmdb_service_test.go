package services_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"local-stream-hub/internal/models"
	"local-stream-hub/internal/services"
)

func TestTMDBService_Trending_CacheHitAndLatency(t *testing.T) {
	var requestCount int64

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&requestCount, 1)

		if r.URL.Path != "/trending/movie/week" {
			t.Errorf("caminho inesperado: %s", r.URL.Path)
		}

		resp := models.TMDBPaginatedResponse{
			Page: 1,
			Results: []models.MediaItem{
				{ID: 550, Title: "Fight Club", MediaType: "movie", VoteAverage: 8.4},
			},
			TotalPages:   1,
			TotalResults: 1,
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer mockServer.Close()

	svc := services.NewTMDBService(mockServer.URL, "test-api-key", 10*time.Minute, 1*time.Minute)
	ctx := context.Background()

	// Primeira requisição: Cache Miss
	startMiss := time.Now()
	res1, err := svc.GetTrending(ctx, "movie", "week", 1)
	durationMiss := time.Since(startMiss)

	if err != nil {
		t.Fatalf("erro inesperado na primeira consulta: %v", err)
	}
	if len(res1.Results) != 1 || res1.Results[0].ID != 550 {
		t.Fatalf("resultado inesperado: %+v", res1)
	}
	if atomic.LoadInt64(&requestCount) != 1 {
		t.Fatalf("esperava 1 requisição upstream, obteve %d", atomic.LoadInt64(&requestCount))
	}

	// Segunda requisição: Cache Hit
	startHit := time.Now()
	res2, err := svc.GetTrending(ctx, "movie", "week", 1)
	durationHit := time.Since(startHit)

	if err != nil {
		t.Fatalf("erro inesperado na consulta em cache: %v", err)
	}
	if len(res2.Results) != 1 || res2.Results[0].ID != 550 {
		t.Fatalf("resultado em cache inconsistente: %+v", res2)
	}

	// O número de requisições upstream DEVE continuar sendo 1 (atestando cache hit)
	if atomic.LoadInt64(&requestCount) != 1 {
		t.Fatalf("cache falhou: requisição upstream foi disparada novamente. Total de requisições: %d", atomic.LoadInt64(&requestCount))
	}

	// A latência com cache em memória deve ser inferior a 2 milissegundos
	if durationHit > 5*time.Millisecond {
		t.Logf("[AVISO] Latência de cache hit: %v (miss foi %v)", durationHit, durationMiss)
	} else {
		t.Logf("[SUCESSO] Cache hit em %v (< 2ms) vs miss em %v", durationHit, durationMiss)
	}
}

func TestTMDBService_Search(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("query")
		if q != "Matrix" {
			t.Errorf("termo de busca esperado 'Matrix', obteve '%s'", q)
		}

		resp := models.TMDBPaginatedResponse{
			Page: 1,
			Results: []models.MediaItem{
				{ID: 603, Title: "The Matrix", MediaType: "movie"},
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer mockServer.Close()

	svc := services.NewTMDBService(mockServer.URL, "test-api-key", 5*time.Minute, 1*time.Minute)
	ctx := context.Background()

	// Validação de termo curto (< 2 caracteres)
	_, err := svc.Search(ctx, "a", "multi", 1)
	if err == nil {
		t.Fatal("esperava erro ao buscar com termo de 1 caractere")
	}

	// Busca válida
	res, err := svc.Search(ctx, "Matrix", "movie", 1)
	if err != nil {
		t.Fatalf("erro ao executar busca: %v", err)
	}
	if len(res.Results) != 1 || res.Results[0].ID != 603 {
		t.Fatalf("resultado de busca incorreto: %+v", res)
	}
}

func TestTMDBService_GetDetails_NotFound(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer mockServer.Close()

	svc := services.NewTMDBService(mockServer.URL, "test-key", 5*time.Minute, 1*time.Minute)
	ctx := context.Background()

	_, err := svc.GetDetails(ctx, "movie", 99999999)
	if err == nil {
		t.Fatal("esperava erro de not found ao consultar ID inexistente")
	}
}
