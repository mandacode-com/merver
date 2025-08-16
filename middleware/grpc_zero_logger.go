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

		if err != nil {
			logger.Error().
				Err(err).
				Str("method", info.FullMethod).
				Dur("duration", duration).
				Str("status", status.Code(err).String()).
				Str("trace", fmt.Sprintf("%+v", err)).
				Str("ip", addr).
				Str("auth_info", authType).
				Msg("gRPC request error")
		} else if resp == nil {
			if logger.GetLevel() == zerolog.WarnLevel {
				logger.Warn().
					Str("method", info.FullMethod).
					Dur("duration", duration).
					Str("status", status.Code(err).String()).
					Str("ip", addr).
					Str("auth_info", authType).
					Msg("gRPC request completed with warning: nil response")
			}
		} else {
			if logger.GetLevel() == zerolog.InfoLevel {
				logger.Info().
					Str("method", info.FullMethod).
					Dur("duration", duration).
					Str("status", status.Code(err).String()).
					Str("ip", addr).
					Str("auth_info", authType).
					Msg("gRPC request completed")
			}
		}

		return resp, err
	}
}
