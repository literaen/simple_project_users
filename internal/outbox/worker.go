package outbox

import (
	"context"
	"log"
	"time"
)

type OutboxWorker struct {
	service *OutBoxService
}

func NewOutboxWorker(service *OutBoxService) *OutboxWorker {
	return &OutboxWorker{service}
}

func (w *OutboxWorker) Start(ctx context.Context, interval time.Duration, batchSize int) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log.Println("outbox worker started")

	for {
		select {
		case <-ctx.Done():
			log.Println("outbox worker stopped")
			return
		case <-ticker.C:
			w.processBatch(ctx, batchSize)
		}
	}
}

func (w *OutboxWorker) processBatch(ctx context.Context, batchSize int) {
	events, err := w.service.GetUnprocessedEvents(ctx, batchSize)
	if err != nil {
		log.Printf("GetUnprocessedEvents error: %v", err)
		return
	}

	for _, event := range events {
		if err := w.service.ProcessEvent(ctx, &event); err != nil {
			log.Printf("event send error: %v", err)
		} else {
			log.Printf("event sent and marked: %v", event.ID)
		}
	}
}
