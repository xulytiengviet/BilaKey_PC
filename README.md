# BilaKey PC 1.3.0 — Ocean Pro Lite

Bộ gõ Windows Unicode ba chế độ **CVNSS4.0 · Telex · VNI**, viết bằng Go và Win32 thuần. Bản 1.3.0 nâng rule oracle CVNSS lên `5.0.0-audit-safe`, giữ đầy đủ 56 nhóm nhập nhằng để kiểm tra và đổi ứng viên, đồng thời sửa các xung đột phím và escape-key của bản 1.2.0.

## Chọn đúng tệp

| Máy Windows | Tệp chạy trực tiếp |
|---|---|
| Intel/AMD 64-bit | `dist/BilaKey-PC-1.3.0-Ocean-Pro-Lite-x64.exe` |
| Intel/AMD 32-bit | `dist/BilaKey-PC-1.3.0-Ocean-Pro-Lite-x86.exe` |
| Snapdragon/Windows ARM | `dist/BilaKey-PC-1.3.0-Ocean-Pro-Lite-arm64.exe` |

Không cần cài đặt. Tệp `.xz` là bản nén để lưu trữ/phân phối, cần giải nén trước khi chạy. Đối chiếu toàn vẹn bằng `SHA256SUMS.txt`.

> Bản phát hành này chưa có chữ ký số Authenticode. Windows SmartScreen có thể cảnh báo ở lần chạy đầu. Chỉ dùng tệp có SHA-256 trùng với bảng bàn giao.

## Phím tắt

| Phím | Tác vụ |
|---|---|
| `Ctrl+Shift+Space` | Bật/tắt tiếng Việt |
| `Ctrl+Shift+1` | Chọn CVNSS4.0 |
| `Ctrl+Shift+2` | Chọn Telex |
| `Ctrl+Shift+3` | Chọn VNI |
| `Ctrl+Shift+0` | Đổi ứng viên cho mã CVNSS nhập nhằng |
| `Tab` trong cửa sổ BilaKey | Chuyển tuần tự ba kiểu gõ |
| `Shift` một lần | Viết hoa một chữ |
| `Shift` hai lần | Bật/tắt khóa viết hoa BilaCaps |
| `Backspace` ngay sau dấu cách | Quay lại từ vừa gõ để sửa |

`Ctrl+Tab` không còn bị BilaKey chiếm; trình duyệt và trình soạn thảo được toàn quyền xử lý phím này.

## Điểm mới 1.3.0

- Sinh bảng Go trực tiếp từ oracle CVNSS4.0 `5.0.0-audit-safe`.
- 758 dòng gốc, 336 patch entries, 56 nhóm nhập nhằng và 56 policy canonical được khóa bằng invariant.
- Không ghi đè reverse-map âm thầm; candidate graph cho phép xem/đổi cách giải.
- Sửa năm collision mã ngắn `ed`, `es`, `od`, `of`, `os`.
- Telex/VNI có escape-key: `ass → as`, `aaa → aa`, `a11 → a1`, `a66 → a6`.
- Hotkey không xung đột Ctrl+Tab; thay bằng nhóm `Ctrl+Shift` có cấu trúc.
- Tạm dừng mặc định trong ô mật khẩu Win32 native.
- Chống chạy nhiều phiên bản, nhận DPI hệ thống, hỗ trợ x86/x64/ARM64.
- Cấu hình ghi atomically vào thư mục cấu hình người dùng; không cần Administrator.
- Không telemetry, không mã ứng dụng truy cập mạng, không tải rule lúc chạy.

## Kiểm chứng phát hành

| Gate | Kết quả |
|---|---:|
| Oracle JavaScript self-test | PASS · 23 kiểm tra |
| Python audit | PASS · 56/56 policy |
| C++ independent checker | PASS · 56 nhóm |
| Go unit/regression tests | PASS |
| Go race detector | PASS |
| Host vet | PASS |
| Windows vet | PASS, trừ heuristic `unsafeptr` tại biên Win32 callback |
| Cross-build x86/x64/ARM64 | PASS |
| Tái lập byte-for-byte x64 | PASS |

Chi tiết xem `REPORT_AUDIT_1.3.0.md`, `docs/BUILD_AUDIT_1.3.0.log` và `SBOM.cdx.json`.

## Build từ nguồn

Yêu cầu Go 1.23 trở lên, Node.js, Python 3, g++ và xz:

```bash
GO_BIN=go scripts/build_release.sh
```

Script kiểm tra oracle, so sánh mã sinh lại, chạy test/race/vet/benchmark, cross-build ba kiến trúc, nén XZ và tạo SHA-256. Runtime phát hành chỉ là một tệp EXE; Node/Python/C++ không đi kèm runtime.

## Giới hạn đã biết

- Chưa có Authenticode certificate.
- Chưa smoke-test GUI/hook/khay hệ thống trên máy Windows thật trong phiên build container này.
- Chế độ bỏ qua mật khẩu phát hiện control Win32 có `ES_PASSWORD`; ô mật khẩu tùy biến trong trình duyệt/Electron cần kiểm thử ứng dụng cụ thể.
- Chưa có chế độ bật/tắt theo từng ứng dụng và chưa hỗ trợ bảng mã cũ như TCVN3/VNI Windows.
- CVNSS4.0 là ánh xạ nhiều-một; policy canonical không thể phục hồi 100% mọi từ nếu thiếu ngữ cảnh. `Ctrl+Shift+0` là đường thoát minh bạch cho trường hợp nhập nhằng.

## Giấy phép

Mã nguồn BilaKey mới: MIT, xem `LICENSE`. Thành phần/ý tưởng bên thứ ba: xem `docs/THIRD_PARTY_NOTICES.md`.
