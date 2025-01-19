package digitalshelfapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (session *Session) CreateCase(args ...string) error {
	err := validateLoggedIn(session)
	if err != nil {
		return err
	}

	if len(args) < 1 {
		return fmt.Errorf("missing case name")
	}

	if session.CurrentLocation == uuid.Nil {
		return fmt.Errorf("no location set - please set a location")
	}

	url := session.BaseURL + "cases"

	type parameters struct {
		Name       string `json:"name"`
		LocationID string `json:"location_id"`
	}

	params := parameters{
		Name:       args[0],
		LocationID: session.CurrentLocation.String(),
	}

	requestBody, err := json.Marshal(params)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+session.Token)

	res, err := session.DSAPIClient.HttpClient.Do(req)
	if err != nil {
		fmt.Println("error making request: ", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusCreated {
		fmt.Println("Case created successfully")
		return nil
	} else {
		return fmt.Errorf("error creating case")
	}
}

func (session *Session) GetCases() error {
	err := validateLoggedIn(session)
	if err != nil {
		return err
	}

	if session.CurrentLocation == uuid.Nil {
		return fmt.Errorf("no location set - please set a location")
	}

	url := session.BaseURL + "locations/" + session.CurrentLocation.String() + "/cases"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+session.Token)

	res, err := session.DSAPIClient.HttpClient.Do(req)
	if err != nil {
		fmt.Println("error making request: ", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		var cases []Case
		err = json.NewDecoder(res.Body).Decode(&cases)
		if err != nil {
			return err
		}

		for _, c := range cases {
			fmt.Printf("Case Name: %s, Case ID: %s\n", c.Name, c.ID)
		}
		return nil
	} else {
		return fmt.Errorf("error getting cases")
	}
}

func (session *Session) ValidateCase(caseID string) error {
	err := validateLoggedIn(session)
	if err != nil {
		return err
	}

	url := session.BaseURL + "cases/" + caseID

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+session.Token)

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

		if c.LocationID != session.CurrentLocation {
			return fmt.Errorf("case is not at the current location")
		}
		return nil
	} else {
		return fmt.Errorf("case not found")
	}
}
