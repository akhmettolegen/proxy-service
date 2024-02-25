package usecase

import "github.com/akhmettolegen/proxy-service/internal/entity"

type Task interface {
	Create(task entity.TaskRequest) (string, error)
	GetById(id string) (entity.Task, error)
}
