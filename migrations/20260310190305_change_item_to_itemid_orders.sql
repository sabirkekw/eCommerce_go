-- +goose Up
ALTER TABLE orders 
RENAME COLUMN item TO item_id;

-- +goose Down
ALTER TABLE orders
RENAME COLUMN item_id TO item;
