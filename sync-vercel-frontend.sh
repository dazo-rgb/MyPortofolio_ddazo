#!/usr/bin/env sh
set -eu
ROOT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
SOURCE="$ROOT_DIR/frontend"
TARGET="$ROOT_DIR/backend/frontend"

if [ ! -f "$SOURCE/index.html" ]; then
  echo "Error: frontend/index.html not found." >&2
  exit 1
fi

mkdir -p "$TARGET"
cp -R "$SOURCE/." "$TARGET/"
echo "Vercel frontend copy synchronized: backend/frontend/"
