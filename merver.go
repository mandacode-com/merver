// Package merver
package merver

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/mandacode-com/merr"
)

type Server interface {
	Run(ctx context.Context) error
	Stop(ctx context.Context) error
}

type ServerGroup struct {
	servers []Server
}

func NewServerGroup(servers ...Server) *ServerGroup {
	return &ServerGroup{
		servers: servers,
	}
}

func (sm *ServerGroup) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(sm.servers))

	// Start all servers
	for _, server := range sm.servers {
		wg.Add(1)
		go func(s Server) {
			defer wg.Done()
			if err := s.Run(ctx); err != nil {
				errCh <- merr.New(merr.ErrInternalServerError, "Failed to run server", err.Error(), err)
			}
		}(server)
	}

	// Handle shutdown signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
	case <-signalChan:
	case err := <-errCh:
		return err
	}

	// Graceful shutdown
	for _, server := range sm.servers {
		if err := server.Stop(ctx); err != nil {
			return merr.New(merr.ErrInternalServerError, "Failed to stop server", err.Error(), err)
		}
	}

	// Ensure all Start routines complete
	wg.Wait()
	return nil
}
