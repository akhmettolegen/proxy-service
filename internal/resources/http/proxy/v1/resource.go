package v1

import (
	"encoding/json"
	"github.com/akhmettolegen/proxy-service/internal/managers"
	"github.com/akhmettolegen/proxy-service/internal/models"
	"github.com/akhmettolegen/proxy-service/internal/models/httperrors"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
)

type ProxyResource struct {
	ProxyManager managers.ProxyManager
}

func (rs ProxyResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {

		r.Post("/", rs.proxyRequest)
	})

	return r
}

// @Tags proxyRequest
// @Description Proxy request
// @Accept  json
// @Produce  json
// @Param body body models.ProxyRequest true "Request"
// @Success 200 {object} models.ProxyResponse
// @Failure 400 {object} httperrors.Response
// @Failure 401 {object} httperrors.Response
// @Failure 500 {object} httperrors.Response
// @Router /proxy [post]
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
