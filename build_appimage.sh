#!/bin/bash

# build_appimage_minimal.sh
# Minimal manual AppImage creation for Retinotopy

set -e

APP_NAME="retinotopy"
VERSION=$(git describe --tags --always || echo "v0.0.0")
ARCH="x86_64"

echo "=== Building ${APP_NAME} ${VERSION} Minimal AppImage ==="

# 1. Build the Go binary (CGO_ENABLED=0 is preferred for portability if possible, 
# but we saw it segfaulted too. Let's try CGO_ENABLED=0 again now that we are doing minimal)
echo "--> Compiling Go binary..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${APP_NAME}_bin retinotopy.go

# 2. Prepare AppDir
echo "--> Preparing AppDir..."
rm -rf AppDir
mkdir -p AppDir/usr/bin
mkdir -p AppDir/usr/lib
mkdir -p AppDir/usr/share/retinotopy/assets

cp ${APP_NAME}_bin AppDir/usr/bin/${APP_NAME}
cp -r assets/* AppDir/usr/share/retinotopy/assets/

# 3. Include ONLY SDL3
SDL3_LIB=$(ldconfig -p | grep "libSDL3.so.0" | head -n 1 | awk '{print $4}')
if [ -z "$SDL3_LIB" ]; then
    echo "Error: libSDL3.so.0 not found on system."
    exit 1
fi
echo "--> Bundling ${SDL3_LIB}"
cp "${SDL3_LIB}" AppDir/usr/lib/

# 4. Create AppRun (the core entry point)
cat > AppDir/AppRun << 'EOF'
#!/bin/sh
HERE="$(dirname "$(readlink -f "${0}")")"
export LD_LIBRARY_PATH="${HERE}/usr/lib:${LD_LIBRARY_PATH}"
exec "${HERE}/usr/bin/retinotopy" "$@"
EOF
chmod +x AppDir/AppRun

# 5. Desktop file and Icon at root (needed by appimagetool)
cp assets/retinotopy.desktop AppDir/
cp assets/icons/icon.png AppDir/retinotopy.png
# Fix icon reference in desktop file to match root icon
sed -i 's/Icon=retinotopy/Icon=retinotopy/' AppDir/retinotopy.desktop

# 6. Download appimagetool
if [ ! -f appimagetool-x86_64.AppImage ]; then
    echo "--> Downloading appimagetool..."
    curl -L https://github.com/AppImage/AppImageKit/releases/download/continuous/appimagetool-x86_64.AppImage -o appimagetool-x86_64.AppImage
    chmod +x appimagetool-x86_64.AppImage
fi

# 7. Build AppImage
echo "--> Running appimagetool..."
ARCH=x86_64 ./appimagetool-x86_64.AppImage AppDir "${APP_NAME}-${VERSION}-linux-${ARCH}.AppImage"

echo "=== Build Complete! ==="
