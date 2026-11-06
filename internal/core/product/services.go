package product

import (
	"Go_Commerce/internal/domain/product"
	"Go_Commerce/internal/infrastructure/logger"
	"Go_Commerce/internal/ports/incoming"
	"Go_Commerce/internal/ports/outgoing"
	"context"
	"strings"
)

type ProductService struct {
	repo outgoing.ProductRepository
	log  logger.Logger
}

func NewProductService(repo outgoing.ProductRepository, log logger.Logger) incoming.ProductService {
	return &ProductService{repo: repo, log: log}
}

var _ incoming.ProductService = (*ProductService)(nil)

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
	if in.SellerID <= 0 || name == "" || in.Price <= 0 || in.Stock < 0 {
		return product.Product{}, product.ErrInvalidInput
	}

	created, err := s.repo.Create(ctx, product.Product{
		SellerID:    in.SellerID,
		Name:        name,
		Description: strings.TrimSpace(in.Description),
		Price:       in.Price,
		Stock:       in.Stock,
	})
	if err != nil {
		s.log.Error("product create failed", "seller_id", in.SellerID, "error", err)
		return product.Product{}, err
	}
	s.log.Info("product created", "product_id", created.ID, "seller_id", created.SellerID, "price_paisa", created.Price)
	return created, nil
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
	if in.Price != nil {
		if *in.Price <= 0 {
			return product.Product{}, product.ErrInvalidInput
		}
		current.Price = *in.Price
	}
	if in.Stock != nil {
		if *in.Stock < 0 {
			return product.Product{}, product.ErrInvalidInput
		}
		current.Stock = *in.Stock
	}

	updated, err := s.repo.Update(ctx, current)
	if err != nil {
		s.log.Error("product update failed", "product_id", id, "seller_id", sellerID, "error", err)
		return product.Product{}, err
	}
	s.log.Info("product updated", "product_id", id, "seller_id", sellerID)
	return updated, nil
}

func (s *ProductService) Delete(ctx context.Context, id int64, sellerID int64) error {
	current, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if current.SellerID != sellerID {
		return product.ErrForbidden
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	s.log.Info("product deleted", "product_id", id, "seller_id", sellerID)
	return nil
}
