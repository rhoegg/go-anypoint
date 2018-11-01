package go_anypoint

import (
	"context"
	"errors"
	"net/http"
)

const loginBasePath = "/accounts/login"

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResult struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	RedirectURL string
}

func (c *Client) Login(ctx context.Context, loginRequest *LoginRequest) (*LoginResult, *Response, error) {
	if loginRequest == nil {
		return nil, nil, errors.New("loginRequest can not be nil")
	}

	path := loginBasePath

	req, err := c.NewRequest(ctx, http.MethodPost, path, loginRequest)
	if err != nil {
		return nil, nil, err
	}

	loginResult := new(LoginResult)
	resp, err := c.Do(ctx, req, loginResult)
	if err != nil {
		return nil, resp, err
	}

	return loginResult, resp, err
}
