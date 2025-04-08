package app

import "context"

type Storage interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
}

type Application interface{}
