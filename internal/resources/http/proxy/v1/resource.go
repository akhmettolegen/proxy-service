package v1

import (
	"encoding/json"
	"github.com/akhmettolegen/test-service/internal/manager/proxy"
	"github.com/akhmettolegen/test-service/internal/models"
	"github.com/akhmettolegen/test-service/internal/models/httperrors"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
)

type ProxyResource struct {
	ProxyManager *proxy.Manager
}

func (rs ProxyResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {

		r.Post("/", rs.proxyRequest)
	})

	return r
}

func (rs ProxyResource) proxyRequest(w http.ResponseWriter, r *http.Request) {

	var req *models.ProxyRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = render.Render(w, r, httperrors.BadRequest(err))
		return
	}

	result, err := rs.ProxyManager.ProxyRequest(req)
	if err != nil {
		_ = render.Render(w, r, httperrors.Internal(err))
		return
	}

	render.JSON(w, r, result)
}