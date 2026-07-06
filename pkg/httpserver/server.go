package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const (
	_defaultAddr            = ":5050"
	_defaultReadTimeout     = 5 * time.Second
	_defaultWriteTimeout    = 5 * time.Second
	_defaultShutdownTimeout = 3 * time.Second
)

type Option func(*Server)

type Server struct {
	ctx context.Context
	eg  *errgroup.Group

	App    *fiber.App
	notify chan error

	address         string
	readTimeout     time.Duration
	writeTimeout    time.Duration
	shutdownTimeout time.Duration

	logger *zap.Logger
}

func NewServer(l *zap.Logger) *Server {
	group, ctx := errgroup.WithContext(context.Background())
	group.SetLimit(1)

	s := &Server{
		ctx:             ctx,
		eg:              group,
		App:             nil,
		address:         _defaultAddr,
		readTimeout:     _defaultReadTimeout,
		writeTimeout:    _defaultWriteTimeout,
		shutdownTimeout: _defaultShutdownTimeout,
		logger:          l,
	}

	app := fiber.New(fiber.Config{
		ReadTimeout:  s.readTimeout,
		WriteTimeout: s.writeTimeout,
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
	})
	s.App = app

	return s
}

func (s *Server) Start() {
	s.logger.Debug("START HTTP SERVER")

	s.eg.Go(func() error {
		err := s.App.Listen(s.address)
		if err != nil {
			s.notify <- fmt.Errorf("failed to start server: %w", err)
			close(s.notify)
			return err
		}

		return nil
	})
}

func (s *Server) WaitForShutdown(l zap.Logger) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	var err error

	select {
	case sig := <-interrupt:
		l.Debug(fmt.Sprintf("app - Run - signal: %s", sig.String()))
	case err = <-s.Notify():
		l.Error(fmt.Sprintf("app - Run - httpServer.Notify: %s", err))
	}

	s.Shutdown()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	var shutdownErrors []error

	err := s.App.ShutdownWithTimeout(s.shutdownTimeout)
	if err != nil && !errors.Is(err, context.Canceled) {
		s.logger.Error(fmt.Sprintf("httpServer shutdown with timeout: %s", err))

		shutdownErrors = append(shutdownErrors, err)
	}

	err = s.eg.Wait()
	if err != nil && !errors.Is(err, context.Canceled) {
		s.logger.Error(fmt.Sprintf("httpServer eg.Wait shutdown: %s", err))

		shutdownErrors = append(shutdownErrors, err)
	}

	s.logger.Debug("Graceful http server shutdown")

	return errors.Join(shutdownErrors...)
}
