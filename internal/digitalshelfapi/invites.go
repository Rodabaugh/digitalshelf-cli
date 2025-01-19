package digitalshelfapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (session *Session) InviteUser(args ...string) error {
	err := validateLoggedIn(session)
	if err != nil {
		return err
	}

	if len(args) < 1 {
		return fmt.Errorf("please specify a user_id")
	}

	userID := args[0]
	if _, err := uuid.Parse(userID); err != nil {
		return fmt.Errorf("invalid user_id: %v", err)
	}

	if session.CurrentLocation == uuid.Nil {
		return fmt.Errorf("please set a current location first")
	}

	url := session.BaseURL + "locations/" + session.CurrentLocation.String() + "/invites"

	type parameters struct {
		UserID string `json:"user_id"`
	}

	params := parameters{
		UserID: userID,
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
	req.Header.Set("Authorization", "Bearer "+session.Token)

	res, err := session.DSAPIClient.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}

	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("error inviting user: %v", res.Status)
	}

	fmt.Println("User invited successfully")

	return nil
}

func (session *Session) GetUserInvites() error {
	err := validateLoggedIn(session)
	if err != nil {
		return err
	}

	url := session.BaseURL + "/users/" + session.User.ID.String() + "/invites"

	type UserInvite []struct {
		UserID       uuid.UUID `json:"userID"`
		LocationID   uuid.UUID `json:"location_id"`
		LocationName string    `json:"location_name"`
		OwnerID      uuid.UUID `json:"owner_id"`
		InvitedAt    time.Time `json:"invited_at"`
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+session.Token)

	res, err := session.DSAPIClient.HttpClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		var invites UserInvite
		json.NewDecoder(res.Body).Decode(&invites)

		fmt.Println("User Invites:")

		for _, invite := range invites {
			fmt.Printf("Location Name: %s, Location ID: %s, Invited At: %s\n", invite.LocationName, invite.LocationID.String(), invite.InvitedAt.String())
		}
		return nil
	}
	return fmt.Errorf("error getting user invites")
}

func (session *Session) RemoveUserInvite(args ...string) error {
	err := validateLoggedIn(session)
	if err != nil {
		return err
	}

	if session.CurrentLocation == uuid.Nil {
		return fmt.Errorf("please set a current location first")
	}

	if len(args) < 1 {
		return fmt.Errorf("missing invite ID")
	}

	inviteID := args[0]
	if _, err := uuid.Parse(inviteID); err != nil {
		return fmt.Errorf("invalid invite ID: %v", err)
	}

	url := session.BaseURL + "locations/" + session.CurrentLocation.String() + "/invites/" + inviteID

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+session.Token)

	res, err := session.DSAPIClient.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNoContent {
		fmt.Println("Invite removed successfully")
		return nil
	}

	return fmt.Errorf("error removing invite: %v", res.Status)
}
