# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Go implementation of the HCP 7 Tesla retinotopic mapping experiment. Presents visual stimuli (checkerboard bars, rotating wedges, expanding/contracting circles) to subjects for brain visual cortex mapping. Uses SDL3 via the [goxpyriment](https://github.com/chrplr/goxpyriment) framework.

## Build & Run

```bash
# Run directly (subject 0, run 1)
go run retinotopy.go -s 0 -r 1

# Build executable
go build -o retinotopy

# Cross-compile for all platforms
./build.sh           # outputs to build/
goreleaser build --snapshot --clean  # uses .goreleaser.yaml

# CGO_ENABLED=0 is required — SDL3 is loaded dynamically via purego (no C compiler needed)
```

**CLI flags:** `-r <run_id 1-6>`, `-s <subject_id>`, `-age`, `-gender`, `-handedness`, `-screen-width <cm>`, `-viewing-distance <cm>`, `-refresh-rate <Hz>`, `-fullscreen`, `-display <id, 0 = primary>`. All of them pre-fill the info dialog, and `-headless` (from goxpyriment) skips the dialog and returns those values directly. Note that `GetParticipantInfo` prefers a value cached in `~/.cache/goxpyriment/last_session.json` over a flag-supplied default (every field except `subject_id`), so interactively a flag only sets what the dialog starts from. Under `-headless` there is no dialog to correct it, so flags passed explicitly on the command line override the cache — see the `flag.Visit` block in `main()` — and each override is logged.

`-lead-in <seconds>` is the one flag that is not a dialog field: each stimulus
order begins with a fixation-only baseline (16 s for the bars, 22 s for the
wedges and rings) before the first pattern, and `trimLeadIn` drops all but the
requested seconds of it, trimming `MaskOrder`/`PatternOrder`/`DotOrder` together
so they stay in step. Default 2 s; a negative value keeps the CSV's own lead-in,
which is what a real scanning session needs.

**Run types (1–6):** bars (×2), CCW wedge, CW wedge, expanding circles, contracting circles.

## Architecture

Single-package application (`main`) in `retinotopy.go`:

- **`Retinotopy` struct** — holds all experiment state: screen, renderer, textures, loaded stimuli, timing config, data logger.
- **`LoadStimuli()`** — loads PNG patterns and masks from `assets/patterns/` and `assets/masks/` (embedded via `//go:embed`), fixation grid, and CSV stimulus orders from `assets/StimuliOrder/`.
- **`Run()`** — main loop: iterates frames, blends pattern + mask textures via `updateCombinedTexture()`, presents on screen, logs timing, detects color-change response (attention task).
- **`Instructions()`** — shows per-run instruction text and waits for experimenter keypress.

**goxpyriment** provides: screen/renderer lifecycle, clock/timing, input handling, and CSV data logging (results saved to `$HOME/goxpy_data/`, each run writing a `.csv` plus an `-info.txt` sidecar; override with `exp.SetOutputDirectory()`).

**`Stim-Generators/`** — standalone utilities for regenerating PNG assets (bar checkerboards, wedges, eccentric checkerboards). Not imported by main; run independently if assets need to be regenerated.

## Assets

Embedded assets in `assets/`:
- `patterns/` — 100 checkerboard phase patterns
- `masks/{swippingBars,rotatingWedge,expendingCircles}/` — mask sequences per run type
- `StimuliOrder/maskOrderRetinotopy.csv` + per-subject CSVs — frame-by-frame stimulus order
- `Inconsolata.ttf`, `fixationGrid.png`, `icons/`

## Browser (WebAssembly) build

`GOOS=js GOARCH=wasm` is a supported target, deployed to GitHub Pages by
`.github/workflows/pages.yml`. Two constraints, both verified by breaking them:

- **Upstream go-sdl3 cannot build for js** (`undefined: sdl.WaitAnimationFrame`).
  The build needs `chrplr/go-sdl3-wasm` (branch `wasm-render-fixes`), applied as
  a `replace` **in CI only** — committing it breaks the native release builds.
- **`control.GetParticipantInfo` panics in the browser** (`not implemented on
  js`, from `SDL_GetCurrentRenderOutputSize`): the dialog opens a second SDL
  window and shuts SDL down afterwards. `browser_js.go` therefore forces
  `-headless`, which takes GetParticipantInfo's early return (no SDL at all),
  and synthesizes the flags from the page URL — goxpyriment's own
  `platformPrepareFlags` does this but only inside `NewExperimentFromFlags`,
  which this program does not use.

- **`sdl.CreateSurfaceFrom` is broken on js** in the fork (`Cannot convert a
  BigInt value to a number`, from `syscall/js.Value.Call`). `loadTextureFromBytes`
  therefore uploads into a `CreateTexture` + `Update` texture rather than going
  through a surface — the path the combined stimulus texture already used. Avoid
  reintroducing `CreateSurfaceFrom`.

The bundle is ~135 MB because every stimulus PNG is embedded. Verified running in
Chrome: dialog skipped, URL parameters applied, stimuli loaded, fixation grid and
run loop rendering.

## Release / CI

`.goreleaser.yaml` and `.github/workflows/build.yml` handle multi-platform releases (Linux amd64/arm64, macOS amd64/arm64, Windows amd64/arm64) including deb/rpm packages, macOS DMG, and Windows NSIS installer. Triggered on `v*` tags.
