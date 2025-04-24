package outbox

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/literaen/simple_project/users/internal/config"

	"github.com/literaen/simple_project/pkg/kafka"

	"gorm.io/gorm"
)

type OutBoxService struct {
	repo OutBoxRepository
	kfw  *kafka.KFW
}

func NewOutBoxService(config *config.Config, repo OutBoxRepository) *OutBoxService {
	return &OutBoxService{
		repo: repo,
		kfw:  kafka.NewKafkaWriter(config.KAFKA_BROKERS, "users.events"),
	}
}

func (s *OutBoxService) AddEvent(tx *gorm.DB, eventType string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	event := OutboxEvent{
		EventType: eventType,
		Payload:   string(data),
		Processed: false,
	}

	if err := s.repo.CreateEvent(context.TODO(), tx, &event); err != nil {
		return fmt.Errorf("failed to insert outbox event: %v", err)
	}

	return nil
}

func (s *OutBoxService) GetUnprocessedEvents(ctx context.Context, batchSize int) ([]OutboxEvent, error) {
	return s.repo.GetUnprocessedEvents(ctx, batchSize)
}

func (s *OutBoxService) ProcessEvent(ctx context.Context, event *OutboxEvent) error {
	if err := s.kfw.WriteString(ctx, []byte(event.EventType), event.Payload); err != nil {
		return fmt.Errorf("error writing to kafka: %v", err)
	}

	// Помечаем как обработанное в транзакции
	return s.repo.MarkEventAsProcessed(ctx, event.ID)
}
