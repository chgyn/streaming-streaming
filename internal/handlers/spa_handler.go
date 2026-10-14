package handlers

import (
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// RegisterSPARoutes registra o roteamento de arquivos estáticos embutidos e o fallback para SPA client-side routing.
func RegisterSPARoutes(r *gin.Engine, distFS fs.FS) {
	fileServer := http.FileServer(http.FS(distFS))

	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// Se a requisição for para a API e deu 404, não retorna o HTML da SPA
		if strings.HasPrefix(path, "/api/") || path == "/api" {
			c.JSON(http.StatusNotFound, gin.H{"error": "API route not found"})
			return
		}

		// Remove barra inicial para consulta ao FS
		cleanPath := strings.TrimPrefix(path, "/")
		if cleanPath == "" {
			cleanPath = "index.html"
		}

		// Checa se o arquivo existe fisicamente no subdiretório dist
		f, err := distFS.Open(cleanPath)
		if err == nil {
			_ = f.Close()
			fileServer.ServeHTTP(c.Writer, c.Request)
			return
		}

		// Fallback para SPA client-side routing
		indexData, err := fs.ReadFile(distFS, "index.html")
		if err != nil {
			c.String(http.StatusInternalServerError, "index.html not found in embedded fs")
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", indexData)
	})
}
