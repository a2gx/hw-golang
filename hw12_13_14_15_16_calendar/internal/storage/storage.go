package storage

import (
	"fmt"
	"strings"

	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/app"
	storage_memory "github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/storage/memory"
	storage_sql "github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/pkg/logger"
)

type Options struct {
	StorageType string
	Logg        *logger.Logger
}

func New(opts Options) (app.Storage, error) {
	switch strings.ToLower(opts.StorageType) {
	case "memory":
		return storage_memory.New(opts.Logg), nil
	case "sql":
		return storage_sql.New(opts.Logg), nil
	default:
		return nil, fmt.Errorf("unknown storage type: %s", opts.StorageType)
	}
}
