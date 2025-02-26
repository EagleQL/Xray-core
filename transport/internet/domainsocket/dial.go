// +build !windows
// +build !wasm

package domainsocket

import (
	"context"

	"github.com/eagleql/xray-core/common"
	"github.com/eagleql/xray-core/common/net"
	"github.com/eagleql/xray-core/transport/internet"
	"github.com/eagleql/xray-core/transport/internet/tls"
	"github.com/eagleql/xray-core/transport/internet/xtls"
)

func Dial(ctx context.Context, dest net.Destination, streamSettings *internet.MemoryStreamConfig) (internet.Connection, error) {
	settings := streamSettings.ProtocolSettings.(*Config)
	addr, err := settings.GetUnixAddr()
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUnix("unix", nil, addr)
	if err != nil {
		return nil, newError("failed to dial unix: ", settings.Path).Base(err).AtWarning()
	}

	if config := tls.ConfigFromStreamSettings(streamSettings); config != nil {
		return tls.Client(conn, config.GetTLSConfig(tls.WithDestination(dest))), nil
	} else if config := xtls.ConfigFromStreamSettings(streamSettings); config != nil {
		return xtls.Client(conn, config.GetXTLSConfig(xtls.WithDestination(dest))), nil
	}

	return conn, nil
}

func init() {
	common.Must(internet.RegisterTransportDialer(protocolName, Dial))
}
