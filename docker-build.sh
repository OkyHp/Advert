#!/usr/bin/env bash
set -euo pipefail

cd /app

OUT_DIR="${OUT_DIR:-build}"
BIN_NAME="${BIN_NAME:-app.so}"
PKG="${PKG:-.}"
BUILD_DEBUG="${BUILD_DEBUG:-false}"

CGO_ENABLED="${CGO_ENABLED:-1}"
GOOS="${GOOS:-linux}"
GOARCH="${GOARCH:-amd64}"
LDFLAGS="${LDFLAGS:-"-s -w"}"

mkdir -p "${OUT_DIR}"

echo "========================================"
echo "Go:   $(go version)"
echo "GCC:  $(gcc --version | head -n 1)"
echo "ARCH: $(uname -m)"
echo "CGO:  enabled (CGO_ENABLED=1)"
echo "OUT:  ${OUT_DIR}/${BIN_NAME}"
echo "PKG:  ${PKG}"
echo "DBG:  ${BUILD_DEBUG}"
echo "========================================"

if [[ -f go.mod ]]; then
  go mod download
fi

if [[ "${BUILD_DEBUG}" == "true" || "${BUILD_DEBUG}" == "1" ]]; then
  echo "[BUILD] DEBUG (tags=debug)"
  CGO_ENABLED=$CGO_ENABLED GOOS=$GOOS GOARCH=$GOARCH go build -buildmode=c-shared -tags debug -o "${OUT_DIR}/${BIN_NAME}" ${PKG}
else
  echo "[BUILD] RELEASE (!debug)"
  CGO_ENABLED=$CGO_ENABLED GOOS=$GOOS GOARCH=$GOARCH go build -buildmode=c-shared -ldflags "$LDFLAGS" -o "${OUT_DIR}/${BIN_NAME}" ${PKG}
fi

echo ""
echo "[DONE] Output files:"
ls -lah "${OUT_DIR}"
