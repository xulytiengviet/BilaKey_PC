#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

VERSION="2.0.0"
GO_BIN="${GO_BIN:-go}"
PYTHON_BIN="${PYTHON_BIN:-python3}"
GOFMT_BIN="${GOFMT_BIN:-$(command -v gofmt)}"
ORACLE="reference/cvnss4_0_converter.pro.v5_1.bilakey_core.js"
RULE_DIR="reference/cvnss"
export GOFLAGS="${GOFLAGS:-} -buildvcs=false"
export SOURCE_DATE_EPOCH="1783760400"

mkdir -p dist docs .build
rm -f "dist/BilaKey-PC-${VERSION}-CVNSS-Core-x64.exe" \
  "dist/BilaKey-PC-${VERSION}-CVNSS-Core-x64.exe.xz" \
  "dist/BilaKey-PC-${VERSION}-CVNSS-Core-x86.exe" \
  "dist/BilaKey-PC-${VERSION}-CVNSS-Core-x86.exe.xz" \
  "dist/BilaKey-PC-${VERSION}-CVNSS-Core-arm64.exe" \
  "dist/BilaKey-PC-${VERSION}-CVNSS-Core-arm64.exe.xz"

node - "$ORACLE" <<'NODE'
const path = require("path");
const converter = require(path.resolve(process.argv[2]));
const result = converter.selfTest();
console.log(JSON.stringify({version: converter.VERSION, audit: converter.audit(), selfTest: result}));
if (!result.ok) process.exit(1);
NODE

rm -rf .build/generated
mkdir -p .build/generated
"$PYTHON_BIN" tools/python/generate_cvnss_go.py "$RULE_DIR" --out-dir .build/generated
"$GOFMT_BIN" -w .build/generated/*.go
rm -rf .build/committed
mkdir -p .build/committed
cp internal/core/cvnss_generated_*.go .build/committed/
diff -ru .build/generated .build/committed
"$PYTHON_BIN" tools/python/audit_cvnss.py "$RULE_DIR" --out docs/cvnss_collision_audit.tsv

g++ -std=c++17 -O2 -Wall -Wextra -Werror tools/cpp/collision_stats.cpp -o .build/collision_stats
.build/collision_stats docs/cvnss_collision_audit.tsv

test -z "$("$GO_BIN" fmt ./...)"
"$GO_BIN" test ./... -count=1
CGO_ENABLED=1 "$GO_BIN" test -race ./... -count=1
"$GO_BIN" vet ./...
# Win32 callbacks deliver LPARAM as uintptr by API contract. Go 1.26's generic
# unsafeptr heuristic cannot prove that lifetime, so the Windows pass disables
# only that analyzer; host vet and race remain fully enabled.
GOOS=windows GOARCH=amd64 "$GO_BIN" vet -unsafeptr=false ./...
"$GO_BIN" test ./internal/core -run '^$' -bench . -benchmem -count=1

COMMON_LDFLAGS="-s -w -H windowsgui -buildid="
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 "$GO_BIN" build -trimpath -ldflags "$COMMON_LDFLAGS" \
  -o "dist/BilaKey-PC-${VERSION}-CVNSS-Core-x64.exe" ./cmd/bilakey
CGO_ENABLED=0 GOOS=windows GOARCH=386 "$GO_BIN" build -trimpath -ldflags "$COMMON_LDFLAGS" \
  -o "dist/BilaKey-PC-${VERSION}-CVNSS-Core-x86.exe" ./cmd/bilakey
CGO_ENABLED=0 GOOS=windows GOARCH=arm64 "$GO_BIN" build -trimpath -ldflags "$COMMON_LDFLAGS" \
  -o "dist/BilaKey-PC-${VERSION}-CVNSS-Core-arm64.exe" ./cmd/bilakey

xz -9e -f -k "dist/BilaKey-PC-${VERSION}-CVNSS-Core-x64.exe"
xz -9e -f -k "dist/BilaKey-PC-${VERSION}-CVNSS-Core-x86.exe"
xz -9e -f -k "dist/BilaKey-PC-${VERSION}-CVNSS-Core-arm64.exe"
xz -t dist/BilaKey-PC-${VERSION}-CVNSS-Core-*.xz

file dist/BilaKey-PC-${VERSION}-CVNSS-Core-*.exe
sha256sum dist/BilaKey-PC-${VERSION}-CVNSS-Core-*.exe \
  dist/BilaKey-PC-${VERSION}-CVNSS-Core-*.xz > SHA256SUMS.txt
wc -c dist/BilaKey-PC-${VERSION}-CVNSS-Core-*.exe \
  dist/BilaKey-PC-${VERSION}-CVNSS-Core-*.xz
