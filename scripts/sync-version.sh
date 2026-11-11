#!/usr/bin/env bash
# Write scripts/git-version.sh into mygo.json "version". Prints the version.
set -euo pipefail
cd "$(dirname "$0")/.."

VER="$(./scripts/git-version.sh)"
python3 - <<PY
import json, pathlib
ver = "${VER}"
path = pathlib.Path("mygo.json")
data = json.loads(path.read_text(encoding="utf-8"))
data["version"] = ver
path.write_text(json.dumps(data, indent=2) + "\n", encoding="utf-8")
print(ver)
PY
