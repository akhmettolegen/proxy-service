package v1

import (
	"encoding/json"
	"errors"
	"github.com/akhmettolegen/proxy-service/internal/usecase"
	"github.com/akhmettolegen/proxy-service/pkg/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
)

type TaskHandler struct {
	TaskUseCase usecase.Task
	logger      logger.Interface
}

func NewTaskHandler(taskUseCase usecase.Task, logger logger.Interface) *TaskHandler {
	return &TaskHandler{
		TaskUseCase: taskUseCase,
		logger:      logger,
	}
}

func (h TaskHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Post("/", h.create)
		r.Get("/{id}", h.getById)
	})

	return r
}

// @Tags Task
// @Description Create task
// @Accept  json
// @Produce  json
// @Param CreateRequest body taskCreateRequest true "Request"
// @Success 200 {object} taskCreateResponse
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /task [post]
func (h TaskHandler) create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req taskCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error(err, "http - v1 - create")
		_ = render.Render(w, r, errResponse(http.StatusBadRequest, "invalid request body"))
		return
	}

	if err := req.validate(); err != nil {
		h.logger.Error(err, "http - v1 - create")
		_ = render.Render(w, r, errResponse(http.StatusBadRequest, "invalid request body"))
		return
	}

	result, err := h.TaskUseCase.Create(ctx, toTaskRequest(req), h.logger)
	if err != nil {
		h.logger.Error(err, "http - v1 - create")
		_ = render.Render(w, r, errResponse(http.StatusInternalServerError, "internal error"))
		return
	}

	render.JSON(w, r, toCreateResponse(result))
}

// @Tags Task
// @Description Get task by id
// @Accept  json
// @Produce  json
// @Param id path string true "Task id"
// @Success 200 {object} taskByIdResponse
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /task/{id} [get]
func (h TaskHandler) getById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	if id == "" {
		h.logger.Error(errors.New("id is empty"), "http - v1 - getById")
		_ = render.Render(w, r, errResponse(http.StatusBadRequest, "invalid request body"))
		return
	}

	result, err := h.TaskUseCase.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, usecase.ErrTaskNotFound) {
			_ = render.Render(w, r, errResponse(http.StatusNotFound, "task not found"))
			return
		}

		h.logger.Error("http - v1 - getById")
		_ = render.Render(w, r, errResponse(http.StatusInternalServerError, "internal error"))
		return
	}

	render.JSON(w, r, toTaskByIdResponse(result))
}
