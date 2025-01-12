package main

import (
	"fmt"

	"github.com/Rodabaugh/digitalshelf-cli/internal/digitalshelfapi"
)

func commandHelp(session *digitalshelfapi.Session, args ...string) error {
	availableCommands := getCommands()

	for _, cmd := range availableCommands {
		fmt.Printf("- %s : %s\n", cmd.name, cmd.description)
	}

	return nil
}
