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

**CLI flags:** `-s <subject_id>`, `-r <run_id 1-6>`, `-d` (dev mode), `--scaling <float>`, `-assets <path>`

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

## Release / CI

`.goreleaser.yaml` and `.github/workflows/build.yml` handle multi-platform releases (Linux amd64/arm64, macOS amd64/arm64, Windows amd64/arm64) including deb/rpm packages, macOS DMG, and Windows NSIS installer. Triggered on `v*` tags.
