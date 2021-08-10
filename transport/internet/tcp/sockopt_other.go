// +build !linux,!freebsd

package tcp

import (
	"github.com/eagleql/xray-core/common/net"
	"github.com/eagleql/xray-core/transport/internet"
)

func GetOriginalDestination(conn internet.Connection) (net.Destination, error) {
	return net.Destination{}, nil
}
