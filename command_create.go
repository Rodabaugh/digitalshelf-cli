package main

import (
	"fmt"

	"github.com/Rodabaugh/digitalshelf-cli/internal/digitalshelfapi"
)

func commandCreate(session *digitalshelfapi.Session, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("please specify what you want to create")
	}
	if len(args) == 1 {
		return fmt.Errorf("please specify what you want to create")
	}
	switch args[0] {
	case "location":
		return session.CreateLocation(args[1])
	case "case":
		return session.CreateCase(args[1])
	case "shelf":
		if len(args) < 3 {
			return fmt.Errorf("please specify a shelf name and case ID")
		}
		return session.CreateShelf(args[1], args[2])
	default:
		return fmt.Errorf("unknown create command: %s", args[0])
	}
}
