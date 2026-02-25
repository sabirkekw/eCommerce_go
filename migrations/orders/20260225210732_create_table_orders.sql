-- +goose Up
CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    item VARCHAR(255) NOT NULL,
    quantity INT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS orders
