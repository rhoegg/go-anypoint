package go_anypoint

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	defaultBaseURL = "https://anypoint.mulesoft.com"
	mediaType = "application/json"
)

type Client struct {
	client *http.Client
	accessToken string

	BaseURL *url.URL
	Username string
	Password string

	BusinessGroup BusinessGroupService
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
	if body != nil {
		err = json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", mediaType)

	return req, nil
}

func (c *Client) DoAuthenticated(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	token, err := c.getAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorizationx", fmt.Sprintf("Bearer %v", token))
	return c.Do(ctx, req, v)
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

	err = CheckResponse(resp)
	if (err != nil) {
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
			if err != nil {
				return nil, err
			}
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err != nil {
				return nil, err
			}
		}
	}

	return response, err
}

func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	data, _ := ioutil.ReadAll(r.Body)

	return errors.New(fmt.Sprintf("(%v) %v", r.StatusCode, string(data)))
}

func (c *Client) getAccessToken(ctx context.Context) (string, error) {
	if c.accessToken == "" {
		result, _, err := c.Login(ctx, &LoginRequest{Username: c.Username, Password: c.Password})
		if err != nil {
			return "", err
		}
		c.accessToken = result.AccessToken
	}
	return c.accessToken, nil
}

func DoRequestWithClient(ctx context.Context, client *http.Client, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	return client.Do(req)
}