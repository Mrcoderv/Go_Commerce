package repositories

import (
	"Go_Commerce/internal/adapters/outgoing/postgres"
	"Go_Commerce/internal/domain/order"
	"Go_Commerce/internal/ports/outgoing"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

var _ outgoing.OrderRepository = (*OrderRepository)(nil)

func (r *OrderRepository) Create(ctx context.Context, o order.Order) (order.Order, error) {
	now := time.Now()
	if o.CreatedAt.IsZero() {
		o.CreatedAt = now
	}

	if o.Status == "" {
		o.Status = order.StatusPending
	}

	q := postgres.Querier(ctx, r.db)
	err := q.QueryRowContext(ctx,
		`INSERT INTO orders (buyer_id, status, total_amount, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		o.BuyerID, o.Status, o.Total, o.CreatedAt, now,
	).Scan(&o.ID)
	if err != nil {
		return o, fmt.Errorf("create order: %w", err)
	}

	for i := range o.Items {
		o.Items[i].OrderID = o.ID

		err = q.QueryRowContext(ctx,
			`INSERT INTO order_items (order_id, product_id, seller_id, product_name, unit_price, quantity)
			 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
			o.ID, o.Items[i].ProductID, o.Items[i].SellerID,
			o.Items[i].ProductName, o.Items[i].UnitPrice, o.Items[i].Quantity,
		).Scan(&o.Items[i].ID)
		if err != nil {
			return o, fmt.Errorf("create order item: %w", err)
		}
	}

	return o, nil
}

func (r *OrderRepository) GetByID(ctx context.Context, id int64) (order.Order, error) {
	var o order.Order
	var status string

	err := postgres.Querier(ctx, r.db).QueryRowContext(ctx,
		`SELECT id, buyer_id, status, total_amount, created_at FROM orders WHERE id = $1`,
		id,
	).Scan(&o.ID, &o.BuyerID, &status, &o.Total, &o.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return order.Order{}, order.ErrNotFound
	}
	if err != nil {
		return order.Order{}, fmt.Errorf("get order by id %d: %w", id, err)
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
	rows, err := postgres.Querier(ctx, r.db).QueryContext(ctx,
		`SELECT id, order_id, product_id, seller_id, product_name, unit_price, quantity
		 FROM order_items WHERE order_id = $1 ORDER BY id ASC`,
		orderID,
	)
	if err != nil {
		return nil, fmt.Errorf("load items: %w", err)
	}
	defer rows.Close()

	var items []order.Item

	for rows.Next() {
		var item order.Item

		if err := rows.Scan(
			&item.ID, &item.OrderID, &item.ProductID, &item.SellerID,
			&item.ProductName, &item.UnitPrice, &item.Quantity,
		); err != nil {
			return nil, fmt.Errorf("scan item: %w", err)
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("load items: %w", err)
	}
	if items == nil {
		items = []order.Item{}
	}
	return items, nil
}

func (r *OrderRepository) ListByBuyer(ctx context.Context, buyerID int64, limit, offset int) ([]order.Order, error) {
	rows, err := postgres.Querier(ctx, r.db).QueryContext(ctx,
		`SELECT id, buyer_id, status, total_amount, created_at
		 FROM orders WHERE buyer_id = $1
		 ORDER BY created_at DESC, id DESC LIMIT $2 OFFSET $3`,
		buyerID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list orders by buyer: %w", err)
	}
	defer rows.Close()

	return r.scanOrders(ctx, rows)
}

func (r *OrderRepository) ListBySeller(ctx context.Context, sellerID int64, limit, offset int) ([]order.Order, error) {
	rows, err := postgres.Querier(ctx, r.db).QueryContext(ctx,
		`SELECT DISTINCT o.id, o.buyer_id, o.status, o.total_amount, o.created_at
		 FROM orders o
		 JOIN order_items oi ON oi.order_id = o.id
		 WHERE oi.seller_id = $1
		 ORDER BY o.created_at DESC, o.id DESC LIMIT $2 OFFSET $3`,
		sellerID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list orders by seller: %w", err)
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
			return nil, fmt.Errorf("scan order: %w", err)
		}

		o.Status = order.Status(status)

		var err error
		o.Items, err = r.loadItems(ctx, o.ID)
		if err != nil {
			return nil, err
		}

		orders = append(orders, o)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("scan orders: %w", err)
	}
	if orders == nil {
		orders = []order.Order{}
	}
	return orders, nil
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, id int64, status order.Status) error {
	response, err := postgres.Querier(ctx, r.db).ExecContext(ctx,
		`UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2`,
		status, id,
	)
	if err != nil {
		return fmt.Errorf("error updating order status %d: %w", id, err)
	}
	n, err := response.RowsAffected()

	if err != nil {
		return fmt.Errorf("error updating order status %d: %w", id, err)
	}
	if n == 0 {
		return order.ErrNotFound
	}
	return nil
}
