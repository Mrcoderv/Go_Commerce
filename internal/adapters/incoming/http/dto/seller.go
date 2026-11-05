package dto

import "Go_Commerce/internal/domain/seller"

type RegisterSellerRequest struct {
	Body struct {
		Name    string `json:"name" minLength:"1"`
		Email   string `json:"email" format:"email"`
		Phone   string `json:"phone,omitempty"`
		Address string `json:"address,omitempty"`
	}
}

type GetSellerRequest struct {
	ID int64 `path:"id" minimum:"1"`
}

type ListSellersRequest struct {
	Limit  int `query:"limit" default:"20" minimum:"1" maximum:"100"`
	Offset int `query:"offset" default:"0" minimum:"0"`
}

type UpdateSellerRequest struct {
	ID   int64 `path:"id" minimum:"1"`
	Body struct {
		Name    *string `json:"name,omitempty" minLength:"1"`
		Email   *string `json:"email,omitempty" format:"email"`
		Phone   *string `json:"phone,omitempty"`
		Address *string `json:"address,omitempty"`
	}
}

type DeleteSellerRequest struct {
	ID int64 `path:"id" minimum:"1"`
}

type SellerResponse struct {
	Body seller.Seller
}

type ListSellersResponse struct {
	Body []seller.Seller
}
