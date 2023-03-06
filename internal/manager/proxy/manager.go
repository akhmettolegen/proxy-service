package proxy

import (
	"context"
	"encoding/json"
	"errors"
	httpCli "github.com/akhmettolegen/test-service/internal/client/http"
	"github.com/akhmettolegen/test-service/internal/models"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Manager struct {
	ctx context.Context
}

func NewManager(ctx context.Context) *Manager {
	return &Manager{
		ctx: ctx,
	}
}

func (m *Manager) ProxyRequest(req *models.ProxyRequest) (*models.ProxyResponse, error) {
	reqByte, err := json.Marshal(&req.Body)
	if err != nil {
		return nil, err
	}

	resp, err := httpCli.Request(m.ctx, req.Method, req.Url, req.Headers, reqByte)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		errTxt := "ProxyRequest error: code=" + strconv.Itoa(resp.StatusCode) + " message=" + resp.Status
		rawBody, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			errResponse := new(ErrorResponse)
			if err = json.Unmarshal(rawBody, errResponse); err == nil && errResponse.Error != nil {
				errTxt = "ProxyRequest error:" + errResponse.Error.Message
			}
		}
		log.Println("[ERROR]", errTxt)
		return nil, errors.New(errTxt)
	}

	taskId := uuid.NewString()

	return &models.ProxyResponse{
		Id:      taskId,
		Status:  resp.Status,
		Headers: resp.Header,
		Length:  resp.ContentLength,
	}, nil
}

type ErrorResponse struct {
	Error *ErrorStatus `json:"error"`
}

type ErrorStatus struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}
