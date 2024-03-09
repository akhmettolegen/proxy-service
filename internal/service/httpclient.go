package service

import (
	"bytes"
	"context"
	"net/http"
)

type Client struct {
	ctx     context.Context
	httpCli *http.Client
}

func NewClient(ctx context.Context, cli *http.Client) *Client {
	return &Client{
		ctx:     ctx,
		httpCli: cli,
	}
}

func (c *Client) Request(method string, url string, headers map[string]string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req = req.WithContext(c.ctx)

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	return c.httpCli.Do(req)
}
