package main

import (
	"fmt"

	"github.com/Rodabaugh/digitalshelf-cli/internal/digitalshelfapi"
)

func commandUpdate(session *digitalshelfapi.Session, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("please specify what you want to update")
	}
	if len(args) == 1 {
		return fmt.Errorf("please specify what you want to update")
	}
	switch args[0] {
	case "movie":
		if len(args) < 3 {
			return fmt.Errorf("please specify a movie ID and new shelf ID")
		}
		return session.UpdateMovieShelf(args[1], args[2])
	default:
		return fmt.Errorf("unknown update command: %s", args[0])
	}
}
