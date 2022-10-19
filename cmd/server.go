package cmd

import (
	"context"
	"errors"
	logger "github.com/rs/zerolog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"test/handler"
	"time"
)

var Log = logger.New(os.Stdout)

// Server is a light wrapper around http.Server.
type Server struct {
	*http.Server
	address  string
	listener net.Listener
	// Channel used to signal server has shutdown
	serverShutdown chan struct{}
}

// StartNew creates a Server and starts it.
func StartNew() error {
	http.HandleFunc("/api", handler.GetDescription)
	http.HandleFunc("/", handler.ErrorHandler)
	srv, err := New(&http.Server{})
	if err != nil {
		return err
	}

	return srv.Start()
}

// New creates a new Server, built off of a base http.Server.
func New(s *http.Server) (*Server, error) {
	// Default server timout, in seconds
	const defaultSrvTimeout = 15 * time.Second
	const port = "8081"

	var (
		srv *Server
		err error
	)

	// ensure timeouts are set
	if s.ReadTimeout == 0 {
		s.ReadTimeout = defaultSrvTimeout
	}

	if s.WriteTimeout == 0 {
		s.WriteTimeout = defaultSrvTimeout
	}

	listener, err := net.Listen("tcp", net.JoinHostPort("0.0.0.0", port))
	if err != nil {
		Log.Error().Err(err).Msgf("%s unavailable", port)
	}
	srv = &Server{
		s,
		port,
		listener,
		make(chan struct{}),
	}

	return srv, nil
}

// Start begins serving, and listens for termination signals to shutdown gracefully.
func (srv *Server) Start() error {
	var err error

	go srv.shutdown()

	Log.Log().
		Str("address", srv.address).
		Int("pid", os.Getpid()).
		Msg("Server listening")

	err = srv.Serve(srv.listener)
	if err != nil && err != http.ErrServerClosed {
		return errors.New("server failed to start")
	}

	<-srv.serverShutdown

	return nil
}

// Shutdown server gracefully on SIGINT or SIGTERM.
func (srv *Server) shutdown() {
	// Block until signal is received
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	// Allow up to thirty seconds for server operations to finish before
	// canceling them.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		Log.Error().
			Err(err).
			Msg("Server shutdown error")
	}

	Log.Log().Msg("Server shutdown")

	// Close channel to signal shutdown is complete
	close(srv.serverShutdown)
}
