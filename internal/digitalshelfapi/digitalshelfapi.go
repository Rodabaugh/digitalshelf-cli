package digitalshelfapi

import (
	"net/http"

	"github.com/google/uuid"
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
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}
