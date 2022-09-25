// Copyright 2022 Walter Scheper <walter.scheper@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid

import (
	"fmt"
	"net"
)

const (
	ipv4Len = 4
	ipv6Len = 16
)

var (
	ipv4Gen = SliceOfN(Byte(), ipv4Len, ipv4Len)
	ipv6Gen = SliceOfN(Byte(), ipv6Len, ipv6Len)
)

type ipGen struct {
	ipv6 bool
}

func (g *ipGen) String() string {
	return fmt.Sprintf("IP(ipv6=%v)", g.ipv6)
}

func (g *ipGen) value(t *T) net.IP {
	var gen *Generator[[]byte]
	if g.ipv6 {
		gen = ipv6Gen
	} else {
		gen = ipv4Gen
	}

	b := gen.Draw(t, g.String())
	return net.IP(b)
}

func IPv4() *Generator[net.IP] {
	return newGenerator[net.IP](&ipGen{
		ipv6: false,
	})
}

func IPv6() *Generator[net.IP] {
	return newGenerator[net.IP](&ipGen{
		ipv6: true,
	})
}
