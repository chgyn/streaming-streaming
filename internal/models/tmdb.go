package models

// MediaItem representa um item resumido de filme ou série
type MediaItem struct {
	ID            int64    `json:"id"`
	Title         string   `json:"title,omitempty"`
	Name          string   `json:"name,omitempty"` // Séries no TMDB utilizam "name"
	OriginalTitle string   `json:"original_title,omitempty"`
	OriginalName  string   `json:"original_name,omitempty"`
	Overview      string   `json:"overview"`
	PosterPath    string   `json:"poster_path"`
	BackdropPath  string   `json:"backdrop_path"`
	MediaType     string   `json:"media_type"`
	ReleaseDate   string   `json:"release_date,omitempty"`
	FirstAirDate  string   `json:"first_air_date,omitempty"`
	VoteAverage   float64  `json:"vote_average"`
	VoteCount     int64    `json:"vote_count"`
	Popularity    float64  `json:"popularity"`
	GenreIDs      []int64  `json:"genre_ids,omitempty"`
}

// TMDBPaginatedResponse padroniza listagens paginadas do TMDB
type TMDBPaginatedResponse struct {
	Page         int         `json:"page"`
	Results      []MediaItem `json:"results"`
	TotalPages   int         `json:"total_pages"`
	TotalResults int         `json:"total_results"`
}

// TMDBGenre descreve um gênero
type TMDBGenre struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// TMDBEpisode descreve um episódio em séries
type TMDBEpisode struct {
	ID            int64   `json:"id"`
	EpisodeNumber int     `json:"episode_number"`
	SeasonNumber  int     `json:"season_number"`
	Name          string  `json:"name"`
	Overview      string  `json:"overview"`
	StillPath     string  `json:"still_path"`
	VoteAverage   float64 `json:"vote_average"`
	AirDate       string  `json:"air_date,omitempty"`
}

// TMDBSeason descreve uma temporada
type TMDBSeason struct {
	ID           int64         `json:"id"`
	SeasonNumber int           `json:"season_number"`
	Name         string        `json:"name"`
	Overview     string        `json:"overview"`
	PosterPath   string        `json:"poster_path"`
	EpisodeCount int           `json:"episode_count"`
	AirDate      string        `json:"air_date,omitempty"`
	Episodes     []TMDBEpisode `json:"episodes,omitempty"`
}

// TMDBDetailResponse reúne os detalhes completos de um título
type TMDBDetailResponse struct {
	ID               int64        `json:"id"`
	Title            string       `json:"title,omitempty"`
	Name             string       `json:"name,omitempty"`
	Overview         string       `json:"overview"`
	PosterPath       string       `json:"poster_path"`
	BackdropPath     string       `json:"backdrop_path"`
	MediaType        string       `json:"media_type"`
	Genres           []TMDBGenre  `json:"genres"`
	VoteAverage      float64      `json:"vote_average"`
	VoteCount        int64        `json:"vote_count"`
	ReleaseDate      string       `json:"release_date,omitempty"`
	FirstAirDate     string       `json:"first_air_date,omitempty"`
	Status           string       `json:"status"`
	Tagline          string       `json:"tagline"`
	Runtime          int          `json:"runtime,omitempty"`
	NumberOfSeasons  int          `json:"number_of_seasons,omitempty"`
	NumberOfEpisodes int          `json:"number_of_episodes,omitempty"`
	Seasons          []TMDBSeason `json:"seasons,omitempty"`
}
