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
	if len(args) < 2 {
		return fmt.Errorf("please specify what you want to add and a barcode")
	}

	var shelfID uuid.UUID
	var err error

	if len(args) == 2 {
		if session.CurrentShelf != uuid.Nil {
			shelfID = session.CurrentShelf
		} else {
			fmt.Print("Please enter a shelf ID: ")
			var shelfIDStr string
			fmt.Scanln(&shelfIDStr)
			shelfID, err = uuid.Parse(shelfIDStr)
			if err != nil {
				return fmt.Errorf("invalid shelf ID: %v", err)
			}
		}
	} else {
		shelfID, err = uuid.Parse(args[1])
		if err != nil {
			return fmt.Errorf("invalid shelf ID: %v", err)
		}
	}

	switch args[0] {
	case "movie":
		return addMovie(session, shelfID, args[1])
	case "show":
		return addShow(session, shelfID, args[1])
	case "book":
		return addBook(session, shelfID, args[1])
	case "music":
		return addMusic(session, shelfID, args[1])
	case "movie_bulk":
		return benchmarkCreateMovie(session, shelfID, args[1])
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
		fmt.Printf("Format: %s\n", movie.Format)
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

func getMovieDetails() (digitalshelfapi.Movie, error) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Entering New Movie\n----------------------------\n")
	var title, genre, actors, writer, director, releaseDateStr, format string
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
	fmt.Print("Format: ")
	scanner.Scan()
	format = scanner.Text()

	movie := digitalshelfapi.Movie{
		Title:       title,
		Genre:       genre,
		Actors:      actors,
		Writer:      writer,
		Director:    director,
		ReleaseDate: releaseDate,
		Format:      format,
	}

	return movie, nil
}

func addShow(session *digitalshelfapi.Session, shelfID uuid.UUID, barcode string) error {
	show, err := session.LookupShowBarcode(barcode)
	if err == nil {
		fmt.Printf("Show found!\n\n")
		fmt.Printf("Title: %s\n", show.Title)
		fmt.Printf("Season: %s\n", show.Season)
		fmt.Printf("Genre: %s\n", show.Genre)
		fmt.Printf("Actors: %s\n", show.Actors)
		fmt.Printf("Writer: %s\n", show.Writer)
		fmt.Printf("Director: %s\n", show.Director)
		fmt.Printf("Format: %s\n", show.Format)
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
	var title, season, genre, actors, writer, director, releaseDateStr, format string
	fmt.Print("Title: ")
	scanner.Scan()
	title = scanner.Text()
	fmt.Print("Season: ")
	scanner.Scan()
	season = scanner.Text()
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
	fmt.Print("Format: ")
	scanner.Scan()
	format = scanner.Text()

	show := digitalshelfapi.Show{
		Title:       title,
		Season:      season,
		Genre:       genre,
		Actors:      actors,
		Writer:      writer,
		Director:    director,
		ReleaseDate: releaseDate,
		Format:      format,
	}

	return show, nil
}

func addBook(session *digitalshelfapi.Session, shelfID uuid.UUID, barcode string) error {
	book, err := session.LookupBookBarcode(barcode)
	if err == nil {
		fmt.Printf("Book found!\n\n")
		fmt.Printf("Title: %s\n", book.Title)
		fmt.Printf("Author: %s\n", book.Author)
		fmt.Printf("Genre: %s\n", book.Genre)
		fmt.Printf("Publication Date: %s\n", book.PublicationDate)
		fmt.Println("Do you want to add this book to the shelf? (y/n)")
		var answer string
		fmt.Scanln(&answer)
		if answer != "y" {
			return fmt.Errorf("book not added to the shelf")
		}
	}

	if err != nil && err.Error() == "book not found" {
		fmt.Printf("This barcode does not exist in the database. Please enter it manually.\n\n")
		book, err = getBookDetails()
		if err != nil {
			return fmt.Errorf("error getting book details: %v", err)
		}
		book.Barcode = barcode
	}

	err = session.AddBook(shelfID, book)
	if err != nil {
		return fmt.Errorf("error adding book to shelf: %v", err)
	}

	return err
}

func getBookDetails() (digitalshelfapi.Book, error) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Entering New Book\n----------------------------\n")
	var title, author, genre, publicationDateStr string
	fmt.Print("Title: ")
	scanner.Scan()
	title = scanner.Text()
	fmt.Print("Author: ")
	scanner.Scan()
	author = scanner.Text()
	fmt.Print("Genre: ")
	scanner.Scan()
	genre = scanner.Text()
	fmt.Print("Publication Date (YYYY-MM-DD): ")
	scanner.Scan()
	publicationDateStr = scanner.Text()
	publicationDate, err := time.Parse("2006-01-02", publicationDateStr)
	if err != nil {
		return digitalshelfapi.Book{}, fmt.Errorf("error parsing publication date: %v", err)
	}

	book := digitalshelfapi.Book{
		Title:           title,
		Author:          author,
		Genre:           genre,
		PublicationDate: publicationDate,
	}

	return book, nil
}

func addMusic(session *digitalshelfapi.Session, shelfID uuid.UUID, barcode string) error {
	music, err := session.LookupMusicBarcode(barcode)
	if err == nil {
		fmt.Printf("Music found!\n\n")
		fmt.Printf("Title: %s\n", music.Title)
		fmt.Printf("Artist: %s\n", music.Artist)
		fmt.Printf("Genre: %s\n", music.Genre)
		fmt.Printf("Format: %s\n", music.Format)
		fmt.Printf("Release Date: %s\n", music.ReleaseDate)
		fmt.Println("Do you want to add this music to the shelf? (y/n)")
		var answer string
		fmt.Scanln(&answer)
		if answer != "y" {
			return fmt.Errorf("music not added to the shelf")
		}
	}

	if err != nil && err.Error() == "music not found" {
		fmt.Printf("This barcode does not exist in the database. Please enter it manually.\n\n")
		music, err = getMusicDetails()
		if err != nil {
			return fmt.Errorf("error getting music details: %v", err)
		}
		music.Barcode = barcode
	}

	err = session.AddMusic(shelfID, music)
	if err != nil {
		return fmt.Errorf("error adding music to shelf: %v", err)
	}

	return err
}

func getMusicDetails() (digitalshelfapi.Music, error) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Entering New Music\n----------------------------\n")
	var title, artist, genre, format, releaseDateStr string
	fmt.Print("Title: ")
	scanner.Scan()
	title = scanner.Text()
	fmt.Print("Artist: ")
	scanner.Scan()
	artist = scanner.Text()
	fmt.Print("Genre: ")
	scanner.Scan()
	genre = scanner.Text()
	fmt.Print("Format: ")
	scanner.Scan()
	format = scanner.Text()
	fmt.Print("Release Date (YYYY-MM-DD): ")
	scanner.Scan()
	releaseDateStr = scanner.Text()
	releaseDate, err := time.Parse("2006-01-02", releaseDateStr)
	if err != nil {
		return digitalshelfapi.Music{}, fmt.Errorf("error parsing release date: %v", err)
	}

	music := digitalshelfapi.Music{
		Title:       title,
		Artist:      artist,
		Genre:       genre,
		Format:      format,
		ReleaseDate: releaseDate,
	}

	return music, nil
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
