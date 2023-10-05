package sigpairadmin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type UserId = int32

type Client struct {
	BaseUrl    string
	HTTPClient *http.Client
	adminToken string
}

type CreateUserPayload struct {
	Name string `json:"name"`
}
type UserTokenPayload struct {
	UserId   UserId `json:"user_id"`
	Lifetime uint32 `json:"lifetime"`
}
type UserTokenRes struct {
	UserToken string `json:"token"`
}

type CreateUserRes struct {
	UserId int32 `json:"user_id"`
}

// NewClient creates a new client for the given base-url and admin-token
func NewClient(baseUrl string, adminToken string) *Client {
	return &Client{
		BaseUrl:    baseUrl,
		HTTPClient: &http.Client{},
		adminToken: adminToken,
	}
}

func (c *Client) createPostRequest(route string, jsonData []byte) ([]byte, error) {

	req, err := http.NewRequest("POST", c.BaseUrl+route, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.adminToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned error: %v", res.Status)
	}

	return io.ReadAll(res.Body)
}

// CreateUser creates a new user with the given name
//
//	Returns the `UserId` of the created user
func (c *Client) CreateUser(userName string) (UserId, error) {
	payload := CreateUserPayload{
		Name: userName,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return -1, err
	}
	body, err := c.createPostRequest("/v1/create-user", jsonData)
	if err != nil {
		return -1, err
	}
	var response CreateUserRes
	if err := json.Unmarshal(body, &response); err != nil {
		return -1, fmt.Errorf("failed to parse response: %v", err)
	}

	return response.UserId, nil
}

// GenUserToken generates a new user token for the given user-id and lifetime (seconds)
//
// Params:
//   - userId: the id of the user to generate the token for
//   - lifetime: the lifetime of the token in seconds
//
// # Returns:
//   - the generated user-token
func (c *Client) GenUserToken(userId UserId, lifetime uint32) (string, error) {
	payload := UserTokenPayload{
		UserId:   userId,
		Lifetime: lifetime,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	res, err := c.createPostRequest("/v1/user-token", jsonData)
	if err != nil {
		return "", err
	}

	var response UserTokenRes
	if err := json.Unmarshal(res, &response); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	return response.UserToken, nil
}
