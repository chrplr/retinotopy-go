// Copyright (2026) Christophe Pallier <christophe@pallier.org>
// Distributed under the GNU General Public License v3.

//go:build js

package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"strings"
	"syscall/js"
)

// prepareFlags turns the page URL's query string into command-line arguments,
// and forces -headless. Call it after every flag is registered and immediately
// before flag.Parse.
//
// goxpyriment does the same translation in control.platformPrepareFlags, but
// that is unexported and runs only inside NewExperimentFromFlags — a
// constructor this program deliberately does not use, because it assembles its
// own dialog field set. Without a local copy, nothing on the page could set the
// subject, the run, or the monitor geometry: os.Args is empty in the browser.
//
// -headless is forced because the setup dialog cannot run here at all. It opens
// a second SDL window and shuts SDL down when it closes, neither of which a
// single-canvas page supports, and its renderer path calls
// SDL_GetCurrentRenderOutputSize, still a panic-stub on js — the program aborts
// before its first frame. Skipping the dialog takes GetParticipantInfo's early
// return, which touches no SDL whatsoever, and the values arrive from the URL
// instead.
func prepareFlags() {
	args := os.Args[:1:1]

	loc := js.Global().Get("location")
	if !loc.IsUndefined() {
		raw := strings.TrimPrefix(loc.Get("search").String(), "?")
		for _, pair := range strings.Split(raw, "&") {
			if pair == "" {
				continue
			}
			key, val, hasVal := strings.Cut(pair, "=")
			k, err := url.QueryUnescape(key)
			if err != nil || k == "" {
				continue
			}
			// Appended unconditionally below; a second copy would be harmless
			// but confusing in the logged argument list.
			if k == "headless" {
				continue
			}
			// Unknown keys (analytics parameters, a stray fragment) are skipped
			// rather than passed on: flag.Parse treats an unrecognised flag as
			// fatal, which in a browser means a blank canvas and no explanation.
			if flag.Lookup(k) == nil {
				log.Printf("browser: ignoring unknown URL parameter %q", k)
				continue
			}
			if !hasVal || val == "" {
				args = append(args, "-"+k)
				continue
			}
			v, err := url.QueryUnescape(val)
			if err != nil {
				log.Printf("browser: ignoring malformed URL parameter %q", pair)
				continue
			}
			args = append(args, "-"+k+"="+v)
		}
	}

	os.Args = append(args, "-headless")
}
