package digitalshelfapi

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Client -
type Client struct {
	HttpClient http.Client
}

type Session struct {
	DSAPIClient     Client
	Platform        string
	BaseURL         string
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

type Movie struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Genre       string    `json:"genre"`
	Actors      string    `json:"actors"`
	Writer      string    `json:"writer"`
	Director    string    `json:"director"`
	Barcode     string    `json:"barcode"`
	ShelfID     uuid.UUID `json:"shelf_id"`
	ReleaseDate time.Time `json:"release_date"`
}

type Show struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Season      int       `json:"season"`
	Genre       string    `json:"genre"`
	Actors      string    `json:"actors"`
	Writer      string    `json:"writer"`
	Director    string    `json:"director"`
	Barcode     string    `json:"barcode"`
	ShelfID     uuid.UUID `json:"shelf_id"`
	ReleaseDate time.Time `json:"release_date"`
}

type Book struct {
	ID              uuid.UUID `json:"id"`
	Title           string    `json:"title"`
	Author          string    `json:"author"`
	Genre           string    `json:"genre"`
	Barcode         string    `json:"barcode"`
	ShelfID         uuid.UUID `json:"shelf_id"`
	PublicationDate time.Time `json:"publication_date"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Music struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Artist      string    `json:"artist"`
	Genre       string    `json:"genre"`
	Barcode     string    `json:"barcode"`
	Format      string    `json:"format"`
	ShelfID     uuid.UUID `json:"shelf_id"`
	ReleaseDAte time.Time `json:"release_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
