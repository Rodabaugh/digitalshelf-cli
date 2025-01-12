package main

import (
	"fmt"

	"github.com/Rodabaugh/digitalshelf-cli/internal/digitalshelfapi"
)

func commandGet(session *digitalshelfapi.Session, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("please specify what you want to get")
	}
	switch args[0] {
	case "locations":
		return session.GetUserLocations()
	default:
		return fmt.Errorf("unknown get command: %s", args[0])
	}
}
