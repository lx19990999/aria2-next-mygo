#!/usr/bin/env bash
# Detached GUI launch for Linux / macOS.
set -euo pipefail
cd "$(dirname "$0")"

./build.sh
pkill -f '[a]ria2-next-mygo$' >/dev/null 2>&1 || true
nohup ./aria2-next-mygo >/dev/null 2>&1 &
disown || true
echo "aria2-next-mygo started"
