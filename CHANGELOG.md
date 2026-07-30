# 📝 Changelog

## 2.0.0 — CVNSS Core Edition · 2026-07-30

### CVNSS4.0 trở thành lõi

- Đưa CVNSS4.0 thành chế độ mặc định và trung tâm của UI, tài liệu, kiểm thử và release pipeline.
- Xác định Telex/VNI là adapter tương thích.
- Nâng oracle lên `5.1.0-bilakey-core`, giữ nguyên 758 base rows, 336 patch entries, 56 ambiguity policies và 5 critical collisions.
- Thêm resolver nhận biết âm đầu, candidate scoring và audit details.
- Sửa mất dấu thanh ở `qyl/qyz/qys/qyj/qyr`.
- Loại các dạng sai cấu trúc như `vyệt`, `tyếng`, `ngyêng` khỏi vòng candidate thông thường nhưng vẫn giữ trong audit trail.
- Sửa spell checker từng từ chặn nhầm toàn bộ từ bắt đầu bằng `qu`.

### Chất lượng và công cụ

- Thêm `AuditCVNSS`, `TransformText` và CLI decode/inspect/audit đa nền tảng.
- Thêm regression, mixed-text, property/fuzz smoke tests.
- Thêm CI: oracle JS, generator parity, Python audit, C++ checker, Go test/race/vet/fuzz và cross-build x86/x64/ARM64.
- Chuẩn hóa module về `github.com/xulytiengviet/BilaKey_PC`.
- Bổ sung architecture, core spec, release gates, migration, security và contributing guide.
- Đưa logo chính thức vào README và nhận diện ứng dụng.

### Ghi nhận

- Phát triển: Long Ngo.
- Hỗ trợ CVNSS4.0: NNC Trần Tư Bình và cộng đồng CVNSS4.0/BilaKey.

## 1.3.0 — Ocean Pro Lite · 2026-07-11

- Sinh bảng Go từ oracle `5.0.0-audit-safe`.
- Giữ candidate graph 56 nhóm, 56 canonical policies và sửa 5 short-code collisions.
- Hotkey có cấu trúc, escape-key Telex/VNI, password gate Win32, single instance, ARM64 và build tái lập.

## 1.2.0 — Ocean Lite

- UI Ocean, x86/x64/ARM64, cấu hình atomic và các cải tiến Win32.
