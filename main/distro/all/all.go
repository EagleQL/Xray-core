package all

import (
	// The following are necessary as they register handlers in their init functions.

	// Required features. Can't remove unless there is replacements.
	_ "github.com/eagleql/xray-core/app/dispatcher"
	_ "github.com/eagleql/xray-core/app/proxyman/inbound"
	_ "github.com/eagleql/xray-core/app/proxyman/outbound"

	// Default commander and all its services. This is an optional feature.
	_ "github.com/eagleql/xray-core/app/commander"
	_ "github.com/eagleql/xray-core/app/log/command"
	_ "github.com/eagleql/xray-core/app/proxyman/command"
	_ "github.com/eagleql/xray-core/app/stats/command"

	// Other optional features.
	_ "github.com/eagleql/xray-core/app/dns"
	_ "github.com/eagleql/xray-core/app/dns/fakedns"
	_ "github.com/eagleql/xray-core/app/log"
	_ "github.com/eagleql/xray-core/app/policy"
	_ "github.com/eagleql/xray-core/app/reverse"
	_ "github.com/eagleql/xray-core/app/router"
	_ "github.com/eagleql/xray-core/app/stats"

	// Inbound and outbound proxies.
	_ "github.com/eagleql/xray-core/proxy/blackhole"
	_ "github.com/eagleql/xray-core/proxy/dns"
	_ "github.com/eagleql/xray-core/proxy/dokodemo"
	_ "github.com/eagleql/xray-core/proxy/freedom"
	_ "github.com/eagleql/xray-core/proxy/http"
	_ "github.com/eagleql/xray-core/proxy/mtproto"
	_ "github.com/eagleql/xray-core/proxy/shadowsocks"
	_ "github.com/eagleql/xray-core/proxy/socks"
	_ "github.com/eagleql/xray-core/proxy/trojan"
	_ "github.com/eagleql/xray-core/proxy/vless/inbound"
	_ "github.com/eagleql/xray-core/proxy/vless/outbound"
	_ "github.com/eagleql/xray-core/proxy/vmess/inbound"
	_ "github.com/eagleql/xray-core/proxy/vmess/outbound"

	// Transports
	_ "github.com/eagleql/xray-core/transport/internet/domainsocket"
	_ "github.com/eagleql/xray-core/transport/internet/grpc"
	_ "github.com/eagleql/xray-core/transport/internet/http"
	_ "github.com/eagleql/xray-core/transport/internet/kcp"
	_ "github.com/eagleql/xray-core/transport/internet/quic"
	_ "github.com/eagleql/xray-core/transport/internet/tcp"
	_ "github.com/eagleql/xray-core/transport/internet/tls"
	_ "github.com/eagleql/xray-core/transport/internet/udp"
	_ "github.com/eagleql/xray-core/transport/internet/websocket"
	_ "github.com/eagleql/xray-core/transport/internet/xtls"

	// Transport headers
	_ "github.com/eagleql/xray-core/transport/internet/headers/http"
	_ "github.com/eagleql/xray-core/transport/internet/headers/noop"
	_ "github.com/eagleql/xray-core/transport/internet/headers/srtp"
	_ "github.com/eagleql/xray-core/transport/internet/headers/tls"
	_ "github.com/eagleql/xray-core/transport/internet/headers/utp"
	_ "github.com/eagleql/xray-core/transport/internet/headers/wechat"
	_ "github.com/eagleql/xray-core/transport/internet/headers/wireguard"

	// JSON & TOML & YAML
	_ "github.com/eagleql/xray-core/main/json"
	_ "github.com/eagleql/xray-core/main/toml"
	_ "github.com/eagleql/xray-core/main/yaml"

	// Load config from file or http(s)
	_ "github.com/eagleql/xray-core/main/confloader/external"

	// Commands
	_ "github.com/eagleql/xray-core/main/commands/all"
)
