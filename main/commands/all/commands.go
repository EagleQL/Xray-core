package all

import (
	"github.com/eagleql/xray-core/main/commands/all/api"
	"github.com/eagleql/xray-core/main/commands/all/tls"
	"github.com/eagleql/xray-core/main/commands/base"
)

// go:generate go run github.com/eagleql/xray-core/common/errors/errorgen

func init() {
	base.RootCommand.Commands = append(
		base.RootCommand.Commands,
		api.CmdAPI,
		//cmdConvert,
		tls.CmdTLS,
		cmdUUID,
	)
}
