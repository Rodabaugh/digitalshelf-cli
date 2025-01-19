package digitalshelfapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (session *Session) CreateUser(args ...string) error {
	err := validateLoggedIn(session)
	if err == nil {
		return err
	}

	if len(args) < 3 {
		return fmt.Errorf("please specify name, email, and password")
	}

	name := args[0]
	email := args[1]
	password := args[2]

	if name == "" {
		return fmt.Errorf("name is required")
	}
	if email == "" {
		return fmt.Errorf("email is required")
	}
	if password == "" {
		return fmt.Errorf("password is required")
	}

	url := session.BaseURL + "users"

	type parameters struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	params := parameters{
		Name:     name,
		Email:    email,
		Password: password,
	}

	reqBody, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("error marshalling request body: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := session.DSAPIClient.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}

	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("error creating user: %v", res.Status)
	}

	fmt.Println("User created successfully")
	fmt.Println("Please login with the new user credentials")

	return nil
}

func (session *Session) SearchUsers(args ...string) error {
	err := validateLoggedIn(session)
	if err != nil {
		return err
	}

	if len(args) < 1 {
		return fmt.Errorf("please specify an email address to search for")
	}

	email := args[0]
	url := session.BaseURL + "search/users?email=" + email

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Add("Authorization", "Bearer "+session.Token)
	res, err := session.DSAPIClient.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error searching for user: %v", res.Status)
	}

	resultUser := User{}

	err = json.NewDecoder(res.Body).Decode(&resultUser)
	if err != nil {
		return fmt.Errorf("error decoding response: %v", err)
	}

	fmt.Printf("User found: %s (%s) ID: %s\n", resultUser.Name, resultUser.Email, resultUser.ID)
	return nil
}
