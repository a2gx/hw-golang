package app

import "context"

type Storage interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
}

type Server interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type Application interface{}
