package anypointplatform

import (
	"context"
	"fmt"
	"net/http"
)

const bgBasePath = "accounts/api/organizations"

type BusinessGroupService interface {
	Get(context.Context, string) (*BusinessGroup, *Response, error)
	Create(context.Context, *BusinessGroupCreateRequest) (*BusinessGroup, *Response, error)
	CreateWithName(context.Context, string) (*BusinessGroup, error)
	Delete(context.Context, string) (*Response, error)
}

type BusinessGroupServiceOp struct {
	client *Client
}

type BusinessGroup struct {
	ID       string
	Name     string
	ClientID string
}

type BusinessGroupCreateRequest struct {
	Name                 string
	OwnerID              string
	ParentID string `json:parentOrganizationId`
}

func (s *BusinessGroupServiceOp) Get(ctx context.Context, id string) (*BusinessGroup, *Response, error) {
	path := fmt.Sprintf("%s/%s", bgBasePath, id)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	bg := new(BusinessGroup)
	resp, err := s.client.DoAuthenticated(ctx, req, bg)
	if err != nil {
		return nil, resp, err
	}
	return bg, resp, err
}

func (s *BusinessGroupServiceOp) Create(ctx context.Context, createRequest *BusinessGroupCreateRequest) (*BusinessGroup, *Response, error) {
	path := bgBasePath

	if createRequest.OwnerID == "" {
		p, _, err := s.client.Profile.Get(ctx)
		if err != nil {
			return nil, nil, err
		}
		createRequest.OwnerID = p.ID
	}

	if createRequest.ParentID == "" {
		p, _, err := s.client.Profile.Get(ctx)
		if err != nil {
			return nil, nil, err
		}
		createRequest.ParentID = p.OrganizationID
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	bg := new(BusinessGroup)
	resp, err := s.client.Do(ctx, req, bg)
	if err != nil {
		return nil, resp, err
	}

	return bg, resp, err
}

func (s *BusinessGroupServiceOp) CreateWithName(ctx context.Context, name string) (*BusinessGroup, error) {
	req := &BusinessGroupCreateRequest{Name: name}
	resp, _, err := s.Create(ctx, req)
	return resp, err
}

func (s *BusinessGroupServiceOp) Delete(ctx context.Context, id string) (*Response, error) {
	path := fmt.Sprintf("%s/%s", bgBasePath, id)
	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
