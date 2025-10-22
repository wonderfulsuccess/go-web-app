#!/usr/bin/env bash
set -euo pipefail

command -v zip >/dev/null 2>&1 || { echo "zip command not found" >&2; exit 1; }

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}" )" && pwd)"
FRONT_DIR="$ROOT_DIR/front"
BACK_DIR="$ROOT_DIR/back"
DIST_DIR="$BACK_DIR/webserver/dist"
RELEASE_DIR="$ROOT_DIR/release"
STAMP="$(date +"%Y%m%d-%H%M%S")"

mkdir -p "$RELEASE_DIR"
rm -rf "$RELEASE_DIR"/*

echo "==> Building front-end assets"
(
  cd "$FRONT_DIR"
  npm install
  npm run build
)

if [ ! -d "$DIST_DIR" ]; then
  echo "Front-end build output not found at $DIST_DIR" >&2
  exit 1
fi

TMP_ROOT="$(mktemp -d)"
cleanup() {
  rm -rf "$TMP_ROOT"
}
trap cleanup EXIT

declare -a TARGETS=(
  "darwin arm64 macos-arm64"
  "darwin amd64 macos-amd64"
  "windows amd64 windows-amd64"
  "linux amd64 linux-amd64"
)

for entry in "${TARGETS[@]}"; do
  read -r GOOS GOARCH LABEL <<<"${entry}"
  echo "==> Building backend for ${LABEL}"

  PACKAGE_DIR="$TMP_ROOT/${LABEL}"
  mkdir -p "$PACKAGE_DIR"

  BINARY_NAME="go-web-app"
  if [ "$GOOS" = "windows" ]; then
    BINARY_NAME+=".exe"
  fi

  build_env=(
    "GOOS=$GOOS"
    "GOARCH=$GOARCH"
    "CGO_ENABLED=1"
  )
  if [ "$GOOS" = "windows" ]; then
    build_env+=("CC=x86_64-w64-mingw32-gcc")
  elif [ "$GOOS" = "linux" ]; then
    build_env+=("CC=x86_64-linux-musl-gcc")
  fi

  (
    cd "$BACK_DIR"
    env "${build_env[@]}" go build -trimpath -ldflags="-s -w" -o "$PACKAGE_DIR/$BINARY_NAME" ./app
  )

  mkdir -p "$PACKAGE_DIR/webserver"
  cp -a "$DIST_DIR" "$PACKAGE_DIR/webserver/"

  ZIP_NAME="${LABEL}-${STAMP}.zip"
  (
    cd "$PACKAGE_DIR"
    zip -rq "$RELEASE_DIR/$ZIP_NAME" .
  )
  echo "   -> Created $RELEASE_DIR/$ZIP_NAME"
done

echo "==> Release artifacts are in $RELEASE_DIR"
