package digitalshelfapi

import (
	"net/http"
)

// Client -
type Client struct {
	HttpClient http.Client
}

type Session struct {
	DSAPIClient  Client
	Base_url     string
	User         User
	Token        string
	RefreshToken string
}

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
