package http

import "github.com/akhmettolegen/test-service/internal/manager/proxy"

type APIServerOption func(srv *APIServer)

func WithProxyManager(proxyManager *proxy.Manager) APIServerOption {
	return func(srv *APIServer) {
		srv.proxyManager = proxyManager
	}
}
