package v1

import (
	"errors"
	"github.com/akhmettolegen/proxy-service/internal/entity"
)

type taskCreateRequest struct {
	Method  string            `json:"method"`
	Url     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}

func toTaskRequest(in taskCreateRequest) entity.TaskRequest {
	return entity.TaskRequest{
		Method:  in.Method,
		Url:     in.Url,
		Headers: in.Headers,
	}
}

func (r *taskCreateRequest) validate() error {
	if r.Method == "" {
		return errors.New("method is empty")
	}
	if r.Url == "" {
		return errors.New("url is empty")
	}

	return nil
}

type taskCreateResponse struct {
	Id string `json:"id"`
}

func toCreateResponse(id string) taskCreateResponse {
	return taskCreateResponse{Id: id}
}

type taskByIdResponse struct {
	Id             string              `json:"id"`
	Status         string              `json:"status"`
	HttpStatusCode int                 `json:"httpStatusCode"`
	Headers        map[string][]string `json:"headers"`
	Length         int64               `json:"length"`
}

func toTaskByIdResponse(in entity.Task) taskByIdResponse {
	return taskByIdResponse{
		Id:             in.Id,
		Status:         in.Status,
		HttpStatusCode: in.HttpStatusCode,
		Headers:        in.Headers,
		Length:         in.Length,
	}
}
