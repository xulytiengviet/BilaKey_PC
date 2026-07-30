<div align="center">
  <img src="assets/brand/bilakey-logo.svg" alt="Logo BilaKey — chữ B trắng trên nền xanh đại dương" width="168" />

# 🌊 BilaKey PC 2.0.0

### Bộ gõ Windows lấy **CVNSS4.0 làm lõi**, không chỉ là một tùy chọn

[![Version](https://img.shields.io/badge/version-2.0.0-0756d8?style=for-the-badge)](VERSION)
[![License](https://img.shields.io/badge/license-MIT-00a86b?style=for-the-badge)](LICENSE)
[![Platform](https://img.shields.io/badge/Windows-x86%20%7C%20x64%20%7C%20ARM64-0078d4?style=for-the-badge&logo=windows)](BUILD.md)
[![Go](https://img.shields.io/badge/Go-1.23%2B-00add8?style=for-the-badge&logo=go)](go.mod)
[![CVNSS4.0](https://img.shields.io/badge/CVNSS4.0-Core-173ea5?style=for-the-badge)](docs/CVNSS_CORE_SPEC.md)
[![Privacy](https://img.shields.io/badge/telemetry-none-5427c7?style=for-the-badge)](SECURITY.md)

**Nhanh · Nhẹ · Unicode · Offline · Kiểm toán được · Mã nguồn mở**

</div>

---

## ✨ BilaKey 2.0.0 là gì?

BilaKey PC là bộ gõ tiếng Việt Unicode cho Windows, được thiết kế lại theo kiến trúc **CVNSS4.0 Core**:

```text
Phím người dùng
      │
      ▼
Composition State Machine
      │
      ▼
CVNSS4.0 Core Resolver
  ├─ bảng quy tắc được sinh từ oracle
  ├─ candidate graph 56 nhóm nhập nhằng
  ├─ resolver nhận biết âm đầu
  ├─ kiểm tra chính tả và văn bản hỗn hợp
  └─ audit / inspect / regression / fuzz
      │
      ├──────────────► Unicode tiếng Việt
      │
      ├─ Telex adapter ─┐
      └─ VNI adapter ───┴─ lớp tương thích, không thay thế lõi
```

Trong phiên bản 2.0.0, **CVNSS4.0 là chế độ mặc định, trung tâm của giao diện, cấu hình, kiểm thử và quy trình phát hành**. Telex và VNI vẫn có mặt để hỗ trợ chuyển tiếp cho người dùng phổ thông, nhưng được xác định rõ là các adapter tương thích.

## 🚀 Điểm mới nổi bật

| Năng lực | BilaKey 2.0.0 |
|---|---|
| 🧠 **CVNSS4.0 Core** | Lõi mặc định; Telex/VNI là adapter |
| 🔎 **Candidate graph** | Giữ đủ 56 nhóm nhập nhằng, không ghi đè reverse-map âm thầm |
| 🧭 **Resolver theo âm đầu** | Chọn ứng viên dựa trên cả onset + rime + chính tả |
| ✅ **Sửa lỗi họ `qu + uy`** | `qyl/qyz/qys/qyj/qyr → quỳ/quỷ/quỹ/quý/quỵ` |
| 🛡️ **Mixed-text safe** | Không tự ý làm hỏng `OpenAI`, `GitHub`, URL và định danh mã nguồn khi bật kiểm tra chính tả |
| ⚡ **Streaming composition** | Cập nhật từ sau từng phím, Backspace phục hồi và đổi ứng viên tại chỗ |
| 🧪 **Kiểm thử nhiều lớp** | Oracle JS · Python audit · C++ checker · Go unit/regression/race/fuzz |
| 🧰 **CLI kiểm toán** | Decode, inspect candidate graph và audit không cần giao diện Windows |
| 🔒 **Riêng tư** | Không telemetry, không tải quy tắc lúc chạy, không phụ thuộc dịch vụ đám mây |
| 🏗️ **Đa kiến trúc** | Cross-build Windows x86, x64 và ARM64 |

## ⌨️ Ví dụ CVNSS4.0

```text
Input CVNSS     Unicode
────────────    ─────────
toiy            tôi
iwy             yêu
vidf            việt
tizb            tiếng
qyl             quỳ
qyj             quý
ses             sẽ
```

Với mã nhiều nghĩa, BilaKey lưu toàn bộ candidate graph và cho phép đổi ứng viên bằng `Ctrl+Shift+0`. Các ứng viên sai cấu trúc chính tả sau một âm đầu được hạ hạng hoặc loại khỏi vòng chọn thông thường, nhưng vẫn còn trong dữ liệu audit.

## 📦 Cài đặt và chạy

### Bản phát hành Windows

Workflow CI tạo artifact Release Candidate từ chính cây nguồn này; bản phát hành chính thức sẽ được đính kèm trong mục **Releases**:

| Thiết bị | Tên artifact |
|---|---|
| Máy Intel/AMD hiện đại | `BilaKey-PC-2.0.0-CVNSS-Core-x64.exe` — khuyến nghị |
| Máy Windows 32-bit cũ | `BilaKey-PC-2.0.0-CVNSS-Core-x86.exe` |
| Windows on ARM / Snapdragon | `BilaKey-PC-2.0.0-CVNSS-Core-arm64.exe` |

Bảng checksum của build RC đã xác minh nằm tại [`SHA256SUMS.txt`](SHA256SUMS.txt). Repository không theo dõi binary trong cây nguồn; artifact được sinh bởi workflow/build script để giảm rủi ro nhầm phiên bản.

> Đây là **Release Candidate chưa ký Authenticode**, dù source audit và cross-build đã PASS. Chỉ dùng artifact từ Actions/Releases có SHA-256 trùng bảng công bố; nhãn stable chờ smoke-test Windows thật và ký số.

### Build từ mã nguồn

Yêu cầu: Go 1.23+, Node.js, Python 3, `g++` và `xz`.

```bash
git clone https://github.com/xulytiengviet/BilaKey_PC.git
cd BilaKey_PC
go test ./...
GO_BIN=go scripts/build_release.sh
```

Script phát hành sẽ kiểm tra oracle, sinh lại bảng Go, so sánh byte-for-byte, chạy audit độc lập, unit/race/vet/benchmark và cross-build ba kiến trúc.

## 🎛️ Phím tắt

| Phím | Tác vụ |
|---|---|
| `Ctrl+Shift+Space` | Bật/tắt BilaKey |
| `Ctrl+Shift+1` | Chọn **CVNSS4.0 Core** |
| `Ctrl+Shift+2` | Chọn Telex adapter |
| `Ctrl+Shift+3` | Chọn VNI adapter |
| `Ctrl+Shift+0` | Đổi ứng viên CVNSS đang nhập nhằng |
| `Shift` một lần | Viết hoa một từ |
| `Shift` hai lần | Bật/tắt BilaCaps |
| `Backspace` sau delimiter | Quay lại từ vừa commit để sửa |

`Ctrl+Tab` được trả hoàn toàn cho trình duyệt, IDE và trình soạn thảo.

## 🧰 CLI dành cho kiểm thử và nghiên cứu

```bash
# Decode một mã CVNSS
go run ./cmd/bilakey-cli qyl

# Chuyển văn bản hỗn hợp an toàn
go run ./cmd/bilakey-cli -text "qyl tizb vidf · OpenAI/GitHub"

# Xem candidate graph và lý do xếp hạng
go run ./cmd/bilakey-cli -inspect vidf

# Xuất audit lõi
go run ./cmd/bilakey-cli -audit
```

## 🏛️ Kiến trúc mã nguồn

```text
BilaKey_PC/
├── cmd/
│   ├── bilakey/          # Ứng dụng Win32 thường trú
│   └── bilakey-cli/      # Decode, inspect và audit đa nền tảng
├── internal/
│   ├── core/             # CVNSS4.0 Core + Telex/VNI adapters
│   ├── hotkey/           # Phân giải phím tắt không xung đột
│   ├── macro/            # Gõ tắt
│   ├── settings/         # Cấu hình atomic
│   ├── typingstate/      # BilaCaps, sentence caps, rollback
│   └── win/              # Win32 hook, tray, UI, Unicode sender
├── reference/            # Oracle CVNSS và bản legacy để đối chiếu
├── tools/                # Generator/audit Python, checker C++
├── scripts/              # Build phát hành tái lập
├── assets/brand/         # Logo chính thức
└── docs/                 # Đặc tả, kiến trúc, release gates
```

Xem chi tiết tại [Kiến trúc 2.0](docs/ARCHITECTURE.md), [Đặc tả CVNSS Core](docs/CVNSS_CORE_SPEC.md) và [Cổng phát hành](docs/RELEASE_GATES.md).

## 🔐 An toàn và quyền riêng tư

BilaKey không triển khai telemetry, quảng cáo, tài khoản người dùng hoặc truy cập mạng ở tầng ứng dụng. Quy tắc CVNSS được nhúng tĩnh; cấu hình được lưu trong thư mục người dùng và ghi theo cơ chế tệp tạm → đổi tên atomic.

Cơ chế hiện tại sử dụng Win32 low-level keyboard hook + Unicode `SendInput`, với clipboard chỉ là fallback tùy chọn. Lộ trình TSF được giữ như một gate kiến trúc sau 2.0; xem [SECURITY.md](SECURITY.md) và [RELEASE_GATES.md](docs/RELEASE_GATES.md).

## 📊 Trạng thái chất lượng

| Gate nguồn 2.0.0 | Trạng thái |
|---|---:|
| Oracle self-test, gồm hồi quy `qu + uy` | PASS |
| 758 dòng gốc · 336 patch · 56 policy | PASS |
| Không silent reverse overwrite | PASS |
| Go unit/regression | PASS |
| Mixed-text safety | PASS |
| Cross-platform CLI | PASS |
| Windows x86/x64/ARM64 cross-build | PASS trong CI/build script |
| Smoke-test 20 ứng dụng Windows thật | Cần pilot |
| Authenticode | Cần chứng thư phát hành |

Mã nguồn và kiến trúc đạt mục tiêu kỹ thuật khoảng **9,0/10**. Mốc **9,2/10 phát hành đại trà** chỉ được công bố sau khi hai gate cuối — Windows real-world matrix và Authenticode — hoàn tất; dự án không báo PASS giả.

## 🤝 Ghi nhận và cộng đồng

- **Phát triển và duy trì:** **Long Ngo**.
- **Hỗ trợ nền tảng CVNSS4.0:** **NNC Trần Tư Bình**, thông qua dự án CVNSS4.0.
- **Đóng góp thực địa:** cộng đồng sử dụng CVNSS4.0 và bộ gõ BilaKey.
- **Cộng đồng Facebook:** [CVNSS4.0 và Bộ gõ BilaKey](https://www.facebook.com/groups/251479779599477).

Tên tác giả, dự án và cộng đồng được ghi nhận nhằm thể hiện nguồn gốc tri thức và sự hỗ trợ; giấy phép của phần mã nguồn BilaKey mới vẫn là MIT. Xem [CREDITS.md](docs/CREDITS.md) và [THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md).

## 🌱 Đóng góp

Lỗi mapping, trường hợp nhập nhằng, tương thích ứng dụng và cải tiến tài liệu đều được hoan nghênh. Mọi thay đổi quy tắc cần kèm golden vector hoặc regression test; xem [CONTRIBUTING.md](CONTRIBUTING.md).

## 📄 Giấy phép

BilaKey PC được phát hành theo **MIT License**.

```text
Copyright (c) 2026 Long Ngo
```

Xem toàn văn tại [LICENSE](LICENSE).
