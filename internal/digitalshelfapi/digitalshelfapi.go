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
	DSAPIClient     Client
	Base_url        string
	User            User
	Token           string
	RefreshToken    string
	CurrentLocation uuid.UUID
}

type User struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

type Case struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	LocationID uuid.UUID `json:"location_id"`
}

type Shelf struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	CaseID uuid.UUID `json:"case_id"`
}
