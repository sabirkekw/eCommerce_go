-- +goose Up
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    quantity INT NOT NULL,
    description VARCHAR(255)
);


-- +goose Down
DROP TABLE IF EXISTS products;
