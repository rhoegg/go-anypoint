package go_anypoint

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

const defaultBaseURL = "https://anypoint.mulesoft.com"

type Client struct {
	client *http.Client

	BaseURL *url.URL

	BusinessGroup 	BusinessGroupService
}

type Response struct {
	*http.Response
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{client: httpClient, BaseURL: baseURL}

	c.BusinessGroup = &BusinessGroupServiceOp{client: c}

	return c
}

func (c *Client) NewRequest(ctx context.Context, method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	buf := new(bytes.Buffer)

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	resp, err := DoRequestWithClient(ctx, c.client, req)
	if err != nil {
		return nil, err
	}

	defer func() {
		if rerr := resp.Body.Close(); err == nil {
			err = rerr
		}
	}()

	response := &Response{Response: resp}

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
		if err != nil {
			return nil, err
		}
	}

	return response, err
}

func DoRequestWithClient(ctx context.Context, client *http.Client, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	return client.Do(req)
}