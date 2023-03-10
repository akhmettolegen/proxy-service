package proxy

import (
	"context"
	"encoding/json"
	"errors"
	restCli "github.com/akhmettolegen/proxy-service/internal/clients"
	"github.com/akhmettolegen/proxy-service/internal/managers"
	"github.com/akhmettolegen/proxy-service/internal/models"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type Manager struct {
	ctx         context.Context
	HttpClient  restCli.HttpClient
	RequestsMap RequestsMap
}

type RequestsMap struct {
	mu     *sync.RWMutex
	reqMap map[string]string
}

func NewManager(ctx context.Context, cli restCli.HttpClient, rMap map[string]string, mu *sync.RWMutex) managers.ProxyManager {
	return &Manager{
		ctx:        ctx,
		HttpClient: cli,
		RequestsMap: RequestsMap{
			mu:     mu,
			reqMap: rMap,
		},
	}
}

func (m *Manager) ProxyRequest(req *models.ProxyRequest) (*models.ProxyResponse, error) {
	reqByte, err := json.Marshal(&req.Body)
	if err != nil {
		return nil, err
	}

	resp, err := m.HttpClient.Request(req.Method, req.Url, req.Headers, reqByte)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		errTxt := "ProxyRequest error: code=" + strconv.Itoa(resp.StatusCode) + " message=" + resp.Status
		rawBody, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			errResponse := new(models.ErrorResponse)
			if err = json.Unmarshal(rawBody, errResponse); err == nil && errResponse.Error != nil {
				errTxt = "ProxyRequest error:" + errResponse.Error.Message
			}
		}
		log.Println("[ERROR]", errTxt)
		return nil, errors.New(errTxt)
	}

	taskId := uuid.NewString()

	go m.AddRequestToMap(taskId, resp.Status)

	return &models.ProxyResponse{
		Id:      taskId,
		Status:  resp.Status,
		Headers: resp.Header,
		Length:  resp.ContentLength,
	}, nil
}

func (m *Manager) AddRequestToMap(taskId, status string) {
	m.RequestsMap.mu.Lock()
	defer m.RequestsMap.mu.Unlock()
	m.RequestsMap.reqMap[taskId] = status
}
