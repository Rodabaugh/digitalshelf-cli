package digitalshelfapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (session *Session) LookupBookBarcode(args ...string) (Book, error) {
	if len(args) < 1 {
		return Book{}, fmt.Errorf("please provide a barcode")
	}

	err := validateLoggedIn(session)
	if err != nil {
		return Book{}, err
	}

	barcode := args[0]

	url := session.BaseURL + "search/book_barcodes/" + barcode

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Book{}, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+session.Token)

	res, err := session.DSAPIClient.HttpClient.Do(req)
	if err != nil {
		return Book{}, fmt.Errorf("error making request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return Book{}, fmt.Errorf("book not found")
	}

	if res.StatusCode != http.StatusOK {
		return Book{}, fmt.Errorf("error looking up book: %v", res.Status)
	}

	var book Book
	if err := json.NewDecoder(res.Body).Decode(&book); err != nil {
		return Book{}, fmt.Errorf("error decoding response: %v", err)
	}
	fmt.Println("Book found")
	fmt.Printf("Title: %s\n", book.Title)
	fmt.Printf("Author: %s\n", book.Author)
	fmt.Printf("Genre: %s\n", book.Genre)
	fmt.Printf("Publication Date: %s\n", book.PublicationDate)

	return book, nil
}

func (session *Session) AddBook(shelfID uuid.UUID, book Book) error {
	err := validateLoggedIn(session)
	if err != nil {
		return err
	}

	type parameters struct {
		Title           string    `json:"title"`
		Author          string    `json:"author"`
		Genre           string    `json:"genre"`
		Barcode         string    `json:"barcode"`
		ShelfID         uuid.UUID `json:"shelf_id"`
		PublicationDate time.Time `json:"publication_date"`
	}

	params := parameters{
		Title:           book.Title,
		Author:          book.Author,
		Genre:           book.Genre,
		Barcode:         book.Barcode,
		ShelfID:         shelfID,
		PublicationDate: book.PublicationDate,
	}

	url := session.BaseURL + "books"

	reqBody, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("error marshalling request: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+session.Token)

	res, err := session.DSAPIClient.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusCreated {
		fmt.Println("Book added successfully")
		return nil
	}

	return fmt.Errorf("error adding book: %v", res.Status)
}

func (session *Session) GetBooks(args ...string) error {
	err := validateLoggedIn(session)
	if err != nil {
		return err
	}

	if len(args) < 1 {
		return fmt.Errorf("please specify a shelf ID")
	}

	shelfID := args[0]

	url := session.BaseURL + "shelves/" + shelfID + "/books"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+session.Token)

	res, err := session.DSAPIClient.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error getting books: %v", res.Status)
	}

	var books []Book
	if err := json.NewDecoder(res.Body).Decode(&books); err != nil {
		return fmt.Errorf("error decoding response: %v", err)
	}

	for _, book := range books {
		fmt.Printf("Title: %s\n", book.Title)
		fmt.Printf("Author: %s\n", book.Author)
		fmt.Printf("Genre: %s\n", book.Genre)
		fmt.Printf("Publication Date: %s\n", book.PublicationDate)
		fmt.Println()
	}

	return nil
}

func (session *Session) GetAllLocationBooks(args ...string) error {
	err := validateLoggedIn(session)
	if err != nil {
		return fmt.Errorf("you must be logged in to do that")
	}

	if len(args) < 1 {
		return fmt.Errorf("please specify a location ID")
	}

	locationID, err := uuid.Parse(args[0])
	if err != nil {
		return fmt.Errorf("invalid location ID: %v", err)
	}

	url := session.BaseURL + "locations/" + locationID.String() + "/books"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+session.Token)

	res, err := session.DSAPIClient.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error getting books: %v", res.Status)
	}

	var books []Book
	if err := json.NewDecoder(res.Body).Decode(&books); err != nil {
		return fmt.Errorf("error decoding response: %v", err)
	}

	for _, book := range books {
		fmt.Printf("ID: %s\n", book.ID)
		fmt.Printf("Title: %s\n", book.Title)
		fmt.Printf("Author: %s\n", book.Author)
		fmt.Printf("Publication Date: %s\n", book.PublicationDate)
		fmt.Println()
	}
	return nil
}

func (session *Session) GetBook(args ...string) error {
	err := validateLoggedIn(session)
	if err != nil {
		return err
	}

	if len(args) < 1 {
		return fmt.Errorf("please specify a book ID")
	}

	bookUUID, err := uuid.Parse(args[0])
	if err != nil {
		return fmt.Errorf("invalid book ID: %v", err)
	}

	url := session.BaseURL + "books/" + bookUUID.String()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+session.Token)

	res, err := session.DSAPIClient.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error getting book: %v", res.Status)
	}

	var book Book
	if err := json.NewDecoder(res.Body).Decode(&book); err != nil {
		return fmt.Errorf("error decoding response: %v", err)
	}

	fmt.Printf("Title: %s\n", book.Title)
	fmt.Printf("Author: %s\n", book.Author)
	fmt.Printf("Genre: %s\n", book.Genre)
	fmt.Printf("Publication Date: %s\n", book.PublicationDate)
	fmt.Printf("Barcode: %s\n", book.Barcode)

	return nil
}

func (session *Session) SearchBooks(args ...string) error {
	err := validateLoggedIn(session)
	if err != nil {
		return err
	}

	if len(args) < 1 {
		return fmt.Errorf("please provide a search query")
	}

	if session.CurrentLocation == uuid.Nil {
		return fmt.Errorf("no location set, please set a location first")
	}

	query := args[0]

	url := session.BaseURL + "search/books"

	reqBody, err := json.Marshal(map[string]string{
		"query":       query,
		"location_id": session.CurrentLocation.String(),
	})
	if err != nil {
		return fmt.Errorf("error marshalling request: %v", err)
	}

	req, err := http.NewRequest("GET", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+session.Token)

	res, err := session.DSAPIClient.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error searching books: %v", res.Status)
	}

	var books []Book
	if err := json.NewDecoder(res.Body).Decode(&books); err != nil {
		return fmt.Errorf("error decoding response: %v", err)
	}

	for _, book := range books {
		fmt.Printf("Title: %s\n", book.Title)
		fmt.Printf("Author: %s\n", book.Author)
		fmt.Printf("Genre: %s\n", book.Genre)
		fmt.Printf("Publication Date: %s\n", book.PublicationDate)
		fmt.Println()
	}

	return nil
}
