-- +goose Up
ALTER TABLE orders ADD COLUMN user_id INT NOT NULL REFERENCES users(id);

-- +goose Down
ALTER TABLE orders DROP COLUMN user_id;
