# 🚦 Cổng phát hành BilaKey 2.0.0

## Gate A — dữ liệu và resolver

- [x] Oracle tự kiểm thử.
- [x] 56/56 ambiguity policy có quyết định rõ.
- [x] Năm short-code collision được audit.
- [x] Hồi quy `qu + uy`.
- [x] Mixed-text safety.
- [x] Candidate details có lý do xếp hạng.

## Gate B — chất lượng mã nguồn

- [x] `gofmt` sạch.
- [x] `go test ./...`.
- [x] `go vet ./...`.
- [x] Race test cho package đa nền tảng.
- [x] Fuzz target cho decoder và inspector.
- [x] CLI audit đa nền tảng.
- [x] SBOM và third-party notices.

## Gate C — build

- [x] Windows x64 cross-build.
- [x] Windows x86 cross-build.
- [x] Windows ARM64 cross-build.
- [x] SHA-256 cho artifact phát hành.
- [ ] Reproducible build được xác minh trong workflow phát hành chính thức.

## Gate D — Windows thực địa

Bắt buộc trước khi gắn nhãn stable đại trà:

- [ ] Windows 10 x64.
- [ ] Windows 11 x64.
- [ ] Windows 11 ARM64.
- [ ] Notepad, Word, Excel, Outlook.
- [ ] Edge, Chrome, Firefox.
- [ ] VS Code/Electron.
- [ ] Teams/Zoom.
- [ ] Windows Terminal.
- [ ] Elevated app, RDP và nhiều mức DPI.
- [ ] Không mất/lặp phím trong stress test.
- [ ] Password fields: Win32, Chromium, Electron, Office.

## Gate E — phát hành đáng tin cậy

- [ ] Authenticode cho EXE/installer.
- [ ] Publisher/version/icon/manifest được nhúng.
- [ ] VirusTotal hoặc quy trình AV nội bộ được lưu hồ sơ.
- [ ] Pilot cộng đồng không còn lỗi Critical/High.

**Quy ước điểm:** nguồn và kiến trúc có thể đạt 9,0/10 sau Gate A–C. Chỉ công bố 9,2/10 stable sau Gate D–E.
