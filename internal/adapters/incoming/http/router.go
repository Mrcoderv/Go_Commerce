package http

import (
	"fmt"
	nethttp "net/http"

	"Go_Commerce/internal/adapters/incoming/http/handlers"
	"Go_Commerce/internal/adapters/incoming/http/middleware"
	"Go_Commerce/internal/infrastructure/config"
	"Go_Commerce/internal/infrastructure/logger"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humamux"
	"github.com/gorilla/mux"
)

func New(
	cfg *config.Config,
	log logger.Logger,
	authH *handlers.AuthHandler,
	buyerH *handlers.BuyerHandler,
	sellerH *handlers.SellerHandler,
	productH *handlers.ProductHandler,
	orderH *handlers.OrderHandler,
) nethttp.Handler {
	r := mux.NewRouter()

	// Logging is outermost so panics recovered below are still logged.
	r.Use(middleware.LoggingMiddleware(log))
	r.Use(middleware.RecoverMiddleware(log))
	r.Use(middleware.AuthMiddleware(cfg))

	adapter := humamux.NewAdapter(r)
	humaCfg := huma.DefaultConfig("Go_Commerce API", "1.0.0")
	humaCfg.Info.Description = `Marketplace API for sellers, buyers, products and orders.

**Testing from this page**
1. Call **POST /auth/register** (or **POST /auth/login**) and copy the ` + "`token`" + ` from the response.
2. Open the authentication control in this page and paste the token (do not type ` + "`Bearer `" + ` yourself).
3. Call the protected endpoints. The same token is sent as ` + "`Authorization: Bearer <token>`" + `.

**Money is integer paisa**: ` + "`49950`" + ` means NPR 499.50. Fractional values like ` + "`499.50`" + ` are rejected.

Seller and buyer identity always comes from the JWT, never from the request body.`
	humaCfg.OpenAPIPath = "/openapi"
	humaCfg.DocsPath = "/docs"
	humaCfg.SchemasPath = "/schemas"
	// Keep the original Huma docs page (Stoplight Elements) instead of Swagger UI.
	humaCfg.DocsRenderer = huma.DocsRendererStoplightElements
	humaCfg.Servers = []*huma.Server{
		{
			URL:         fmt.Sprintf("http://localhost:%d", cfg.AppPort),
			Description: "Local development",
		},
	}
	humaCfg.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"bearerAuth": {
			Type:         "http",
			Scheme:       "bearer",
			BearerFormat: "JWT",
			Description:  "JWT returned by POST /auth/login or POST /auth/register. Paste only the token; the docs page adds the 'Bearer ' prefix.",
		},
	}
	humaCfg.Tags = []*huma.Tag{
		{Name: "Health", Description: "Service health."},
		{Name: "Auth", Description: "Register and log in to obtain a JWT."},
		{Name: "Sellers", Description: "Seller accounts. Identity comes from the JWT."},
		{Name: "Buyers", Description: "Buyer accounts. Identity comes from the JWT."},
		{Name: "Products", Description: "Products. Writes are restricted to the owning seller."},
		{Name: "Orders", Description: "Orders. Buyers and sellers only see their own data."},
	}
	api := huma.NewAPI(humaCfg, adapter)

	// Reusable security requirement for protected operations.
	bearer := []map[string][]string{{"bearerAuth": {}}}

	huma.Register(api, huma.Operation{
		OperationID: "health",
		Method:      nethttp.MethodGet,
		Path:        "/health",
		Summary:     "Check API health",
		Description: "Public endpoint. Returns a simple health status.",
		Tags:        []string{"Health"},
	}, (&handlers.Handler{}).Health)

	huma.Register(api, huma.Operation{
		OperationID:   "auth-register",
		Method:        nethttp.MethodPost,
		Path:          "/auth/register",
		Summary:       "Register a seller or buyer account",
		Description:   "Public endpoint. Creates the seller/buyer record plus its login user and returns a JWT. Use role `seller` or `buyer`.",
		Tags:          []string{"Auth"},
		DefaultStatus: nethttp.StatusCreated,
	}, authH.Register)

	huma.Register(api, huma.Operation{
		OperationID: "auth-login",
		Method:      nethttp.MethodPost,
		Path:        "/auth/login",
		Summary:     "Log in and get a JWT",
		Description: "Public endpoint. Returns a JWT for the given email, password and role.",
		Tags:        []string{"Auth"},
	}, authH.Login)

	huma.Register(api, huma.Operation{
		OperationID:   "register-buyer",
		Method:        nethttp.MethodPost,
		Path:          "/buyers",
		Summary:       "Register a buyer (legacy)",
		Description:   "Public endpoint. Creates a buyer without a login account. Prefer POST /auth/register with role=buyer if you need a token.",
		Tags:          []string{"Buyers"},
		DefaultStatus: nethttp.StatusCreated,
	}, buyerH.Register)

	huma.Register(api, huma.Operation{
		OperationID: "list-buyers",
		Method:      nethttp.MethodGet,
		Path:        "/buyers",
		Summary:     "List buyers",
		Description: "Public endpoint. Paginated with limit/offset.",
		Tags:        []string{"Buyers"},
	}, buyerH.List)

	huma.Register(api, huma.Operation{
		OperationID: "get-buyer",
		Method:      nethttp.MethodGet,
		Path:        "/buyers/{id}",
		Summary:     "Get buyer by ID",
		Description: "Requires a buyer JWT. The path id must match the buyer id in the token.",
		Tags:        []string{"Buyers"},
		Security:    bearer,
	}, buyerH.Get)

	huma.Register(api, huma.Operation{
		OperationID: "update-buyer",
		Method:      nethttp.MethodPatch,
		Path:        "/buyers/{id}",
		Summary:     "Update buyer",
		Description: "Requires a buyer JWT. The path id must match the buyer id in the token.",
		Tags:        []string{"Buyers"},
		Security:    bearer,
	}, buyerH.Update)

	huma.Register(api, huma.Operation{
		OperationID:   "delete-buyer",
		Method:        nethttp.MethodDelete,
		Path:          "/buyers/{id}",
		Summary:       "Delete buyer",
		Description:   "Requires a buyer JWT. Returns 409 if the buyer has existing orders.",
		Tags:          []string{"Buyers"},
		DefaultStatus: nethttp.StatusNoContent,
		Security:      bearer,
	}, buyerH.Delete)

	huma.Register(api, huma.Operation{
		OperationID:   "register-seller",
		Method:        nethttp.MethodPost,
		Path:          "/sellers",
		Summary:       "Register a seller (legacy)",
		Description:   "Public endpoint. Creates a seller without a login account. Prefer POST /auth/register with role=seller if you need a token.",
		Tags:          []string{"Sellers"},
		DefaultStatus: nethttp.StatusCreated,
	}, sellerH.Register)

	huma.Register(api, huma.Operation{
		OperationID: "list-sellers",
		Method:      nethttp.MethodGet,
		Path:        "/sellers",
		Summary:     "List sellers",
		Description: "Public endpoint. Paginated with limit/offset.",
		Tags:        []string{"Sellers"},
	}, sellerH.List)

	huma.Register(api, huma.Operation{
		OperationID: "get-seller",
		Method:      nethttp.MethodGet,
		Path:        "/sellers/{id}",
		Summary:     "Get seller by ID",
		Description: "Requires a seller JWT. The path id must match the seller id in the token.",
		Tags:        []string{"Sellers"},
		Security:    bearer,
	}, sellerH.GetByID)

	huma.Register(api, huma.Operation{
		OperationID: "update-seller",
		Method:      nethttp.MethodPatch,
		Path:        "/sellers/{id}",
		Summary:     "Update seller",
		Description: "Requires a seller JWT. The path id must match the seller id in the token.",
		Tags:        []string{"Sellers"},
		Security:    bearer,
	}, sellerH.Update)

	huma.Register(api, huma.Operation{
		OperationID:   "delete-seller",
		Method:        nethttp.MethodDelete,
		Path:          "/sellers/{id}",
		Summary:       "Delete seller",
		Description:   "Requires a seller JWT. The path id must match the seller id in the token.",
		Tags:          []string{"Sellers"},
		DefaultStatus: nethttp.StatusNoContent,
		Security:      bearer,
	}, sellerH.Delete)

	huma.Register(api, huma.Operation{
		OperationID:   "create-product",
		Method:        nethttp.MethodPost,
		Path:          "/products",
		Summary:       "Create a product",
		Description:   "Requires a seller JWT. The seller id is taken from the token; do not send seller_id in the body. Price is integer paisa.",
		Tags:          []string{"Products"},
		DefaultStatus: nethttp.StatusCreated,
		Security:      bearer,
	}, productH.Create)

	huma.Register(api, huma.Operation{
		OperationID: "list-products",
		Method:      nethttp.MethodGet,
		Path:        "/products",
		Summary:     "List products",
		Description: "Public endpoint. Optionally filter by seller_id; paginated with limit/offset.",
		Tags:        []string{"Products"},
	}, productH.List)

	huma.Register(api, huma.Operation{
		OperationID: "get-product",
		Method:      nethttp.MethodGet,
		Path:        "/products/{id}",
		Summary:     "Get a product by ID",
		Description: "Public endpoint.",
		Tags:        []string{"Products"},
	}, productH.GetByID)

	huma.Register(api, huma.Operation{
		OperationID: "update-product",
		Method:      nethttp.MethodPatch,
		Path:        "/products/{id}",
		Summary:     "Update a product",
		Description: "Requires a seller JWT. Only the seller who owns the product may update it.",
		Tags:        []string{"Products"},
		Security:    bearer,
	}, productH.Update)

	huma.Register(api, huma.Operation{
		OperationID:   "delete-product",
		Method:        nethttp.MethodDelete,
		Path:          "/products/{id}",
		Summary:       "Delete a product",
		Description:   "Requires a seller JWT. Only the owner may delete it. Returns 409 if the product appears in orders.",
		Tags:          []string{"Products"},
		DefaultStatus: nethttp.StatusNoContent,
		Security:      bearer,
	}, productH.Delete)

	huma.Register(api, huma.Operation{
		OperationID:   "place-order",
		Method:        nethttp.MethodPost,
		Path:          "/orders",
		Summary:       "Place an order",
		Description:   "Requires a buyer JWT. The buyer id is taken from the token; do not send buyer_id. Product name, seller and unit price are snapshotted. Fails with 409 if stock is insufficient.",
		Tags:          []string{"Orders"},
		DefaultStatus: nethttp.StatusCreated,
		Security:      bearer,
	}, orderH.PlaceOrder)

	huma.Register(api, huma.Operation{
		OperationID: "get-order",
		Method:      nethttp.MethodGet,
		Path:        "/orders/{id}",
		Summary:     "Get order by ID",
		Description: "Requires a JWT for the buyer who placed the order or a seller who has an item in it.",
		Tags:        []string{"Orders"},
		Security:    bearer,
	}, orderH.GetOrder)

	huma.Register(api, huma.Operation{
		OperationID: "update-order-status",
		Method:      nethttp.MethodPatch,
		Path:        "/orders/{id}/status",
		Summary:     "Update order status",
		Description: "Requires a seller JWT for a seller participating in the order. Legal transitions: pending->confirmed/cancelled, confirmed->shipped/cancelled, shipped->delivered, delivered->returned.",
		Tags:        []string{"Orders"},
		Security:    bearer,
	}, orderH.UpdateStatus)

	huma.Register(api, huma.Operation{
		OperationID: "cancel-order",
		Method:      nethttp.MethodPost,
		Path:        "/orders/{id}/cancel",
		Summary:     "Cancel an order",
		Description: "Requires the buyer JWT that placed the order. Restores stock in the same transaction.",
		Tags:        []string{"Orders"},
		Security:    bearer,
	}, orderH.Cancel)

	huma.Register(api, huma.Operation{
		OperationID: "list-buyer-orders",
		Method:      nethttp.MethodGet,
		Path:        "/buyers/{id}/orders",
		Summary:     "List orders for a buyer",
		Description: "Requires a buyer JWT. The path id must match the buyer id in the token.",
		Tags:        []string{"Orders"},
		Security:    bearer,
	}, orderH.ListByBuyer)

	huma.Register(api, huma.Operation{
		OperationID: "list-seller-orders",
		Method:      nethttp.MethodGet,
		Path:        "/sellers/{id}/orders",
		Summary:     "List orders for a seller",
		Description: "Requires a seller JWT. The path id must match the seller id in the token. Each order only includes this seller's items and subtotal.",
		Tags:        []string{"Orders"},
		Security:    bearer,
	}, orderH.ListBySeller)

	return r
}
