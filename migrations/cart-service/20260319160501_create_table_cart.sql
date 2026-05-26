-- +goose Up
CREATE TABLE IF NOT EXISTS cart (
    user_id     INT     NOT NULL,
    product_id  INT     NOT NULL,
    quantity    INT     NOT NULL CHECK (quantity > 0),
    description VARCHAR(255),

    PRIMARY KEY (user_id, product_id)
);

-- +goose Down
DROP TABLE IF EXISTS cart;
