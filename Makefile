.PHONY: proto order-proto sso-proto

order-proto:
	protoc --proto_path=proto/order-service/ \
	       --go_out=pkg/api/order --go_opt=paths=source_relative \
	       --go-grpc_out=pkg/api/order --go-grpc_opt=paths=source_relative \
	       proto/order-service/order.proto

sso-proto:
	protoc --proto_path=proto/sso-service/ \
	       --go_out=pkg/api/sso --go_opt=paths=source_relative \
	       --go-grpc_out=pkg/api/sso --go-grpc_opt=paths=source_relative \
	       proto/sso-service/sso.proto

proto: order-proto sso-proto

migrate-up:
	goose postgres "host=localhost user=postgres password=postgres dbname=postgres sslmode=disable" -dir migrations up

migrate-down:
	goose postgres "host=localhost user=postgres password=postgres dbname=postgres sslmode=disable" -dir migrations down

run:
	go run sso-service/cmd/app/main.go &
	go run order-service/cmd/app/main.go &
	wait