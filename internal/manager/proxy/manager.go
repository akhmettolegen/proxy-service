package proxy

import (
	"context"
	"github.com/akhmettolegen/test-service/internal/models"
)

type Manager struct {
	ctx context.Context
}

func NewManager(ctx context.Context) *Manager {
	return &Manager{
		ctx: ctx,
	}
}

func (m *Manager) ProxyRequest(req *models.ProxyRequest) (*models.ProxyResponse, error) {
	return &models.ProxyResponse{
		Id:      "responseId",
		Status:  "NEW",
		Headers: nil,
		Length:  10,
	}, nil
}
