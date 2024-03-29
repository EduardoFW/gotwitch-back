package twitch

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	"api.gotwitch.tk/models"
	"api.gotwitch.tk/settings"
	"github.com/google/go-querystring/query"
)

// Rate limit
type RateLimit struct {
	Limit     int       `json:"limit"`
	Remaining int       `json:"remaining"`
	Reset     time.Time `json:"reset"`
}

// Services controller instance
type TwitchService struct {
	tokenManager TokenManager `inject:""`
	rateLimit    *RateLimit
}

type RateLimitError struct {
	When time.Time
}

func (e *RateLimitError) Error() string {
	return "Rate limit exceeded, please wait until " + e.When.String()
}

func (t *TwitchService) GetTwitchToken() (string, error) {
	return t.tokenManager.GetToken()
}

func (t *TwitchService) processHeaders(res *http.Header) {
	t.rateLimit = &RateLimit{}
	if res.Get("RateLimit-Limit") != "" {
		t.rateLimit.Limit, _ = strconv.Atoi(res.Get("RateLimit-Limit"))
	}
	if res.Get("RateLimit-Remaining") != "" {
		t.rateLimit.Remaining, _ = strconv.Atoi(res.Get("RateLimit-Remaining"))
	}
	if res.Get("RateLimit-Reset") != "" {
		timestamp, err := strconv.Atoi(res.Get("RateLimit-Reset"))
		if err != nil {
			return
		}
		t.rateLimit.Reset = time.Unix(int64(timestamp), 0)
	}
}

type GetStreamListParams struct {
	After     string   `url:"after"`
	Before    string   `url:"before"`
	First     int      `url:"first"`
	GameId    []string `url:"game_id"`
	Language  []string `url:"language"`
	UserId    []string `url:"user_id"`
	UserLogin []string `url:"user_login"`
}

type Pagination struct {
	Cursor string `json:"cursor"`
}

type StreamResponse struct {
	Data       []models.Stream `json:"data"`
	Pagination Pagination      `json:"pagination"`
}

func (t *TwitchService) GetStreamList(params *GetStreamListParams) (*StreamResponse, error) {
	var streamResponse StreamResponse

	canRequest, when := t.HasRateLimit()
	if !canRequest {
		return nil, &RateLimitError{When: *when}
	}

	if params.First == 0 {
		params.First = 100
	}

	values, _ := query.Values(params)

	url := "https://api.twitch.tv/helix/streams?" + values.Encode()

	res, err := t.executeRequest(url)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(res.Body).Decode(&streamResponse)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	return &streamResponse, nil
}

func (t *TwitchService) HasRateLimit() (bool, *time.Time) {
	if t.rateLimit == nil { // Possible first request so no rate limit set yet
		return true, nil
	}
	if t.rateLimit.Remaining == 0 && t.rateLimit.Reset.After(time.Now()) { // Rate limit exceeded and reset time is in the future
		return false, &t.rateLimit.Reset
	}
	return true, nil
}

type SearchCategoriesParams struct {
	Query string `url:"query"`
	After string `url:"after"`
	First int    `url:"first"`
}

type CategoryResponse struct {
	Data       []models.Category `json:"data"`
	Pagination Pagination        `json:"pagination"`
}

func (t *TwitchService) SearchCategories(params *SearchCategoriesParams) (*CategoryResponse, error) {
	var categoryResponse CategoryResponse

	canRequest, when := t.HasRateLimit()
	if !canRequest {
		return nil, &RateLimitError{When: *when}
	}

	if params.First == 0 {
		params.First = 100
	}

	values, _ := query.Values(params)

	url := "https://api.twitch.tv/helix/search/categories?" + values.Encode()

	res, err := t.executeRequest(url)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(res.Body).Decode(&categoryResponse)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	return &categoryResponse, nil
}

type ListGamesParams struct {
	After string `url:"after"`
	First int    `url:"first"`
}

func (t *TwitchService) ListGames(params *ListGamesParams) (*CategoryResponse, error) {
	var categoryResponse CategoryResponse

	canRequest, when := t.HasRateLimit()
	if !canRequest {
		return nil, &RateLimitError{When: *when}
	}

	if params.First == 0 {
		params.First = 100
	}

	values, _ := query.Values(params)

	url := "https://api.twitch.tv/helix/games/top?" + values.Encode()

	res, err := t.executeRequest(url)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(res.Body).Decode(&categoryResponse)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	return &categoryResponse, nil
}

func (t *TwitchService) executeRequest(url string) (*http.Response, error) {
	client := http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	token, err := t.GetTwitchToken()
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"Authorization": []string{"Bearer " + token},
		"Client-Id":     []string{settings.ServerSettings.TwitchClientID},
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	t.processHeaders(&res.Header)

	return res, nil
}

var singleInstanceTwitchService *TwitchService
var onceTwitchService sync.Once

func GetTwitchService() *TwitchService {
	onceTwitchService.Do(func() {
		singleInstanceTwitchService = &TwitchService{}
	})
	return singleInstanceTwitchService
}
