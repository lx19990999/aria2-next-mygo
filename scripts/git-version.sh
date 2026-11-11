#!/usr/bin/env bash
# Print the app version (no leading "v"): VERSION > GITHUB_REF_NAME > git tag > dev
set -euo pipefail

if [[ -n "${VERSION:-}" ]]; then
  echo "${VERSION#v}"
  exit 0
fi

if [[ "${GITHUB_REF_NAME:-}" == v* ]]; then
  echo "${GITHUB_REF_NAME#v}"
  exit 0
fi

tag=""
if tag="$(git describe --tags --exact-match HEAD 2>/dev/null)"; then
  :
elif tag="$(git describe --tags --abbrev=0 2>/dev/null)"; then
  :
else
  echo "dev"
  exit 0
fi

echo "${tag#v}"
