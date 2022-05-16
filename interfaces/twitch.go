package interfaces

type TwitchToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type Stream struct {
	Id           string   `json:"id"`
	UserId       string   `json:"user_id"`
	UserLogin    string   `json:"user_login"`
	UserName     string   `json:"user_name"`
	GameId       string   `json:"game_id"`
	GameName     string   `json:"game_name"`
	Type         string   `json:"type"`
	Title        string   `json:"title"`
	ViewerCount  int      `json:"viewer_count"`
	StartedAt    string   `json:"started_at"`
	Language     string   `json:"language"`
	ThumbnailUrl string   `json:"thumbnail_url"`
	TagIds       []string `json:"tag_ids"`
	IsMature     bool     `json:"is_mature"`
}

type Pagination struct {
	Cursor string `json:"cursor"`
}

type StreamResponse struct {
	Data       []Stream   `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type Category struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	BoxArtUrl string `json:"box_art_url"`
}

type CategoryResponse struct {
	Data       []Category `json:"data"`
	Pagination Pagination `json:"pagination"`
}
