# Build BilaKey PC 2.5.0

## Build đầy đủ có audit

```bash
GO_BIN=go scripts/build_release.sh
```

Toolchain tối thiểu: Go 1.23+, Node.js 22+, Python 3.12+, g++ 13+ và xz. Bản RC này đã được xác minh với Go 1.23.2, Node.js 22.16, Python 3.13.5, g++ 14.2 và xz 5.8.1.

## Build nhanh x64

```bash
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 \
go build -buildvcs=false -trimpath \
  -ldflags='-s -w -H windowsgui -buildid=' \
  -o dist/BilaKey-PC-2.5.0-CVNSS-Core-x64.exe ./cmd/bilakey
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
