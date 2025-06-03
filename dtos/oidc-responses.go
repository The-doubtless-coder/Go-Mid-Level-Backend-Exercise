package dtos

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

func LoginWithPasswordGrant(username string, password string) (*TokenResponse, error) {
	keycloakTokenURL := os.Getenv("KEYCLOAK_TOKEN_URL") // e.g. http://localhost:8080/realms/master/protocol/openid-connect/token

	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("client_id", os.Getenv("KEYCLOAK_CUSTOMER_CLIENT_ID"))         // e.g., go-customer
	data.Set("client_secret", os.Getenv("KEYCLOAK_CUSTOMER_CLIENT_SECRET")) // your confidential client secret
	data.Set("username", username)
	data.Set("password", password)

	req, err := http.NewRequest("POST", keycloakTokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("invalid login credentials or Keycloak rejected the request")
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}
