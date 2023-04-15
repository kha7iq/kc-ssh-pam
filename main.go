package main

import (
	"log"
	"os"

	"github.com/kha7iq/kc-ssh-pam/internal/auth"
	"github.com/kha7iq/kc-ssh-pam/internal/conf"
	"golang.org/x/oauth2"
)

func main() {
	c, err := conf.LoadConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	providerEndpoint := c.Endpoint + "/realms/" + c.Realm
	username := os.Getenv("PAM_USER")
	var otp string

	// Analyze the input from stdIn and split the password if it containcts "/"  return otp and pass
	password, otp, err := auth.ReadPasswordWithOTP()
	if err != nil {
		log.Fatal(err)
	}
	// Get provider configuration
	provider, err := auth.GetOIDCProvider(providerEndpoint)
	if err != nil {
		log.Fatalf("Failed to retrieve provider configuration: %v\n", err)
	}

	// Create a new OAuth2 config with the provided client configuration and scopes
	oauth2Config := oauth2.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{"profile", "email", c.ClientScope},
	}

	// Retrieve an OIDC token using the password grant type
	accessToken, err := auth.RequestJWT(username, password, otp, c.ClientScope, oauth2Config)
	if err != nil {
		log.Fatalf("Failed to retrieve token: %v\n", err)
		os.Exit(2)
	}

	// Verify the token signature and retrieve the ID token
	if err := auth.VerifyToken(accessToken, c.ClientID, c.ClientSecret, c.Realm, c.Endpoint); err != nil {
		// handle the error
		log.Fatal(err)
		os.Exit(3)
	}
	log.Println("Token acquired and verified Sucessfully.")
}
