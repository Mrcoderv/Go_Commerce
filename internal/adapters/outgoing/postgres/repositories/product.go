package repositories

import (
	"context"
	"database/sql"

	"Go_Commerce/internal/domain/product"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(ctx context.Context, p product.Product) error {
	query := `INSERT INTO products (seller_id, product_name, product_description, unit_price, stock) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return r.db.QueryRowContext(ctx, query, p.SellerID, p.Name, p.Description, p.UnitPrice, p.Stock).Scan(&p.ID)
}

func (r *ProductRepository) GetByID(ctx context.Context, id int64) (product.Product, error) {
	p := product.Product{}
	query := `SELECT id, seller_id, product_name, product_description, unit_price, stock FROM products WHERE id = $1 AND deleted_at IS NULL`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&p.ID, &p.SellerID, &p.Name, &p.Description, &p.UnitPrice, &p.Stock)
	if err != nil {
		return product.Product{}, err
	}
	return p, nil
}

func (r *ProductRepository) GetAll(ctx context.Context) ([]product.Product, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, seller_id, product_name, product_description, unit_price, stock FROM products WHERE deleted_at IS NULL`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []product.Product
	for rows.Next() {
		p := product.Product{}
		if err := rows.Scan(&p.ID, &p.SellerID, &p.Name, &p.Description, &p.UnitPrice, &p.Stock); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

func (r *ProductRepository) Update(ctx context.Context, p product.Product) error {
	query := `UPDATE products SET seller_id = $1, product_name = $2, product_description = $3, unit_price = $4, stock = $5, updated_at = NOW() WHERE id = $6`
	_, err := r.db.ExecContext(ctx, query, p.SellerID, p.Name, p.Description, p.UnitPrice, p.Stock, p.ID)
	return err
}

func (r *ProductRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `UPDATE products SET deleted_at = NOW() WHERE id = $1`, id)
	return err
}
