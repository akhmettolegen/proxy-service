package service

import (
	"net/http"
)

type Service interface {
	Request(method string, url string, headers map[string]string, body []byte) (*http.Response, error)
}
