package outbox

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OutboxEvent struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	EventType string
	Payload   string
	Processed bool
	CreatedAt time.Time
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&OutboxEvent{})
}
