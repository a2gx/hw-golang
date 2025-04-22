package servergrpc

import (
	"context"
	"time"

	"github.com/alxbuylov/hw-golang/hw12_13_14_15_calendar/pkg/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

func loggingInterceptor(logg *logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()

		// Получаем информацию о клиенте
		p, ok := peer.FromContext(ctx)
		remoteAddr := "unknown"
		if ok {
			remoteAddr = p.Addr.String()
		}

		resp, err := handler(ctx, req)
		latency := time.Since(start)

		logg.Info(
			"gRPC request processed",
			"remote_addr", remoteAddr,
			"time", start.Format("02/Jan/2006:15:04:05 -0700"),
			"method", info.FullMethod,
			"latency_ms", latency.Milliseconds(),
			"error", err != nil,
		)

		return resp, err
	}
}
