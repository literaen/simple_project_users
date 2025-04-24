package outbox

import (
	"context"

	"github.com/google/uuid"
	"github.com/literaen/simple_project/pkg/postgres"
	"gorm.io/gorm"
)

type OutBoxRepository interface {
	// Сохраняет событие в транзакции
	CreateEvent(ctx context.Context, tx *gorm.DB, event *OutboxEvent) error

	// Получает необроботанные события (пачками)
	GetUnprocessedEvents(ctx context.Context, limit int) ([]OutboxEvent, error)

	// Помечает событие как обработанное
	MarkEventAsProcessed(ctx context.Context, id uuid.UUID) error
}

type outboxRepository struct {
	gdb *postgres.GDB
}

func NewOutBoxRepository(gdb *postgres.GDB) OutBoxRepository {
	return &outboxRepository{gdb: gdb}
}

func (r *outboxRepository) CreateEvent(ctx context.Context, tx *gorm.DB, event *OutboxEvent) error {
	return tx.WithContext(ctx).Create(event).Error
}

func (r *outboxRepository) GetUnprocessedEvents(ctx context.Context, limit int) ([]OutboxEvent, error) {
	var events []OutboxEvent
	err := r.gdb.DB.WithContext(ctx).
		Where("processed = ?", false).
		Order("created_at").
		Limit(limit).
		Find(&events).Error
	return events, err
}

func (r *outboxRepository) MarkEventAsProcessed(ctx context.Context, id uuid.UUID) error {
	return r.gdb.DB.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).
			Model(&OutboxEvent{}).
			Where("id = ?", id).
			Update("processed", true).Error
	})
}
