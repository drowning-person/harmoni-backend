package app

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"harmoni/internal/types/iface"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go.uber.org/zap"
)

// Application is the main struct of the application
type Application struct {
	logger  *zap.SugaredLogger
	id      string
	name    string
	version string
	runners []iface.Executor
	signals []os.Signal
}

// Option application support option
type Option func(application *Application)

// NewApp creates a new Application
func NewApp(logger *zap.SugaredLogger, ops ...Option) *Application {
	app := &Application{
		logger: logger,
	}
	for _, op := range ops {
		op(app)
	}
	// default random id
	if len(app.id) == 0 {
		bytes := make([]byte, 24)
		_, _ = rand.Read(bytes)
		app.id = hex.EncodeToString(bytes)
	}
	// default accept signals
	if len(app.signals) == 0 {
		app.signals = []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT}
	}
	return app
}

// WithID application add id
func WithID(id string) func(application *Application) {
	return func(application *Application) {
		application.id = id
	}
}

// WithName application add name
func WithName(name string) func(application *Application) {
	return func(application *Application) {
		application.name = name
	}
}

// WithVersion application add version
func WithVersion(version string) func(application *Application) {
	return func(application *Application) {
		application.version = version
	}
}

// WithServer application add server
func WithServer(runners ...iface.Executor) func(application *Application) {
	return func(application *Application) {
		application.runners = runners
	}
}

// WithSignals application add listen signals
func WithSignals(signals []os.Signal) func(application *Application) {
	return func(application *Application) {
		application.signals = signals
	}
}

func (app *Application) recover() {
	if err := recover(); err != nil {
		app.logger.Error(err)
	}
}

// Run application run
func (app *Application) Run(ctx context.Context) error {
	if len(app.runners) == 0 {
		return nil
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, app.signals...)
	errCh := make(chan error, 1)

	for _, s := range app.runners {
		go func(srv iface.Executor) {
			defer app.recover()
			if err := srv.Start(); err != nil {
				app.logger.Errorf("failed to start server, err: %s", err)
				errCh <- err
			}
		}(s)
	}

	select {
	case err := <-errCh:
		_ = app.Stop()
		return err
	case <-ctx.Done():
		return app.Stop()
	case <-quit:
		return app.Stop()
	}
}

// Stop application stop
func (app *Application) Stop() error {
	wg := sync.WaitGroup{}
	for _, s := range app.runners {
		wg.Add(1)
		go func(srv iface.Executor) {
			defer wg.Done()
			defer app.recover()
			if err := srv.Shutdown(); err != nil {
				app.logger.Errorf("failed to stop server, err: %s", err)
			}
		}(s)
	}
	// wait all server graceful shutdown
	wg.Wait()
	return nil
}
