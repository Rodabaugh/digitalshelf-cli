package digitalshelfapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (session *Session) GetUserLocations(args ...string) error {
	err := validateLoggedIn(session)
	if err != nil {
		return err
	}

	url := session.Base_url + "users/" + session.User.ID.String() + "/locations"

	type locationMembership []struct {
		LocationID   string    `json:"location_id"`
		LocationName string    `json:"location_name"`
		OwnerID      string    `json:"owner_id"`
		JoinedAt     time.Time `json:"joined_at"`
	}

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
		var locations locationMembership
		err = json.NewDecoder(res.Body).Decode(&locations)
		if err != nil {
			return err
		}
		fmt.Println("Locations:")
		for _, location := range locations {
			fmt.Printf("Location ID: %s, Location Name: %s\n", location.LocationID, location.LocationName)
		}
		return nil
	} else {
		return fmt.Errorf("error getting user locations")
	}
}

func (session *Session) CreateLocation(args ...string) error {
	err := validateLoggedIn(session)
	if err != nil {
		return err
	}

	url := session.Base_url + "locations"
	type parameters struct {
		Name    string    `json:"name"`
		OwnerID uuid.UUID `json:"owner_id"`
	}
	params := parameters{
		Name:    args[0],
		OwnerID: session.User.ID,
	}

	type response struct {
		ID        uuid.UUID `json:"id"`
		Name      string    `json:"name"`
		OwnerID   uuid.UUID `json:"owner_id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
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
	if res.StatusCode == http.StatusCreated {
		var response response
		err = json.NewDecoder(res.Body).Decode(&response)
		if err != nil {
			return err
		}
		fmt.Printf("Location created successfully, ID: %s\n", response.ID)
		return nil
	} else {
		return fmt.Errorf("error creating location")
	}
}

func (session *Session) JoinLocaion(args ...string) error {
	err := validateLoggedIn(session)
	if err != nil {
		return err
	}

	url := session.Base_url + "locations/" + args[0] + "/members"
	type parameters struct {
		UserID uuid.UUID `json:"user_id"`
	}

	params := parameters{
		UserID: session.User.ID,
	}
	request_data, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("error marshalling request data: %v", err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(request_data))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	res, err := session.DSAPIClient.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusCreated {
		fmt.Println("Joined successfully")
		return nil
	} else if res.StatusCode == http.StatusNotFound {
		return fmt.Errorf("location not found")
	} else {
		return fmt.Errorf("error adding member to location")
	}
}
