package memory

import (
	"context"
	"sync"

	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/app"
)

type Storage struct {
	// TODO
	mu sync.RWMutex //nolint:unused
}

var _ app.Storage = &Storage{} // check the interface

func New() *Storage {
	return &Storage{}
}

// Close implements app.Storage.
func (s *Storage) Close(ctx context.Context) error {
	panic("unimplemented")
}

// Connect implements app.Storage.
func (s *Storage) Connect(ctx context.Context) error {
	panic("unimplemented")
}
