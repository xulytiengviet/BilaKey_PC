# Build BilaKey PC 2.0.0

## Build đầy đủ có audit

```bash
GO_BIN=go scripts/build_release.sh
```

Toolchain đã xác minh cho bản bàn giao: Go 1.26.4, Node.js 24, Python 3, GCC/g++ 13, xz.

## Build nhanh x64

```bash
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 \
go build -buildvcs=false -trimpath \
  -ldflags='-s -w -H windowsgui -buildid=' \
  -o dist/BilaKey-PC-2.0.0-CVNSS-Core-x64.exe ./cmd/bilakey
```

Không dùng `-gcflags=all=-B`: cờ đó loại bỏ bounds-check và không phù hợp với bản phát hành ưu tiên an toàn.

## Cấu trúc gate

1. Oracle JS self-test.
2. Sinh lại Go table và `cmp` với source đã commit.
3. Python candidate/policy audit.
4. C++ checker độc lập.
5. Go fmt/test/race/vet/benchmark.
6. Cross-build x86/x64/ARM64.
7. XZ và SHA-256.
