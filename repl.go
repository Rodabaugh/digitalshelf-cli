package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Rodabaugh/digitalshelf-cli/internal/digitalshelfapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*digitalshelfapi.Session, ...string) error
}

func startRepl(session *digitalshelfapi.Session) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		print("digitalshelf > ")
		scanner.Scan()
		input := scanner.Text()
		cleanedInput := cleanInput(input)

		if len(cleanedInput) == 0 {
			continue
		}

		commandName := cleanedInput[0]
		args := []string{}

		if len(cleanedInput) > 1 {
			args = cleanedInput[1:]
		}

		command, exists := getCommands()[commandName]
		if !exists {
			fmt.Println("Invalid command")
			continue
		}

		err := command.callback(session, args...)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the DigitalShelf",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Show available commands",
			callback:    commandHelp,
		},
		"login": {
			name:        "login",
			description: "Login to your DigitalShelf account",
			callback:    commandLogin,
		},
		"create": {
			name:        "create",
			description: "Create a new item/location",
			callback:    commandCreate,
		},
		"join": {
			name:        "join",
			description: "Join a location",
			callback:    commandJoin,
		},
		"get": {
			name:        "get",
			description: "Get information about an item/location",
			callback:    commandGet,
		},
	}
}
