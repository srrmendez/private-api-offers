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
		Name:       "Secondary Offer",
		Pattern:    "/v1/secondary",
		HandleFunc: getSecondaryOffers,
		Method:     http.MethodGet,
		ShouldLog:  true,
	},
	{
		Name:       "Get Offer",
		Pattern:    "/v1/{id}",
		HandleFunc: getOffer,
		Method:     http.MethodGet,
		ShouldLog:  true,
	},
	{
		Name:       "Search Offers",
		Pattern:    "/v1/",
		HandleFunc: searchOffers,
		Method:     http.MethodGet,
		ShouldLog:  true,
	},
	{
		Name:       "Sync Offers",
		Pattern:    "/v1/",
		HandleFunc: createOffers,
		Method:     http.MethodPost,
		ShouldLog:  true,
	},
}
