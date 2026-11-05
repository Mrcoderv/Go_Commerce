# Design decisions

1. **Money in paisa (int64)**  
   Store prices and totals as `BIGINT` in paisa. No floats. Avoids rounding errors.

2. **Soft/hard deletes**  
   Sellers/products use soft delete (`deleted_at`) and reads filter `deleted_at IS NULL`. Buyers use hard delete; if they have orders, return conflict.

3. **SellerID on OrderItem**  
   Keep `seller_id` on `order_items` so an order can mix sellers and we can list orders by seller.

4. **Price snapshot**  
   Copy `product_name` and `unit_price` into `order_items` at order time so price changes don't affect history.

5. **Order status transitions**  
   Enforce legal transitions in `domain/order/status.go`: pending->confirmed/cancelled, confirmed->shipped/cancelled, shipped->delivered. Others invalid.
