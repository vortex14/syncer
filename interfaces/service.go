package interfaces

import (
	"context"
	"time"
)

// Batch is a batch of items.
type Batch []Item

// Item is some abstract item.
type Item struct{}

// Service defines external service that can process batches of items.
type Service interface {
	GetLimits() (n uint64, p time.Duration)
	Process(ctx context.Context, batch Batch) error
}
