#!/usr/bin/env bash

set -eu

TLD_URL=https://data.iana.org/TLD/tlds-alpha-by-domain.txt
TLD="$(curl -fsSL https://data.iana.org/TLD/tlds-alpha-by-domain.txt | grep -v '^#')"

cat > tld.go << EOF
// Copyright 2020 Walter Scheper <walter.scheper@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid

import "strings"

// sourced from https://data.iana.org/TLD/tlds-alpha-by-domain.txt
// Version $(date +%Y%m%d00), Last Updated $(date --utc)
const tldsByAlpha = \`
${TLD}
\`

var tlds = strings.Split(strings.TrimSpace(tldsByAlpha), "\n")
EOF
