all:
	go mod tidy && go mod download && \
	go run cmd/app/main.go