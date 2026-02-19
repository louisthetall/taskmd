#!/usr/bin/env bash
#
# Build an MCPB (MCP Bundle) for taskmd.
#
# Usage: ./scripts/build-mcpb.sh <binary-path> <version> <goos> <goarch> <output-dir>
#
# Example:
#   ./scripts/build-mcpb.sh ./dist/taskmd-darwin-arm64 0.1.0 darwin arm64 ./dist

set -euo pipefail

if [[ $# -ne 5 ]]; then
    echo "Usage: $0 <binary-path> <version> <goos> <goarch> <output-dir>" >&2
    exit 1
fi

BINARY_PATH="$1"
VERSION="$2"
GOOS="$3"
GOARCH="$4"
OUTPUT_DIR="$5"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
MANIFEST_TEMPLATE="$PROJECT_ROOT/apps/cli/mcpb/manifest.template.json"

if [[ ! -f "$BINARY_PATH" ]]; then
    echo "Error: binary not found: $BINARY_PATH" >&2
    exit 1
fi

if [[ ! -f "$MANIFEST_TEMPLATE" ]]; then
    echo "Error: manifest template not found: $MANIFEST_TEMPLATE" >&2
    exit 1
fi

# Determine binary name
BINARY_NAME="taskmd"
if [[ "$GOOS" == "windows" ]]; then
    BINARY_NAME="taskmd.exe"
fi

BUNDLE_NAME="taskmd-v${VERSION}-${GOOS}-${GOARCH}.mcpb"

# Create temp staging directory
STAGING_DIR=$(mktemp -d)
trap 'rm -rf "$STAGING_DIR"' EXIT

# Copy binary into server/
mkdir -p "$STAGING_DIR/server"
cp "$BINARY_PATH" "$STAGING_DIR/server/$BINARY_NAME"
chmod +x "$STAGING_DIR/server/$BINARY_NAME"

# Generate manifest.json from template
sed "s/VERSION_PLACEHOLDER/$VERSION/g" "$MANIFEST_TEMPLATE" > "$STAGING_DIR/manifest.json"

# For windows, update entry_point and mcp_config.command to use .exe
if [[ "$GOOS" == "windows" ]]; then
    # Use a temp file for portability
    TMP_MANIFEST=$(mktemp)
    sed 's|"entry_point": "server/taskmd"|"entry_point": "server/taskmd.exe"|g' "$STAGING_DIR/manifest.json" \
        | sed 's|"command": "${__dirname}/server/taskmd"|"command": "${__dirname}/server/taskmd.exe"|g' \
        > "$TMP_MANIFEST"
    mv "$TMP_MANIFEST" "$STAGING_DIR/manifest.json"
fi

# Update platform compatibility based on target OS
if [[ "$GOOS" == "darwin" ]]; then
    PLATFORM="darwin"
elif [[ "$GOOS" == "linux" ]]; then
    PLATFORM="linux"
elif [[ "$GOOS" == "windows" ]]; then
    PLATFORM="win32"
fi
TMP_MANIFEST=$(mktemp)
sed "s/\"darwin\"/\"$PLATFORM\"/g" "$STAGING_DIR/manifest.json" > "$TMP_MANIFEST"
mv "$TMP_MANIFEST" "$STAGING_DIR/manifest.json"

# Create the .mcpb (ZIP archive)
mkdir -p "$OUTPUT_DIR"
(cd "$STAGING_DIR" && zip -r "$BUNDLE_NAME" manifest.json server/)
mv "$STAGING_DIR/$BUNDLE_NAME" "$OUTPUT_DIR/$BUNDLE_NAME"

echo "Created $OUTPUT_DIR/$BUNDLE_NAME"
