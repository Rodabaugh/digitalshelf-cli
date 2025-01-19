package main

import (
	"fmt"

	"github.com/Rodabaugh/digitalshelf-cli/internal/digitalshelfapi"
)

func commandRemove(session *digitalshelfapi.Session, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("please specify what you want to remove")
	}
	switch args[0] {
	case "member":
		if len(args) < 2 {
			return fmt.Errorf("please specify a user ID")
		}
		return session.RemoveLocationMember(args[1])
	default:
		return fmt.Errorf("unknown remove command: %s", args[0])
	}
}
