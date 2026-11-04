package http

import (
	"Go_Commerce/internal/adapters/incoming/http/handlers"
	"Go_Commerce/internal/infrastructure/config"

	nethttp "net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humamux"
	"github.com/gorilla/mux"
)

func New(cfg *config.Config, buyerH *handlers.BuyerHandler, sellerHandler *handlers.SellerHandler, orderH *handlers.OrderHandler) nethttp.Handler {
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
		OperationID:   "register-buyer",
		Method:        nethttp.MethodPost,
		Path:          "/buyers",
		Summary:       "Register a buyer",
		DefaultStatus: nethttp.StatusCreated,
	}, buyerH.Register)

	huma.Register(api, huma.Operation{
		OperationID: "get-buyer",
		Method:      nethttp.MethodGet,
		Path:        "/buyers/{id}",
		Summary:     "Get buyer by ID",
	}, buyerH.Get)

	huma.Register(api, huma.Operation{
		OperationID: "update-buyer",
		Method:      nethttp.MethodPatch,
		Path:        "/buyers/{id}",
		Summary:     "Update a buyer (only the fields you send)",
	}, buyerH.Update)

	huma.Register(api, huma.Operation{
		OperationID:   "delete-buyer",
		Method:        nethttp.MethodDelete,
		Path:          "/buyers/{id}",
		Summary:       "Delete a buyer",
		DefaultStatus: nethttp.StatusNoContent,
	}, buyerH.Delete)

	huma.Register(api, huma.Operation{
		OperationID:   "register-seller",
		Method:        nethttp.MethodPost,
		Path:          "/sellers",
		Summary:       "Register a seller",
		DefaultStatus: nethttp.StatusCreated,
		Tags:          []string{"Sellers"},
	}, sellerHandler.Register)

	return r
}
