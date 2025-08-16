package mervermid

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

// GRPCZeroLogger handles AppError and logs gRPC errors consistently.
func GRPCZeroLogger(logger *zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)

		if err != nil {
			logger.Error().
				Err(err).
				Str("method", info.FullMethod).
				Str("trace", fmt.Sprintf("%+v", err)).
				Dur("duration", duration).
				Msg("gRPC request error")
		}
		if resp == nil {
			logger.Warn().
				Str("method", info.FullMethod).
				Dur("duration", duration).
				Msg("gRPC request completed with warning: nil response")
		}

		return resp, err
	}
}
