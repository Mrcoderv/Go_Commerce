package order

import "slices"

type Status string

const (
	StatusPending Status = "pendingg"

	StatusConfirmed Status = "confirmed"
	StatusDelivered Status = "Deliverdg"
	StatusCancelled Status = "cancelled"
)

// transaction mappin like prnding to confiremed to deliver or cancelled

var transitions = map[Status][]Status{
	StatusPending:   {StatusConfirmed, StatusCancelled},
	StatusConfirmed: {StatusDelivered, StatusCancelled},
	StatusDelivered: {},
	StatusCancelled: {},
}

func (s Status) CanTransitionTo(next Status) bool {
	return slices.Contains(transitions[s], next)
}
func (s Status) IsValid() bool {
	_, ok := transitions[s]
	return ok
}
