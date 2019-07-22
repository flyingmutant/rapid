// Copyright 2019 Gregory Petrosyan <gregory.petrosyan@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

/*
Package rapid implements utilities for property-based testing.

Rapid checks that properties you define hold for a large number
of automatically generated test cases. If a failure is found, rapid
fails the current test and presents an automatically minimized
version of the failing test case.

Please note that rapid is alpha software; the documentation
is very incomplete, unclear and probably full of grammatical errors.
*/
package rapid
