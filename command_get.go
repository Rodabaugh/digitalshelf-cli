package main

import (
	"fmt"

	"github.com/Rodabaugh/digitalshelf-cli/internal/digitalshelfapi"
	"github.com/google/uuid"
)

func commandGet(session *digitalshelfapi.Session, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("please specify what you want to get")
	}
	switch args[0] {
	case "locations":
		return session.GetUserLocations()
	case "invites":
		return session.GetUserInvites()
	case "cases":
		return session.GetCases()
	case "shelves":
		if len(args) < 2 {
			return fmt.Errorf("please specify a case ID")
		}
		return session.GetShelves(args[1])
	case "movies":
		if len(args) < 2 {
			return fmt.Errorf("please specify a shelf ID")
		}
		return session.GetMovies(args[1])
	case "movie":
		if len(args) < 2 {
			return fmt.Errorf("please specify a movie ID")
		}
		return session.GetMovie(args[1])
	case "shows":
		if len(args) < 2 {
			return fmt.Errorf("please specify a shelf ID")
		}
		return session.GetShows(args[1])
	case "show":
		if len(args) < 2 {
			return fmt.Errorf("please specify a show ID")
		}
		return session.GetShow(args[1])
	case "books":
		if len(args) < 2 {
			return fmt.Errorf("please specify a shelf ID")
		}
		return session.GetBooks(args[1])
	case "book":
		if len(args) < 2 {
			return fmt.Errorf("please specify a book ID")
		}
		return session.GetBook(args[1])
	case "location":
		return getLocation(session, args[1:]...)
	default:
		return fmt.Errorf("unknown get command: %s", args[0])
	}
}

func getLocation(session *digitalshelfapi.Session, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("please specify what you want to get for this location")
	}
	if session.CurrentLocation == uuid.Nil {
		return fmt.Errorf("please set a location first")
	}

	switch args[0] {
	case "movies":
		return session.GetAllLocationMovies(session.CurrentLocation.String())
	case "shows":
		return session.GetAllLocationShows(session.CurrentLocation.String())
	case "books":
		return session.GetAllLocationBooks(session.CurrentLocation.String())
	default:
		return fmt.Errorf("unknown get command: %s", args[0])
	}
}
