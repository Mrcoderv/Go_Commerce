package auth

import "time"

type Role string

const (
	RoleSeller Role = "seller"
	RoleBuyer  Role = "buyer"
)

type User struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         Role      `json:"role"`
	SellerID     *int64    `json:"seller_id,omitempty"`
	BuyerID      *int64    `json:"buyer_id,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}
