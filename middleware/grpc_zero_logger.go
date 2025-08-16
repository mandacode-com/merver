package mervermid

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
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
		peer, ok := peer.FromContext(ctx)

		var addr string
		var authType string

		if ok {
			addr = peer.Addr.String()
			if ua, ok := peer.AuthInfo.(interface{ AuthType() string }); ok {
				authType = ua.AuthType()
			} else {
				authType = "no auth info"
			}
		} else {
			addr = "unknown"
			authType = "unknown"
		}

		level := logger.GetLevel()

		if err != nil && level <= zerolog.ErrorLevel {
			event := logger.Error()

			event.Err(err).
				Str("method", info.FullMethod).
				Dur("duration", duration).
				Str("status", status.Code(err).String()).
				Str("ip", addr).
				Str("auth_info", authType)

			if level <= zerolog.TraceLevel {
				event.Str("trace", fmt.Sprintf("%+v", err))
			}

			event.Msg("gRPC request error")
		} else if resp == nil && level <= zerolog.WarnLevel {
			event := logger.Warn()

			event.Str("method", info.FullMethod).
				Dur("duration", duration).
				Str("status", status.Code(err).String()).
				Str("ip", addr).
				Str("auth_info", authType)

			event.Msg("gRPC request completed with warning: nil response")
		} else if level <= zerolog.InfoLevel {
			event := logger.Info()
			event.Str("method", info.FullMethod).
				Dur("duration", duration).
				Str("status", status.Code(err).String()).
				Str("ip", addr).
				Str("auth_info", authType)

			event.Msg("gRPC request completed")
		}

		return resp, err
	}
}
