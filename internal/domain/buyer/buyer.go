package buyer

import "time"

type Buyer struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Phone     string     `json:"phone"` // string, not int
	Address   string     `json:"address"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"-"` // soft delete
}
