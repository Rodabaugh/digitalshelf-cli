package digitalshelfapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (session *Session) LookupMovieBarcode(args ...string) (Movie, error) {
	if len(args) < 1 {
		return Movie{}, fmt.Errorf("please provide a barcode")
	}

	err := validateLoggedIn(session)
	if err != nil {
		return Movie{}, err
	}

	barcode := args[0]

	url := session.BaseURL + "search/movie_barcodes/" + barcode

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Movie{}, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+session.Token)

	res, err := session.DSAPIClient.HttpClient.Do(req)
	if err != nil {
		return Movie{}, fmt.Errorf("error making request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return Movie{}, fmt.Errorf("movie not found")
	}

	if res.StatusCode != http.StatusOK {
		return Movie{}, fmt.Errorf("error looking up movie: %v", res.Status)
	}

	var movie Movie
	if err := json.NewDecoder(res.Body).Decode(&movie); err != nil {
		return Movie{}, fmt.Errorf("error decoding response: %v", err)
	}
	fmt.Println("Movie found")
	fmt.Printf("Title: %s\n", movie.Title)
	fmt.Printf("Genre: %s\n", movie.Genre)
	fmt.Printf("Actors: %s\n", movie.Actors)
	fmt.Printf("Writer: %s\n", movie.Writer)
	fmt.Printf("Director: %s\n", movie.Director)
	fmt.Printf("Format: %s\n", movie.Format)
	fmt.Printf("Release Date: %s\n", movie.ReleaseDate)

	return movie, nil
}

func (session *Session) AddMovie(shelfID uuid.UUID, movie Movie) error {
	err := validateLoggedIn(session)
	if err != nil {
		return err
	}

	type parameters struct {
		Title       string    `json:"title"`
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
		Title:       movie.Title,
		Genre:       movie.Genre,
		Actors:      movie.Actors,
		Writer:      movie.Writer,
		Director:    movie.Director,
		Barcode:     movie.Barcode,
		Format:      movie.Format,
		ShelfID:     shelfID,
		ReleaseDate: movie.ReleaseDate,
	}

	url := session.BaseURL + "movies"

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
		fmt.Println("Movie added successfully")
		return nil
	}

	return fmt.Errorf("error adding movie: %v", res.Status)
}

func (session *Session) GetMovies(args ...string) error {
	err := validateLoggedIn(session)
	if err != nil {
		return err
	}

	if len(args) < 1 {
		return fmt.Errorf("please specify a shelf ID")
	}

	shelfID := args[0]

	url := session.BaseURL + "shelves/" + shelfID + "/movies"

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
		return fmt.Errorf("error getting movies: %v", res.Status)
	}

	var movies []Movie
	if err := json.NewDecoder(res.Body).Decode(&movies); err != nil {
		return fmt.Errorf("error decoding response: %v", err)
	}

	for _, movie := range movies {
		fmt.Printf("ID: %s\n", movie.ID)
		fmt.Printf("Title: %s\n", movie.Title)
		fmt.Printf("Genre: %s\n", movie.Genre)
		fmt.Printf("Actors: %s\n", movie.Actors)
		fmt.Printf("Writer: %s\n", movie.Writer)
		fmt.Printf("Director: %s\n", movie.Director)
		fmt.Printf("Format: %s\n", movie.Format)
		fmt.Printf("Release Date: %s\n", movie.ReleaseDate)
		fmt.Println()
	}

	return nil
}

func (session *Session) GetAllLocationMovies(args ...string) error {
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

	url := session.BaseURL + "locations/" + locationID.String() + "/movies"

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
		return fmt.Errorf("error getting movies: %v", res.Status)
	}

	var movies []Movie
	if err := json.NewDecoder(res.Body).Decode(&movies); err != nil {
		return fmt.Errorf("error decoding response: %v", err)
	}

	for _, movie := range movies {
		fmt.Printf("ID: %s\n", movie.ID)
		fmt.Printf("Title: %s\n", movie.Title)
		fmt.Printf("Release Date: %s\n", movie.ReleaseDate)
		fmt.Printf("Format: %s\n", movie.Format)
		fmt.Println()
	}
	return nil
}

func (session *Session) GetMovie(args ...string) error {
	err := validateLoggedIn(session)
	if err != nil {
		return err
	}

	if len(args) < 1 {
		return fmt.Errorf("please specify a movie ID")
	}

	movieUUID, err := uuid.Parse(args[0])
	if err != nil {
		return fmt.Errorf("invalid movie ID: %v", err)
	}

	url := session.BaseURL + "movies/" + movieUUID.String()

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
		return fmt.Errorf("error getting movie: %v", res.Status)
	}

	var movie Movie
	if err := json.NewDecoder(res.Body).Decode(&movie); err != nil {
		return fmt.Errorf("error decoding response: %v", err)
	}

	fmt.Printf("Title: %s\n", movie.Title)
	fmt.Printf("Genre: %s\n", movie.Genre)
	fmt.Printf("Actors: %s\n", movie.Actors)
	fmt.Printf("Writer: %s\n", movie.Writer)
	fmt.Printf("Director: %s\n", movie.Director)
	fmt.Printf("Release Date: %s\n", movie.ReleaseDate)
	fmt.Printf("Barcode: %s\n", movie.Barcode)
	fmt.Printf("Format: %s\n", movie.Format)

	return nil
}

func (session *Session) SearchMovies(args ...string) error {
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

	url := session.BaseURL + "search/movies"

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
		return fmt.Errorf("error searching movies: %v", res.Status)
	}

	var movies []Movie
	if err := json.NewDecoder(res.Body).Decode(&movies); err != nil {
		return fmt.Errorf("error decoding response: %v", err)
	}

	for _, movie := range movies {
		fmt.Printf("Title: %s\n", movie.Title)
		fmt.Printf("Genre: %s\n", movie.Genre)
		fmt.Printf("Actors: %s\n", movie.Actors)
		fmt.Printf("Writer: %s\n", movie.Writer)
		fmt.Printf("Director: %s\n", movie.Director)
		fmt.Printf("Release Date: %s\n", movie.ReleaseDate)
		fmt.Printf("Format: %s\n", movie.Format)
		fmt.Println()
	}

	return nil
}

func (session *Session) UpdateMovieShelf(args ...string) error {
	err := validateLoggedIn(session)
	if err != nil {
		return err
	}

	if len(args) < 2 {
		return fmt.Errorf("please provide a movie ID and a shelf ID")
	}

	movieID, err := uuid.Parse(args[0])
	if err != nil {
		return fmt.Errorf("invalid movie ID: %v", err)
	}

	shelfID, err := uuid.Parse(args[1])
	if err != nil {
		return fmt.Errorf("invalid shelf ID: %v", err)
	}

	err = session.validateShelf(shelfID.String())
	if err != nil {
		return err
	}

	url := session.BaseURL + "movies/" + movieID.String()

	reqBody, err := json.Marshal(map[string]string{"shelf_id": shelfID.String()})
	if err != nil {
		return fmt.Errorf("error marshalling request: %v", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+session.Token)

	res, err := session.DSAPIClient.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		fmt.Println("Movie shelf updated successfully")
		return nil
	}

	return fmt.Errorf("error updating movie shelf: %v", res.Status)
}
