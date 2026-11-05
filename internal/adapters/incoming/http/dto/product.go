package dto

import "Go_Commerce/internal/domain/product"

type CreateProductRequest struct {
	SellerIDHeader int64 `header:"X-Seller-ID"`
	Body           struct {
		SellerID    int64   `json:"seller_id,omitempty"`
		Name        string  `json:"name" minLength:"1"`
		Description string  `json:"description,omitempty"`
		Price       float64 `json:"price,omitempty" doc:"Price (plain number or in NPR)"`
		UnitPrice   int64   `json:"unit_price,omitempty" doc:"Price in paisa"`
		Stock       int     `json:"stock" minimum:"0"`
	}
}

type GetProductRequest struct {
	ID int64 `path:"id" minimum:"1"`
}

type ListProductsRequest struct {
	SellerID int64 `query:"seller_id,omitempty"`
	Limit    int   `query:"limit" default:"20" minimum:"1" maximum:"100"`
	Offset   int   `query:"offset" default:"0" minimum:"0"`
}

type UpdateProductRequest struct {
	ID             int64 `path:"id" minimum:"1"`
	SellerIDHeader int64 `header:"X-Seller-ID"`
	Body           struct {
		SellerID    int64    `json:"seller_id,omitempty"`
		Name        *string  `json:"name,omitempty"`
		Description *string  `json:"description,omitempty"`
		Price       *float64 `json:"price,omitempty"`
		UnitPrice   *int64   `json:"unit_price,omitempty"`
		Stock       *int     `json:"stock,omitempty"`
	}
}

type DeleteProductRequest struct {
	ID             int64 `path:"id" minimum:"1"`
	SellerIDHeader int64 `header:"X-Seller-ID"`
	SellerIDQuery  int64 `query:"seller_id"`
}

type ProductResponse struct {
	Body product.Product
}

type ListProductsResponse struct {
	Body []product.Product
}
