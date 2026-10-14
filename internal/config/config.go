package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config armazena todas as variáveis de configuração da aplicação.
type Config struct {
	Port                string
	TMDBAPIKey          string
	TMDBBaseURL         string
	CacheTTLMinutes     int
	CacheCleanupMinutes int
	SuperflixBaseURL    string
	GinMode             string
}

// LoadConfig carrega as configurações das variáveis de ambiente e arquivo .env opcional.
func LoadConfig() *Config {
	// Carrega arquivo .env se existir (ignora erro se não existir em containers/produção)
	if err := godotenv.Load(); err != nil {
		log.Println("[INFO] Arquivo .env não encontrado, utilizando variáveis de ambiente do sistema")
	}

	port := getEnv("PORT", "8080")
	tmdbKey := getEnv("TMDB_API_KEY", "")
	if tmdbKey == "" {
		// Suporte alternativo a TMDB_TOKEN
		tmdbKey = getEnv("TMDB_TOKEN", "")
	}

	tmdbBaseURL := getEnv("TMDB_BASE_URL", "https://api.themoviedb.org/3")
	superflixBaseURL := getEnv("SUPERFLIX_BASE_URL", "https://superflixapi.quest")
	ginMode := getEnv("GIN_MODE", "release")

	cacheTTL := getEnvAsInt("CACHE_TTL_MINUTES", 60)
	cacheCleanup := getEnvAsInt("CACHE_CLEANUP_MINUTES", 10)

	return &Config{
		Port:                port,
		TMDBAPIKey:          tmdbKey,
		TMDBBaseURL:         tmdbBaseURL,
		CacheTTLMinutes:     cacheTTL,
		CacheCleanupMinutes: cacheCleanup,
		SuperflixBaseURL:    superflixBaseURL,
		GinMode:             ginMode,
	}
}

func getEnv(key, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists && val != "" {
		return val
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	valStr := getEnv(key, "")
	if valStr == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultVal
	}
	return val
}
