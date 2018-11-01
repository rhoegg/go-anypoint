package go_anypoint

import (
	"context"
	"fmt"
	"net/http"
)

const bgBasePath = "accounts/api/organizations"

type BusinessGroupService interface {
	Get(context.Context, string) (*BusinessGroup, *Response, error)
}

type BusinessGroupServiceOp struct {
	client *Client
}

type BusinessGroup struct {
	Name string
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
