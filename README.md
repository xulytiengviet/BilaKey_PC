<div align="center">
  <img src="assets/brand/bilakey-logo.svg" alt="Logo BilaKey — chữ B trắng trên nền xanh đại dương" width="168" />

# 🌊 BilaKey PC 2.5.0

### Bộ gõ Windows chỉ còn **2 kiểu gõ: CVNSS4.0 và VNI/Telex**

<a href="https://github.com/xulytiengviet/BilaKey_PC/releases/download/2.5.0/BilaKey-PC-2.5.0-CVNSS-Core-x64.exe">
  <img src="https://img.shields.io/badge/T%E1%BA%A2I_NGAY-Windows_x64-0756d8?style=for-the-badge&logo=windows&logoColor=white" alt="Tải ngay BilaKey PC 2.5.0 cho Windows x64" />
</a>

**Bấm nút trên để tải `.exe` và chạy ngay — portable, không cần cài đặt.**

[![ARM64](https://img.shields.io/badge/T%E1%BA%A3i-ARM64%20%7C%20Snapdragon-173ea5?style=flat-square&logo=windows)](https://github.com/xulytiengviet/BilaKey_PC/releases/download/2.5.0/BilaKey-PC-2.5.0-CVNSS-Core-arm64.exe)
[![x86](https://img.shields.io/badge/T%E1%BA%A3i-Windows%2032--bit-3158a8?style=flat-square&logo=windows)](https://github.com/xulytiengviet/BilaKey_PC/releases/download/2.5.0/BilaKey-PC-2.5.0-CVNSS-Core-x86.exe)
[![ZIP](https://img.shields.io/badge/T%E1%BA%A3i-G%C3%B3i%20Windows%20%C4%91%E1%BA%A7y%20%C4%91%E1%BB%A7-5427c7?style=flat-square&logo=github)](https://github.com/xulytiengviet/BilaKey_PC/releases/download/2.5.0/BilaKey-PC-2.5.0-Windows.zip)
[![SHA](https://img.shields.io/badge/Ki%E1%BB%83m_tra-SHA--256-6b7280?style=flat-square&logo=github)](https://github.com/xulytiengviet/BilaKey_PC/releases/download/2.5.0/SHA256SUMS-2.5.0.txt)
[![Release](https://img.shields.io/badge/Xem-Release%202.5.0-00a86b?style=flat-square&logo=github)](https://github.com/xulytiengviet/BilaKey_PC/releases/tag/2.5.0)

[![Version](https://img.shields.io/badge/version-2.5.0-0756d8?style=for-the-badge)](VERSION)
[![License](https://img.shields.io/badge/license-MIT-00a86b?style=for-the-badge)](LICENSE)
[![Platform](https://img.shields.io/badge/Windows-x86%20%7C%20x64%20%7C%20ARM64-0078d4?style=for-the-badge&logo=windows)](BUILD.md)
[![CVNSS4.0](https://img.shields.io/badge/CVNSS4.0-Core-173ea5?style=for-the-badge)](docs/CVNSS_CORE_SPEC.md)
[![Privacy](https://img.shields.io/badge/telemetry-none-5427c7?style=for-the-badge)](SECURITY.md)

**Nhanh · Nhẹ · Unicode · Offline · Kiểm toán được · Mã nguồn mở MIT**

</div>

---

## ✨ Thay đổi quan trọng trong 2.5.0

BilaKey không còn bắt người dùng chọn riêng Telex hoặc VNI. Giao diện và lõi nhập liệu chỉ còn hai kiểu:

| Kiểu gõ | Vai trò |
|---|---|
| 🧠 **CVNSS4.0** | Lõi trung tâm, chế độ mặc định của BilaKey |
| 🔁 **VNI/Telex** | Engine hợp nhất tự nhận cả phím chữ Telex và phím số VNI |

Khi chọn **VNI/Telex**, người dùng có thể:

- gõ hoàn toàn theo Telex;
- gõ hoàn toàn theo VNI;
- đổi từ Telex sang VNI giữa các từ mà không chuyển chế độ;
- kết hợp phím tạo chữ của một kiểu với phím đặt dấu của kiểu còn lại trong cùng một từ.

```text
Telex       tieengs   → tiếng
VNI         tieng61   → tiếng
Telex       ddoongf   → đồng
VNI         d9ong62   → đồng
Kết hợp     vieet5    → việt
Kết hợp     d9oongf   → đồng
```

Cấu hình cũ `Telex`, `VNI` hoặc `Telex/VNI` được chuyển tự động sang `VNI/Telex`; người dùng không cần xóa cấu hình.

## 🏗️ Kiến trúc nhập liệu 2.5

```text
Phím người dùng
      │
      ▼
Composition State Machine
      │
      ├───────────── CVNSS4.0 Core
      │                ├─ candidate graph 56 nhóm
      │                ├─ resolver nhận biết âm đầu
      │                └─ audit / regression / fuzz
      │
      └───────────── VNI/Telex Unified Engine
                       ├─ phím dấu Telex: s f r x j z
                       ├─ phím tạo chữ Telex: aa aw ee oo ow uw dd
                       ├─ phím dấu VNI: 0 1 2 3 4 5
                       ├─ phím tạo chữ VNI: 6 7 8 9
                       └─ một đầu ra Unicode tiếng Việt
```

## 📦 Tải trực tiếp đã xác minh

| Thiết bị | Liên kết tải trực tiếp |
|---|---|
| **Windows x64 — khuyến nghị** | [BilaKey-PC-2.5.0-CVNSS-Core-x64.exe](https://github.com/xulytiengviet/BilaKey_PC/releases/download/2.5.0/BilaKey-PC-2.5.0-CVNSS-Core-x64.exe) |
| **Windows ARM64 / Snapdragon** | [BilaKey-PC-2.5.0-CVNSS-Core-arm64.exe](https://github.com/xulytiengviet/BilaKey_PC/releases/download/2.5.0/BilaKey-PC-2.5.0-CVNSS-Core-arm64.exe) |
| **Windows x86 32-bit** | [BilaKey-PC-2.5.0-CVNSS-Core-x86.exe](https://github.com/xulytiengviet/BilaKey_PC/releases/download/2.5.0/BilaKey-PC-2.5.0-CVNSS-Core-x86.exe) |
| **Gói đầy đủ Windows ZIP** | [BilaKey-PC-2.5.0-Windows.zip](https://github.com/xulytiengviet/BilaKey_PC/releases/download/2.5.0/BilaKey-PC-2.5.0-Windows.zip) |
| **Bảng kiểm tra SHA-256** | [SHA256SUMS-2.5.0.txt](https://github.com/xulytiengviet/BilaKey_PC/releases/download/2.5.0/SHA256SUMS-2.5.0.txt) |
| **Trang phát hành** | [Release BilaKey PC 2.5.0](https://github.com/xulytiengviet/BilaKey_PC/releases/tag/2.5.0) |

Quy trình phát hành chỉ hoàn tất khi workflow xác minh đủ **ba tệp EXE, một gói ZIP và bảng checksum SHA-256**.

> Các tệp portable hiện chưa ký Authenticode. Windows SmartScreen có thể cảnh báo ở lần chạy đầu; chỉ tải từ repository/release chính thức và đối chiếu SHA-256.

## 🎛️ Phím tắt

| Phím | Tác vụ |
|---|---|
| `Ctrl+Shift+Space` | Bật/tắt BilaKey |
| `Ctrl+Shift+1` | Chọn **CVNSS4.0** |
| `Ctrl+Shift+2` | Chọn **VNI/Telex** |
| `Ctrl+Shift+0` | Đổi ứng viên CVNSS đang nhập nhằng |
| `Shift` một lần | Viết hoa một từ |
| `Shift` hai lần | Bật/tắt BilaCaps |
| `Backspace` sau delimiter | Quay lại từ vừa commit để sửa |

`Ctrl+Shift+3`, `Ctrl+Tab`, Alt và phím Windows được trả cho ứng dụng đang dùng.

## 🧰 CLI

```bash
# CVNSS4.0
go run ./cmd/bilakey-cli -method cvnss qyl

# VNI hoặc Telex đều dùng cùng một method
go run ./cmd/bilakey-cli -method vni-telex tieengs
go run ./cmd/bilakey-cli -method vni-telex tieng61

# Audit CVNSS core
go run ./cmd/bilakey-cli -audit
```

## 🔨 Build từ mã nguồn

```bash
git clone https://github.com/xulytiengviet/BilaKey_PC.git
cd BilaKey_PC
go test ./...
GO_BIN=go scripts/build_release.sh
```

Yêu cầu: Go 1.23+, Node.js 22+, Python 3.12+, `g++` và `xz`.

## 🔐 An toàn và quyền riêng tư

BilaKey không có telemetry, quảng cáo, tài khoản người dùng hoặc network runtime. Quy tắc được nhúng tĩnh; cấu hình lưu trong thư mục người dùng theo cơ chế ghi tệp tạm rồi đổi tên atomic.

## 📊 Trạng thái kiểm chứng

| Gate | Trạng thái |
|---|---:|
| CVNSS oracle, candidate graph và policy | PASS |
| VNI/Telex: vectors Telex thuần | PASS |
| VNI/Telex: vectors VNI thuần | PASS |
| VNI/Telex: vectors kết hợp | PASS |
| Di chuyển cấu hình Telex/VNI cũ | PASS |
| Go unit/regression/race/vet/fuzz | PASS trong CI |
| Windows x86/x64/ARM64 cross-build | PASS trong CI |
| Release có 3 EXE + ZIP + SHA-256 | Bắt buộc trước khi workflow kết thúc |
| Authenticode | Chờ chứng thư phát hành |

## 🤝 Ghi nhận và cộng đồng

- **Phát triển và duy trì:** **Long Ngo**.
- **Hỗ trợ nền tảng CVNSS4.0:** **NNC Trần Tư Bình**, thông qua dự án CVNSS4.0.
- **Cộng đồng:** [CVNSS4.0 và Bộ gõ BilaKey](https://www.facebook.com/groups/251479779599477).

## 📄 Giấy phép

BilaKey PC được phát hành theo **MIT License**.

```text
Copyright (c) 2026 Long Ngo
```
