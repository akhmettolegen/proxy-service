package service

import (
	"context"
	"net/http"
)

type Service interface {
	Request(ctx context.Context, method string, url string, headers map[string]string, body []byte) (*http.Response, error)
}
