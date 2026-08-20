package health

import (
	"net/http"

	"github.com/ThatSoftwareCompany/template-go-api/internal/platform/httpserver"
)

func RegisterRoutes(mux *http.ServeMux, service *Service) {
	controller := NewController(service)
	mux.Handle("/api/v1/health", httpserver.OnlyMethods(http.MethodGet)(http.HandlerFunc(controller.Handle)))
}
