package nats

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/cenkalti/backoff/v4"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

type options struct {
	NATSServerOptions *server.Options
	EnableLogging     bool
}

type Option func(*options)

func WithNATSServerOptions(natsServerOptions *server.Options) Option {
	return func(o *options) {
		o.NATSServerOptions = natsServerOptions
	}
}

func WithLogging() Option {
	return func(o *options) {
		o.EnableLogging = true
	}
}

type Server struct {
	NatsServer *server.Server
}

func New(ctx context.Context, opts ...Option) (*Server, error) {
	options := &options{}
	for _, o := range opts {
		o(options)
	}
	options.NATSServerOptions.DontListen = true

	ns, err := server.NewServer(options.NATSServerOptions)
	if err != nil {
		return nil, err
	}

	if options.EnableLogging {
		ns.ConfigureLogger()
	}

	go func() {
		<-ctx.Done()
		ns.Shutdown()
	}()

	ns.Start()

	return &Server{NatsServer: ns}, nil
}

func (n *Server) Close() error {
	if n.NatsServer != nil && n.NatsServer.Running() {
		n.NatsServer.Shutdown()
	}
	return nil
}

func (n *Server) WaitForServer() {
	b := backoff.NewExponentialBackOff()

	for {
		d := b.NextBackOff()
		ready := n.NatsServer.ReadyForConnections(d)
		if ready {
			break
		}

		slog.Info(fmt.Sprintf("NATS server not ready, waited %s, retrying...", d))
	}
}

func (n *Server) Client() (*nats.Conn, error) {
	opts := []nats.Option{
		nats.InProcessServer(n.NatsServer),
	}
	return nats.Connect(n.NatsServer.ClientURL(), opts...)
}
