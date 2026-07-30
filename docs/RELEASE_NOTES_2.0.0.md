# 🌊 BilaKey PC 2.0.0 — CVNSS4.0 Core Edition

BilaKey PC 2.0.0 là bản portable dành cho Windows, lấy **CVNSS4.0 làm lõi mặc định**. Telex và VNI được giữ như lớp tương thích.

## ⬇️ Chọn đúng tệp

- **Windows x64 — khuyến nghị:** `BilaKey-PC-2.0.0-CVNSS-Core-x64.exe`
- **Windows ARM64 / Snapdragon:** `BilaKey-PC-2.0.0-CVNSS-Core-arm64.exe`
- **Windows x86 32-bit:** `BilaKey-PC-2.0.0-CVNSS-Core-x86.exe`
- **Gói đầy đủ ba kiến trúc:** `BilaKey-PC-2.0.0-Windows.zip`
- **Kiểm tra toàn vẹn:** `SHA256SUMS-2.0.0.txt`

## 🚀 Cách dùng

1. Tải đúng tệp `.exe` theo kiến trúc máy.
2. Chạy trực tiếp, không cần cài đặt và không cần quyền Administrator.
3. Dùng `Ctrl+Shift+Space` để bật/tắt BilaKey.
4. Dùng `Ctrl+Shift+1` để chọn **CVNSS4.0 Core**.
5. Dùng `Ctrl+Shift+0` để đổi ứng viên khi mã CVNSS có nhiều cách giải.

## 🧠 Nâng cấp chính

- CVNSS4.0 trở thành trung tâm của engine, giao diện, kiểm thử và quy trình phát hành.
- Resolver nhận biết âm đầu và chính tả, sửa đúng `qyl/qyz/qys/qyj/qyr → quỳ/quỷ/quỹ/quý/quỵ`.
- Giữ candidate graph 56 nhóm nhập nhằng và không ghi đè reverse-map âm thầm.
- Bảo vệ văn bản hỗn hợp như `OpenAI`, `GitHub`, URL và định danh mã nguồn.
- Hỗ trợ Windows x86, x64 và ARM64.
- Không telemetry, không tải rule qua mạng và không phụ thuộc dịch vụ đám mây.

## 🔐 Lưu ý an toàn

Bản 2.0.0 hiện là bản portable **chưa ký Authenticode**. Windows SmartScreen có thể cảnh báo ở lần chạy đầu. Hãy chỉ tải từ release chính thức này và đối chiếu SHA-256 với `SHA256SUMS-2.0.0.txt`.

## 🤝 Ghi nhận

- **Phát triển và duy trì:** Long Ngo.
- **Hỗ trợ nền tảng CVNSS4.0:** NNC Trần Tư Bình qua dự án CVNSS4.0.
- **Cộng đồng:** CVNSS4.0 và Bộ gõ BilaKey — https://www.facebook.com/groups/251479779599477

Mã nguồn mới của BilaKey PC được phát hành theo giấy phép **MIT**.
