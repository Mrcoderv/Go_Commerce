package order

type Status string

const (
	StatusPending Status = "pendingg"

	StatusConfirmed Status = "confirmed"
	StatusDelivered Status = "Deliverdg"
	StatusCancelled Status = "cancelled"
)
