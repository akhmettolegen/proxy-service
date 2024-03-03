package usecase

import (
	"context"
	"github.com/akhmettolegen/proxy-service/internal/entity"
	"github.com/akhmettolegen/proxy-service/pkg/logger"
)

type Task interface {
	Create(ctx context.Context, task entity.TaskRequest, l logger.Interface) (string, error)
	GetById(ctx context.Context, id string) (entity.Task, error)
}
