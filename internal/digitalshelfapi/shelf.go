package digitalshelfapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (session *Session) CreateShelf(args ...string) error {
	err := validateLoggedIn(session)
	if err != nil {
		return err
	}

	if session.CurrentLocation == uuid.Nil {
		return fmt.Errorf("please set a current location first")
	}

	if len(args) < 2 {
		return fmt.Errorf("please specify a case ID and shelf name")
	}

	type parameters struct {
		Name   string    `json:"name"`
		CaseID uuid.UUID `json:"case_id"`
	}

	caseID := args[0]
	shelfName := args[1]

	caseUUID, err := uuid.Parse(caseID)
	if err != nil {
		return fmt.Errorf("invalid case ID")
	}

	url := session.BaseURL + "shelves"

	params := parameters{
		Name:   shelfName,
		CaseID: caseUUID,
	}

	reqBody, err := json.Marshal(params)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := session.DSAPIClient.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusCreated {
		fmt.Println("Shelf created successfully")
		return nil
	} else {
		return fmt.Errorf("error creating shelf")
	}
}

func (session *Session) GetShelves(args ...string) error {
	err := validateLoggedIn(session)
	if err != nil {
		return err
	}

	if len(args) < 1 {
		return fmt.Errorf("please specify a case ID")
	}

	caseID := args[0]
	err = session.ValidateCase(caseID)
	if err != nil {
		return err
	}

	url := session.BaseURL + "cases/" + caseID + "/shelves"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	res, err := session.DSAPIClient.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		var shelves []Shelf
		err = json.NewDecoder(res.Body).Decode(&shelves)
		if err != nil {
			return err
		}

		for _, shelf := range shelves {
			fmt.Printf("Name: %s, ID: %s\n", shelf.Name, shelf.ID)
		}
		return nil
	} else {
		return fmt.Errorf("error getting shelves")
	}
}

func (session *Session) validateShelf(shelfID string) error {
	err := validateLoggedIn(session)
	if err != nil {
		return err
	}

	url := session.BaseURL + "shelves/" + shelfID

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	res, err := session.DSAPIClient.HttpClient.Do(req)
	if err != nil {
		fmt.Println("error making request: ", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		var c Case
		err = json.NewDecoder(res.Body).Decode(&c)
		if err != nil {
			return err
		}
		return nil
	} else {
		return fmt.Errorf("shelf not found")
	}
}
