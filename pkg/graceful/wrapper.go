package graceful

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KennyChenFight/golib/loglib"
	"go.uber.org/zap"
)

func Wrapper(logger *loglib.Logger, fn func(ctx context.Context) error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan error, 1)

	go func() {
		done <- fn(ctx)
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-done:
		return
	case <-shutdown:
		cancel()
		timeoutCtx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		select {
		case <-done:
			return
		case <-timeoutCtx.Done():
			logger.Error("shutdown timeout", zap.Error(timeoutCtx.Err()))
			return
		}
	}
}
