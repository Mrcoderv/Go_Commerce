package handlers

import (
	"context"

	"Go_Commerce/internal/adapters/incoming/http/middleware"
	coreauth "Go_Commerce/internal/core/auth"
	domainauth "Go_Commerce/internal/domain/auth"
	"Go_Commerce/internal/domain/order"
)

func claimsForRole(ctx context.Context, role domainauth.Role) (*coreauth.Claims, error) {
	claims, ok := middleware.ClaimsFromContext(ctx)
	if !ok {
		return nil, domainauth.ErrUnauthorized
	}
	if claims.Role != role {
		return nil, domainauth.ErrForbidden
	}
	return claims, nil
}

func sellerClaims(ctx context.Context) (*coreauth.Claims, error) {
	claims, err := claimsForRole(ctx, domainauth.RoleSeller)
	if err != nil || claims.SellerID == nil {
		if err != nil {
			return nil, err
		}
		return nil, domainauth.ErrUnauthorized
	}
	return claims, nil
}

func buyerClaims(ctx context.Context) (*coreauth.Claims, error) {
	claims, err := claimsForRole(ctx, domainauth.RoleBuyer)
	if err != nil || claims.BuyerID == nil {
		if err != nil {
			return nil, err
		}
		return nil, domainauth.ErrUnauthorized
	}
	return claims, nil
}

func orderAccessClaims(ctx context.Context) (*coreauth.Claims, error) {
	claims, ok := middleware.ClaimsFromContext(ctx)
	if !ok {
		return nil, domainauth.ErrUnauthorized
	}
	if claims.Role != domainauth.RoleBuyer && claims.Role != domainauth.RoleSeller {
		return nil, domainauth.ErrForbidden
	}
	if claims.Role == domainauth.RoleBuyer && claims.BuyerID == nil {
		return nil, domainauth.ErrUnauthorized
	}
	if claims.Role == domainauth.RoleSeller && claims.SellerID == nil {
		return nil, domainauth.ErrUnauthorized
	}
	return claims, nil
}

func canAccessOrder(claims *coreauth.Claims, o order.Order) bool {
	if claims.Role == domainauth.RoleBuyer {
		return claims.BuyerID != nil && *claims.BuyerID == o.BuyerID
	}
	if claims.Role == domainauth.RoleSeller && claims.SellerID != nil {
		for _, item := range o.Items {
			if item.SellerID == *claims.SellerID {
				return true
			}
		}
	}
	return false
}
