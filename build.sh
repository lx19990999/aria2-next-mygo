#!/usr/bin/env bash
# Cross-platform build for aria2-next-mygo (Windows / Linux / macOS).
set -euo pipefail
cd "$(dirname "$0")"

OUT="${OUT:-aria2-next-mygo}"
LDFLAGS="${LDFLAGS:-}"
VER="$(./scripts/git-version.sh)"
VERSION_X="-X main.appVersion=${VER}"

case "$(uname -s 2>/dev/null || echo unknown)" in
  MINGW*|MSYS*|CYGWIN*)
    OUT="${OUT%.exe}.exe"
    if [[ -z "${LDFLAGS}" ]]; then
      LDFLAGS="-H windowsgui"
    fi
    ;;
esac

if [[ "$(go env GOOS)" == "windows" && -z "${LDFLAGS}" ]]; then
  OUT="${OUT%.exe}.exe"
  LDFLAGS="-H windowsgui"
fi

if [[ -n "${LDFLAGS}" ]]; then
  LDFLAGS="${LDFLAGS} ${VERSION_X}"
else
  LDFLAGS="${VERSION_X}"
fi

echo "building ${OUT} (version=${VER} GOOS=$(go env GOOS) GOARCH=$(go env GOARCH) ldflags=${LDFLAGS})"
go build -ldflags "${LDFLAGS}" -o "${OUT}" .
echo "ok: ${OUT}"
