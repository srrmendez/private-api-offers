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
	{
		Name:       "Create Commercial System Service Order",
		Pattern:    "/v1/service",
		HandleFunc: createServiceOrder,
		Method:     http.MethodPost,
		ShouldLog:  true,
	},
	{
		Name:       "Get Order",
		Pattern:    "/v1/{id}",
		HandleFunc: getOrder,
		Method:     http.MethodGet,
		ShouldLog:  true,
	},
	{
		Name:       "Update Order",
		Pattern:    "/v1/{id}",
		HandleFunc: updateOrder,
		Method:     http.MethodPut,
		ShouldLog:  true,
	},
	{
		Name:       "Search Orders",
		Pattern:    "/v1/",
		HandleFunc: searchOrders,
		Method:     http.MethodGet,
		ShouldLog:  true,
	},
	{
		Name:       "Create Order",
		Pattern:    "/v1/",
		HandleFunc: createOrder,
		Method:     http.MethodPost,
		ShouldLog:  true,
	},
}
