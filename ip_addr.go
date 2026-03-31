// Copyright 2022 Walter Scheper <walter.scheper@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid

import (
	"fmt"
	"net/netip"
)

const (
	ipv4Len = 4
	ipv6Len = 16
)

type ipGen struct {
	ipv6 bool
}

func (g *ipGen) String() string {
	return fmt.Sprintf("IP(ipv6=%v)", g.ipv6)
}

func (g *ipGen) value(t *T) netip.Addr {
	var b []byte
	if g.ipv6 {
		b = SliceOfN(Byte(), ipv6Len, ipv6Len).Draw(t, g.String())
	} else {
		b = SliceOfN(Byte(), ipv4Len, ipv4Len).Draw(t, g.String())
	}

	addr, _ := netip.AddrFromSlice(b)
	return addr
}

func IPv4() *Generator[netip.Addr] {
	return newGenerator[netip.Addr](&ipGen{
		ipv6: false,
	})
}

func IPv6() *Generator[netip.Addr] {
	return newGenerator[netip.Addr](&ipGen{
		ipv6: true,
	})
}
