package anypointplatform

import (
	"context"
	"fmt"
	"net/http"
)

const orgBasePath = "accounts/api/organizations"

type MasterOrganizationService interface {
	Get(context.Context) (*MasterOrganization, *Response, error)
}

type MasterOrganizationServiceOp struct {
	client *Client
}

type MasterOrganization struct {
	ID       string
	Name     string
	ClientID string
}

func (s *MasterOrganizationServiceOp) Get(ctx context.Context) (*MasterOrganization, *Response, error) {
	p, _, err := s.client.Profile.Get(ctx)
	if err != nil {
		return nil, nil, err
	}
	path := fmt.Sprintf("%s/%s", orgBasePath, p.OrganizationID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	org := new(MasterOrganization)
	resp, err := s.client.DoAuthenticated(ctx, req, org)
	if err != nil {
		return nil, resp, err
	}
	return org, resp, err
}

