FROM golang:1.25-alpine AS builder
WORKDIR /app/
COPY go.mod go.sum ./

RUN go mod download
COPY . .
RUN go build -o main ./products-service/cmd/app/

FROM alpine:latest AS final
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/config/products-service/ /root/config/products-service/
CMD ["./main"]
