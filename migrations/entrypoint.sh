#!/bin/bash

goose postgres "$SSO_DB_STRING" -dir /migrations/sso-service/ up
goose postgres "$ORDERS_DB_STRING" -dir /migrations/order-service/ up
goose postgres "$PRODUCTS_DB_STRING" -dir /migrations/products-service/ up
goose postgres "$CART_DB_STRING" -dir /migrations/cart-service/ up