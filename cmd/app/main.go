package main

import "github.com/akhmettolegen/proxy-service/internal/server"

// @title Proxy API
// @version 1.0
//
// @BasePath /api/v1
// @Description HTTP server for proxying HTTP-requests to 3rd-party services

func main() {
	server.Start()
}
