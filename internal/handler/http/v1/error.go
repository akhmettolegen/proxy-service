package v1

import (
	"github.com/go-chi/render"
	"net/http"
)

type response struct {
	ErrMessage     string `json:"error"`
	HTTPStatusCode int    `json:"-"`
}

func (e *response) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func errResponse(code int, msg string) render.Renderer {
	return &response{
		ErrMessage:     msg,
		HTTPStatusCode: code,
	}
}
