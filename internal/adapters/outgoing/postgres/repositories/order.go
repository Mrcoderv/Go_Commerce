package repositories

import (
	"context"
	"database/sql"
	"time"

	"Go_Commerce/internal/domain/order"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(ctx context.Context, o order.Order) (order.Order, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return o, err
	}
	defer tx.Rollback()

	now := time.Now()
	if o.CreatedAt.IsZero() {
		o.CreatedAt = now
	}

	if o.Status == "" {
		o.Status = order.Status("pending")
	}

	err = tx.QueryRowContext(ctx,
		`INSERT INTO orders (buyer_id, status, total_amount, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		o.BuyerID, o.Status, o.Total, o.CreatedAt, now,
	).Scan(&o.ID)
	if err != nil {
		return o, err
	}

	for i := range o.Items {
		o.Items[i].OrderID = o.ID

		err = tx.QueryRowContext(ctx,
			`INSERT INTO order_items (order_id, product_id, seller_id, product_name, unit_price, quantity)
			 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
			o.ID, o.Items[i].ProductID, o.Items[i].SellerID,
			o.Items[i].ProductName, o.Items[i].UnitPrice, o.Items[i].Quantity,
		).Scan(&o.Items[i].ID)
		if err != nil {
			return o, err
		}
	}

	if err = tx.Commit(); err != nil {
		return o, err
	}

	return o, nil
}

func (r *OrderRepository) GetByID(ctx context.Context, id int64) (order.Order, error) {
	var o order.Order
	var status string

	err := r.db.QueryRowContext(ctx,
		`SELECT id, buyer_id, status, total_amount, created_at FROM orders WHERE id = $1`,
		id,
	).Scan(&o.ID, &o.BuyerID, &status, &o.Total, &o.CreatedAt)
	if err != nil {
		return o, err
	}

	o.Status = order.Status(status)
	items, err := r.loadItems(ctx, o.ID)
	if err != nil {
		return o, err
	}
	o.Items = items

	return o, nil
}

func (r *OrderRepository) loadItems(ctx context.Context, orderID int64) ([]order.Item, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, order_id, product_id, seller_id, product_name, unit_price, quantity
		 FROM order_items WHERE order_id = $1`,
		orderID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []order.Item

	for rows.Next() {
		var item order.Item

		if err := rows.Scan(
			&item.ID, &item.OrderID, &item.ProductID, &item.SellerID,
			&item.ProductName, &item.UnitPrice, &item.Quantity,
		); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *OrderRepository) ListByBuyer(ctx context.Context, buyerID int64, limit, offset int) ([]order.Order, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, buyer_id, status, total_amount, created_at
		 FROM orders WHERE buyer_id = $1
		 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		buyerID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanOrders(ctx, rows)
}

func (r *OrderRepository) ListBySeller(ctx context.Context, sellerID int64, limit, offset int) ([]order.Order, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT DISTINCT o.id, o.buyer_id, o.status, o.total_amount, o.created_at
		 FROM orders o
		 JOIN order_items oi ON oi.order_id = o.id
		 WHERE oi.seller_id = $1
		 ORDER BY o.created_at DESC LIMIT $2 OFFSET $3`,
		sellerID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanOrders(ctx, rows)
}

func (r *OrderRepository) scanOrders(ctx context.Context, rows *sql.Rows) ([]order.Order, error) {
	var orders []order.Order

	for rows.Next() {
		var o order.Order
		var status string

		if err := rows.Scan(&o.ID, &o.BuyerID, &status, &o.Total, &o.CreatedAt); err != nil {
			return nil, err
		}

		o.Status = order.Status(status)

		var err error
		o.Items, err = r.loadItems(ctx, o.ID)
		if err != nil {
			return nil, err
		}

		orders = append(orders, o)
	}

	return orders, rows.Err()
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, id int64, status order.Status) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2`,
		status, id,
	)
	return err
}