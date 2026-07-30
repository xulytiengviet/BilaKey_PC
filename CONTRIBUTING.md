# 🤝 Đóng góp cho BilaKey

## Quy trình

1. Tạo issue mô tả input, output hiện tại, output mong đợi và ứng dụng Windows liên quan.
2. Tạo nhánh từ `main`.
3. Mỗi thay đổi mapping/resolver phải kèm regression test hoặc golden vector.
4. Chạy:

```bash
go fmt ./...
go test ./... -count=1
go vet ./...
```

5. Không thêm telemetry, network runtime hoặc dependency nặng nếu chưa có thảo luận kiến trúc.
6. Pull request phải ghi rõ ảnh hưởng đến CVNSS core, Telex/VNI adapter và tương thích Windows.

## Thay đổi oracle

Không sửa trực tiếp `internal/core/cvnss_generated.go`. Sửa dữ liệu mô-đun trong `reference/cvnss/` hoặc generator, sau đó chạy:

```bash
python3 tools/python/generate_cvnss_go.py \
  reference/cvnss \
  --out-dir internal/core
gofmt -w internal/core/cvnss_generated_*.go
```

Mọi thay đổi invariant phải được review chuyên môn và cập nhật đặc tả.
