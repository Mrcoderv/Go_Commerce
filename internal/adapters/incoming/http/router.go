package http

import (
	"Go_Commerce/internal/adapters/incoming/http/handlers"
	"Go_Commerce/internal/infrastructure/config"

	nethttp "net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humamux"
	"github.com/gorilla/mux"
)

// New builds a Huma-powered HTTP handler. It registers GET /health (returns 200)
// and exposes the OpenAPI spec at /openapi.json plus the docs UI at /docs.
func New(cfg *config.Config) nethttp.Handler {
	r := mux.NewRouter()
	adapter := humamux.NewAdapter(r)

	humaCfg := huma.DefaultConfig("Go Commerce API", "1.0.0")
	humaCfg.OpenAPIPath = "/openapi"
	humaCfg.DocsPath = "/docs"

	api := huma.NewAPI(humaCfg, adapter)

	huma.Register(api, huma.Operation{
		OperationID: "health",
		Method:      nethttp.MethodGet,
		Path:        "/health",
		Summary:     "Check API health",
	}, (&handlers.Handler{}).Health)

	return r
}
