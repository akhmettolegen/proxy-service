package service

import (
	"bytes"
	"context"
	"net/http"
)

type Client struct {
	httpCli *http.Client
}

func NewClient(cli *http.Client) *Client {
	return &Client{
		httpCli: cli,
	}
}

func (c *Client) Request(ctx context.Context, method string, url string, headers map[string]string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	resp, err := c.httpCli.Do(req)

	if err != nil {
		return nil, err
	}

	return resp, err
}
