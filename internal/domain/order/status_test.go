package order_test

import (
	"Go_Commerce/internal/domain/order"
	"testing"
)

func TestStatusTransitions(t *testing.T) {
	tests := []struct {
		from     order.Status
		to       order.Status
		expected bool
	}{
		// Legal transitions
		{order.StatusPending, order.StatusConfirmed, true},
		{order.StatusPending, order.StatusCancelled, true},
		{order.StatusConfirmed, order.StatusShipped, true},
		{order.StatusConfirmed, order.StatusCancelled, true},
		{order.StatusShipped, order.StatusDelivered, true},

		// Illegal transitions
		{order.StatusPending, order.StatusShipped, false},
		{order.StatusPending, order.StatusDelivered, false},
		{order.StatusConfirmed, order.StatusPending, false},
		{order.StatusConfirmed, order.StatusDelivered, false},
		{order.StatusShipped, order.StatusCancelled, false},
		{order.StatusShipped, order.StatusPending, false},
		{order.StatusDelivered, order.StatusCancelled, false},
		{order.StatusDelivered, order.StatusPending, false},
		{order.StatusCancelled, order.StatusConfirmed, false},
		{order.StatusCancelled, order.StatusPending, false},
	}

	for _, tt := range tests {
		got := tt.from.CanTransitionTo(tt.to)
		if got != tt.expected {
			t.Errorf("CanTransitionTo(%s -> %s) = %v; want %v", tt.from, tt.to, got, tt.expected)
		}
	}
}

func TestStatusValidity(t *testing.T) {
	valid := []order.Status{
		order.StatusPending,
		order.StatusConfirmed,
		order.StatusShipped,
		order.StatusDelivered,
		order.StatusCancelled,
	}

	for _, s := range valid {
		if !s.IsValid() {
			t.Errorf("expected %s to be valid", s)
		}
	}

	invalid := []order.Status{
		"",
		"unknown",
		"pendingg",
	}

	for _, s := range invalid {
		if s.IsValid() {
			t.Errorf("expected %s to be invalid", s)
		}
	}
}
