package main

import (
	"fmt"

	"github.com/Rodabaugh/digitalshelf-cli/internal/digitalshelfapi"
)

func commandSearch(session *digitalshelfapi.Session, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("please specify what you want to search")
	}
	switch args[0] {
	case "users":
		if len(args) < 2 {
			return fmt.Errorf("please specify an email address")
		}
		return session.SearchUsers(args[1])
	case "movies":
		if len(args) < 2 {
			return fmt.Errorf("please specify a search term")
		}
		return session.SearchMovies(args[1])
	case "shows":
		if len(args) < 2 {
			return fmt.Errorf("please specify a search term")
		}
		return session.SearchShows(args[1])
	case "books":
		if len(args) < 2 {
			return fmt.Errorf("please specify a search term")
		}
		return session.SearchBooks(args[1])
	case "music":
		if len(args) < 2 {
			return fmt.Errorf("please specify a search term")
		}
		return session.SearchMusic(args[1])
	default:
		return fmt.Errorf("unknown search command: %s", args[0])
	}
}
