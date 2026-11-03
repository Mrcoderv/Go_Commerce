package seller

import (
	"context"
	"strings"
	"Go_Commerce/internal/domain/seller"
	"Go_Commerce/internal/ports/incoming"
	"Go_Commerce/internal/ports/outgoing"
)
type SellerService struct {
	repo outgoing.SellerRepository
}
func NewSellerService(repo outgoing.SellerRepository) *SellerService {
	return &SellerService{repo: repo}
}