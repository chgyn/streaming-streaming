package models

// Channel representa a entidade de um canal de TV ao vivo da SuperFlixAPI.
type Channel struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	LogoURL     string `json:"logo_url"`
	EmbedURL    string `json:"embed_url"`
	Category    string `json:"category"`
	IsActive    bool   `json:"is_active"`
}

// SuperflixChannelResponse mapeia a resposta envelopada retornada pelo endpoint /lista da Superflix.
type SuperflixChannelResponse struct {
	Success bool      `json:"success"`
	Data    []Channel `json:"data"`
	Total   int       `json:"total"`
}

// ChannelCategoryItem representa uma categoria de canais.
type ChannelCategoryItem struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}
