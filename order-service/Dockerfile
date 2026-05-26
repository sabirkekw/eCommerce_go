FROM golang:1.25 AS builder
WORKDIR /app/

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main ./order-service/cmd/app/

# Use a glibc-based runtime so the confluent-kafka-go binary runs correctly
FROM debian:bookworm-slim AS final
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/config/order-service/ /root/config/order-service/
CMD ["./main"]
