package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"Go_Commerce/internal/adapters/outgoing/postgres"
	"Go_Commerce/internal/domain/product"
	"Go_Commerce/internal/ports/outgoing"

	"github.com/jackc/pgx/v5/pgconn"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// Compile-time check: the build fails if a method is missing.
var _ outgoing.ProductRepository = (*ProductRepository)(nil)

const productColumns = `id, seller_id, product_name, product_description, unit_price, stock`

func (r *ProductRepository) Create(ctx context.Context, p product.Product) (product.Product, error) {
	const q = `INSERT INTO products (seller_id, product_name, product_description, unit_price, stock)
	           VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := postgres.Querier(ctx, r.db).
		QueryRowContext(ctx, q, p.SellerID, p.Name, p.Description, p.UnitPrice, p.Stock).
		Scan(&p.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" { // seller does not exist
			return product.Product{}, product.ErrInvalidInput
		}
		return product.Product{}, fmt.Errorf("create product: %w", err)
	}
	return p, nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id int64) (product.Product, error) {
	const q = `SELECT ` + productColumns + ` FROM products WHERE id = $1 AND deleted_at IS NULL`

	var p product.Product
	err := postgres.Querier(ctx, r.db).
		QueryRowContext(ctx, q, id).
		Scan(&p.ID, &p.SellerID, &p.Name, &p.Description, &p.UnitPrice, &p.Stock)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return product.Product{}, product.ErrNotFound
		}
		return product.Product{}, fmt.Errorf("get product: %w", err)
	}
	return p, nil
}

func (r *ProductRepository) List(ctx context.Context, sellerID *int64, limit, offset int) ([]product.Product, error) {
	const q = `SELECT ` + productColumns + ` FROM products
	           WHERE deleted_at IS NULL AND ($1::bigint IS NULL OR seller_id = $1::bigint)
	           ORDER BY id DESC LIMIT $2 OFFSET $3`

	rows, err := postgres.Querier(ctx, r.db).QueryContext(ctx, q, sellerID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list products: %w", err)
	}
	defer rows.Close()

	products := []product.Product{}
	for rows.Next() {
		var p product.Product
		if err := rows.Scan(&p.ID, &p.SellerID, &p.Name, &p.Description, &p.UnitPrice, &p.Stock); err != nil {
			return nil, fmt.Errorf("scan product: %w", err)
		}
		products = append(products, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list products: %w", err)
	}
	return products, nil
}

func (r *ProductRepository) Update(ctx context.Context, p product.Product) (product.Product, error) {
	const q = `UPDATE products
	           SET product_name = $1, product_description = $2, unit_price = $3, stock = $4, updated_at = NOW()
	           WHERE id = $5 AND deleted_at IS NULL
	           RETURNING ` + productColumns

	var out product.Product
	err := postgres.Querier(ctx, r.db).
		QueryRowContext(ctx, q, p.Name, p.Description, p.UnitPrice, p.Stock, p.ID).
		Scan(&out.ID, &out.SellerID, &out.Name, &out.Description, &out.UnitPrice, &out.Stock)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return product.Product{}, product.ErrNotFound
		}
		return product.Product{}, fmt.Errorf("update product: %w", err)
	}
	return out, nil
}

func (r *ProductRepository) Delete(ctx context.Context, id int64) error {
	const q = `UPDATE products SET deleted_at = NOW(), updated_at = NOW()
	           WHERE id = $1 AND deleted_at IS NULL`

	res, err := postgres.Querier(ctx, r.db).ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("delete product: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete product: %w", err)
	}
	if n == 0 {
		return product.ErrNotFound
	}
	return nil
}

// DecreaseStock is ONE statement: the check and the decrease cannot be split
// by another buyer. Zero rows means not enough stock (or no such product).
func (r *ProductRepository) DecreaseStock(ctx context.Context, id int64, qty int) error {
	const q = `UPDATE products SET stock = stock - $1, updated_at = NOW()
	           WHERE id = $2 AND stock >= $1 AND deleted_at IS NULL`

	res, err := postgres.Querier(ctx, r.db).ExecContext(ctx, q, qty, id)
	if err != nil {
		return fmt.Errorf("decrease stock: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("decrease stock: %w", err)
	}
	if n == 0 {
		return product.ErrInsufficientStock
	}
	return nil
}

func (r *ProductRepository) IncreaseStock(ctx context.Context, id int64, qty int) error {
	const q = `UPDATE products SET stock = stock + $1, updated_at = NOW()
	           WHERE id = $2 AND deleted_at IS NULL`

	res, err := postgres.Querier(ctx, r.db).ExecContext(ctx, q, qty, id)
	if err != nil {
		return fmt.Errorf("increase stock: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("increase stock: %w", err)
	}
	if n == 0 {
		return product.ErrNotFound
	}
	return nil
}