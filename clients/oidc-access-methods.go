package clients

import (
	"Savannah_Screening_Test/dtos"
	"bytes"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func GetKeyCloakAdminToken() (string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file:: Using os files instead")
	}
	log.Println("Loading keycloak backend app details")

	clientID := os.Getenv("OIDC_CLIENT_ID")
	clientSecret := os.Getenv("OIDC_CLIENT_SECRET")
	grantType := os.Getenv("OIDC_GRANT_TYPE")

	data := url.Values{
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"grant_type":    {grantType},
	}

	resp, err := http.PostForm("http://localhost:8080/realms/master/protocol/openid-connect/token", data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
	}

	json.NewDecoder(resp.Body).Decode(&result)
	//fmt.Println("ADMIN TOKEN" + result.AccessToken)
	return result.AccessToken, nil
}

func CreateUserInKeycloak(request dtos.SignUpRequest, token string) (string, error) {
	url := "http://localhost:8080/admin/realms/master/users"
	userData := map[string]interface{}{
		"username":  request.Email,
		"email":     request.Email,
		"enabled":   true,
		"firstName": request.Name,
		"credentials": []map[string]interface{}{
			{
				"type":      "password",
				"value":     request.Password,
				"temporary": false,
			},
		},
	}

	jsonData, _ := json.Marshal(userData)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		_, err = io.ReadAll(resp.Body)
		fmt.Println(resp.StatusCode)       //logging to check error
		return "", errors.New(err.Error()) //custom message
	}

	location := resp.Header.Get("Location")
	parts := strings.Split(location, "/")
	userID := parts[len(parts)-1] //use ID for customer in my local dd
	return userID, nil
}

func AssignRoleToUser(userID, roleName, token string) error {
	// 1. Get Role ID
	roleURL := fmt.Sprintf("http://localhost:8080/admin/realms/master/roles/%s", roleName)
	req, _ := http.NewRequest("GET", roleURL, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var role struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	json.NewDecoder(resp.Body).Decode(&role)

	// 2. Assign Role
	assignURL := fmt.Sprintf("http://localhost:8080/admin/realms/master/users/%s/role-mappings/realm", userID)
	rolePayload := []map[string]string{{"id": role.ID, "name": role.Name}}
	body, _ := json.Marshal(rolePayload)
	req, _ = http.NewRequest("POST", assignURL, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func LoginWithPasswordGrant(username, password string) (*dtos.TokenResponse, error) {
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

	var tokenResp dtos.TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

var cachedKey *rsa.PublicKey
var lastFetched time.Time

func GetKeycloakPublicKey() (*rsa.PublicKey, error) {
	if cachedKey != nil && time.Since(lastFetched) < 5*time.Minute {
		return cachedKey, nil
	}

	jwksURL := os.Getenv("KEYCLOAK_JWKS_URL")
	if jwksURL == "" {
		jwksURL = "http://localhost:8080/realms/master/protocol/openid-connect/certs"
	}

	resp, err := http.Get(jwksURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var jwks struct {
		Keys []struct {
			N   string `json:"n"`
			E   string `json:"e"`
			Kty string `json:"kty"`
			Alg string `json:"alg"`
			Use string `json:"use"`
		} `json:"keys"`
	}

	if err := json.Unmarshal(body, &jwks); err != nil {
		return nil, err
	}

	if len(jwks.Keys) == 0 {
		return nil, errors.New("no keys found in JWKS")
	}

	pubKey, err := parseRSAPublicKeyFromJWKS(jwks.Keys[0].N, jwks.Keys[0].E)
	if err != nil {
		return nil, err
	}

	cachedKey = pubKey
	lastFetched = time.Now()
	return cachedKey, nil
}

func parseRSAPublicKeyFromJWKS(nStr, eStr string) (*rsa.PublicKey, error) {
	nb, err := base64.RawURLEncoding.DecodeString(nStr)
	if err != nil {
		return nil, err
	}
	eb, err := base64.RawURLEncoding.DecodeString(eStr)
	if err != nil {
		return nil, err
	}

	n := new(big.Int).SetBytes(nb)
	e := 0
	for _, b := range eb {
		e = e<<8 + int(b)
	}

	return &rsa.PublicKey{N: n, E: e}, nil
}
