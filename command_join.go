package main

import (
	"fmt"

	"github.com/Rodabaugh/digitalshelf-cli/internal/digitalshelfapi"
)

func commandJoin(session *digitalshelfapi.Session, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("please specify a location ID")
	}

	return session.JoinLocaion(args[0])
}
