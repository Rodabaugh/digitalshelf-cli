package main

import "fmt"

func commandHelp() error {
	availableCommands := getCommands()

	for _, cmd := range availableCommands {
		fmt.Printf("- %s : %s\n", cmd.name, cmd.description)
	}

	return nil
}
