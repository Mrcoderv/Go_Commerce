package product

import (
	"Go_Commerce/internal/domain/product"
	"Go_Commerce/internal/ports/incoming"
	"Go_Commerce/internal/ports/outgoing"
	"context"
	"strings"
)

type ProductService struct {
	repo outgoing.ProductRepository
}

func NewProductService(repo outgoing.ProductRepository) incoming.ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetByID(ctx context.Context, id int64) (product.Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ProductService) List(ctx context.Context, f incoming.ProductFilter) ([]product.Product, error) {
	limit, offset := f.Limit, f.Offset
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.List(ctx, f.SellerID, limit, offset)
}
func (s *ProductService) Create(ctx context.Context, in incoming.CreateProductInput) (product.Product, error) {
	name := strings.TrimSpace(in.Name)
	if in.SellerID <= 0 || name == "" || in.UnitPrice <= 0 || in.Stock < 0 {
		return product.Product{}, product.ErrInvalidInput
	}

	return s.repo.Create(ctx, product.Product{
		SellerID:    in.SellerID,
		Name:        name,
		Description: strings.TrimSpace(in.Description),
		UnitPrice:   in.UnitPrice,
		Stock:       in.Stock,
	})
}
func (s *ProductService) CheckAvailability(ctx context.Context, productID int64, qty int) (product.Product, error) {
	if qty <= 0 {
		return product.Product{}, product.ErrInvalidInput
	}
	p, err := s.repo.GetByID(ctx, productID)
	if err != nil {
		return product.Product{}, err
	}
	if p.Stock < qty {
		return product.Product{}, product.ErrInsufficientStock
	}
	return p, nil
}

func (s *ProductService) ReduceStock(ctx context.Context, productID int64, qty int) error {
	if qty <= 0 {
		return product.ErrInvalidInput
	}
	return s.repo.DecreaseStock(ctx, productID, qty)
}
func (s *ProductService) Update(ctx context.Context, id int64, sellerID int64, in incoming.UpdateProductInput) (product.Product, error) {
	current, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return product.Product{}, err
	}
	if current.SellerID != sellerID {
		return product.Product{}, product.ErrForbidden
	}

	if in.Name != nil {
		name := strings.TrimSpace(*in.Name)
		if name == "" {
			return product.Product{}, product.ErrInvalidInput
		}
		current.Name = name
	}
	if in.Description != nil {
		current.Description = strings.TrimSpace(*in.Description)
	}
	if in.UnitPrice != nil {
		if *in.UnitPrice <= 0 {
			return product.Product{}, product.ErrInvalidInput
		}
		current.UnitPrice = *in.UnitPrice
	}
	if in.Stock != nil {
		if *in.Stock < 0 {
			return product.Product{}, product.ErrInvalidInput
		}
		current.Stock = *in.Stock
	}

	return s.repo.Update(ctx, current)
}

func (s *ProductService) Delete(ctx context.Context, id int64, sellerID int64) error {
	current, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if current.SellerID != sellerID {
		return product.ErrForbidden
	}
	return s.repo.Delete(ctx, id)
}
