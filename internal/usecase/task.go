package usecase

import (
	"context"
	"encoding/json"
	"github.com/akhmettolegen/proxy-service/internal/entity"
	"github.com/akhmettolegen/proxy-service/internal/repo"
	"github.com/akhmettolegen/proxy-service/internal/service"
	"github.com/akhmettolegen/proxy-service/pkg/logger"
	"github.com/google/uuid"
	"net/http"
)

type TaskUseCase struct {
	repo       repo.TaskRepo
	httpClient service.Service
}

func New(repo repo.TaskRepo, cli service.Service) *TaskUseCase {
	return &TaskUseCase{repo: repo, httpClient: cli}
}

func (u *TaskUseCase) Create(ctx context.Context, req entity.TaskRequest, l logger.Interface) (string, error) {
	taskId := uuid.NewString()

	go u.processRequest(ctx, taskId, req, l)

	return taskId, nil
}

func (u *TaskUseCase) processRequest(ctx context.Context, taskId string, req entity.TaskRequest, l logger.Interface) {
	task := entity.Task{
		Id:     taskId,
		Status: entity.TaskStatusInProcess,
	}

	if err := u.repo.Store(ctx, task); err != nil {
		l.Error("store task error %v", err)
		return
	}

	reqByte, err := json.Marshal(&req.Body)
	if err != nil {
		l.Error("marshal request error, %v", err)
		return
	}

	res, err := u.httpClient.Request(ctx, req.Method, req.Url, req.Headers, reqByte)
	if err != nil {
		l.Error("http request error %v", err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK || res.StatusCode == http.StatusCreated {
		task.Status = entity.TaskStatusDone
		task.Headers = res.Header
		task.Length = res.ContentLength
		task.HttpStatusCode = res.StatusCode
	} else {
		task.Status = entity.TaskStatusError
	}

	if err := u.repo.Store(ctx, task); err != nil {
		l.Error("store task error %v", err)
		return
	}
}

func (u *TaskUseCase) GetById(ctx context.Context, id string) (entity.Task, error) {
	return u.repo.GetById(ctx, id)
}
