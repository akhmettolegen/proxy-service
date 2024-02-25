package usecase

import (
	"github.com/akhmettolegen/proxy-service/internal/entity"
	"github.com/akhmettolegen/proxy-service/internal/repo"
	"github.com/google/uuid"
)

type TaskUseCase struct {
	repo repo.TaskRepo
}

func New(repo repo.TaskRepo) *TaskUseCase {
	return &TaskUseCase{repo: repo}
}

func (u *TaskUseCase) Store() (string, error) {
	taskId := uuid.NewString()

	return taskId, nil
}

func (u *TaskUseCase) GetById(id string) (entity.Task, error) {
	return u.repo.GetById(id)
}
