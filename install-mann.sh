#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BIN_NAME="mann"
INSTALL_DIR="/opt/4rji/bin"

echo "==> Compiling ${BIN_NAME}..."
cd "$SCRIPT_DIR"
go build -o "$BIN_NAME" .

echo "==> Installing to ${INSTALL_DIR} (requires sudo)..."
sudo mkdir -p "$INSTALL_DIR"
sudo install -m 755 "$BIN_NAME" "$INSTALL_DIR/$BIN_NAME"

echo "==> Done: ${INSTALL_DIR}/${BIN_NAME}"
