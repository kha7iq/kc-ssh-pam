package auth

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Nerzal/gocloak/v13"
	"github.com/coreos/go-oidc"
	"github.com/golang-jwt/jwt"
)

type Token struct {
	// AccessToken is the token that authorizes and authenticates
	// the requests.
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	Expiry       time.Time `json:"expiry,omitempty"`
}

func VerifyToken(aToken, cID, cSecret, providerRealm, providerUrl string) error {

	// Set up the Keycloak client
	client := gocloak.NewClient(providerUrl)

	// Set up the JWT token parser
	parser := jwt.Parser{}

	// Parse the access token
	token, _, err := parser.ParseUnverified(aToken, jwt.MapClaims{})
	if err != nil {
		return fmt.Errorf("Error parsing access token:", err)

	}

	// Get the token's claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("Error getting token claims")
	}

	// Get the token's expiration time
	expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)

	// Check if the token is expired
	if time.Now().After(expirationTime) {
		return fmt.Errorf("Access token has expired")
	}

	// Verify the access token with Keycloak
	_, err = client.RetrospectToken(context.Background(), aToken, cID, cSecret, providerRealm)
	if err != nil {
		return fmt.Errorf("Access token verification failed:", err)
	}

	return nil
}

func RequestJWT(username, password, otp, tokenUrl, clientid, clientsecret, clientscope string) (string, error) {
	// Create a new HTTP client
	client := &http.Client{}

	// Encode the username, password (with OTP appended)
	urlV := url.Values{}
	urlV.Add("grant_type", "password")
	urlV.Add("client_id", clientid)
	urlV.Add("client_secret", clientsecret)
	urlV.Add("username", username)
	urlV.Add("password", password)

	if len(clientscope) > 0 {
		urlV.Add("scop", clientscope)
	}

	if len(otp) > 0 {
		urlV.Add("totp", otp)
	}

	// Send a POST request to the token endpoint
	req, err := http.NewRequest("POST", tokenUrl, strings.NewReader(urlV.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	// Check if the response status code is not 200 OK
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}
	// Check if the response body is nil before decoding it
	if resp.Body == nil {
		return "", fmt.Errorf("response body is nil")
	}
	// Decode the JSON response
	var tokenResp Token
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", err
	}
	// Check if the token response is valid
	if tokenResp.AccessToken == "" {
		return "", fmt.Errorf("access token is empty")
	}

	return tokenResp.AccessToken, nil
}

func ReadPasswordWithOTP() (string, string, error) {
	var password string
	var otp string

	stdinScanner := bufio.NewScanner(os.Stdin)
	if stdinScanner.Scan() {
		pass := strings.Trim(stdinScanner.Text(), "\x00")

		// Extract the password and OTP from the input string
		if strings.Contains(pass, "/") {
			creds := strings.Split(pass, "/")
			password = creds[0]
			otp = creds[1]
		} else {
			password = pass
			otp = ""
		}
	}

	// Check for errors during input
	if err := stdinScanner.Err(); err != nil {
		return "", "", err
	}

	return password, otp, nil
}

func GetOIDCProvider(providerEndpoint string) (*oidc.Provider, error) {
	// Create a context to perform the OIDC discovery process
	ctx := context.Background()

	// Retrieve provider configuration
	provider, err := oidc.NewProvider(ctx, providerEndpoint)
	if err != nil {
		return nil, err
	}

	return provider, nil
}
