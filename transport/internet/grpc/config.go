package grpc

import (
	"github.com/eagleql/xray-core/common"
	"github.com/eagleql/xray-core/transport/internet"
)

const protocolName = "grpc"

func init() {
	common.Must(internet.RegisterProtocolConfigCreator(protocolName, func() interface{} {
		return new(Config)
	}))
}
