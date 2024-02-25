package usecase

import "github.com/akhmettolegen/proxy-service/internal/entity"

type Task interface {
	Store(task entity.Task) (string, error)
	GetById(id string) (entity.Task, error)
}
