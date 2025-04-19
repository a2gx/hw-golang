package storage

import (
	"fmt"
	"strings"

	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/app"
	storagememory "github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/storage/memory"
	storagesql "github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/pkg/logger"
)

type Options struct {
	StorageType string
	DatabaseDNS string
	Logg        *logger.Logger
}

func New(opts Options) (app.Storage, error) {
	switch strings.ToLower(opts.StorageType) {
	case "memory":
		return storagememory.New(opts.Logg), nil
	case "sql":
		return storagesql.New(opts.Logg, opts.DatabaseDNS), nil
	default:
		return nil, fmt.Errorf("unknown storage type: %s", opts.StorageType)
	}
}
