package anypointplatform

import (
	"context"
	"net/http"
)

const profileBasePath = "accounts/api/profile"

//go:generate pegomock generate --use-experimental-model-gen -o mock_profile_test.go --package anypointplatform_test github.com/rhoegg/go-anypoint/anypointplatform ProfileService
type ProfileService interface {
	Get(context.Context) (*Profile, *Response, error)
}

type ProfileServiceOp struct {
	client *Client
}

type Profile struct {
	ID             string
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