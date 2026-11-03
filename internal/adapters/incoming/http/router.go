package http

import (
	"Go_Commerce/internal/adapters/incoming/http/handlers"
	"Go_Commerce/internal/infrastructure/config"

	nethttp "net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humamux"
	"github.com/gorilla/mux"
)

func New(cfg *config.Config) nethttp.Handler {
	r := mux.NewRouter()
	adapter := humamux.NewAdapter(r)

	humaCfg := huma.DefaultConfig("this is the go defult  ", "99.99")
	humaCfg.OpenAPIPath = "/openapi"
	humaCfg.DocsPath = "/docs"

	api := huma.NewAPI(humaCfg, adapter)

	huma.Register(api, huma.Operation{
		OperationID: "health",
		Method:      nethttp.MethodGet,
		Path:        "/health",
		Summary:     "Check API health",
	}, (&handlers.Handler{}).Health)

	huma.Register(api, huma.Operation{
		OperationID: "get-order",
		Method:      nethttp.MethodGet,
		Path:        "/orders/{id}",
		Summary:     "Get order by ID",
	}, (&handlers.OrderHandler{}).GetOrder)
	//huma.Register (api, huma.Operation{
	//	OperationID: "create-order",
	//	Method:      nethttp.MethodGet,
	//	Path:        "/orders/{id}",
	//	Summary:     "Post order by ID",
	//}, (&handlers.OrderHandler{}).PostOrder)
	//
	//
	//

	return r
}
