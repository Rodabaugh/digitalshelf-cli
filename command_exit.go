package main

import (
	"os"

	"github.com/Rodabaugh/digitalshelf-cli/internal/digitalshelfapi"
)

func commandExit(session *digitalshelfapi.Session, args ...string) error {
	os.Exit(0)
	return nil
}
