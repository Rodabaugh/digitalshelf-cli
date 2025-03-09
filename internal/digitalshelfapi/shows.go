package digitalshelfapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (session *Session) LookupShowBarcode(args ...string) (Show, error) {
	if len(args) < 1 {
		return Show{}, fmt.Errorf("please provide a barcode")
	}

	err := validateLoggedIn(session)
	if err != nil {
		return Show{}, err
	}

	barcode := args[0]

	url := session.BaseURL + "search/show_barcodes/" + barcode

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Show{}, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+session.Token)

	res, err := session.DSAPIClient.HttpClient.Do(req)
	if err != nil {
		return Show{}, fmt.Errorf("error making request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return Show{}, fmt.Errorf("show not found")
	}

	if res.StatusCode != http.StatusOK {
		return Show{}, fmt.Errorf("error looking up show: %v", res.Status)
	}

	var show Show
	if err := json.NewDecoder(res.Body).Decode(&show); err != nil {
		return Show{}, fmt.Errorf("error decoding response: %v", err)
	}
	fmt.Println("Show found")
	fmt.Printf("Title: %s\n", show.Title)
	fmt.Printf("Season: %s\n", show.Season)
	fmt.Printf("Genre: %s\n", show.Genre)
	fmt.Printf("Actors: %s\n", show.Actors)
	fmt.Printf("Writer: %s\n", show.Writer)
	fmt.Printf("Director: %s\n", show.Director)
	fmt.Printf("Format: %s\n", show.Format)
	fmt.Printf("Release Date: %s\n", show.ReleaseDate)

	return show, nil
}

func (session *Session) AddShow(shelfID uuid.UUID, show Show) error {
	err := validateLoggedIn(session)
	if err != nil {
		return err
	}

	type parameters struct {
		Title       string    `json:"title"`
		Season      string    `json:"season"`
		Genre       string    `json:"genre"`
		Actors      string    `json:"actors"`
		Writer      string    `json:"writer"`
		Director    string    `json:"director"`
		Barcode     string    `json:"barcode"`
		Format      string    `json:"format"`
		ShelfID     uuid.UUID `json:"shelf_id"`
		ReleaseDate time.Time `json:"release_date"`
	}

	params := parameters{
		Title:       show.Title,
		Season:      show.Season,
		Genre:       show.Genre,
		Actors:      show.Actors,
		Writer:      show.Writer,
		Director:    show.Director,
		Barcode:     show.Barcode,
		Format:      show.Format,
		ShelfID:     shelfID,
		ReleaseDate: show.ReleaseDate,
	}

	url := session.BaseURL + "shows"

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
		fmt.Println("Show added successfully")
		return nil
	}

	return fmt.Errorf("error adding show: %v", res.Status)
}

func (session *Session) GetShows(args ...string) error {
	err := validateLoggedIn(session)
	if err != nil {
		return err
	}

	if len(args) < 1 {
		return fmt.Errorf("please specify a shelf ID")
	}

	shelfID := args[0]

	url := session.BaseURL + "shelves/" + shelfID + "/shows"

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
		return fmt.Errorf("error getting shows: %v", res.Status)
	}

	var shows []Show
	if err := json.NewDecoder(res.Body).Decode(&shows); err != nil {
		return fmt.Errorf("error decoding response: %v", err)
	}

	for _, show := range shows {
		fmt.Printf("Title: %s\n", show.Title)
		fmt.Printf("Season: %s\n", show.Season)
		fmt.Printf("Genre: %s\n", show.Genre)
		fmt.Printf("Actors: %s\n", show.Actors)
		fmt.Printf("Writer: %s\n", show.Writer)
		fmt.Printf("Director: %s\n", show.Director)
		fmt.Printf("Format: %s\n", show.Format)
		fmt.Printf("Release Date: %s\n", show.ReleaseDate)
		fmt.Println()
	}

	return nil
}

func (session *Session) GetAllLocationShows(args ...string) error {
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

	url := session.BaseURL + "locations/" + locationID.String() + "/shows"

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
		return fmt.Errorf("error getting shows: %v", res.Status)
	}

	var shows []Show
	if err := json.NewDecoder(res.Body).Decode(&shows); err != nil {
		return fmt.Errorf("error decoding response: %v", err)
	}

	for _, show := range shows {
		fmt.Printf("ID: %s\n", show.ID)
		fmt.Printf("Title: %s\n", show.Title)
		fmt.Printf("Season: %s\n", show.Season)
		fmt.Printf("Format: %s\n", show.Format)
		fmt.Printf("Release Date: %s\n", show.ReleaseDate)
		fmt.Println()
	}
	return nil
}

func (session *Session) GetShow(args ...string) error {
	err := validateLoggedIn(session)
	if err != nil {
		return err
	}

	if len(args) < 1 {
		return fmt.Errorf("please specify a show ID")
	}

	showUUID, err := uuid.Parse(args[0])
	if err != nil {
		return fmt.Errorf("invalid show ID: %v", err)
	}

	url := session.BaseURL + "shows/" + showUUID.String()

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
		return fmt.Errorf("error getting show: %v", res.Status)
	}

	var show Show
	if err := json.NewDecoder(res.Body).Decode(&show); err != nil {
		return fmt.Errorf("error decoding response: %v", err)
	}

	fmt.Printf("Title: %s\n", show.Title)
	fmt.Printf("Season: %s\n", show.Season)
	fmt.Printf("Genre: %s\n", show.Genre)
	fmt.Printf("Actors: %s\n", show.Actors)
	fmt.Printf("Writer: %s\n", show.Writer)
	fmt.Printf("Director: %s\n", show.Director)
	fmt.Printf("Release Date: %s\n", show.ReleaseDate)
	fmt.Printf("Barcode: %s\n", show.Barcode)
	fmt.Printf("Format: %s\n", show.Format)

	return nil
}

func (session *Session) SearchShows(args ...string) error {
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

	url := session.BaseURL + "search/shows"

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
		return fmt.Errorf("error searching shows: %v", res.Status)
	}

	var shows []Show
	if err := json.NewDecoder(res.Body).Decode(&shows); err != nil {
		return fmt.Errorf("error decoding response: %v", err)
	}

	for _, show := range shows {
		fmt.Printf("Title: %s\n", show.Title)
		fmt.Printf("Season: %s\n", show.Season)
		fmt.Printf("Genre: %s\n", show.Genre)
		fmt.Printf("Actors: %s\n", show.Actors)
		fmt.Printf("Writer: %s\n", show.Writer)
		fmt.Printf("Director: %s\n", show.Director)
		fmt.Printf("Format: %s\n", show.Format)
		fmt.Printf("Release Date: %s\n", show.ReleaseDate)
		fmt.Println()
	}

	return nil
}
