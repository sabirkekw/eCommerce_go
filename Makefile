.PHONY: proto

sso-proto:
	protoc --proto_path=proto/sso-service/ \
	       --go_out=proto/sso-service --go_opt=paths=source_relative \
	       --go-grpc_out=proto/sso-service --go-grpc_opt=paths=source_relative \
	       proto/sso-service/sso.proto