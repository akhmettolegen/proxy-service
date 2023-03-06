package managers

import "github.com/akhmettolegen/proxy-service/internal/models"

type ProxyManager interface {
	ProxyRequest(req *models.ProxyRequest) (*models.ProxyResponse, error)
}
