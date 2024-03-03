package repo

import (
	"context"
	"github.com/akhmettolegen/proxy-service/internal/entity"
)

type TaskRepo interface {
	Store(ctx context.Context, task entity.Task) error
	GetById(ctx context.Context, id string) (entity.Task, error)
}
