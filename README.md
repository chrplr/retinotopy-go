# Retinotopy Experiment (Go Implementation)

![](retino-stim_small.png)

This app runs the **HCP Retinotopic Mapping experiment** — a visual stimulus protocol used in brain imaging (fMRI) studies to map how the visual cortex is organized. Subjects watch a series of flickering checkerboard patterns (bars, rotating wedges, expanding/contracting rings) while occasionally pressing a key when a central dot changes color. The recorded timing data is used to reconstruct a map of the visual field on the brain.

It is also a working example of how to use the [goxpyriment](https://github.com/chrplr/goxpyriment) library to build standalone experiment executables in Go.

The experiment is based on:

> Benson, N. C., Jamison, K. W., Arcaro, M. J., Vu, A. T., Glasser, M. F., Coalson, T. S., Van Essen, D. C., Yacoub, E., Ugurbil, K., Winawer, J., & Kay, K. (2018). The Human Connectome Project 7 Tesla retinotopy dataset: Description and population receptive field analysis. *Journal of Vision*, 18(13), 23. https://doi.org/10.1167/18.13.23

> **Warning: Stimulus timing has not been independently validated yet. Use with caution in actual experiments.**

*Found a bug? Please report it at <https://github.com/chrplr/retinotopy-go/issues>*

Christophe Pallier, 05/03/2026 [![DOI](https://zenodo.org/badge/1173914957.svg)](https://doi.org/10.5281/zenodo.18887912)

---

## 1. Installation from binaries



Download the latest installer for your operating system directly below. All required files are bundled — no extra software needed.

### Windows
1. Download [retinotopy-windows-x86_64-setup.exe](https://github.com/chrplr/retinotopy-go/releases/latest/download/retinotopy-windows-x86_64-setup.exe) (Intel/AMD) or [retinotopy-windows-arm64-setup.exe](https://github.com/chrplr/retinotopy-go/releases/latest/download/retinotopy-windows-arm64-setup.exe) (ARM).
2. Double-click the file and follow the installer steps.
3. A **Retinotopy** shortcut will appear in your Start Menu and on your Desktop.

### macOS
1. Download [retinotopy-macos-arm64-app.zip](https://github.com/chrplr/retinotopy-go/releases/latest/download/retinotopy-macos-arm64-app.zip) (M1/M2/M3/M4) or [retinotopy-macos-x86_64-app.zip](https://github.com/chrplr/retinotopy-go/releases/latest/download/retinotopy-macos-x86_64-app.zip) (Intel).
2. Extract the archive and drag **Retinotopy.app** to your **Applications** folder (or anywhere you like).
3. Double-click to launch.

> [!WARNING]
> macOS may show a security warning the first time you open the app. See [macOS installation and security](https://chrplr.github.io/note-about-macos-unsigned-apps) for an explanation and step-by-step instructions to bypass it.

### Linux
1. Download [retinotopy-linux-x86_64.AppImage](https://github.com/chrplr/retinotopy-go/releases/latest/download/retinotopy-linux-x86_64.AppImage) (Intel/AMD) or [retinotopy-linux-aarch64.AppImage](https://github.com/chrplr/retinotopy-go/releases/latest/download/retinotopy-linux-aarch64.AppImage) (ARM).
2. Make it executable: right-click the file → **Properties** → **Permissions** → check **"Allow executing file as program"**.
3. Double-click to run.

---

## 2. Alternative: ZIP archives and system packages

If you prefer not to use an installer, download a raw binary archive from the [Releases page](https://github.com/chrplr/retinotopy-go/releases).

- **Windows ZIP:** Extract and run `retinotopy.exe` (SDL3.dll is included in the archive).
- **Linux packages:** `.deb` (Ubuntu/Debian) and `.rpm` (Fedora/RHEL) are available on the Releases page.

---

## 3. Running the Experiment

### Controls

| Action | Effect |
| :--- | :--- |
| **Any key or mouse click** | Register response when the fixation dot changes color |
| **ESC** | Exit and save data |

### Choosing a run

A setup dialog opens first, where the subject ID, run and monitor geometry are
entered. Flags on the command line pre-fill it, so a session can be prepared in
one line and confirmed with a click:

```
retinotopy -s 1 -r 3
```

The six available runs are:

| `-r` | Name | Stimulus |
| :--- | :--- | :--- |
| 1 | RETBAR1 | Swiping bars (first pass) |
| 2 | RETBAR2 | Swiping bars (second pass) |
| 3 | RETCCW | Counter-clockwise rotating wedge |
| 4 | RETCW | Clockwise rotating wedge |
| 5 | RETEXP | Expanding circles |
| 6 | RETCON | Contracting circles |

### All command-line options

Each flag sets the starting value of one dialog field; all of them can still be
changed in the dialog itself.

| Flag | Description | Default |
| :--- | :--- | :--- |
| `-r <1-6>` | Run number (see above) | `1` |
| `-s <id>` | Subject ID (selects the per-subject stimulus order, and names the results file) | *(empty)* |
| `-age <text>` | Participant age | *(empty)* |
| `-gender <text>` | Participant gender (M / F / NB) | *(empty)* |
| `-handedness <R\|L>` | Participant handedness | `R` |
| `-screen-width <cm>` | Screen width in centimetres | `30` |
| `-viewing-distance <cm>` | Viewing distance in centimetres | `50` |
| `-refresh-rate <Hz>` | Display refresh rate | `60` |
| `-display <id>` | Monitor to open on, `0` = primary | `0` |
| `-fullscreen` | Start fullscreen; `-fullscreen=false` gives a 1024×768 window | `true` |
| `-headless` | Skip the dialog and run with the values above | off |
| `-lead-in <s>` | Seconds of the run's opening blank period to keep; negative keeps the CSV's own | `2` |

`-lead-in` is not a dialog field. Every stimulus order opens with a
fixation-only baseline — 16 s for the bars, 22 s for the wedges and rings —
before the first checkerboard appears, which is a long wait when checking the
stimuli, so those frames are trimmed to 2 s by default. **Pass `-lead-in=-1` for
a scanning session**: the dropped frames are baseline the analysis expects.

Screen width and viewing distance are what set the stimulus size: the
checkerboard is scaled so its maximum eccentricity is 15° of visual angle,
clamped so that it still fits on screen. The resulting scaling and eccentricity
are recorded alongside the data, so measure them rather than leaving the
defaults.

**Remembered between sessions.** Every field except the subject ID is cached in
`~/.cache/goxpyriment/last_session.json` and restored ahead of a flag's default,
so the dialog reopens where the last session left it. Under `-headless` there is
no dialog to correct that, so any flag given explicitly on the command line
overrides the cached value, and each override is logged.

---

## 4. Building from Source

Only needed if you want to modify the code or compile it yourself.

### Prerequisites

1. **Install Go:** [go.dev/doc/install](https://go.dev/doc/install)

### Build steps

```bash
git clone https://github.com/chrplr/retinotopy-go.git
cd retinotopy-go
go mod download
go run retinotopy.go -s 0 -r 1     # run directly
go build -o retinotopy              # build an executable
```

To cross-compile for all platforms at once:

```bash
./build.sh
```

### Browser (WebAssembly) build

The experiment also builds for `GOOS=js GOARCH=wasm` and is deployed to GitHub
Pages by [`.github/workflows/pages.yml`](.github/workflows/pages.yml) on every
push to `main`.

Two things make it work, and both are easy to trip over:

- **It needs the `chrplr/go-sdl3-wasm` fork**, not upstream go-sdl3, which lacks
  the js bindings (building against upstream fails with `undefined:
  sdl.WaitAnimationFrame`). The `replace` lives in CI only and **must never be
  committed** — it would break the native release builds for everyone.
- **The setup dialog cannot run in a browser.** It opens a second SDL window and
  shuts SDL down when it closes, which a single-canvas page does not support, so
  the browser build forces `-headless` and takes its settings from the page URL
  instead (`?s=7&r=3&screen-width=40`). See `browser_js.go`.

A third trap is already worked around: `SDL_CreateSurfaceFrom`'s js binding in
the fork throws `Cannot convert a BigInt value to a number`, which killed the
run while loading the fixation grid. `loadTextureFromBytes` now uploads straight
into a texture instead of building a surface first.

To build and serve the bundle locally:

```bash
git clone -b wasm-render-fixes https://github.com/chrplr/go-sdl3-wasm ../go-sdl3-wasm
go mod edit -replace github.com/Zyko0/go-sdl3=../go-sdl3-wasm   # NEVER commit this
(cd ../go-sdl3-wasm && go run ./cmd/wasmsdl serve \
    -html "$PWD/../retinotopy-go/web/index.html" "$PWD/../retinotopy-go")
# → http://localhost:8080/
go mod edit -dropreplace github.com/Zyko0/go-sdl3                # before committing
```

Check `git diff go.mod` before committing — a leftover `replace` points the
published module at a path that exists only on your machine.

> **The bundle is about 135 MB** (112 MB gzipped), nearly all of it `main.wasm`:
> the 100 pattern and 1680 mask PNGs are embedded in the binary. That is a long
> first load for a participant, and it is the price of embedding the stimuli.
> Browser frame timing is also not the desktop build's — use a native binary for
> anything whose onsets matter, and treat the browser build as a demonstration.

---

## See also

- <https://osf.io/bw9ec/overview> (original Matlab version)
- <https://github.com/Goffaux-Lab/psychopy-retinotopy>
- <https://github.com/hiroshiban/Retinotopy>
- <https://github.com/egaffincahn/RetinotopicMapping>

---

## License and Authorship

Developed by [Christophe Pallier](https://github.com/chrplr) (2026), porting a previous Python/[Expyriment](http://expyriment.org) version with the help of Gemini.

Distributed under the [GNU General Public License v3](https://www.gnu.org/licenses/gpl-3.0.html).
