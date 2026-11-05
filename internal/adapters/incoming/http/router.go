package http

import (
	"Go_Commerce/internal/adapters/incoming/http/handlers"
	"Go_Commerce/internal/adapters/incoming/http/middleware"
	"Go_Commerce/internal/infrastructure/config"
	nethttp "net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humamux"
	"github.com/gorilla/mux"
)

func New(
	cfg *config.Config,
	authH *handlers.AuthHandler,
	buyerH *handlers.BuyerHandler,
	sellerH *handlers.SellerHandler,
	productH *handlers.ProductHandler,
	orderH *handlers.OrderHandler,
) nethttp.Handler {
	r := mux.NewRouter()

	r.Use(middleware.RecoverMiddleware)
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.AuthMiddleware(cfg))

	adapter := humamux.NewAdapter(r)
	humaCfg := huma.DefaultConfig("Go_Commerce API", "1.0.0")
	humaCfg.OpenAPIPath = "/openapi"
	humaCfg.DocsPath = "/docs"
	api := huma.NewAPI(humaCfg, adapter)

	huma.Register(api, huma.Operation{OperationID: "health", Method: nethttp.MethodGet, Path: "/health", Summary: "Check API health"}, (&handlers.Handler{}).Health)

	huma.Register(api, huma.Operation{OperationID: "auth-register", Method: nethttp.MethodPost, Path: "/auth/register", Summary: "Register seller or buyer", DefaultStatus: nethttp.StatusCreated, Tags: []string{"Auth"}}, authH.Register)
	huma.Register(api, huma.Operation{OperationID: "auth-login", Method: nethttp.MethodPost, Path: "/auth/login", Summary: "Login seller or buyer", Tags: []string{"Auth"}}, authH.Login)

	huma.Register(api, huma.Operation{OperationID: "register-buyer", Method: nethttp.MethodPost, Path: "/buyers", Summary: "Register a buyer (legacy)", DefaultStatus: nethttp.StatusCreated, Tags: []string{"Buyers"}}, buyerH.Register)
	huma.Register(api, huma.Operation{OperationID: "list-buyers", Method: nethttp.MethodGet, Path: "/buyers", Summary: "List buyers", Tags: []string{"Buyers"}}, buyerH.List)
	huma.Register(api, huma.Operation{OperationID: "get-buyer", Method: nethttp.MethodGet, Path: "/buyers/{id}", Summary: "Get buyer by ID", Tags: []string{"Buyers"}}, buyerH.Get)
	huma.Register(api, huma.Operation{OperationID: "update-buyer", Method: nethttp.MethodPatch, Path: "/buyers/{id}", Summary: "Update buyer", Tags: []string{"Buyers"}}, buyerH.Update)
	huma.Register(api, huma.Operation{OperationID: "delete-buyer", Method: nethttp.MethodDelete, Path: "/buyers/{id}", Summary: "Delete buyer", DefaultStatus: nethttp.StatusNoContent, Tags: []string{"Buyers"}}, buyerH.Delete)

	huma.Register(api, huma.Operation{OperationID: "register-seller", Method: nethttp.MethodPost, Path: "/sellers", Summary: "Register a seller (legacy)", DefaultStatus: nethttp.StatusCreated, Tags: []string{"Sellers"}}, sellerH.Register)
	huma.Register(api, huma.Operation{OperationID: "list-sellers", Method: nethttp.MethodGet, Path: "/sellers", Summary: "List sellers", Tags: []string{"Sellers"}}, sellerH.List)
	huma.Register(api, huma.Operation{OperationID: "get-seller", Method: nethttp.MethodGet, Path: "/sellers/{id}", Summary: "Get seller by ID", Tags: []string{"Sellers"}}, sellerH.GetByID)
	huma.Register(api, huma.Operation{OperationID: "update-seller", Method: nethttp.MethodPatch, Path: "/sellers/{id}", Summary: "Update seller", Tags: []string{"Sellers"}}, sellerH.Update)
	huma.Register(api, huma.Operation{OperationID: "delete-seller", Method: nethttp.MethodDelete, Path: "/sellers/{id}", Summary: "Delete seller", DefaultStatus: nethttp.StatusNoContent, Tags: []string{"Sellers"}}, sellerH.Delete)

	huma.Register(api, huma.Operation{OperationID: "create-product", Method: nethttp.MethodPost, Path: "/products", Summary: "Create a product (seller auth)", DefaultStatus: nethttp.StatusCreated, Tags: []string{"Products"}}, productH.Create)
	huma.Register(api, huma.Operation{OperationID: "list-products", Method: nethttp.MethodGet, Path: "/products", Summary: "List products", Tags: []string{"Products"}}, productH.List)
	huma.Register(api, huma.Operation{OperationID: "get-product", Method: nethttp.MethodGet, Path: "/products/{id}", Summary: "Get a product by ID", Tags: []string{"Products"}}, productH.GetByID)
	huma.Register(api, huma.Operation{OperationID: "update-product", Method: nethttp.MethodPatch, Path: "/products/{id}", Summary: "Update a product (seller auth)", Tags: []string{"Products"}}, productH.Update)
	huma.Register(api, huma.Operation{OperationID: "delete-product", Method: nethttp.MethodDelete, Path: "/products/{id}", Summary: "Delete a product (seller auth)", DefaultStatus: nethttp.StatusNoContent, Tags: []string{"Products"}}, productH.Delete)

	huma.Register(api, huma.Operation{OperationID: "place-order", Method: nethttp.MethodPost, Path: "/orders", Summary: "Place an order (buyer auth)", DefaultStatus: nethttp.StatusCreated, Tags: []string{"Orders"}}, orderH.PlaceOrder)
	huma.Register(api, huma.Operation{OperationID: "get-order", Method: nethttp.MethodGet, Path: "/orders/{id}", Summary: "Get order by ID", Tags: []string{"Orders"}}, orderH.GetOrder)
	huma.Register(api, huma.Operation{OperationID: "update-order-status", Method: nethttp.MethodPatch, Path: "/orders/{id}/status", Summary: "Update order status", Tags: []string{"Orders"}}, orderH.UpdateStatus)
	huma.Register(api, huma.Operation{OperationID: "cancel-order", Method: nethttp.MethodPost, Path: "/orders/{id}/cancel", Summary: "Cancel an order", Tags: []string{"Orders"}}, orderH.Cancel)
	huma.Register(api, huma.Operation{OperationID: "list-buyer-orders", Method: nethttp.MethodGet, Path: "/buyers/{id}/orders", Summary: "List orders for a buyer", Tags: []string{"Orders"}}, orderH.ListByBuyer)
	huma.Register(api, huma.Operation{OperationID: "list-seller-orders", Method: nethttp.MethodGet, Path: "/sellers/{id}/orders", Summary: "List orders for a seller", Tags: []string{"Orders"}}, orderH.ListBySeller)

	return r
}
