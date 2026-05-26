-- +goose Up
CREATE TABLE IF NOT EXISTS products(
    id SERIAL PRIMARY KEY,
    product_name VARCHAR(255) NOT NULL UNIQUE,
    quantity INTEGER,
    description VARCHAR(255)
);

-- +goose Down
DROP TABLE IF EXISTS products;
