package main

import (
	"fmt"

	"github.com/Rodabaugh/digitalshelf-cli/internal/digitalshelfapi"
)

func commandSet(session *digitalshelfapi.Session, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("please specify what you want to set")
	}
	switch args[0] {
	case "location":
		return session.SetCurrentLocation(args[1])
	case "shelf":
		return session.SetCurrentShelf(args[1])
	default:
		return fmt.Errorf("unknown set command: %s", args[0])
	}
}
