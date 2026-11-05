package order

import "slices"

type Status string

const (
	StatusPending   Status = "pending"
	StatusConfirmed Status = "confirmed"
	StatusShipped   Status = "shipped"
	StatusDelivered Status = "delivered"
	StatusCancelled Status = "cancelled"
)

var transitions = map[Status][]Status{
	StatusPending:   {StatusConfirmed, StatusCancelled},
	StatusConfirmed: {StatusShipped, StatusCancelled},
	StatusShipped:   {StatusDelivered},
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
