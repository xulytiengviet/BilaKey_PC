# Di chuyển lên BilaKey PC 2.5.0

## Thay đổi kiểu gõ

BilaKey 2.5.0 chỉ còn hai lựa chọn trong giao diện:

1. **CVNSS4.0** — lõi trung tâm của BilaKey.
2. **VNI/Telex** — một engine hợp nhất nhận cả quy ước VNI và Telex.

Cấu hình cũ có `Telex`, `VNI`, `Telex/VNI` hoặc `VNI/Telex` được chuẩn hóa tự động thành `VNI/Telex`. Không cần xóa `config.json`.

## Phím tắt

- `Ctrl+Shift+1`: CVNSS4.0.
- `Ctrl+Shift+2`: VNI/Telex.
- `Ctrl+Shift+3`: được trả lại cho ứng dụng đang dùng.

## Ví dụ

| Cách gõ | Chuỗi | Kết quả |
|---|---|---|
| Telex | `tieengs` | `tiếng` |
| VNI | `tieng61` | `tiếng` |
| Telex | `ddoongf` | `đồng` |
| VNI | `d9ong62` | `đồng` |
| Kết hợp | `vieet5` | `việt` |
| Kết hợp | `d9oongf` | `đồng` |
