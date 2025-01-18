package main

import (
	"fmt"

	"github.com/Rodabaugh/digitalshelf-cli/internal/digitalshelfapi"
)

func commandInvite(session *digitalshelfapi.Session, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("please specify a user ID")
	}

	return session.InviteUser(args[0])
}
