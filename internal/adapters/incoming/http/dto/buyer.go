package dto

import "Go_Commerce/internal/domain/buyer"

type RegisterBuyerRequest struct { // dto
	Body struct {
		Name    string `json:"name" minLength:"1"`
		Email   string `json:"email" format:"email"`
		Phone   string `json:"phone,omitempty"`
		Address string `json:"address,omitempty"`
	}
}

type GetBuyerRequest struct {
	ID int64 `path:"id" minimum:"1"`
}

type ListBuyersRequest struct {
	Limit  int `query:"limit" default:"20" minimum:"1" maximum:"100"`
	Offset int `query:"offset" default:"0" minimum:"0"`
}

type UpdateBuyerRequest struct {
	ID   int64 `path:"id" minimum:"1"`
	Body struct {
		Name    *string `json:"name,omitempty" doc:"New name"`
		Phone   *string `json:"phone,omitempty" doc:"New phone"`
		Address *string `json:"address,omitempty" doc:"New address"`
	}
}

type DeleteBuyerRequest struct {
	ID int64 `path:"id" minimum:"1"`
}

type BuyerResponse struct {
	Body buyer.Buyer
}

type ListBuyersResponse struct {
	Body []buyer.Buyer
}
