package repo

import "github.com/akhmettolegen/proxy-service/internal/entity"

type TaskRepo interface {
	Store(task entity.Task) error
	GetById(id string) (entity.Task, error)
}
