package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/akhmettolegen/proxy-service/internal/entity"
	"github.com/akhmettolegen/proxy-service/internal/repo"
	"github.com/akhmettolegen/proxy-service/internal/service"
	"github.com/akhmettolegen/proxy-service/pkg/logger"
	"github.com/google/uuid"
	"net/http"
	"sync"
)

var ErrTaskNotFound = errors.New("task not found")

type TaskUseCase struct {
	repo       repo.TaskRepo
	httpClient service.Service
	swg        *sync.WaitGroup
}

func New(repo repo.TaskRepo, cli service.Service, swg *sync.WaitGroup) *TaskUseCase {
	return &TaskUseCase{repo: repo, httpClient: cli, swg: swg}
}

func (u *TaskUseCase) Create(ctx context.Context, req entity.TaskRequest, l logger.Interface) (string, error) {
	taskId := uuid.NewString()

	go u.processRequest(ctx, taskId, req, l)

	return taskId, nil
}

func (u *TaskUseCase) processRequest(ctx context.Context, taskId string, req entity.TaskRequest, l logger.Interface) {
	u.swg.Add(1)
	defer u.swg.Done()

	task := entity.Task{
		Id:     taskId,
		Status: entity.TaskStatusInProcess,
	}

	if err := u.repo.Store(ctx, task); err != nil {
		l.Error("store task error %v", err)
		return
	}

	var (
		reqByte []byte
		err     error
	)
	if req.Body != nil {
		reqByte, err = json.Marshal(&req.Body)
		if err != nil {
			l.Error("marshal request error, %v", err)
			return
		}
	}

	res, err := u.httpClient.Request(req.Method, req.Url, req.Headers, reqByte)
	if err != nil {
		l.Error("http request error %v", err)
		return
	}
	defer res.Body.Close()

	task.Length = res.ContentLength
	task.Headers = res.Header
	task.HttpStatusCode = res.StatusCode

	if res.StatusCode == http.StatusOK || res.StatusCode == http.StatusCreated {
		task.Status = entity.TaskStatusDone
	} else {
		task.Status = entity.TaskStatusError
	}

	if err := u.repo.Store(ctx, task); err != nil {
		l.Error("store task error %v", err)
		return
	}
}

func (u *TaskUseCase) GetById(ctx context.Context, id string) (entity.Task, error) {
	task, err := u.repo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrTaskNotFound) {
			return entity.Task{}, ErrTaskNotFound
		}

		return entity.Task{}, err
	}

	return task, nil
}
