package handlers

import (
	"errors"

	"Go_Commerce/internal/domain/auth"
	"Go_Commerce/internal/domain/buyer"
	"Go_Commerce/internal/domain/order"
	"Go_Commerce/internal/domain/product"
	"Go_Commerce/internal/domain/seller"
	"Go_Commerce/internal/infrastructure/logger"

	"github.com/danielgtaylor/huma/v2"
)

// MapError maps domain sentinel errors to appropriate Huma HTTP errors and
// logs the failure through the configured logger.
func MapError(log logger.Logger, err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, seller.ErrNotFound),
		errors.Is(err, buyer.ErrNotFound),
		errors.Is(err, product.ErrNotFound),
		errors.Is(err, order.ErrNotFound),
		errors.Is(err, auth.ErrNotFound):
		log.Warn("request rejected", "status", 404, "error", err)
		return huma.Error404NotFound(err.Error())

	case errors.Is(err, seller.ErrInvalidInput),
		errors.Is(err, buyer.ErrInvalidInput),
		errors.Is(err, product.ErrInvalidInput),
		errors.Is(err, order.ErrInvalidInput),
		errors.Is(err, auth.ErrInvalidInput):
		log.Warn("request rejected", "status", 400, "error", err)
		return huma.Error400BadRequest(err.Error())

	case errors.Is(err, auth.ErrUnauthorized),
		errors.Is(err, auth.ErrInvalidCreds):
		log.Warn("request rejected", "status", 401, "error", err)
		return huma.Error401Unauthorized(err.Error())

	case errors.Is(err, product.ErrForbidden),
		errors.Is(err, auth.ErrForbidden):
		log.Warn("request rejected", "status", 403, "error", err)
		return huma.Error403Forbidden(err.Error())

	case errors.Is(err, seller.ErrEmailAlreadyExists),
		errors.Is(err, buyer.ErrEmailAlreadyExists),
		errors.Is(err, buyer.ErrConflict),
		errors.Is(err, product.ErrConflict),
		errors.Is(err, product.ErrInsufficientStock),
		errors.Is(err, order.ErrInvalidTransition),
		errors.Is(err, auth.ErrEmailExists):
		log.Warn("request rejected", "status", 409, "error", err)
		return huma.Error409Conflict(err.Error())

	default:
		log.Error("unhandled error", "error", err)
		return huma.Error500InternalServerError("internal server error")
	}
}
