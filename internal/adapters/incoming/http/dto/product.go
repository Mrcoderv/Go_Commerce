package dto

import "Go_Commerce/internal/domain/product"

type CreateProductRequest struct {
	Body struct {
		Name        string `json:"name" minLength:"1" doc:"Product name"`
		Description string `json:"description,omitempty" doc:"Optional description"`
		Price       int64  `json:"price" minimum:"1" doc:"Price in paisa (fixed money, e.g. 49950 = NPR 499.50)"`
		Stock       int    `json:"stock" minimum:"0" doc:"Available stock"`
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
	ID   int64 `path:"id" minimum:"1" doc:"Product ID"`
	Body struct {
		Name        *string `json:"name,omitempty" doc:"New product name"`
		Description *string `json:"description,omitempty" doc:"New description"`
		Price       *int64  `json:"price,omitempty" minimum:"1" doc:"New price in paisa"`
		Stock       *int    `json:"stock,omitempty" doc:"New stock"`
	}
}

type DeleteProductRequest struct {
	ID int64 `path:"id" minimum:"1" doc:"Product ID"`
}

type ProductResponse struct {
	Body product.Product
}

type ListProductsResponse struct {
	Body []product.Product
}
