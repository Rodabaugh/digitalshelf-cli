package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/Rodabaugh/digitalshelf-cli/internal/digitalshelfapi"
	"github.com/google/uuid"
)

func commandAdd(session *digitalshelfapi.Session, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("please specify what you want to set")
	}
	switch args[0] {
	case "movie":
		if len(args) < 3 {
			return fmt.Errorf("please specify a shelf ID and movie barcode")
		}
		shelfID, err := uuid.Parse(args[1])
		if err != nil {
			return fmt.Errorf("invalid shelf ID: %v", err)
		}

		return addMovie(session, shelfID, args[2])
	case "show":
		if len(args) < 3 {
			return fmt.Errorf("please specify a shelf ID and show barcode")
		}
		shelfID, err := uuid.Parse(args[1])
		if err != nil {
			return fmt.Errorf("invalid shelf ID: %v", err)
		}
		return addShow(session, shelfID, args[2])
	case "moviebulk":
		if len(args) < 3 {
			return fmt.Errorf("please specify a shelf ID and movie barcode")
		}
		shelfID, err := uuid.Parse(args[1])
		if err != nil {
			return fmt.Errorf("invalid shelf ID: %v", err)
		}

		return benchmarkCreateMovie(session, shelfID, args[2])
	default:
		return fmt.Errorf("unknown add command: %s", args[0])
	}
}

func addMovie(session *digitalshelfapi.Session, shelfID uuid.UUID, barcode string) error {
	movie, err := session.LookupMovieBarcode(barcode)
	if err == nil {
		fmt.Printf("Movie found!\n\n")
		fmt.Printf("Title: %s\n", movie.Title)
		fmt.Printf("Genre: %s\n", movie.Genre)
		fmt.Printf("Actors: %s\n", movie.Actors)
		fmt.Printf("Writer: %s\n", movie.Writer)
		fmt.Printf("Director: %s\n", movie.Director)
		fmt.Printf("Release Date: %s\n", movie.ReleaseDate)
		fmt.Println("Do you want to add this movie to the shelf? (y/n)")
		var answer string
		fmt.Scanln(&answer)
		if answer != "y" {
			return fmt.Errorf("movie not added to the shelf")
		}
	}

	if err != nil && err.Error() == "movie not found" {
		fmt.Printf("This barcode does not exist in the database. Please enter it manually.\n\n")
		movie, err = getMovieDetails()
		if err != nil {
			return fmt.Errorf("error getting movie details: %v", err)
		}
		movie.Barcode = barcode
	}

	err = session.AddMovie(shelfID, movie)
	if err != nil {
		return fmt.Errorf("error adding movie to shelf: %v", err)
	}

	return err
}

// There is an issue with this fuction that does not handle spaces in the input. This needs to be fixed.
func getMovieDetails() (digitalshelfapi.Movie, error) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Entering New Movie\n----------------------------\n")
	var title, genre, actors, writer, director, releaseDateStr string
	fmt.Print("Title: ")
	scanner.Scan()
	title = scanner.Text()
	fmt.Print("Genre: ")
	scanner.Scan()
	genre = scanner.Text()
	fmt.Print("Actors: ")
	scanner.Scan()
	actors = scanner.Text()
	fmt.Print("Writer: ")
	scanner.Scan()
	writer = scanner.Text()
	fmt.Print("Director: ")
	scanner.Scan()
	director = scanner.Text()
	fmt.Print("Release Date (YYYY-MM-DD): ")
	scanner.Scan()
	releaseDateStr = scanner.Text()
	releaseDate, err := time.Parse("2006-01-02", releaseDateStr)
	if err != nil {
		return digitalshelfapi.Movie{}, fmt.Errorf("error parsing release date: %v", err)
	}

	movie := digitalshelfapi.Movie{
		Title:       title,
		Genre:       genre,
		Actors:      actors,
		Writer:      writer,
		Director:    director,
		ReleaseDate: releaseDate,
	}

	return movie, nil
}

func addShow(session *digitalshelfapi.Session, shelfID uuid.UUID, barcode string) error {
	show, err := session.LookupShowBarcode(barcode)
	if err == nil {
		fmt.Printf("Show found!\n\n")
		fmt.Printf("Title: %s\n", show.Title)
		fmt.Printf("Season: %d\n", show.Season)
		fmt.Printf("Genre: %s\n", show.Genre)
		fmt.Printf("Actors: %s\n", show.Actors)
		fmt.Printf("Writer: %s\n", show.Writer)
		fmt.Printf("Director: %s\n", show.Director)
		fmt.Printf("Release Date: %s\n", show.ReleaseDate)
		fmt.Println("Do you want to add this show to the shelf? (y/n)")
		var answer string
		fmt.Scanln(&answer)
		if answer != "y" {
			return fmt.Errorf("show not added to the shelf")
		}
	}

	if err != nil && err.Error() == "show not found" {
		fmt.Printf("This barcode does not exist in the database. Please enter it manually.\n\n")
		show, err = getShowDetails()
		if err != nil {
			return fmt.Errorf("error getting show details: %v", err)
		}
		show.Barcode = barcode
	}

	err = session.AddShow(shelfID, show)
	if err != nil {
		return fmt.Errorf("error adding show to shelf: %v", err)
	}

	return err
}

func getShowDetails() (digitalshelfapi.Show, error) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Entering New Show\n----------------------------\n")
	var title, genre, actors, writer, director, releaseDateStr string
	var season int
	fmt.Print("Title: ")
	scanner.Scan()
	title = scanner.Text()
	fmt.Print("Season: ")
	fmt.Scanln(&season)
	fmt.Print("Genre: ")
	scanner.Scan()
	genre = scanner.Text()
	fmt.Print("Actors: ")
	scanner.Scan()
	actors = scanner.Text()
	fmt.Print("Writer: ")
	scanner.Scan()
	writer = scanner.Text()
	fmt.Print("Director: ")
	scanner.Scan()
	director = scanner.Text()
	fmt.Print("Release Date (YYYY-MM-DD): ")
	scanner.Scan()
	releaseDateStr = scanner.Text()
	releaseDate, err := time.Parse("2006-01-02", releaseDateStr)
	if err != nil {
		return digitalshelfapi.Show{}, fmt.Errorf("error parsing release date: %v", err)
	}

	show := digitalshelfapi.Show{
		Title:       title,
		Season:      season,
		Genre:       genre,
		Actors:      actors,
		Writer:      writer,
		Director:    director,
		ReleaseDate: releaseDate,
	}

	return show, nil
}

func benchmarkCreateMovie(session *digitalshelfapi.Session, shelfID uuid.UUID, barcode string) error {
	if session.Platform == "prod" {
		return fmt.Errorf("benchmarking is not allowed in production")
	}

	movie, err := session.LookupMovieBarcode(barcode)
	if err == nil {
		fmt.Printf("Movie found!\n\n")
		fmt.Printf("Title: %s\n", movie.Title)
		fmt.Printf("Genre: %s\n", movie.Genre)
		fmt.Printf("Actors: %s\n", movie.Actors)
		fmt.Printf("Writer: %s\n", movie.Writer)
		fmt.Printf("Director: %s\n", movie.Director)
		fmt.Printf("Release Date: %s\n", movie.ReleaseDate)
		fmt.Println("Do you want to add this movie to the shelf? (y/n)")
		var answer string
		fmt.Scanln(&answer)
		if answer != "y" {
			return fmt.Errorf("movie not added to the shelf")
		}
	}

	if err != nil && err.Error() == "movie not found" {
		fmt.Printf("This barcode does not exist in the database. Please enter it manually.\n\n")
		movie, err = getMovieDetails()
		if err != nil {
			return fmt.Errorf("error getting movie details: %v", err)
		}
		movie.Barcode = barcode
	}

	numberOfMovies := 500000
	startTime := time.Now()
	for i := 0; i < numberOfMovies; i++ {
		err = session.AddMovie(shelfID, movie)
		if err != nil {
			return fmt.Errorf("error adding movie to shelf: %v", err)
		}
		fmt.Printf("Creating movie %v\n", i)
	}

	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	fmt.Printf("Created %v movies in: %v\n", numberOfMovies, elapsedTime)

	return err
}
