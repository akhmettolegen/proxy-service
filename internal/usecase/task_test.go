package usecase_test

import (
	"bytes"
	"context"
	"errors"
	"github.com/akhmettolegen/proxy-service/internal/entity"
	repoMocks "github.com/akhmettolegen/proxy-service/internal/repo/mocks"
	serviceMocks "github.com/akhmettolegen/proxy-service/internal/service/mocks"
	"github.com/akhmettolegen/proxy-service/internal/usecase"
	"github.com/akhmettolegen/proxy-service/pkg/logger"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"sync"
	"testing"
)

var errInternalServErr = errors.New("internal server error")

type test struct {
	name string
	mock func()
	res  interface{}
	err  error
}

func task(t *testing.T) (*usecase.TaskUseCase, *repoMocks.MockTaskRepo, *serviceMocks.MockService, logger.Interface) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	repo := repoMocks.NewMockTaskRepo(mockCtl)
	service := serviceMocks.NewMockService(mockCtl)
	swg := &sync.WaitGroup{}

	task := usecase.New(repo, service, swg)

	log := logger.New("debug")

	return task, repo, service, log
}

func TestCreate(t *testing.T) {
	t.Parallel()

	task, repo, service, log := task(t)

	json := `{"name":"Test Name","full_name":"test full name"}`
	r := io.NopCloser(bytes.NewReader([]byte(json)))

	tests := []test{
		{
			name: "random string result",
			mock: func() {
				repo.EXPECT().Store(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
				service.EXPECT().Request(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&http.Response{
					StatusCode: 200,
					Body:       r,
				}, nil).AnyTimes()
			},
			res: "b0a5e2ec-f59a-4b58-a129-31ebcd52be60",
			err: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()
			got, err := task.Create(context.Background(), entity.TaskRequest{}, log)

			assert.Len(t, tc.res, len(got))
			require.ErrorIs(t, err, tc.err)
		})
	}
}

func TestGetById(t *testing.T) {
	t.Parallel()

	task, repo, _, _ := task(t)

	tests := []test{
		{
			name: "empty result",
			mock: func() {
				repo.EXPECT().GetById(gomock.Any(), gomock.Any()).Return(entity.Task{}, nil)
			},
			res: entity.Task{},
			err: nil,
		},
		{
			name: "result with error",
			mock: func() {
				repo.EXPECT().GetById(gomock.Any(), gomock.Any()).Return(entity.Task{}, errInternalServErr)
			},
			res: entity.Task{},
			err: errInternalServErr,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()
			got, err := task.GetById(context.Background(), "")

			require.Equal(t, tc.res, got)
			require.ErrorIs(t, err, tc.err)
		})
	}
}
