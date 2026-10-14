package web

import (
	"embed"
	"io/fs"
)

//go:embed dist/*
var DistFS embed.FS

// GetFileSystem retorna a subárvore do diretório dist embutido.
func GetFileSystem() fs.FS {
	sub, err := fs.Sub(DistFS, "dist")
	if err != nil {
		panic(err)
	}
	return sub
}
