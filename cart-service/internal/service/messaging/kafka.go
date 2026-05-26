package messaging

import (
	"context"
	"encoding/json"
	"os"
	"strconv"

	models "github.com/sabirkekw/ecommerce_go/cart-service/internal/models/product"
	"github.com/sabirkekw/ecommerce_go/pkg/apierrors"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.uber.org/zap"
)

type KafkaProducer struct {
	logger   *zap.SugaredLogger
	Producer *kafka.Producer
	topic    string
}

func New(logger *zap.SugaredLogger, topic string) *KafkaProducer {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BROKERS"),
	})
	if err != nil {
		panic(err)
	}

	return &KafkaProducer{
		logger:   logger,
		Producer: producer,
		topic:    topic,
	}
}

func (p *KafkaProducer) SendCheckoutMessage(ctx context.Context, userID int32, products []*models.ProductData) error {
	payload, err := SerializeToJSON(userID, products)
	if err != nil {
		p.logger.Errorw("Failed to serialize message", "error", err)
		return apierrors.ErrUnknown
	}

	deliveryChan := make(chan kafka.Event, 1)
	defer close(deliveryChan)

	err = p.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &p.topic,
			Partition: kafka.PartitionAny,
		},
		Key:   []byte(strconv.Itoa(int(userID))),
		Value: payload,
	}, deliveryChan)

	if err != nil {
		p.logger.Errorw("Failed to produce message", "error", err)
		return apierrors.ErrUnknown
	}

	select {
	case e := <-deliveryChan:
		m := e.(*kafka.Message)
		if m.TopicPartition.Error != nil {
			p.logger.Errorw("Failed to deliver message", "error", m.TopicPartition.Error)
			return apierrors.ErrUnknown
		}
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}

func (p *KafkaProducer) Close() {
	p.Producer.Flush(5000)
	p.Producer.Close()
}

func SerializeToJSON(userID int32, products []*models.ProductData) ([]byte, error) {
	type ProductIDQuantity struct {
		ID       int32 `json:"id"`
		Quantity int32 `json:"quantity"`
	}
	type CheckoutMessage struct {
		UserID   int32                `json:"user_id"`
		Products []*ProductIDQuantity `json:"products"`
	}

	var productsData []*ProductIDQuantity
	for _, p := range products {
		productsData = append(productsData, &ProductIDQuantity{
			ID:       p.ID,
			Quantity: p.Quantity,
		})
	}

	message := CheckoutMessage{
		UserID:   userID,
		Products: productsData,
	}
	return json.Marshal(message)
}
