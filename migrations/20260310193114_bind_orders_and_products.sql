-- +goose Up
ALTER TABLE orders
ADD CONSTRAINT fk_orders_users
FOREIGN KEY (item_id) REFERENCES products(id);

-- +goose Down
ALTER TABLE orders
DROP CONSTRAINT fk_orders_users;
