package models

// PlayerResolveRequest define os parâmetros para resolver a URL de streaming
type PlayerResolveRequest struct {
	Type      string `form:"type" binding:"required,oneof=movie tv channel"`
	TMDBID    int64  `form:"tmdb_id" binding:"omitempty"`
	ChannelID string `form:"channel_id" binding:"omitempty"`
	Season    int    `form:"season" binding:"omitempty,gte=0"`
	Episode   int    `form:"episode" binding:"omitempty,gte=0"`
}

// PlayerResolveResponse define a estrutura de resposta para o player
type PlayerResolveResponse struct {
	TMDBID    int64  `json:"tmdb_id,omitempty"`
	ChannelID string `json:"channel_id,omitempty"`
	Type      string `json:"type"`
	EmbedURL  string `json:"embed_url"`
}
