package services_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"local-stream-hub/internal/models"
	"local-stream-hub/internal/services"
)

func TestChannelService_GetChannelsAndCategories(t *testing.T) {
	mockData := models.SuperflixChannelResponse{
		Success: true,
		Total:   3,
		Data: []models.Channel{
			{
				ID:       "ae",
				Name:     "A&E",
				LogoURL:  "https://superflixapi.quest/img/ae.png",
				Category: "Variedades",
				IsActive: true,
			},
			{
				ID:       "sportv",
				Name:     "SporTV",
				LogoURL:  "https://superflixapi.quest/img/sportv.png",
				Category: "Esportes",
				IsActive: true,
			},
			{
				ID:       "espn",
				Name:     "ESPN Brasil",
				LogoURL:  "https://superflixapi.quest/img/espn.png",
				Category: "Esportes",
				IsActive: true,
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/lista" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mockData)
	}))
	defer server.Close()

	svc := services.NewChannelService(server.URL, 5*time.Minute)
	ctx := context.Background()

	// 1. Obter todos os canais
	channels, err := svc.GetChannels(ctx, "", "")
	if err != nil {
		t.Fatalf("erro inesperado ao buscar canais: %v", err)
	}
	if len(channels) != 3 {
		t.Fatalf("esperava 3 canais, obteve %d", len(channels))
	}

	// 2. Filtrar por categoria "Esportes"
	sports, err := svc.GetChannels(ctx, "Esportes", "")
	if err != nil {
		t.Fatalf("erro ao filtrar por categoria: %v", err)
	}
	if len(sports) != 2 {
		t.Fatalf("esperava 2 canais de esportes, obteve %d", len(sports))
	}

	// 3. Filtrar por busca de texto "A&E"
	aeChan, err := svc.GetChannels(ctx, "", "a&e")
	if err != nil {
		t.Fatalf("erro ao buscar texto: %v", err)
	}
	if len(aeChan) != 1 || aeChan[0].ID != "ae" {
		t.Fatalf("esperava canal A&E, obteve: %+v", aeChan)
	}

	// 4. Obter categorias únicas ordenadas
	categories, err := svc.GetCategories(ctx)
	if err != nil {
		t.Fatalf("erro ao buscar categorias: %v", err)
	}
	if len(categories) != 2 {
		t.Fatalf("esperava 2 categorias distintas ('Esportes', 'Variedades'), obteve %d", len(categories))
	}
	if categories[0] != "Esportes" || categories[1] != "Variedades" {
		t.Fatalf("ordem de categorias inesperada: %v", categories)
	}
}

func TestChannelService_ResilienceAndRateLimit(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts == 1 {
			mockData := models.SuperflixChannelResponse{
				Success: true,
				Total:   1,
				Data: []models.Channel{
					{ID: "globo", Name: "TV Globo", Category: "Canais Abertos", IsActive: true},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(mockData)
			return
		}
		// Segunda requisição simula rate limit
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"error":true,"message":"Muitas requisições","retry_after":10}`))
	}))
	defer server.Close()

	svc := services.NewChannelService(server.URL, 5*time.Minute)
	ctx := context.Background()

	// Primeira chamada popula o cache
	channels, err := svc.GetChannels(ctx, "", "")
	if err != nil || len(channels) != 1 {
		t.Fatalf("primeira chamada falhou: %v, total: %d", err, len(channels))
	}

	// Segunda chamada deve usar o cache in-memory e não falhar
	channels2, err := svc.GetChannels(ctx, "", "")
	if err != nil || len(channels2) != 1 {
		t.Fatalf("segunda chamada em cache falhou: %v", err)
	}
}
