package server

import (
	"fmt"

	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/server/grpc"
	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/internal/server/http"

	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/pkg/logger"
)

type Options struct {
	ServerType string
	Logg       *logger.Logger
	App        *app.App
	AddrHttp   string
	AddrGrpc   string
}

func New(opts Options) (app.Server, error) {
	switch opts.ServerType {
	case "http":
		return serverhttp.New(opts.AddrHttp, opts.Logg, opts.App), nil
	case "grpc":
		return servergrpc.New(opts.AddrGrpc, opts.Logg, opts.App), nil
	default:
		return nil, fmt.Errorf("unknown server type: %s", opts.ServerType)
	}
}
