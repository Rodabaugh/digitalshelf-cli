package digitalshelfapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (session *Session) LookupMusicBarcode(args ...string) (Music, error) {
	if len(args) < 1 {
		return Music{}, fmt.Errorf("please provide a barcode")
	}

	err := validateLoggedIn(session)
	if err != nil {
		return Music{}, err
	}

	barcode := args[0]

	url := session.BaseURL + "search/music_barcodes/" + barcode

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Music{}, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+session.Token)

	res, err := session.DSAPIClient.HttpClient.Do(req)
	if err != nil {
		return Music{}, fmt.Errorf("error making request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return Music{}, fmt.Errorf("music not found")
	}

	if res.StatusCode != http.StatusOK {
		return Music{}, fmt.Errorf("error looking up music: %v", res.Status)
	}

	var music Music
	if err := json.NewDecoder(res.Body).Decode(&music); err != nil {
		return Music{}, fmt.Errorf("error decoding response: %v", err)
	}
	fmt.Println("Music found")
	fmt.Printf("Title: %s\n", music.Title)
	fmt.Printf("Artist: %s\n", music.Artist)
	fmt.Printf("Genre: %s\n", music.Genre)
	fmt.Printf("Format: %s\n", music.Format)
	fmt.Printf("Release Date: %s\n", music.ReleaseDate)

	return music, nil
}

func (session *Session) AddMusic(shelfID uuid.UUID, music Music) error {
	err := validateLoggedIn(session)
	if err != nil {
		return err
	}

	type parameters struct {
		Title       string    `json:"title"`
		Artist      string    `json:"artist"`
		Genre       string    `json:"genre"`
		Barcode     string    `json:"barcode"`
		Format      string    `json:"format"`
		ShelfID     uuid.UUID `json:"shelf_id"`
		ReleaseDate time.Time `json:"release_date"`
	}

	params := parameters{
		Title:       music.Title,
		Artist:      music.Artist,
		Genre:       music.Genre,
		Barcode:     music.Barcode,
		Format:      music.Format,
		ShelfID:     shelfID,
		ReleaseDate: music.ReleaseDate,
	}

	url := session.BaseURL + "music"

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
		fmt.Println("Music added successfully")
		return nil
	}

	return fmt.Errorf("error adding music: %v", res.Status)
}

func (session *Session) GetMusic(args ...string) error {
	err := validateLoggedIn(session)
	if err != nil {
		return err
	}

	if len(args) < 1 {
		return fmt.Errorf("please specify a shelf ID")
	}

	shelfID := args[0]

	url := session.BaseURL + "shelves/" + shelfID + "/music"

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
		return fmt.Errorf("error getting music: %v", res.Status)
	}

	var musicList []Music
	if err := json.NewDecoder(res.Body).Decode(&musicList); err != nil {
		return fmt.Errorf("error decoding response: %v", err)
	}

	for _, music := range musicList {
		fmt.Printf("Title: %s\n", music.Title)
		fmt.Printf("Artist: %s\n", music.Artist)
		fmt.Printf("Genre: %s\n", music.Genre)
		fmt.Printf("Format: %s\n", music.Format)
		fmt.Printf("Release Date: %s\n", music.ReleaseDate)
		fmt.Println()
	}

	return nil
}

func (session *Session) GetAllLocationMusic(args ...string) error {
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

	url := session.BaseURL + "locations/" + locationID.String() + "/music"

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
		return fmt.Errorf("error getting music: %v", res.Status)
	}

	var musicList []Music
	if err := json.NewDecoder(res.Body).Decode(&musicList); err != nil {
		return fmt.Errorf("error decoding response: %v", err)
	}

	for _, music := range musicList {
		fmt.Printf("ID: %s\n", music.ID)
		fmt.Printf("Title: %s\n", music.Title)
		fmt.Printf("Artist: %s\n", music.Artist)
		fmt.Printf("Genre: %s\n", music.Genre)
		fmt.Printf("Format: %s\n", music.Format)
		fmt.Printf("Release Date: %s\n", music.ReleaseDate)
		fmt.Println()
	}
	return nil
}

func (session *Session) GetMusicByID(args ...string) error {
	err := validateLoggedIn(session)
	if err != nil {
		return err
	}

	if len(args) < 1 {
		return fmt.Errorf("please specify a music ID")
	}

	musicUUID, err := uuid.Parse(args[0])
	if err != nil {
		return fmt.Errorf("invalid music ID: %v", err)
	}

	url := session.BaseURL + "music/" + musicUUID.String()

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
		return fmt.Errorf("error getting music: %v", res.Status)
	}

	var music Music
	if err := json.NewDecoder(res.Body).Decode(&music); err != nil {
		return fmt.Errorf("error decoding response: %v", err)
	}

	fmt.Printf("Title: %s\n", music.Title)
	fmt.Printf("Artist: %s\n", music.Artist)
	fmt.Printf("Genre: %s\n", music.Genre)
	fmt.Printf("Release Date: %s\n", music.ReleaseDate)
	fmt.Printf("Format: %s\n", music.Format)
	fmt.Printf("Barcode: %s\n", music.Barcode)

	return nil
}

func (session *Session) SearchMusic(args ...string) error {
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

	url := session.BaseURL + "search/music"

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
		return fmt.Errorf("error searching music: %v", res.Status)
	}

	var musicList []Music
	if err := json.NewDecoder(res.Body).Decode(&musicList); err != nil {
		return fmt.Errorf("error decoding response: %v", err)
	}

	for _, music := range musicList {
		fmt.Printf("Title: %s\n", music.Title)
		fmt.Printf("Artist: %s\n", music.Artist)
		fmt.Printf("Genre: %s\n", music.Genre)
		fmt.Printf("Format: %s\n", music.Format)
		fmt.Printf("Release Date: %s\n", music.ReleaseDate)
		fmt.Println()
	}

	return nil
}
