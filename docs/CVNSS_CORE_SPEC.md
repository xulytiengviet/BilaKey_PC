# 🧠 CVNSS4.0 Core Specification — BilaKey 2.0.0

## Mục tiêu

Đặc tả này định nghĩa hành vi runtime của CVNSS4.0 trong BilaKey. Oracle bảng quy tắc và resolver runtime là hai lớp riêng:

- **Oracle:** bảo tồn dữ liệu, candidate graph và policy lịch sử.
- **Resolver:** dùng âm đầu và chính tả để chọn cách hiển thị phù hợp trong một bộ gõ thời gian thực.

## Invariant dữ liệu

| Thuộc tính | Giá trị khóa |
|---|---:|
| Base rows | 758 |
| Patch entries | 336 |
| Code rimes | 810 |
| Ambiguity groups | 56 |
| Canonical policies | 56 |
| Critical short-code collisions | 5 |
| Silent reverse overwrite | 0 |

## Quy tắc lựa chọn

1. Tách code thành âm đầu và code-rime bằng longest match.
2. Lấy tất cả reverse candidates; không chỉ lấy một map value.
3. Ghép mỗi candidate thành từ Quốc ngữ.
4. Đánh giá cấu trúc chính tả sau âm đầu.
5. Ưu tiên policy oracle khi không có bằng chứng ngữ cảnh mạnh hơn.
6. Với `q`/`qu`, ưu tiên ứng viên giữ `u` làm glide không mang dấu.
7. Sắp xếp ổn định; thứ tự nguồn là tie-breaker cuối cùng.
8. Candidate UI thông thường chỉ xoay vòng kết quả chính tả hợp lệ và không trùng.
9. Audit details vẫn giữ các đường giải thô bị loại khỏi UI.

## Hồi quy bắt buộc

```text
qyl → quỳ
qyz → quỷ
qys → quỹ
qyj → quý
qyr → quỵ
vidf → việt, không ưu tiên vyệt
tizb → tiếng, không ưu tiên tyếng
wizy → nghiêng, không ưu tiên ngyêng
```

## Văn bản hỗn hợp

Khi `SpellCheck` và `AutoRestoreWrongKey` cùng bật, một token bị biến đổi nhưng không tạo thành âm tiết tiếng Việt hợp lệ phải được trả về raw input. Đây là hàng rào bảo vệ tên sản phẩm, URL, email và mã nguồn.

## Tương thích

Telex/VNI không được phép sửa candidate graph CVNSS hoặc thay đổi default mode. Chúng chỉ triển khai `Engine.Transform` qua adapter tương ứng.
