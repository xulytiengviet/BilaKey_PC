# 🔄 Chuyển từ BilaKey 1.3 sang 2.0

## Không đổi

- Cấu hình người dùng vẫn nằm trong thư mục `BilaKeyPC`.
- CVNSS4.0 vẫn là mặc định.
- Hotkey `Ctrl+Shift+Space` và `Ctrl+Shift+1/2/3/0` được giữ.
- Macro TSV và BilaCaps tiếp tục tương thích.

## Thay đổi hành vi

- `qyl/qyz/qys/qyj/qyr` được giải đúng thành `quỳ/quỷ/quỹ/quý/quỵ`.
- Candidate được xếp hạng sau khi ghép âm đầu; một số dạng như `vyệt`, `tyếng`, `ngyêng` không còn xuất hiện trong vòng chọn thông thường.
- Giao diện nhấn mạnh CVNSS4.0 Core; Telex/VNI được ghi rõ là adapter.
- Module Go chuyển về đường dẫn chính thức `github.com/xulytiengviet/BilaKey_PC`.

## Tệp thực thi cũ

Không chép đè 2.0 lên 1.3 trong lúc 1.3 đang chạy. Thoát từ menu tray, sao lưu `config.json` và macro TSV, sau đó chạy bản mới.
