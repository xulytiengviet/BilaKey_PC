# 🔐 Chính sách an toàn

## Báo cáo lỗ hổng

Không đăng công khai dữ liệu nhạy cảm, PoC khai thác hoặc thông tin cá nhân. Hãy mở một GitHub Security Advisory riêng tư khi chức năng đó khả dụng, hoặc liên hệ maintainer qua hồ sơ GitHub của dự án.

Thông tin nên có:

- phiên bản và SHA-256;
- Windows/kiến trúc;
- ứng dụng đích;
- bước tái hiện tối thiểu;
- tác động;
- log không chứa nội dung mật khẩu hoặc clipboard riêng tư.

## Mô hình đe dọa

BilaKey là phần mềm nhập liệu toàn cục, vì vậy các vùng nhạy cảm gồm keyboard hook, Unicode injection, clipboard fallback, macro và cấu hình khởi động.

Cam kết 2.0:

- không telemetry;
- không tải rule lúc chạy;
- clipboard fallback mặc định tắt;
- macro mặc định tắt;
- tạm dừng trên Win32 `ES_PASSWORD` mặc định bật;
- cấu hình ghi với quyền người dùng và cơ chế atomic;
- build phát hành phải có checksum, SBOM và Authenticode trước nhãn stable.

## Giới hạn đã biết

Control mật khẩu tùy biến của Chromium/Electron không phải lúc nào cũng mang `ES_PASSWORD`. Vì thế bản chưa qua password matrix không được xem là hoàn tất Gate D.
