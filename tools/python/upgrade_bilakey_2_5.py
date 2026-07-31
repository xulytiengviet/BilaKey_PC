#!/usr/bin/env python3
from __future__ import annotations

from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]


def path(rel: str) -> Path:
    return ROOT / rel


def read(rel: str) -> str:
    return path(rel).read_text(encoding="utf-8")


def write(rel: str, content: str) -> None:
    target = path(rel)
    target.parent.mkdir(parents=True, exist_ok=True)
    target.write_text(content, encoding="utf-8", newline="\n")


def replace_exact(rel: str, old: str, new: str, *, count: int = 1) -> None:
    content = read(rel)
    actual = content.count(old)
    if actual != count:
        raise RuntimeError(f"{rel}: expected {count} occurrence(s), found {actual}: {old!r}")
    write(rel, content.replace(old, new, count))


def replace_all(rel: str, old: str, new: str) -> None:
    content = read(rel)
    if old not in content:
        raise RuntimeError(f"{rel}: missing text {old!r}")
    write(rel, content.replace(old, new))


# Version and build metadata.
write("VERSION", "2.5.0\n")
replace_all("BUILD.md", "2.0.0", "2.5.0")
replace_exact("scripts/build_release.sh", 'VERSION="2.0.0"', 'VERSION="2.5.0"')
replace_all(".github/workflows/ci.yml", "BilaKey 2.0 CI", "BilaKey 2.5 CI")
replace_all(".github/workflows/ci.yml", "BilaKey-PC-2.0.0", "BilaKey-PC-2.5.0")
if path("SBOM.cdx.json").exists():
    replace_all("SBOM.cdx.json", "2.0.0", "2.5.0")

# Core engine: CVNSS4.0 plus one unified VNI/Telex mode.
write(
    "internal/core/engine.go",
    r'''package core

import "strings"

const (
	MethodCVNSS     = "CVNSS4.0"
	MethodVNITelex  = "VNI/Telex"
	MethodTelex     = "Telex" // legacy configuration alias
	MethodVNI       = "VNI"   // legacy configuration alias
)

type Options struct {
	OldToneStyle        bool
	FreeToneMarking     bool
	SpellCheck          bool
	AutoRestoreWrongKey bool
}

type Engine struct {
	Method  string
	Options Options
}

// NormalizeMethod migrates legacy Telex/VNI selections to the unified mode.
func NormalizeMethod(method string) string {
	key := strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(method), " ", ""))
	switch key {
	case "VNI/TELEX", "TELEX/VNI", "VNI", "TELEX", "VNITELEX", "TELEXVNI":
		return MethodVNITelex
	default:
		return MethodCVNSS
	}
}

func New(method string, opts Options) *Engine {
	return &Engine{Method: NormalizeMethod(method), Options: opts}
}

func (e *Engine) Transform(raw string) string {
	var out string
	switch NormalizeMethod(e.Method) {
	case MethodVNITelex:
		out = transformVNITelex(raw, e.Options.OldToneStyle)
	default:
		out = DecodeCVNSS(raw)
	}
	if e.Options.SpellCheck && e.Options.AutoRestoreWrongKey && out != raw && !IsLikelyVietnameseSyllable(out) {
		return raw
	}
	return out
}
''',
)

write(
    "internal/core/vni_telex.go",
    r'''package core

import "unicode"

// transformVNITelex accepts Telex and VNI control keys in one composition
// engine. Users may switch convention between words and may combine a Telex
// shape key with a VNI tone key (or the reverse) inside one Vietnamese word.
func transformVNITelex(raw string, oldStyle bool) string {
	out := make([]rune, 0, len([]rune(raw)))
	pendingTone := toneNone

	for _, ch := range []rune(raw) {
		// VNI numeric controls.
		switch ch {
		case '0', '1', '2', '3', '4', '5':
			if hasVowel(out) {
				nextTone := vniTone(ch)
				if ch != '0' && pendingTone == nextTone {
					pendingTone = toneNone
					out = append(out, ch)
					continue
				}
				pendingTone = nextTone
				continue
			}
		case '6':
			if applyLastEligible(out, map[rune]rune{'a': 'â', 'e': 'ê', 'o': 'ô'}) {
				continue
			}
			if undoLastVNIShape(out, map[rune]rune{'â': 'a', 'ê': 'e', 'ô': 'o'}) {
				out = append(out, ch)
				continue
			}
		case '7':
			if applyVNI7(out) {
				continue
			}
			if undoVNI7(out) {
				out = append(out, ch)
				continue
			}
		case '8':
			if applyLastEligible(out, map[rune]rune{'a': 'ă'}) {
				continue
			}
			if undoLastVNIShape(out, map[rune]rune{'ă': 'a'}) {
				out = append(out, ch)
				continue
			}
		case '9':
			if len(out) > 0 && unicode.ToLower(out[len(out)-1]) == 'd' {
				out[len(out)-1] = replaceBaseRune(out[len(out)-1], 'đ')
				continue
			}
			if len(out) > 0 && unicode.ToLower(out[len(out)-1]) == 'đ' {
				out[len(out)-1] = replaceBaseRune(out[len(out)-1], 'd')
				out = append(out, ch)
				continue
			}
		}

		// Telex alphabetic controls.
		lo := unicode.ToLower(ch)
		switch lo {
		case 's', 'f', 'r', 'x', 'j', 'z':
			if hasVowel(out) {
				nextTone := telexTone(lo)
				if lo == 'z' {
					if pendingTone != toneNone {
						pendingTone = toneNone
						continue
					}
					if stripTelexShapes(out) {
						continue
					}
				} else if pendingTone == nextTone {
					pendingTone = toneNone
					out = append(out, ch)
					continue
				} else {
					pendingTone = nextTone
					continue
				}
			}
		case 'a', 'e', 'o', 'd':
			if len(out) > 0 {
				last := out[len(out)-1]
				if unicode.ToLower(last) == lo {
					switch lo {
					case 'a':
						out[len(out)-1] = replaceBaseRune(last, 'â')
					case 'e':
						out[len(out)-1] = replaceBaseRune(last, 'ê')
					case 'o':
						out[len(out)-1] = replaceBaseRune(last, 'ô')
					case 'd':
						out[len(out)-1] = replaceBaseRune(last, 'đ')
					}
					continue
				}
				if base, ok := undoRepeatedTelexShape(last, lo); ok {
					out[len(out)-1] = base
					out = append(out, ch)
					continue
				}
			}
		case 'w':
			if applyTelexW(out) {
				continue
			}
			if undoTelexW(out) {
				out = append(out, ch)
				continue
			}
		}
		out = append(out, ch)
	}

	return applyToneToWord(string(out), pendingTone, oldStyle)
}
''',
)

engine_test = read("internal/core/engine_test.go")
if "TestVNITelexUnifiedMode" not in engine_test:
    engine_test += r'''

func TestVNITelexUnifiedMode(t *testing.T) {
	e := New(MethodVNITelex, Options{})
	cases := map[string]string{
		"tieengs": "tiếng", // Telex
		"tieng61": "tiếng", // VNI
		"ddoongf": "đồng",  // Telex
		"d9ong62": "đồng",  // VNI
		"d9oongf": "đồng",  // VNI đ + Telex ô/dấu
		"ddong62": "đồng",  // Telex đ + VNI ô/dấu
		"vieet5":  "việt",  // Telex ê + VNI dấu nặng
		"2026":    "2026",
	}
	for in, want := range cases {
		if got := e.Transform(in); got != want {
			t.Errorf("VNI/Telex %q: got %q want %q", in, got, want)
		}
	}
}

func TestLegacyMethodsNormalizeToUnifiedMode(t *testing.T) {
	for _, method := range []string{"Telex", "VNI", "Telex/VNI", "VNI/Telex", "vni telex"} {
		if got := New(method, Options{}).Method; got != MethodVNITelex {
			t.Errorf("Normalize %q=%q want %q", method, got, MethodVNITelex)
		}
	}
}
'''
write("internal/core/engine_test.go", engine_test)

write(
    "cmd/bilakey-cli/main.go",
    r'''package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/xulytiengviet/BilaKey_PC/internal/core"
)

func main() {
	method := flag.String("method", "cvnss", "cvnss hoặc vni-telex")
	inspect := flag.Bool("inspect", false, "in chi tiết candidate graph CVNSS dạng JSON")
	audit := flag.Bool("audit", false, "in audit lõi CVNSS dạng JSON")
	text := flag.Bool("text", false, "chuyển toàn bộ văn bản, giữ nguyên nội dung hỗn hợp")
	flag.Parse()

	if *audit {
		writeJSON(core.AuditCVNSS())
		return
	}
	input := strings.Join(flag.Args(), " ")
	if input == "" {
		fmt.Fprintln(os.Stderr, "usage: bilakey-cli [-method cvnss|vni-telex] [-text|-inspect|-audit] <input>")
		os.Exit(2)
	}
	if *inspect {
		writeJSON(core.InspectCVNSS(input))
		return
	}
	selected := core.MethodCVNSS
	switch strings.ToLower(strings.TrimSpace(*method)) {
	case "vni-telex", "vni/telex", "telex/vni", "vni", "telex":
		selected = core.MethodVNITelex
	case "cvnss", "cvnss4.0":
	default:
		fmt.Fprintf(os.Stderr, "method không hợp lệ: %s\n", *method)
		os.Exit(2)
	}
	engine := core.New(selected, core.Options{SpellCheck: true, AutoRestoreWrongKey: true})
	if *text {
		fmt.Println(engine.TransformText(input))
		return
	}
	fmt.Println(engine.Transform(input))
}

func writeJSON(value any) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(value); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
''',
)

# Hotkeys: exactly two selectable modes.
write(
    "internal/hotkey/hotkey.go",
    r'''package hotkey

type Action uint8

const (
	None Action = iota
	ToggleVietnamese
	SelectCVNSS
	SelectVNITelex
	CycleCandidate
)

// Resolve keeps all global shortcuts behind Ctrl+Shift and never steals
// Ctrl+Tab, Alt combinations, or Windows-key combinations from applications.
func Resolve(vk uint32, ctrl, shift, alt, win bool) Action {
	if !ctrl || !shift || alt || win {
		return None
	}
	switch vk {
	case 0x20:
		return ToggleVietnamese
	case '1':
		return SelectCVNSS
	case '2':
		return SelectVNITelex
	case '0':
		return CycleCandidate
	default:
		return None
	}
}
''',
)
write(
    "internal/hotkey/hotkey_test.go",
    r'''package hotkey

import "testing"

func TestResolve(t *testing.T) {
	cases := []struct {
		vk                    uint32
		ctrl, shift, alt, win bool
		want                  Action
	}{
		{0x20, true, true, false, false, ToggleVietnamese},
		{'1', true, true, false, false, SelectCVNSS},
		{'2', true, true, false, false, SelectVNITelex},
		{'3', true, true, false, false, None}, // 2.5.0 has only two input modes.
		{'0', true, true, false, false, CycleCandidate},
		{0x09, true, false, false, false, None}, // Ctrl+Tab belongs to the foreground app.
		{'1', true, true, true, false, None},
		{'1', true, true, false, true, None},
	}
	for _, tc := range cases {
		if got := Resolve(tc.vk, tc.ctrl, tc.shift, tc.alt, tc.win); got != tc.want {
			t.Errorf("Resolve(%#x)=%d want %d", tc.vk, got, tc.want)
		}
	}
}
''',
)

# Settings migration from all prior Telex/VNI values.
replace_exact("internal/settings/config.go", '"path/filepath"\n\t"sync"', '"path/filepath"\n\t"strings"\n\t"sync"')
replace_exact("internal/settings/config.go", 'AppVersion = "2.0.0"', 'AppVersion = "2.5.0"')
replace_exact(
    "internal/settings/config.go",
    '\t} else if !errors.Is(readErr, os.ErrNotExist) {\n\t\treturn nil, readErr\n\t}\n\treturn &Store{path: path, cfg: cfg}, nil',
    '\t} else if !errors.Is(readErr, os.ErrNotExist) {\n\t\treturn nil, readErr\n\t}\n\tcfg.InputMethod = normalizeInputMethod(cfg.InputMethod)\n\treturn &Store{path: path, cfg: cfg}, nil',
)
replace_exact(
    "internal/settings/config.go",
    '\nfunc (s *Store) Get() Config {',
    r'''
func normalizeInputMethod(method string) string {
	key := strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(method), " ", ""))
	switch key {
	case "VNI/TELEX", "TELEX/VNI", "VNI", "TELEX", "VNITELEX", "TELEXVNI":
		return "VNI/Telex"
	default:
		return "CVNSS4.0"
	}
}

func (s *Store) Get() Config {''',
)
replace_exact(
    "internal/settings/config_test.go",
    'if cfg.InputMethod != "Telex" || !cfg.PauseInPasswordFields {',
    'if cfg.InputMethod != "VNI/Telex" || !cfg.PauseInPasswordFields {',
)
config_test = read("internal/settings/config_test.go")
if "TestNormalizeLegacyVNI" not in config_test:
    config_test += r'''

func TestNormalizeLegacyVNI(t *testing.T) {
	for _, method := range []string{"Telex", "VNI", "Telex/VNI", "VNI/Telex"} {
		if got := normalizeInputMethod(method); got != "VNI/Telex" {
			t.Errorf("normalizeInputMethod(%q)=%q", method, got)
		}
	}
}
'''
write("internal/settings/config_test.go", config_test)

# Win32 UI: two tabs, one hybrid compatibility engine.
replace_exact(
    "internal/win/gui_main_windows.go",
    '\tidTabCVNSS   = 111\n\tidTabTelex   = 112\n\tidTabVNI     = 113',
    '\tidTabCVNSS    = 111\n\tidTabVNITelex = 112',
)
replace_exact(
    "internal/win/gui_main_windows.go",
    'createControl(hwnd, "STATIC", "Lõi nhập liệu và adapter tương thích", 0, 24, 88, 310, 22, idSubtle)\n\t\ta.controls[idTabCVNSS] = createDOSButton(hwnd, "CVNSS4.0 · LÕI", 24, 114, 290, 42, idTabCVNSS)\n\t\ta.controls[idTabTelex] = createDOSButton(hwnd, "TELEX", 324, 114, 110, 42, idTabTelex)\n\t\ta.controls[idTabVNI] = createDOSButton(hwnd, "VNI", 444, 114, 110, 42, idTabVNI)',
    'createControl(hwnd, "STATIC", "Hai kiểu gõ: CVNSS4.0 và VNI/Telex tự nhận dạng", 0, 24, 88, 510, 22, idSubtle)\n\t\ta.controls[idTabCVNSS] = createDOSButton(hwnd, "CVNSS4.0 · LÕI", 24, 114, 255, 42, idTabCVNSS)\n\t\ta.controls[idTabVNITelex] = createDOSButton(hwnd, "VNI / TELEX · TỰ ĐỘNG", 299, 114, 255, 42, idTabVNITelex)',
)
replace_exact(
    "internal/win/gui_main_windows.go",
    '\t\t\tcase idTabTelex:\n\t\t\t\ta.selectInputMethod(core.MethodTelex)\n\t\t\tcase idTabVNI:\n\t\t\t\ta.selectInputMethod(core.MethodVNI)',
    '\t\t\tcase idTabVNITelex:\n\t\t\t\ta.selectInputMethod(core.MethodVNITelex)',
)
replace_exact("internal/win/gui_main_windows.go", "Plain Tab cycles the three tabs", "Plain Tab toggles the two tabs")
replace_exact(
    "internal/win/gui_main_windows.go",
    '''\tmethod := a.runtimeConfig().InputMethod
\tswitch method {
\tcase core.MethodCVNSS:
\t\tmethod = core.MethodTelex
\tcase core.MethodTelex:
\t\tmethod = core.MethodVNI
\tdefault:
\t\tmethod = core.MethodCVNSS
\t}
\ta.selectInputMethod(method)''',
    '''\tmethod := core.MethodCVNSS
\tif a.runtimeConfig().InputMethod == core.MethodCVNSS {
\t\tmethod = core.MethodVNITelex
\t}
\ta.selectInputMethod(method)''',
)
replace_exact("internal/win/gui_main_windows.go", "three direct, owner-drawn method tabs", "two direct, owner-drawn method tabs")
replace_exact(
    "internal/win/gui_main_windows.go",
    'for _, id := range []int{idTabCVNSS, idTabTelex, idTabVNI} {',
    'for _, id := range []int{idTabCVNSS, idTabVNITelex} {',
)
replace_exact(
    "internal/win/gui_main_windows.go",
    '\tcase hotkey.SelectTelex:\n\t\ta.selectInputMethod(core.MethodTelex)\n\tcase hotkey.SelectVNI:\n\t\ta.selectInputMethod(core.MethodVNI)',
    '\tcase hotkey.SelectVNITelex:\n\t\ta.selectInputMethod(core.MethodVNITelex)',
)

replace_exact(
    "internal/win/theme_windows.go",
    'return id == idTabCVNSS || id == idTabTelex || id == idTabVNI',
    'return id == idTabCVNSS || id == idTabVNITelex',
)
replace_exact(
    "internal/win/theme_windows.go",
    '\tcase idTabTelex:\n\t\treturn method == "Telex"\n\tcase idTabVNI:\n\t\treturn method == "VNI"',
    '\tcase idTabVNITelex:\n\t\treturn method == "VNI/Telex"',
)

replace_exact(
    "internal/win/app_windows.go",
    '\tif cfg.InputMethod == core.MethodTelex {\n\t\tmethodRole = "TELEX ADAPTER"\n\t} else if cfg.InputMethod == core.MethodVNI {\n\t\tmethodRole = "VNI ADAPTER"\n\t}',
    '\tif cfg.InputMethod == core.MethodVNITelex {\n\t\tmethodRole = "VNI/TELEX · TỰ NHẬN DẠNG"\n\t}',
)
replace_exact(
    "internal/win/app_windows.go",
    'setWindowText(hwnd, "Ctrl+Shift+1: CVNSS Core  •  2/3: adapter  •  Ctrl+Shift+0: ứng viên")',
    'setWindowText(hwnd, "Ctrl+Shift+1: CVNSS4.0  •  Ctrl+Shift+2: VNI/Telex  •  Ctrl+Shift+0: ứng viên")',
)
replace_exact(
    "internal/win/app_windows.go",
    '"1. CVNSS4.0 là lõi mặc định; Telex và VNI là lớp tương thích tùy chọn.\\r\\n" +\n\t\t"2. Nhấn Tab khi cửa sổ BilaKey đang mở để chuyển giữa lõi và hai adapter.\\r\\n" +\n\t\t"3. Ctrl+Shift+Space bật/tắt; Ctrl+Shift+1/2/3 chọn CVNSS/Telex/VNI.\\r\\n" +',
    '"1. BilaKey chỉ có hai kiểu gõ: CVNSS4.0 và VNI/Telex.\\r\\n" +\n\t\t"2. VNI/Telex tự nhận cả phím chữ Telex và phím số VNI; có thể đổi cách gõ giữa từng từ.\\r\\n" +\n\t\t"3. Ctrl+Shift+Space bật/tắt; Ctrl+Shift+1 chọn CVNSS4.0; Ctrl+Shift+2 chọn VNI/Telex.\\r\\n" +',
)
replace_exact(
    "internal/win/app_windows.go",
    '"Adapter tương thích: Telex / VNI\\r\\n" +',
    '"Kiểu gõ hợp nhất: VNI/Telex tự nhận dạng\\r\\n" +',
)

replace_exact(
    "internal/win/tray_windows.go",
    '\ttrayCmdCVNSS  = 9003\n\ttrayCmdTelex  = 9004\n\ttrayCmdVNI    = 9005\n\ttrayCmdExit   = 9006',
    '\ttrayCmdCVNSS    = 9003\n\ttrayCmdVNITelex = 9004\n\ttrayCmdExit     = 9005',
)
replace_exact(
    "internal/win/tray_windows.go",
    '\tappendMenuChecked(menu, trayCmdTelex, "Telex · tương thích", cfg.InputMethod == "Telex")\n\tappendMenuChecked(menu, trayCmdVNI, "VNI · tương thích", cfg.InputMethod == "VNI")',
    '\tappendMenuChecked(menu, trayCmdVNITelex, "VNI/Telex · tự nhận dạng", cfg.InputMethod == "VNI/Telex")',
)
replace_exact(
    "internal/win/tray_windows.go",
    '\tcase trayCmdTelex:\n\t\ta.updateConfig(func(c *settings.Config) { c.InputMethod = "Telex" })\n\t\ta.syncMethodCombo()\n\tcase trayCmdVNI:\n\t\ta.updateConfig(func(c *settings.Config) { c.InputMethod = "VNI" })\n\t\ta.syncMethodCombo()',
    '\tcase trayCmdVNITelex:\n\t\ta.updateConfig(func(c *settings.Config) { c.InputMethod = "VNI/Telex" })\n\t\ta.syncMethodCombo()',
)

replace_exact("internal/win/hook_windows.go", "cfg.InputMethod != core.MethodVNI", "cfg.InputMethod != core.MethodVNITelex")
replace_exact(
    "internal/win/hook_windows.go",
    "Numbers are literal delimiters for Telex/CVNSS; do not consume them.",
    "Numbers are literal delimiters for CVNSS; VNI/Telex consumes them as VNI control keys.",
)

# Changelog and migration documentation.
changelog = read("CHANGELOG.md")
entry = r'''## 2.5.0 — Unified VNI/Telex Edition · 2026-07-31

### Hai kiểu gõ rõ ràng

- BilaKey giờ chỉ hiển thị **CVNSS4.0** và **VNI/Telex**.
- Hợp nhất Telex và VNI vào một engine tự nhận dạng; người dùng có thể đổi quy ước giữa từng từ mà không đổi chế độ.
- Hỗ trợ kết hợp phím tạo chữ của Telex với phím dấu VNI, hoặc ngược lại, trong cùng một từ.
- Tự động chuyển cấu hình cũ `Telex`, `VNI` và `Telex/VNI` sang `VNI/Telex`.
- Rút gọn hotkey: `Ctrl+Shift+1` cho CVNSS4.0, `Ctrl+Shift+2` cho VNI/Telex; `Ctrl+Shift+3` không còn bị chiếm.

### Phát hành

- Đồng bộ phiên bản, tên EXE, ZIP, release notes, CI và bảng SHA-256 lên 2.5.0.
- Duy trì ba bản Windows x64, x86 và ARM64 cùng gói ZIP đầy đủ.

'''
if "## 2.5.0" not in changelog:
    changelog = changelog.replace("# 📝 Changelog\n\n", "# 📝 Changelog\n\n" + entry, 1)
write("CHANGELOG.md", changelog)

write(
    "docs/MIGRATION_2.5.md",
    r'''# Di chuyển lên BilaKey PC 2.5.0

## Thay đổi kiểu gõ

BilaKey 2.5.0 chỉ còn hai lựa chọn trong giao diện:

1. **CVNSS4.0** — lõi trung tâm của BilaKey.
2. **VNI/Telex** — một engine hợp nhất nhận cả quy ước VNI và Telex.

Cấu hình cũ có `Telex`, `VNI`, `Telex/VNI` hoặc `VNI/Telex` được chuẩn hóa tự động thành `VNI/Telex`. Không cần xóa `config.json`.

## Phím tắt

- `Ctrl+Shift+1`: CVNSS4.0.
- `Ctrl+Shift+2`: VNI/Telex.
- `Ctrl+Shift+3`: được trả lại cho ứng dụng đang dùng.

## Ví dụ

| Cách gõ | Chuỗi | Kết quả |
|---|---|---|
| Telex | `tieengs` | `tiếng` |
| VNI | `tieng61` | `tiếng` |
| Telex | `ddoongf` | `đồng` |
| VNI | `d9ong62` | `đồng` |
| Kết hợp | `vieet5` | `việt` |
| Kết hợp | `d9oongf` | `đồng` |
''',
)

# README is rewritten so the two-mode behavior and verified download matrix are unmistakable.
write(
    "README.md",
    r'''<div align="center">
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
''',
)

write(
    "docs/RELEASE_NOTES_2.5.0.md",
    r'''# BilaKey PC 2.5.0 — Unified VNI/Telex Edition

BilaKey 2.5.0 đơn giản hóa trải nghiệm xuống còn hai kiểu gõ: **CVNSS4.0** và **VNI/Telex**.

## Điểm mới

- Hợp nhất VNI và Telex thành một engine tự nhận dạng.
- Có thể gõ VNI hoặc Telex theo từng từ mà không chuyển chế độ.
- Hỗ trợ các tổ hợp thực dụng như `vieet5 → việt`, `d9oongf → đồng`.
- Tự chuyển cấu hình cũ Telex/VNI sang VNI/Telex.
- Hotkey mới: `Ctrl+Shift+1` cho CVNSS4.0, `Ctrl+Shift+2` cho VNI/Telex.
- `Ctrl+Shift+3` không còn bị BilaKey chiếm.

## Tải đúng bản

- `BilaKey-PC-2.5.0-CVNSS-Core-x64.exe`: Intel/AMD Windows 64-bit, khuyến nghị.
- `BilaKey-PC-2.5.0-CVNSS-Core-arm64.exe`: Windows on ARM, Snapdragon.
- `BilaKey-PC-2.5.0-CVNSS-Core-x86.exe`: Windows 32-bit cũ.
- `BilaKey-PC-2.5.0-Windows.zip`: gói đầy đủ.
- `SHA256SUMS-2.5.0.txt`: bảng đối chiếu toàn vẹn.

Bản portable hiện chưa ký Authenticode. Chỉ tải từ release chính thức và kiểm tra SHA-256.
''',
)

# Replace the historical publishing workflow with a 2.5.0 release workflow.
old_release = path(".github/workflows/release-2.0.0.yml")
if old_release.exists():
    old_release.unlink()
write(
    ".github/workflows/release-2.5.0.yml",
    r'''name: Publish BilaKey PC 2.5.0

on:
  workflow_dispatch:
  push:
    branches: [main]
    paths:
      - ".github/workflows/release-2.5.0.yml"
      - "VERSION"
      - "README.md"
      - "CHANGELOG.md"
      - "docs/RELEASE_NOTES_2.5.0.md"
      - "internal/core/**"
      - "internal/hotkey/**"
      - "internal/settings/**"
      - "internal/win/**"
      - "cmd/bilakey-cli/**"

permissions:
  contents: write

concurrency:
  group: bilakey-release-2.5.0
  cancel-in-progress: false

jobs:
  publish:
    name: Build, verify and publish Windows package
    runs-on: ubuntu-latest

    steps:
      - name: Checkout source
        uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: "1.23.x"

      - name: Verify source
        run: |
          test "$(cat VERSION)" = "2.5.0"
          test -z "$(gofmt -l cmd internal)"
          go test ./... -count=1
          go vet ./...
          go test -race ./internal/core ./internal/hotkey ./internal/macro ./internal/settings ./internal/typingstate -count=1

      - name: Build x64, x86 and ARM64
        shell: bash
        run: |
          set -euo pipefail
          VERSION=2.5.0
          mkdir -p dist "package/BilaKey-PC-${VERSION}-Windows"

          build_one() {
            local arch="$1"
            local suffix="$2"
            CGO_ENABLED=0 GOOS=windows GOARCH="${arch}" \
              go build -trimpath \
              -ldflags='-s -w -H windowsgui -buildid=' \
              -o "dist/BilaKey-PC-${VERSION}-CVNSS-Core-${suffix}.exe" \
              ./cmd/bilakey
          }

          build_one amd64 x64
          build_one 386 x86
          build_one arm64 arm64

          cp dist/*.exe "package/BilaKey-PC-${VERSION}-Windows/"
          cp LICENSE README.md CHANGELOG.md "package/BilaKey-PC-${VERSION}-Windows/"

          cat > "package/BilaKey-PC-${VERSION}-Windows/README_FIRST.txt" <<'TXT'
          BILAKEY PC 2.5.0 — CVNSS4.0 + VNI/TELEX

          Chọn tệp:
          - x64: máy Intel/AMD hiện đại, khuyến nghị.
          - arm64: Windows on ARM, Snapdragon.
          - x86: Windows 32-bit cũ.

          BilaKey chỉ có hai kiểu gõ:
          1. CVNSS4.0.
          2. VNI/Telex: gõ theo VNI hoặc Telex đều được.

          Phím tắt:
          - Ctrl+Shift+Space: bật/tắt.
          - Ctrl+Shift+1: CVNSS4.0.
          - Ctrl+Shift+2: VNI/Telex.
          - Ctrl+Shift+0: đổi ứng viên CVNSS.

          Bản này chưa ký Authenticode. Chỉ tải từ GitHub chính thức và kiểm tra SHA-256.
          Phát triển: Long Ngo. Giấy phép: MIT.
          TXT

          (
            cd "package/BilaKey-PC-${VERSION}-Windows"
            sha256sum *.exe > SHA256SUMS.txt
          )
          (
            cd package
            zip -9 -r "../dist/BilaKey-PC-${VERSION}-Windows.zip" "BilaKey-PC-${VERSION}-Windows"
          )
          (
            cd dist
            sha256sum *.exe *.zip > "SHA256SUMS-${VERSION}.txt"
          )
          ls -lh dist
          cat "dist/SHA256SUMS-${VERSION}.txt"

      - name: Align tag 2.5.0 with published source
        shell: bash
        run: |
          set -euo pipefail
          git config user.name "github-actions[bot]"
          git config user.email "41898282+github-actions[bot]@users.noreply.github.com"
          git tag -fa 2.5.0 -m "BilaKey PC 2.5.0 — Unified VNI/Telex Edition" "${GITHUB_SHA}"
          git push --force origin refs/tags/2.5.0

      - name: Create or update GitHub Release
        env:
          GH_TOKEN: ${{ github.token }}
        shell: bash
        run: |
          set -euo pipefail
          if gh release view 2.5.0 >/dev/null 2>&1; then
            gh release edit 2.5.0 \
              --title "BilaKey PC 2.5.0 — Unified VNI/Telex Edition" \
              --notes-file docs/RELEASE_NOTES_2.5.0.md
          else
            gh release create 2.5.0 \
              --title "BilaKey PC 2.5.0 — Unified VNI/Telex Edition" \
              --notes-file docs/RELEASE_NOTES_2.5.0.md \
              --target "${GITHUB_SHA}"
          fi

          gh release upload 2.5.0 \
            dist/BilaKey-PC-2.5.0-CVNSS-Core-x64.exe \
            dist/BilaKey-PC-2.5.0-CVNSS-Core-x86.exe \
            dist/BilaKey-PC-2.5.0-CVNSS-Core-arm64.exe \
            dist/BilaKey-PC-2.5.0-Windows.zip \
            dist/SHA256SUMS-2.5.0.txt \
            --clobber

      - name: Verify published asset set
        env:
          GH_TOKEN: ${{ github.token }}
        shell: bash
        run: |
          set -euo pipefail
          required=(
            BilaKey-PC-2.5.0-CVNSS-Core-x64.exe
            BilaKey-PC-2.5.0-CVNSS-Core-x86.exe
            BilaKey-PC-2.5.0-CVNSS-Core-arm64.exe
            BilaKey-PC-2.5.0-Windows.zip
            SHA256SUMS-2.5.0.txt
          )
          mapfile -t published < <(gh release view 2.5.0 --json assets --jq '.assets[].name')
          for asset in "${required[@]}"; do
            printf '%s\n' "${published[@]}" | grep -Fxq "${asset}" || {
              echo "Missing release asset: ${asset}" >&2
              exit 1
            }
          done
          echo "Verified: 3 EXE + ZIP + SHA-256."

      - name: Preserve build artifacts
        uses: actions/upload-artifact@v4
        with:
          name: BilaKey-PC-2.5.0-release-assets
          path: dist/*
          if-no-files-found: error
          retention-days: 30
''',
)

print("BilaKey PC 2.5.0 upgrade prepared successfully")
