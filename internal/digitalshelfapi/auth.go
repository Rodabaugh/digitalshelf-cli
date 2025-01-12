package digitalshelfapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (session *Session) Authenticate(args ...string) error {
	type response struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		Email        string `json:"email"`
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	type parameters struct {
		Email    string
		Password string
	}

	url := session.Base_url + "login"
	email := args[0]
	if email == "" {
		return errors.New("email is required")
	}
	password := args[1]
	if password == "" {
		return errors.New("password is required")
	}

	params := parameters{
		Email:    email,
		Password: password,
	}

	request_data, err := json.Marshal(params)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(request_data))
	if err != nil {
		return err
	}

	res, err := session.DSAPIClient.HttpClient.Do(req)
	if err != nil {
		fmt.Println("error making request: ", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		var response response
		err = json.NewDecoder(res.Body).Decode(&response)
		if err != nil {
			return err
		}
		session.User.ID = response.ID
		session.User.Email = response.Email
		session.User.Name = response.Name
		session.Token = response.Token
		session.RefreshToken = response.RefreshToken
		return nil
	}
	if res.StatusCode == 401 {
		return errors.New("invalid email or password")
	}

	return errors.New("error authenticating")
}
