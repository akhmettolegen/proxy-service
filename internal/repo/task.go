package repo

import (
	"context"
	"errors"
	"github.com/akhmettolegen/proxy-service/internal/entity"
	"sync"
)

var ErrTaskNotFound = errors.New("task not found")

type TaskStorage struct {
	storage map[string]entity.Task
	mu      *sync.RWMutex
}

func New(storage map[string]entity.Task, mu *sync.RWMutex) *TaskStorage {
	return &TaskStorage{
		storage: storage,
		mu:      mu,
	}
}

func (r *TaskStorage) Store(ctx context.Context, task entity.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.storage[task.Id] = task

	return nil
}

func (r *TaskStorage) GetById(ctx context.Context, id string) (entity.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	task, ok := r.storage[id]
	if !ok {
		return entity.Task{}, ErrTaskNotFound
	}

	return task, nil
}
