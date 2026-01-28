package http

import (
	"fmt"
	"net/http"
)

type Client struct {
	BaseURL string
}

func NewClient(baseURL string) Client {
	return Client{
		BaseURL: baseURL,
	}
}

func (c *Client) newRequest(method string, path string, body any) (*http.Request, error) {
	bodyReader, err := MarshalJsonToBody(body)

	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s%s", c.BaseURL, path)

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	return req, nil

}

func (c *Client) Post(path string, body any) (*http.Request, error) {
	return c.newRequest("POST", path, body)
}

func (c *Client) Get(path string, body any) (*http.Request, error) {
	return c.newRequest("GET", path, body)
}

func (c *Client) Delete(path string, body any) (*http.Request, error) {
	return c.newRequest("DELETE", path, body)
}

func (c *Client) Put(path string, body any) (*http.Request, error) {
	return c.newRequest("PUT", path, body)
}

func (c *Client) Patch(path string, body any) (*http.Request, error) {
	return c.newRequest("PATCH", path, body)
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	return client.Do(req)
}
