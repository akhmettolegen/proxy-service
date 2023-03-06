package managers

import "github.com/akhmettolegen/test-service/internal/models"

type ProxyManager interface {
	ProxyRequest(req *models.ProxyRequest) (*models.ProxyResponse, error)
}
