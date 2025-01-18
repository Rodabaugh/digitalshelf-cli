package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Rodabaugh/digitalshelf-cli/internal/digitalshelfapi"
)

func commandRegister(session *digitalshelfapi.Session, _ ...string) error {
	scanner := bufio.NewScanner(os.Stdin)
	var name, email string

	fmt.Print("Enter your name: ")
	scanner.Scan()
	name = scanner.Text()

	fmt.Print("Enter your email: ")
	scanner.Scan()
	email = scanner.Text()

	var newPassword, confirmPassword string

	fmt.Print("Enter your new password: ")
	scanner.Scan()
	newPassword = scanner.Text()

	fmt.Print("Confrim your new password: ")
	scanner.Scan()
	confirmPassword = scanner.Text()

	if newPassword != confirmPassword {
		return fmt.Errorf("passwords do not match")
	}

	return session.CreateUser(name, email, newPassword)
}

func commandLogin(session *digitalshelfapi.Session, args ...string) error {
	var email, password string

	fmt.Print("Enter your email: ")
	fmt.Scanln(&email)

	fmt.Print("Enter your password: ")
	fmt.Scanln(&password)

	err := session.Authenticate(email, password)
	if err != nil {
		return err
	}
	fmt.Printf("Welcome, %s\n", session.User.Name)
	return nil
}

func commandLogout(session *digitalshelfapi.Session, args ...string) error {
	if len(args) == 0 {
		return session.Logout()
	}
	switch args[0] {
	case "all":
		return session.RevokeAllSessions()
	default:
		return fmt.Errorf("unknown logout command: %s", args[0])
	}
}

func commandChangePassword(session *digitalshelfapi.Session, args ...string) error {
	var newPassword, confirmPassword string

	fmt.Print("Enter your new password: ")
	fmt.Scanln(&newPassword)

	fmt.Print("Confrim your new password: ")
	fmt.Scanln(&confirmPassword)

	if newPassword != confirmPassword {
		return fmt.Errorf("passwords do not match")
	}

	return session.ChangePassword(newPassword)
}
