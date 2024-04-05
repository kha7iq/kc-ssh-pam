package users

import (
	"errors"
	"context"

	"github.com/kha7iq/kc-ssh-pam/internal/conf"
	"github.com/Nerzal/gocloak/v13"
)

func IsTOTPConfigured(username string, c conf.Config) (bool, error) {
	keycloakAddr := c.Endpoint
	realm := c.Realm
	realmAPIUsername := c.RealmAPIUsername
	realmAPIPassword := c.RealmAPIPassword
    	client := gocloak.NewClient(keycloakAddr)
    	ctx := context.Background()

	// Request an API user Access Token
    	token, err := client.LoginAdmin(ctx, realmAPIUsername, realmAPIPassword, realm)
    	if err != nil {
		return false, err 
    	}
	
	// Query users by SSH username
    	userCandidates, err := client.GetUsers(ctx, token.AccessToken, realm, gocloak.GetUsersParams{Username: &username})
	if err != nil {
		return false, err
	}
    	if len(userCandidates) > 0 {
		return *userCandidates[0].Totp, nil 	// Note that usernames must be unique within the realm!
    	}

	return false, errors.New("User not found!") 
}
