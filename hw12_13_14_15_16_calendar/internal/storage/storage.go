package storage

import (
	"fmt"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/pkg/logger"

	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/storage/postgres"
)

type Options struct {
	StorageType string
	Logg        *logger.Logger
}

func New(opts Options) (app.Storage, error) {
	switch opts.StorageType {
	case "memory":
		return memory.New(opts.Logg), nil
	case "postgres":
		return postgres.New(), nil
	default:
		return nil, fmt.Errorf("unknown storage type: %s", opts.StorageType)
	}
}
