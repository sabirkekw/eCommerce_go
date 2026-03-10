.PHONY: proto order-proto sso-proto

order-proto:
	protoc --proto_path=proto/order-service/ \
	       --go_out=pkg/api/order --go_opt=paths=source_relative \
	       --go-grpc_out=pkg/api/order --go-grpc_opt=paths=source_relative \
	       proto/order-service/order.proto

sso-proto:
	protoc -I . \
		--go_out . --go_opt paths=source_relative \
		--go-grpc_out . --go-grpc_opt paths=source_relative \
		pkg/api/sso/sso.proto

	protoc -I . --grpc-gateway_out . \
		--grpc-gateway_opt paths=source_relative \
		pkg/api/sso/sso.proto

products-proto:
	protoc -I . \
		--go_out . --go_opt paths=source_relative \
		--go-grpc_out . --go-grpc_opt paths=source_relative \
		pkg/api/products/products.proto

	protoc -I . --grpc-gateway_out . \
		--grpc-gateway_opt paths=source_relative \
		pkg/api/products/products.proto

migrate-up:
	goose postgres "host=localhost user=postgres password=postgres dbname=postgres sslmode=disable" -dir migrations up

migrate-down:
	goose postgres "host=localhost user=postgres password=postgres dbname=postgres sslmode=disable" -dir migrations down

run:
	go run sso-service/cmd/app/main.go &
	go run order-service/cmd/app/main.go &
	wait