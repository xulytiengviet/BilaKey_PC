# THIRD-PARTY NOTICES

## 1. Bamboo Core

BilaKey PC 1.3.0 sử dụng lại **ý tưởng kiến trúc, tên quy tắc, mô hình xử lý Telex/VNI và cấu trúc kiểm tra chính tả** từ mã nguồn Bamboo Core do người dùng cung cấp trong một phiên làm việc trước.

Bản được cung cấp là một Rust port của Bamboo Core gốc, phát hành theo giấy phép MIT. BilaKey PC không liên kết runtime Rust; phần chạy thực tế được viết lại bằng Go thuần.

Bản quyền ghi trong giấy phép nguồn:

- Copyright (C) 2018 Luong Thanh Lam <ltlam93@gmail.com>
- Copyright (C) 2024 nguien <nguyen10t2lhp@gmail.com>

Toàn văn giấy phép được sao chép tại `docs/BAMBOO_CORE_MIT_LICENSE.txt`.

## 2. CVNSS4.0 converter 5.0.0-audit-safe

Bộ ánh xạ CVNSS4.0 trong `internal/core/cvnss_generated.go` được sinh từ tệp chuyển đổi do người dùng cung cấp: `reference/cvnss4_0_converter.pro.v5_0.audit_safe.js`.

Runtime BilaKey PC không cần Node.js hoặc JavaScript. Mọi ánh xạ cần thiết đã được chuyển thành Go map tĩnh. Candidate graph giữ 56 nhóm nhập nhằng; 56 policy canonical và 5 critical collision được kiểm tra công khai trong báo cáo audit.

## 3. BilaKey PC

Phát triển: **Long Ngo, 2026**.

Phần mã nguồn mới của BilaKey PC được phát hành theo MIT License, xem `LICENSE`.
