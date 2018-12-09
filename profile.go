package go_anypoint

import (
	"context"
	"net/http"
)

const profileBasePath = "accounts/api/profile"

type ProfileService interface {
	Get(context.Context) (*Profile, *Response, error)
	GetID(context.Context) (string, error)
}

type ProfileServiceOp struct {
	client *Client
}

type Profile struct {
	ID string
	OrganizationID string
}

func (s *ProfileServiceOp) Get(ctx context.Context) (*Profile, *Response, error) {
	path := profileBasePath

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	profile := new(Profile)
	resp, err := s.client.DoAuthenticated(ctx, req, profile)
	if err != nil {
		return nil, resp, err
	}
	return profile, resp, err
}

func (s *ProfileServiceOp) GetID(ctx context.Context) (string, error) {
	profile, _, err := s.Get(ctx)
	if err != nil {
		return "", err
	}

	return profile.ID, err
}