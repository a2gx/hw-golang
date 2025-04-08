package postgres

import (
	"context"

	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/app"
)

type Storage struct { // TODO
}

func New() *Storage {
	return &Storage{}
}

var _ app.Storage = &Storage{} // check the interface

func (s *Storage) Connect(ctx context.Context) error {
	// TODO
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	// TODO
	return nil
}
