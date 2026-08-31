// Copyright (2026) Christophe Pallier <christophe@pallier.org>
// Distributed under the GNU General Public License v3.

//go:build !js

package main

// prepareFlags is a no-op on desktop: the flags come from the real command
// line, and the setup dialog works. See browser_js.go for what it does in the
// browser, where neither is true.
func prepareFlags() {}
