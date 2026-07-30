# 🧰 Công cụ kiểm toán và sinh mã

- `python/generate_cvnss_go.py`: sinh bảng Go từ oracle và khóa invariant.
- `python/audit_cvnss.py`: audit candidate graph/policy độc lập.
- `cpp/collision_stats.cpp`: kiểm tra TSV bằng implementation C++ độc lập.
- `legacy/v1.3.0/`: công cụ tạo báo cáo và thử nghiệm lịch sử, không nằm trong release pipeline 2.0.

Runtime BilaKey không phụ thuộc Python, Node.js hoặc C++; các công cụ này chỉ chạy khi build/audit.
