package digitalshelfapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

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
