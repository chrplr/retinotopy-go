# Retinotopy Experiment (Go Implementation)

Welcome! This repository contains a Go implementation of the **HCP Retinotopic Mapping experiment** described in: 

> Benson, N. C., Jamison, K. W., Arcaro, M. J., Vu, A. T., Glasser, M. F., Coalson, T. S., Van Essen, D. C., Yacoub, E., Ugurbil, K., Winawer, J., & Kay, K. (2018). The Human Connectome Project 7 Tesla retinotopy dataset: Description and population receptive field analysis. *Journal of Vision*, 18(13), 23. https://doi.org/10.1167/18.13.23

The original stimulation program, available at <https://osf.io/bw9ec/overview>, was written in Matlab.

The present version, a complete rewrite, relies on [goxpyriment](https://github.com/chrplr/goxpyriment), itself based on [go-sdl3](https://github.com/Zyko0/go-sdl3).

**Warning: The timing of presentation has not been checked yet.**

*If you find issues, please report them on <https://github.com/chrplr/retinotopy-go>*


Christophe Pallier 05/03/2026


---

## 1. Using Pre-compiled Binaries (Fastest Start)

If you don't want to install the Go programming language, you can use the pre-compiled files in the `build/` folder.

### Step A: Install the SDL3 Library
Even with a pre-compiled file, your computer needs the **SDL3 shared library** to show graphics.

**Note:** You only need the *runtime* library (core SDL3). You do **NOT** need the development headers, nor do you need the extra `SDL3_image` or `SDL3_ttf` libraries, as these functions are handled internally by Go.

- **Ubuntu/Debian:**
  ```bash
  sudo apt install libsdl3-0
  ```
- **macOS (Homebrew):**
  ```bash
  brew install sdl3
  ```
- **Windows:**
  Download `SDL3.dll` from the [SDL GitHub Releases](https://github.com/libsdl-org/SDL/releases) (look for the `SDL3-3.x.x-win32-x64.zip` or similar) and place it in the `build/` folder alongside the `.exe` files.

### Step B: Choose the right file
Look in the `build/` folder for the file matching your computer:

- **Windows:** `retinotopy-windows-amd64.exe` (Standard PCs) or `retinotopy-windows-arm64.exe` (Surface Pro X, etc.)
- **macOS:** `retinotopy-darwin-arm64` (Apple Silicon M1/M2/M3) or `retinotopy-darwin-amd64` (Intel Macs)
- **Linux:** `retinotopy-linux-amd64` (Standard PCs) or `retinotopy-linux-arm64` (Raspberry Pi, etc.)

### Step C: Run the Experiment
1.  Open your **Terminal** or **PowerShell**.
2.  Navigate to the project folder.
3.  Run the command (replacing `<filename>` with your file):
    - **Windows:** `.\build\<filename>.exe -s 0 -r 1`
    - **macOS:**
      1.  Remove the security "quarantine" flag: `xattr -d com.apple.quarantine build/<filename>`
      2.  Make it executable: `chmod +x build/<filename>`
      3.  Run it: `./build/<filename> -s 0 -r 1`
    - **Linux:**
      1.  Make it executable: `chmod +x build/<filename>`
      2.  Run it: `./build/<filename> -s 0 -r 1`

---

## 2. Building from Source (For Developers)

If you want to modify the code or compile it yourself, follow these steps.

### Prerequisites
1.  **Install Go:** [go.dev/doc/install](https://go.dev/doc/install)
2.  **Install SDL3:** (See Step A above). 
    *Note: Because this project uses `purego`, you do **not** need C compilers or SDL3 development headers (`-dev` packages) to compile.*

### Getting Started
1.  **Clone the Repository:**
    ```bash
    git clone https://github.com/yourusername/retinotopy-go.git
    cd retinotopy-go
    ```
2.  **Download Dependencies:**
    ```bash
    go mod download
    ```

### Running/Building
- **To Run directly:** `go run retinotopy.go -s 0 -r 1`
- **To Build your own executable:** `go build -o my_retinotopy`

---

## 3. Command Line Options

Customize the experiment using these flags:

| Flag | Description | Default |
| :--- | :--- | :--- |
| `-s <id>` | **Subject ID** (used for data logging and stimuli order) | `0` |
| `-r <id>` | **Run ID** (1 to 6, see below) | `1` |
| `-d` | **Development Mode**: Runs in a 900x900 window instead of fullscreen. | `false` |
| `--scaling <f>` | **Scaling Factor**: Adjusts the size of stimuli (e.g., `1.5` for 150% size). | `1.0` |
| `-assets <path>`| **Assets Directory**: Path to the `assets/` folder. | `./assets` |

### Available Runs (`-r`)
1. `RETBAR1` / 2. `RETBAR2` (Swiping Bars)
3. `RETCCW` (Counter-Clockwise Wedge) / 4. `RETCW` (Clockwise Wedge)
5. `RETEXP` (Expanding Circles) / 6. `RETCON` (Contracting Circles)

---

## 4. Controls & Data

-   **ESC:** Exit and save data.
-   **Any Key/Mouse Click:** Press when the center fixation dot changes color.
-   **Data:** Results are saved as `.xpd` files in the `data/` directory with frame-by-frame timing and event logs.

---

## Troubleshooting

- **"SDL3 not found":** Re-check Step A. The library must be installed or the DLL must be in the folder.
- **"Assets not found":** Run the command from the root of the project directory.

---

## See also:

* <https://github.com/Goffaux-Lab/psychopy-retinotopy>
* <https://github.com/hiroshiban/Retinotopy>
* <https://github.com/egaffincahn/RetinotopicMapping>

---


Developed by [Christophe Pallier](https://github.com/chrplr) (2026). (Porting a previous Python using [Expyriment](http://expyriment.org) with the help of Gemini) 
Distributed under the GNU General Public License v3.
