package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sabirkekw/ecommerce_go/order-service/internal/models/order"
	"github.com/sabirkekw/ecommerce_go/order-service/internal/models/products"
	"go.uber.org/zap"
)

type OrderService interface {
	CreateOrder(ctx context.Context, order *order.Order) (int32, error)
	GetOrderByID(ctx context.Context, orderID int32) (*order.Order, error)
	GetOrderByUserID(ctx context.Context, userID int32) ([]*order.Order, error)
	DeleteOrder(ctx context.Context, orderID int32) error
}

type KafkaConsumer struct {
	logger   *zap.SugaredLogger
	Consumer *kafka.Consumer
	Service  OrderService
	group    string
	topic    string
}

// TODO: Implement Kafka consumer to listen for checkout messages and process them in the order service
func NewKafkaConsumer(logger *zap.SugaredLogger, brokers string, group string, topic string, service OrderService) *KafkaConsumer {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"group.id":          group,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		panic(err)
	}

	err = consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		panic(err)
	}

	return &KafkaConsumer{
		logger:   logger,
		Consumer: consumer,
		group:    group,
		topic:    topic,
		Service:  service,
	}
}

func (c *KafkaConsumer) Close() {
	c.Consumer.Close()
}

// TODO: Implement message processing logic to create orders based on the received checkout messages
func (c *KafkaConsumer) Poll() {
	const op = "Order.Messaging.Poll"
	for {
		msg, err := c.Consumer.ReadMessage(-1) // blocks until a message is received
		if err == nil {
			c.logger.Debugw("Consumed message: ", "topic", msg.TopicPartition, "value", string(msg.Value), "op", op)

			// Deserialize full checkout message (user_id + products)
			var checkout products.CheckoutMessage
			err := DeserializeFromJSON(msg.Value, &checkout)
			if err != nil {
				c.logger.Errorw("Failed to deserialize message", "error", err, "op", op)
				continue
			}

			// Fallback: if key is set, prefer it for logging, but use JSON user_id for data
			if msg.Key != nil {
				if keyUserID, errConv := strconv.Atoi(string(msg.Key)); errConv == nil && int32(keyUserID) != checkout.UserID {
					c.logger.Debugw("Kafka key user_id differs from payload user_id", "key_user_id", keyUserID, "payload_user_id", checkout.UserID, "op", op)
				}
			}

			orderID, err := c.Service.CreateOrder(context.Background(), &order.Order{
				UserID: checkout.UserID,
				Products: func() []*order.ProductData {
					var orderProducts []*order.ProductData
					for _, p := range checkout.Products {
						orderProducts = append(orderProducts, &order.ProductData{
							ID:       p.ID,
							Quantity: p.Quantity,
						})
					}
					return orderProducts
				}(),
			})
			if err != nil {
				c.logger.Errorw("Failed to create order", "error", err, "op", op)
			} else {
				c.logger.Debugw("Order created successfully", "order_id", orderID, "op", op)
			}
			_, err = c.Consumer.CommitMessage(msg)
			if err != nil {
				c.logger.Errorw("Failed to commit message", "error", err, "op", op)
			}
			c.logger.Debugw("Message committed successfully", "topic", msg.TopicPartition, "op", op)
			continue

		} else {
			c.logger.Errorw("Consumer error: ", "error", err, "op", op)
		}
	}
}

func DeserializeFromJSON(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		return fmt.Errorf("failed to deserialize message: %w", err)
	}
	return nil
}
