#!/bin/bash

# Build script for Retinotopy Experiment
# Cross-compiles for Linux, macOS, and Windows (amd64 and arm64)

APP_NAME="retinotopy"
SRC_FILE="retinotopy.go"
BUILD_DIR="build"

# Define target platforms: OS/Arch
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
    "windows/arm64"
)

echo "Starting cross-compilation..."
mkdir -p "$BUILD_DIR"

for PLATFORM in "${PLATFORMS[@]}"; do
    # Split OS and ARCH
    IFS="/" read -r GOOS GOARCH <<< "$PLATFORM"
    
    # Set output filename
    OUTPUT_NAME="${APP_NAME}-${GOOS}-${GOARCH}"
    if [ "$GOOS" == "windows" ]; then
        OUTPUT_NAME="${OUTPUT_NAME}.exe"
    fi

    echo "--> Building for ${GOOS}/${GOARCH}..."
    
    # Run build
    # CGO_ENABLED=0 is used because go-sdl3 uses purego for dynamic loading
    env CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH go build -o "${BUILD_DIR}/${OUTPUT_NAME}" "$SRC_FILE"

    if [ $? -ne 0 ]; then
        echo "FAILED: ${GOOS}/${GOARCH}"
    else
        echo "SUCCESS: ${BUILD_DIR}/${OUTPUT_NAME}"
    fi
done

echo ""
echo "Build complete. Binaries are in the '${BUILD_DIR}' directory."
echo "Note: The SDL3 shared library must be available on the target system to run these binaries."
