// Copyright 2020 Walter Scheper <walter.scheper@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid

import (
	"fmt"
	"net"
	"reflect"
)

const (
	ipv4Len = 4
	ipv6Len = 16
)

var (
	ipType = reflect.TypeOf(net.IP{})

	ipv4Gen = SliceOfN(Byte(), ipv4Len, ipv4Len)
	ipv6Gen = SliceOfN(Byte(), ipv6Len, ipv6Len)
)

type ipGen struct {
	ipv6 bool
}

func (g *ipGen) String() string {
	return fmt.Sprintf("IP(ipv6=%v)", g.ipv6)
}

func (*ipGen) type_() reflect.Type {
	return ipType
}

func (g *ipGen) value(t *T) value {
	var gen *Generator
	if g.ipv6 {
		gen = ipv6Gen
	} else {
		gen = ipv4Gen
	}

	b := gen.Draw(t, g.String()).([]byte)
	return net.IP(b)
}

func IPv4() *Generator {
	return newGenerator(&ipGen{
		ipv6: false,
	})
}

func IPv6() *Generator {
	return newGenerator(&ipGen{
		ipv6: true,
	})
}
