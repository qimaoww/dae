/*
 * SPDX-License-Identifier: AGPL-3.0-only
 * Copyright (c) 2022-2025, daeuniverse Organization <dae@v2raya.org>
 */

package outbound

import (
	_ "github.com/qimaoww/outbound/dialer/anytls"
	_ "github.com/qimaoww/outbound/dialer/http"
	_ "github.com/qimaoww/outbound/dialer/hysteria2"
	_ "github.com/qimaoww/outbound/dialer/juicity"
	_ "github.com/qimaoww/outbound/dialer/shadowsocks"
	_ "github.com/qimaoww/outbound/dialer/shadowsocksr"
	_ "github.com/qimaoww/outbound/dialer/socks"
	_ "github.com/qimaoww/outbound/dialer/trojan"
	_ "github.com/qimaoww/outbound/dialer/tuic"
	_ "github.com/qimaoww/outbound/dialer/v2ray"
	_ "github.com/qimaoww/outbound/protocol/anytls"
	_ "github.com/qimaoww/outbound/protocol/hysteria2"
	_ "github.com/qimaoww/outbound/protocol/juicity"
	_ "github.com/qimaoww/outbound/protocol/shadowsocks"
	_ "github.com/qimaoww/outbound/protocol/trojanc"
	_ "github.com/qimaoww/outbound/protocol/tuic"
	_ "github.com/qimaoww/outbound/protocol/vless"
	_ "github.com/qimaoww/outbound/protocol/vmess"
	_ "github.com/qimaoww/outbound/transport/simpleobfs"
	_ "github.com/qimaoww/outbound/transport/tls"
	_ "github.com/qimaoww/outbound/transport/ws"
)
