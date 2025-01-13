package main

import (
	"fmt"
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
	fmt.Printf("Entering New Movie\n----------------------------\n")
	var title, genre, actors, writer, director, releaseDateStr string
	fmt.Print("Title: ")
	fmt.Scanln(&title)
	fmt.Print("Genre: ")
	fmt.Scanln(&genre)
	fmt.Print("Actors: ")
	fmt.Scanln(&actors)
	fmt.Print("Writer: ")
	fmt.Scanln(&writer)
	fmt.Print("Director: ")
	fmt.Scanln(&director)
	fmt.Print("Release Date (YYYY-MM-DD): ")
	fmt.Scanln(&releaseDateStr)
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
