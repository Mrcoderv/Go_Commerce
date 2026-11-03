package product

import (
	"context"
	"Go_Commerce/internal/domain/product"
	"Go_Commerce/internal/ports/incoming"
	"Go_Commerce/internal/ports/outgoing"
)

type ProductService struct {
	repo outgoing.ProductRepository
}

func NewProductService(repo outgoing.ProductRepository) * ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) ListAll(ctx context.Context) ([]product.Product, error) {
	return s.repo.GetAll(ctx)
}