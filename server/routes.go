package server

import (
	"net/http"

	pkgHttp "github.com/srrmendez/services-interface-tools/pkg/http"
)

var Routes = pkgHttp.Routes{
	{
		Name:       "Health Check",
		Pattern:    "/health-check/",
		HandleFunc: healthCheck,
		Method:     http.MethodGet,
		ShouldLog:  true,
	},
}
