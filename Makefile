run:
	go mod tidy && go mod download && \
	go run cmd/app/main.go

mock:
	go install github.com/golang/mock/mockgen@latest
	mockgen -source ./internal/usecase/usecase.go -package usecase_test > ./internal/usecase/mocks/mocks.go
	mockgen -source ./internal/repo/repo.go -package repo_test > ./internal/repo/mocks/mocks.go
	mockgen -source ./internal/service/service.go -package service_test > ./internal/service/mocks/mocks.go

test:
	go test -v -cover -race ./internal/...