# THIRD-PARTY NOTICES

## 1. Bamboo Core

BilaKey tham khảo **ý tưởng kiến trúc, tên quy tắc, mô hình xử lý Telex/VNI và cấu trúc kiểm tra chính tả** từ Bamboo Core. Phần runtime hiện tại được viết lại bằng Go, không liên kết runtime Rust.

Bản quyền được ghi trong nguồn tham chiếu:

- Copyright (C) 2018 Luong Thanh Lam <ltlam93@gmail.com>
- Copyright (C) 2024 nguien <nguyen10t2lhp@gmail.com>

Toàn văn giấy phép được sao chép tại `docs/BAMBOO_CORE_MIT_LICENSE.txt`.

## 2. Dự án CVNSS4.0

Bảng quy tắc và candidate graph trong `internal/core/cvnss_generated.go` được sinh từ oracle:

`reference/cvnss4_0_converter.pro.v5_1.bilakey_core.js`

BilaKey 2.0 bổ sung lớp resolver nhận biết âm đầu và các hồi quy dành cho hoạt động bộ gõ. Dữ liệu giữ 758 dòng gốc, 336 patch entries, 56 nhóm nhập nhằng, 56 policy canonical và 5 critical collision.

Ghi nhận hỗ trợ nền tảng:

- **NNC Trần Tư Bình**, thông qua dự án CVNSS4.0.
- **Cộng đồng CVNSS4.0 và Bộ gõ BilaKey**: <https://www.facebook.com/groups/251479779599477>

## 3. BilaKey PC

Phát triển và duy trì: **Long Ngo, 2026**.

Phần mã nguồn mới của BilaKey PC được phát hành theo MIT License, xem `LICENSE`.
