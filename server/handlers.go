package server

import (
	"net/http"

	pkgHttp "github.com/srrmendez/services-interface-tools/pkg/http"
)

// healthCheck godoc
// @Tags HealthCheck
// @Accept  json
// @Produce json
// @Success 200
// @Router /health-check/ [get]
func healthCheck(w http.ResponseWriter, r *http.Request) {
	pkgHttp.JsonResponse(w, map[string]string{"status": "Running"}, http.StatusOK)
}
