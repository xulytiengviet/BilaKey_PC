# 🏛️ Kiến trúc BilaKey PC 2.0.0

## 1. Nguyên tắc thiết kế

BilaKey 2.0.0 đảo chiều quan hệ của phiên bản cũ: CVNSS4.0 không còn là một nhánh trong `switch`, mà là **nguồn chân lý của hệ thống nhập liệu**. Telex và VNI được giữ trong vai trò adapter tương thích, chia sẻ composition state machine, Unicode sender, viết hoa và rollback.

Các nguyên tắc bắt buộc:

1. **Audit trước tiện lợi:** không ghi đè reverse-map âm thầm.
2. **Quyết định có ngữ cảnh:** candidate được xếp hạng theo âm đầu, vần và chính tả.
3. **Offline-first:** runtime không tải rule hoặc từ điển qua mạng.
4. **Fail-safe:** từ không hợp lệ được trả lại raw input khi bật spell restore.
5. **Một nguồn quy tắc:** oracle → generator → bảng Go → golden tests.
6. **Minh bạch phát hành:** không tuyên bố production khi chưa qua Windows matrix và ký số.

## 2. Luồng nhập liệu

```text
WH_KEYBOARD_LL
   │
   ├─ secure-field gate
   ├─ hotkey resolver
   └─ composition buffer
          │
          ▼
      core.Engine
          │
          ├─ CVNSS4.0 Core (mặc định)
          │    ├─ split initial/code-rime
          │    ├─ retrieve candidate graph
          │    ├─ onset-aware scoring
          │    ├─ orthographic filter
          │    └─ deterministic selection
          │
          ├─ Telex adapter
          └─ VNI adapter
          │
          ▼
 Unicode replacement → commit snapshot → rollback
```

## 3. Resolver CVNSS4.0

Phiên bản 1.3 chọn canonical ở cấp vần rồi mới ghép âm đầu. Cách này làm mất dấu ở họ `qyl/qyz/qys/qyj/qyr`. Phiên bản 2.0:

- giữ toàn bộ ứng viên thô;
- ghép từng ứng viên với âm đầu;
- kiểm tra cấu trúc chính tả;
- ưu tiên ứng viên bảo toàn thanh điệu của glide `qu`;
- xếp hạng ổn định và công khai lý do qua `InspectCVNSS`;
- di chuyển dấu khỏi `u` glide một cách phòng vệ trước khi bỏ ký tự đó.

## 4. Biên Windows

Ứng dụng hiện dùng Go + Win32 thuần:

- low-level hook tiếp nhận phím;
- `SendInput` phát Unicode;
- clipboard là fallback tùy chọn;
- message loop giữ callback hook ngắn, không làm I/O;
- cấu hình và thay đổi mode được xử lý ngoài đường nóng;
- phát hiện `ES_PASSWORD` cho control Win32 native.

TSF text service là hướng nâng cấp hậu 2.0, không được mô tả như tính năng đã hoàn thành.

## 5. Chuỗi cung ứng quy tắc

```text
reference/cvnss/ + wrapper tương thích JavaScript
             │
             ├─ Node selfTest()
             ├─ Python policy audit
             ├─ C++ independent checker
             └─ Python generator
                    │
                    ▼
           internal/core/cvnss_generated_*.go
                    │
                    ├─ Go regression/golden tests
                    └─ executable x86/x64/ARM64
```

Generator khóa các invariant: 758 base rows, 336 patch entries, 56 ambiguity groups, 56 canonical policies và 5 critical collisions.

## 6. Biên mở rộng

API lõi đáng tin cậy:

- `DecodeCVNSS`
- `InspectCVNSS`
- `CVNSSCandidates`
- `AuditCVNSS`
- `Engine.Transform`
- `Engine.TransformText`

Các ứng dụng khác có thể tái sử dụng lõi mà không cần Win32 shell.
