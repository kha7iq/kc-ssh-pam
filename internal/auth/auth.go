package auth

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

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

type OIDCProviderInfo struct {
	Issuer      string      `json:"issuer"`
	AuthURL     string      `json:"authorization_endpoint"`
	TokenURL    string      `json:"token_endpoint"`
	JWKSURL     string      `json:"jwks_uri"`
	UserInfoURL string      `json:"userinfo_endpoint"`
	Algorithms  []string    `json:"id_token_signing_alg_values_supported"`
	KeySet      oidc.KeySet `json:",omitempty"`
}

// Stole this from github.com/coreos/go-oidc
func unmarshalResp(r *http.Response, body []byte, v interface{}) error {
	err := json.Unmarshal(body, &v)
	if err == nil {
		return nil
	}
	ct := r.Header.Get("Content-Type")
	mediaType, _, parseErr := mime.ParseMediaType(ct)
	if parseErr == nil && mediaType == "application/json" {
		return fmt.Errorf("got Content-Type = application/json, but could not unmarshal as JSON: %v", err)
	}
	return fmt.Errorf("expected Content-Type = application/json, got %q: %v", ct, err)
}

func GetProviderInfo(providerEndpoint string) (*OIDCProviderInfo, error) {
	wellKnown := strings.TrimSuffix(providerEndpoint, "/") + "/.well-known/openid-configuration"
	// Query the oidc provider for the configuration
	req, err := http.NewRequest("GET", wellKnown, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %s", resp.Status, body)
	}

	// Create a pointer to hold the internal oidc config from the provider
	p := new(OIDCProviderInfo)
	// Parse the json body of the http response into the config variable
	err = unmarshalResp(resp, body, p)
	if err != nil {
		return nil, fmt.Errorf("oidc: failed to decode provider discovery object: %v", err)
	}
	ctx := context.Background()
	p.KeySet = oidc.NewRemoteKeySet(ctx, p.JWKSURL)
	return p, nil
}

func (provider *OIDCProviderInfo) VerifyToken(aToken string) error {
	// Verify the JWT Signature before further parsing and processing
	ctx := context.Background()
	_, err := provider.KeySet.VerifySignature(ctx, aToken)
	if err != nil {
		return fmt.Errorf("Access token verification failed: %v", err)
	}

	// Set up the JWT token parser
	parser := jwt.Parser{}

	// Parse the access token
	token, _, err := parser.ParseUnverified(aToken, jwt.MapClaims{})
	if err != nil {
		return fmt.Errorf("Error parsing access token: %v", err)
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

	return nil
}

func RequestJWT(username, password, otp, tokenUrl, clientid, clientsecret, clientscope string) (string, error) {
	// Create a new HTTP client
	client := &http.Client{}

	// Encode the username, password (with OTP appended)
	urlV := url.Values{}
	urlV.Add("grant_type", "password")
	urlV.Add("client_id", clientid)
	if len(clientsecret) > 0 {
		urlV.Add("client_secret", clientsecret)
	}
	urlV.Add("username", username)
	urlV.Add("password", password)

	if len(clientscope) > 0 {
		urlV.Add("scope", clientscope)
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
