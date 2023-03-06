package http

import (
	"github.com/akhmettolegen/test-service/internal/managers"
)

type APIServerOption func(srv *APIServer)

func WithProxyManager(proxyManager managers.ProxyManager) APIServerOption {
	return func(srv *APIServer) {
		srv.proxyManager = proxyManager
	}
}
