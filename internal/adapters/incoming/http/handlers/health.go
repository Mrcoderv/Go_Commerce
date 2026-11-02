package handlers

import (
	"context"
)

type HealthInput struct{}

type HealthOutput struct {
	Body HealthResponse
}

type HealthResponse struct {
	Status string `json:"status" doc:"Health status of the API"`
}

type Handler struct{}

func (h *Handler) Health(ctx context.Context, input *HealthInput) (*HealthOutput, error) {
	return &HealthOutput{
		Body: HealthResponse{
			Status: "health is good",
		},
	}, nil
}
