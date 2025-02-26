package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/eagleql/xray-core/common"
	"github.com/eagleql/xray-core/common/net"
	"github.com/eagleql/xray-core/common/session"
	"github.com/eagleql/xray-core/transport/internet"
	"github.com/eagleql/xray-core/transport/internet/grpc/encoding"
	"github.com/eagleql/xray-core/transport/internet/tls"
)

type Listener struct {
	encoding.UnimplementedGRPCServiceServer
	ctx     context.Context
	handler internet.ConnHandler
	local   net.Addr
	config  *Config
	locker  *internet.FileLocker // for unix domain socket

	s *grpc.Server
}

func (l Listener) Tun(server encoding.GRPCService_TunServer) error {
	tunCtx, cancel := context.WithCancel(l.ctx)
	l.handler(encoding.NewHunkConn(server, cancel))
	<-tunCtx.Done()
	return nil
}

func (l Listener) TunMulti(server encoding.GRPCService_TunMultiServer) error {
	tunCtx, cancel := context.WithCancel(l.ctx)
	l.handler(encoding.NewMultiHunkConn(server, cancel))
	<-tunCtx.Done()
	return nil
}

func (l Listener) Close() error {
	l.s.Stop()
	return nil
}

func (l Listener) Addr() net.Addr {
	return l.local
}

func Listen(ctx context.Context, address net.Address, port net.Port, settings *internet.MemoryStreamConfig, handler internet.ConnHandler) (internet.Listener, error) {
	grpcSettings := settings.ProtocolSettings.(*Config)
	var listener *Listener
	if port == net.Port(0) { // unix
		listener = &Listener{
			handler: handler,
			local: &net.UnixAddr{
				Name: address.Domain(),
				Net:  "unix",
			},
			config: grpcSettings,
		}
	} else { // tcp
		listener = &Listener{
			handler: handler,
			local: &net.TCPAddr{
				IP:   address.IP(),
				Port: int(port),
			},
			config: grpcSettings,
		}
	}

	listener.ctx = ctx

	config := tls.ConfigFromStreamSettings(settings)

	var s *grpc.Server
	if config == nil {
		s = grpc.NewServer()
	} else {
		s = grpc.NewServer(grpc.Creds(credentials.NewTLS(config.GetTLSConfig(tls.WithNextProto("h2")))))
	}
	listener.s = s

	if settings.SocketSettings != nil && settings.SocketSettings.AcceptProxyProtocol {
		newError("accepting PROXY protocol").AtWarning().WriteToLog(session.ExportIDToError(ctx))
	}

	go func() {
		var streamListener net.Listener
		var err error
		if port == net.Port(0) { // unix
			streamListener, err = internet.ListenSystem(ctx, &net.UnixAddr{
				Name: address.Domain(),
				Net:  "unix",
			}, settings.SocketSettings)
			if err != nil {
				newError("failed to listen on ", address).Base(err).AtError().WriteToLog(session.ExportIDToError(ctx))
				return
			}
			locker := ctx.Value(address.Domain())
			if locker != nil {
				listener.locker = locker.(*internet.FileLocker)
			}
		} else { // tcp
			streamListener, err = internet.ListenSystem(ctx, &net.TCPAddr{
				IP:   address.IP(),
				Port: int(port),
			}, settings.SocketSettings)
			if err != nil {
				newError("failed to listen on ", address, ":", port).Base(err).AtError().WriteToLog(session.ExportIDToError(ctx))
				return
			}
		}

		encoding.RegisterGRPCServiceServerX(s, listener, grpcSettings.ServiceName)

		if err = s.Serve(streamListener); err != nil {
			newError("Listener for gRPC ended").Base(err).WriteToLog()
		}
	}()

	return listener, nil
}

func init() {
	common.Must(internet.RegisterTransportListener(protocolName, Listen))
}
