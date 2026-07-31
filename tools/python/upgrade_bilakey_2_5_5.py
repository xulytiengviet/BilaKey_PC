#!/usr/bin/env python3
from __future__ import annotations

from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]


def p(rel: str) -> Path:
    return ROOT / rel


def read(rel: str) -> str:
    return p(rel).read_text(encoding="utf-8")


def write(rel: str, content: str) -> None:
    target = p(rel)
    target.parent.mkdir(parents=True, exist_ok=True)
    target.write_text(content, encoding="utf-8", newline="\n")


def replace_exact(rel: str, old: str, new: str) -> None:
    content = read(rel)
    count = content.count(old)
    if count != 1:
        raise RuntimeError(f"{rel}: expected exactly one match, found {count}: {old!r}")
    write(rel, content.replace(old, new, 1))


def replace_all(rel: str, old: str, new: str) -> None:
    content = read(rel)
    if old not in content:
        raise RuntimeError(f"{rel}: missing text {old!r}")
    write(rel, content.replace(old, new))


write("VERSION", "2.5.5\n")
replace_exact("internal/settings/config.go", 'AppVersion = "2.5.0"', 'AppVersion = "2.5.5"')
replace_exact("scripts/build_release.sh", 'VERSION="2.5.0"', 'VERSION="2.5.5"')
replace_all("BUILD.md", "2.5.0", "2.5.5")
if p("SBOM.cdx.json").exists():
    replace_all("SBOM.cdx.json", "2.5.0", "2.5.5")

write(
    "internal/win/gui_main_windows.go",
    r'''//go:build windows

package win

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/xulytiengviet/BilaKey_PC/internal/core"
	"github.com/xulytiengviet/BilaKey_PC/internal/hotkey"
	"github.com/xulytiengviet/BilaKey_PC/internal/settings"
)

const (
	classMain    = "BilaKeyPCMainWindow"
	classOptions = "BilaKeyPCOptionsWindow"
	classMacro   = "BilaKeyPCMacroWindow"
)

const (
	idToggle       = 101
	idExit         = 102
	idExpand       = 103
	idHelp         = 104
	idAbout        = 105
	idSettings     = 106
	idStatus       = 107
	idCapsStatus   = 108
	idLogo         = 109
	idHeroTitle    = 110
	idTabCVNSS     = 111
	idTabVNITelex  = 112
	idHeroSubtitle = 113
	idHeroTagline  = 114
	idModeHint     = 115
	idInfoTitle    = 116
	idInfoBody     = 117
	idInfoMeta     = 118

	idOptFreeTone         = 201
	idOptOldTone          = 202
	idOptClipboard        = 203
	idOptSpell            = 204
	idOptRestore          = 205
	idOptFeedback         = 206
	idOptMacro            = 207
	idOptMacroOff         = 208
	idOptMacroTable       = 209
	idOptShowStartup      = 210
	idOptStartWindows     = 211
	idOptVietnamese       = 212
	idOptSave             = 213
	idOptClose            = 214
	idOptAutoCapInitial   = 215
	idOptAutoCapSentence  = 216
	idOptDoubleShiftCaps  = 217
	idOptRestoreDelimiter = 218
	idOptPausePassword    = 219

	idMacroTrigger     = 301
	idMacroReplacement = 302
	idMacroList        = 303
	idMacroAdd         = 304
	idMacroDelete      = 305
	idMacroSave        = 306
	idMacroChoose      = 307
	idMacroDefault     = 308
	idMacroClose       = 309
	idMacroPath        = 310
)

var (
	mainWndProcPtr    = syscall.NewCallback(mainWndProc)
	optionsWndProcPtr = syscall.NewCallback(optionsWndProc)
	macroWndProcPtr   = syscall.NewCallback(macroWndProc)
)

func (a *App) registerClasses() error {
	if err := registerClass(classMain, mainWndProcPtr); err != nil {
		return err
	}
	if err := registerClass(classOptions, optionsWndProcPtr); err != nil {
		return err
	}
	return registerClass(classMacro, macroWndProcPtr)
}

func createWindow(class, title string, style uint32, x, y, w, h int32, parent uintptr, id int) uintptr {
	hinst, _, _ := procGetModuleHandleW.Call(0)
	hwnd, _, _ := procCreateWindowExW.Call(
		0,
		uintptr(unsafe.Pointer(utf16Ptr(class))),
		uintptr(unsafe.Pointer(utf16Ptr(title))),
		uintptr(style),
		uintptr(x), uintptr(y), uintptr(w), uintptr(h),
		parent,
		uintptr(id),
		hinst,
		0,
	)
	return hwnd
}

func createControl(parent uintptr, class, text string, style uint32, x, y, w, h int32, id int) uintptr {
	hwnd := createWindow(class, text, WS_CHILD|WS_VISIBLE|style, x, y, w, h, parent, id)
	if currentApp != nil {
		currentApp.applyControlTheme(hwnd, id)
	}
	return hwnd
}

func createDOSButton(parent uintptr, text string, x, y, w, h int32, id int) uintptr {
	return createControl(parent, "BUTTON", text, BS_OWNERDRAW|WS_TABSTOP, x, y, w, h, id)
}

func (a *App) createMainWindow() error {
	// 1000 × 650 tạo vùng làm việc rộng, tương đương khoảng 1/4 diện tích
	// màn hình Full HD và vẫn vừa trên các máy Windows phổ biến.
	hwnd := createWindow(classMain, "BilaKey PC "+settings.AppVersion, WS_OVERLAPPEDWINDOW, 150, 70, 1000, 650, 0, 0)
	if hwnd == 0 {
		return fmt.Errorf("không thể tạo cửa sổ chính")
	}
	a.mainHwnd = hwnd
	if a.iconSmall != 0 {
		sendMessage(hwnd, WM_SETICON, ICON_SMALL, a.iconSmall)
	}
	if a.iconBig != 0 {
		sendMessage(hwnd, WM_SETICON, ICON_BIG, a.iconBig)
	}
	a.refreshMainStatus()
	return nil
}

func mainWndProc(hwnd uintptr, msg uint32, wParam, lParam uintptr) uintptr {
	a := currentApp
	if a == nil {
		r, _, _ := procDefWindowProcW.Call(hwnd, uintptr(msg), wParam, lParam)
		return r
	}

	if msg == trayCallbackMessage {
		a.handleTrayMessage(lParam)
		return 0
	}
	if msg == methodCycleMessage {
		a.cycleInputMethod()
		return 0
	}
	if msg == hotkeyActionMessage {
		a.applyHotkeyAction(hotkey.Action(wParam))
		return 0
	}

	switch msg {
	case WM_CREATE:
		createControl(hwnd, "STATIC", "B", WS_BORDER, 32, 22, 84, 84, idLogo)
		createControl(hwnd, "STATIC", "BilaKey PC", 0, 140, 18, 430, 44, idHeroTitle)
		createControl(hwnd, "STATIC", "CVNSS4.0 Core, Unicode, riêng tư, kiểm toán được", 0, 140, 64, 720, 28, idHeroSubtitle)
		createControl(hwnd, "STATIC", "Kiểu gõ dấu chữ tiếng Việt", 0, 140, 94, 420, 26, idHeroTagline)

		createControl(hwnd, "STATIC", "HAI KIỂU GÕ · CHỌN MỘT LẦN, DÙNG ỔN ĐỊNH", 0, 28, 128, 620, 22, idModeHint)
		a.controls[idTabCVNSS] = createDOSButton(hwnd, "CVNSS4.0  •  LÕI", 28, 154, 450, 84, idTabCVNSS)
		a.controls[idTabVNITelex] = createDOSButton(hwnd, "VNI / TELEX  •  TỰ ĐỘNG", 492, 154, 450, 84, idTabVNITelex)

		a.controls[idStatus] = createControl(hwnd, "STATIC", "", WS_BORDER, 28, 250, 914, 42, idStatus)
		a.controls[idCapsStatus] = createControl(hwnd, "STATIC", "", WS_BORDER, 28, 302, 914, 38, idCapsStatus)

		createControl(hwnd, "STATIC", "THÔNG TIN", WS_BORDER, 28, 352, 914, 34, idInfoTitle)
		createControl(hwnd, "STATIC",
			"• VNI + Telex hợp nhất: gõ theo VNI hoặc Telex; BilaKey tự nhận dạng và xuất chữ tiếng Việt Unicode.\r\n\r\n• CVNSS4.0: kiểu gõ chuyên dụng lấy CVNSS4.0 làm lõi trung tâm.",
			WS_BORDER, 28, 386, 574, 108, idInfoBody)
		createControl(hwnd, "STATIC",
			"• Tác giả: Long Ngo phát triển\r\n\r\n• Dự án CVNSS4.0\r\n\r\n• Giấy phép: MIT",
			WS_BORDER, 604, 386, 338, 108, idInfoMeta)

		a.controls[idToggle] = createDOSButton(hwnd, "BẬT/TẮT", 28, 514, 160, 50, idToggle)
		createDOSButton(hwnd, "MỞ RỘNG", 206, 514, 174, 50, idExpand)
		createDOSButton(hwnd, "THÔNG TIN", 398, 514, 174, 50, idAbout)
		createDOSButton(hwnd, "HƯỚNG DẪN", 590, 514, 174, 50, idHelp)
		createDOSButton(hwnd, "THOÁT", 782, 514, 160, 50, idExit)

		a.refreshMainStatus()
		a.syncMethodCombo()
		return 0

	case WM_DRAWITEM:
		if a.drawDOSButton(lParam) {
			return 1
		}
	case WM_CTLCOLORSTATIC, WM_CTLCOLORBTN, WM_CTLCOLOREDIT, WM_CTLCOLORLISTBOX:
		return a.themeControlColor(msg, wParam, lParam)

	case WM_COMMAND:
		id := int(loword(wParam))
		notify := hiword(wParam)
		if notify == BN_CLICKED {
			switch id {
			case idTabCVNSS:
				a.selectInputMethod(core.MethodCVNSS)
			case idTabVNITelex:
				a.selectInputMethod(core.MethodVNITelex)
			case idToggle:
				a.updateConfig(func(c *settings.Config) { c.Enabled = !c.Enabled })
			case idExit:
				a.exitApp()
			case idExpand:
				a.showOptions()
			case idHelp:
				a.showHelp()
			case idAbout:
				a.showAbout()
			case idSettings:
				a.showSettingsInfo()
			}
			return 0
		}

	case WM_KEYDOWN:
		if uint32(wParam) == VK_TAB {
			a.cycleInputMethod()
			return 0
		}
		handled := true
		switch uint32(wParam) {
		case 0x70: // F1
			a.showHelp()
		case 0x71: // F2
			a.updateConfig(func(c *settings.Config) { c.Enabled = !c.Enabled })
		case 0x72: // F3
			a.showOptions()
		case 0x73: // F4
			a.showAbout()
		case 0x77: // F8
			a.showSettingsInfo()
		case 0x79: // F10
			a.exitApp()
		default:
			handled = false
		}
		if handled {
			return 0
		}

	case WM_CLOSE:
		if a.exiting {
			procDestroyWindow.Call(hwnd)
		} else {
			procShowWindow.Call(hwnd, SW_HIDE)
			a.feedback("BilaKey vẫn đang chạy · icon B ở khay hệ thống")
		}
		return 0
	case WM_DESTROY:
		procPostQuitMessage.Call(0)
		return 0
	}
	r, _, _ := procDefWindowProcW.Call(hwnd, uintptr(msg), wParam, lParam)
	return r
}

func (a *App) selectInputMethod(method string) {
	if a.runtimeConfig().InputMethod == method {
		return
	}
	a.updateConfig(func(c *settings.Config) { c.InputMethod = method })
	a.syncMethodCombo()
}

func (a *App) cycleInputMethod() {
	method := core.MethodCVNSS
	if a.runtimeConfig().InputMethod == core.MethodCVNSS {
		method = core.MethodVNITelex
	}
	a.selectInputMethod(method)
}

func (a *App) syncMethodCombo() {
	for _, id := range []int{idTabCVNSS, idTabVNITelex, idToggle} {
		if hwnd := a.controls[id]; hwnd != 0 {
			procInvalidateRect.Call(hwnd, 0, 1)
			procUpdateWindow.Call(hwnd)
		}
	}
}

func (a *App) applyHotkeyAction(action hotkey.Action) {
	switch action {
	case hotkey.ToggleVietnamese:
		a.updateConfig(func(c *settings.Config) { c.Enabled = !c.Enabled })
	case hotkey.SelectCVNSS:
		a.selectInputMethod(core.MethodCVNSS)
	case hotkey.SelectVNITelex:
		a.selectInputMethod(core.MethodVNITelex)
	case hotkey.CycleCandidate:
		a.cycleCVNSSCandidate()
	}
}
''',
)

write(
    "internal/win/theme_windows.go",
    r'''//go:build windows

package win

import "unsafe"

// COLORREF uses 0x00BBGGRR.
const (
	oceanCanvas  = uint32(0x00FFFDFC)
	oceanHero    = uint32(0x00C66B09)
	oceanDeep    = uint32(0x00964A05)
	oceanPanel   = uint32(0x00F7EADC)
	oceanAccent  = uint32(0x00E58B16)
	oceanWhite   = uint32(0x00FFFFFF)
	oceanSoft    = uint32(0x00FFF6EA)
	oceanBorder  = uint32(0x00E3C7AA)
	oceanText    = uint32(0x00472F1C)
	oceanMuted   = uint32(0x00816C59)
	oceanSuccess = uint32(0x005A9B18)
	oceanStatus  = uint32(0x00F0FBF2)
)

func (a *App) initTheme() {
	if a.theme.font != 0 {
		return
	}
	a.theme.bgBrush = createSolidBrush(oceanCanvas)
	a.theme.heroBrush = createSolidBrush(oceanHero)
	a.theme.panelBrush = createSolidBrush(oceanPanel)
	a.theme.editBrush = createSolidBrush(oceanWhite)
	a.theme.cyanBrush = createSolidBrush(oceanAccent)
	a.theme.grayBrush = createSolidBrush(oceanBorder)
	a.theme.softBrush = createSolidBrush(oceanSoft)
	a.theme.statusBrush = createSolidBrush(oceanStatus)
	a.theme.font = createUIFont(FW_NORMAL, 16)
	a.theme.fontSmall = createUIFont(FW_NORMAL, 14)
	a.theme.fontBold = createUIFont(FW_BOLD, 18)
	a.theme.fontMethod = createUIFont(FW_BOLD, 24)
	a.theme.fontTitle = createUIFont(FW_BOLD, 34)
	a.theme.fontLogo = createUIFont(FW_BOLD, 46)
}

func (a *App) destroyTheme() {
	for _, h := range []uintptr{
		a.theme.font,
		a.theme.fontSmall,
		a.theme.fontBold,
		a.theme.fontMethod,
		a.theme.fontTitle,
		a.theme.fontLogo,
		a.theme.bgBrush,
		a.theme.heroBrush,
		a.theme.panelBrush,
		a.theme.editBrush,
		a.theme.cyanBrush,
		a.theme.grayBrush,
		a.theme.softBrush,
		a.theme.statusBrush,
	} {
		if h != 0 {
			procDeleteObject.Call(h)
		}
	}
	a.theme = themeHandles{}
}

func createSolidBrush(color uint32) uintptr {
	h, _, _ := procCreateSolidBrush.Call(uintptr(color))
	return h
}

func createUIFont(weight int, px int) uintptr {
	face := utf16Ptr("Segoe UI")
	h, _, _ := procCreateFontW.Call(
		uintptr(int32(-px)),
		0,
		0,
		0,
		uintptr(weight),
		0,
		0,
		0,
		DEFAULT_CHARSET,
		OUT_DEFAULT_PRECIS,
		CLIP_DEFAULT_PRECIS,
		DEFAULT_QUALITY,
		DEFAULT_PITCH|FF_SWISS,
		uintptr(unsafe.Pointer(face)),
	)
	return h
}

func (a *App) applyControlTheme(hwnd uintptr, id int) {
	if hwnd == 0 {
		return
	}
	font := a.theme.font
	switch id {
	case idLogo:
		font = a.theme.fontLogo
	case idHeroTitle:
		font = a.theme.fontTitle
	case idHeroSubtitle, idHeroTagline, idModeHint, idStatus:
		font = a.theme.fontBold
	case idCapsStatus, idInfoBody, idInfoMeta:
		font = a.theme.fontSmall
	case idInfoTitle:
		font = a.theme.fontBold
	}
	if font != 0 {
		sendMessage(hwnd, WM_SETFONT, font, 1)
	}
}

func (a *App) themeControlColor(msg uint32, hdc, child uintptr) uintptr {
	if hdc == 0 {
		return 0
	}
	procSetBkMode.Call(hdc, TRANSPARENT)
	procSetTextColor.Call(hdc, uintptr(oceanText))

	switch msg {
	case WM_CTLCOLOREDIT, WM_CTLCOLORLISTBOX:
		procSetBkMode.Call(hdc, 2)
		procSetBkColor.Call(hdc, uintptr(oceanWhite))
		procSetTextColor.Call(hdc, uintptr(oceanText))
		return a.theme.editBrush
	case WM_CTLCOLORSTATIC:
		id, _, _ := procGetDlgCtrlID.Call(child)
		switch int(id) {
		case idLogo, idHeroTitle, idHeroSubtitle, idHeroTagline:
			procSetTextColor.Call(hdc, uintptr(oceanWhite))
			return a.theme.heroBrush
		case idStatus:
			procSetTextColor.Call(hdc, uintptr(oceanSuccess))
			return a.theme.statusBrush
		case idCapsStatus:
			procSetTextColor.Call(hdc, uintptr(oceanDeep))
			return a.theme.softBrush
		case idInfoTitle, idInfoBody, idInfoMeta:
			procSetTextColor.Call(hdc, uintptr(oceanText))
			return a.theme.editBrush
		case idModeHint:
			procSetTextColor.Call(hdc, uintptr(oceanDeep))
			return a.theme.bgBrush
		default:
			procSetTextColor.Call(hdc, uintptr(oceanText))
			return a.theme.bgBrush
		}
	case WM_CTLCOLORBTN:
		procSetTextColor.Call(hdc, uintptr(oceanText))
		return a.theme.bgBrush
	}
	return a.theme.bgBrush
}

func (a *App) isMethodTab(id uint32) bool {
	return id == idTabCVNSS || id == idTabVNITelex
}

func (a *App) methodTabActive(id uint32) bool {
	method := a.runtimeConfig().InputMethod
	switch id {
	case idTabCVNSS:
		return method == "CVNSS4.0"
	case idTabVNITelex:
		return method == "VNI/Telex"
	}
	return false
}

func (a *App) drawDOSButton(lParam uintptr) bool {
	if lParam == 0 {
		return false
	}
	dis := (*drawItemStruct)(unsafe.Pointer(lParam))
	if dis.HDC == 0 || dis.HwndItem == 0 {
		return false
	}

	fill := a.theme.editBrush
	border := a.theme.grayBrush
	textColor := oceanDeep
	if a.isMethodTab(dis.CtlID) && a.methodTabActive(dis.CtlID) {
		fill = a.theme.cyanBrush
		border = a.theme.cyanBrush
		textColor = oceanWhite
	}
	if dis.CtlID == idToggle {
		fill = a.theme.cyanBrush
		border = a.theme.cyanBrush
		textColor = oceanWhite
	}
	if dis.ItemState&ODS_SELECTED != 0 {
		fill = a.theme.panelBrush
		border = a.theme.cyanBrush
		textColor = oceanDeep
	}
	if dis.ItemState&ODS_DISABLED != 0 {
		textColor = oceanMuted
	}

	procFillRect.Call(dis.HDC, uintptr(unsafe.Pointer(&dis.RcItem)), fill)
	procFrameRect.Call(dis.HDC, uintptr(unsafe.Pointer(&dis.RcItem)), border)
	procSetBkMode.Call(dis.HDC, TRANSPARENT)
	procSetTextColor.Call(dis.HDC, uintptr(textColor))
	font := a.theme.fontBold
	if a.isMethodTab(dis.CtlID) {
		font = a.theme.fontMethod
	}
	if font != 0 {
		procSelectObject.Call(dis.HDC, font)
	}

	text := getWindowText(dis.HwndItem)
	if text != "" {
		p := utf16Ptr(text)
		procDrawTextW.Call(
			dis.HDC,
			uintptr(unsafe.Pointer(p)),
			^uintptr(0),
			uintptr(unsafe.Pointer(&dis.RcItem)),
			DT_CENTER|DT_VCENTER|DT_SINGLELINE,
		)
	}
	return true
}
''',
)

replace_exact(
    "internal/win/app_windows.go",
    '''type themeHandles struct {
\tfont       uintptr
\tfontBold   uintptr
\tbgBrush    uintptr
\tpanelBrush uintptr
\teditBrush  uintptr
\tcyanBrush  uintptr
\tgrayBrush  uintptr
}''',
    '''type themeHandles struct {
\tfont        uintptr
\tfontSmall   uintptr
\tfontBold    uintptr
\tfontMethod  uintptr
\tfontTitle   uintptr
\tfontLogo    uintptr
\tbgBrush     uintptr
\theroBrush   uintptr
\tpanelBrush  uintptr
\teditBrush   uintptr
\tcyanBrush   uintptr
\tgrayBrush   uintptr
\tsoftBrush   uintptr
\tstatusBrush uintptr
}''',
)

old_refresh = '''func (a *App) refreshMainStatus() {
\tcfg := a.runtimeConfig()
\tstate := "ĐANG BẬT"
\tif !cfg.Enabled {
\t\tstate = "ĐANG TẮT"
\t}
\tmethodRole := "CVNSS4.0 CORE"
\tif cfg.InputMethod == core.MethodVNITelex {
\t\tmethodRole = "VNI/TELEX · TỰ NHẬN DẠNG"
\t}
\tif hwnd := a.controls[idStatus]; hwnd != 0 {
\t\tsetWindowText(hwnd, fmt.Sprintf("  %s  •  %s  •  Unicode  •  Ctrl+Shift+Space: bật/tắt", state, methodRole))
\t}
\tif hwnd := a.controls[idCapsStatus]; hwnd != 0 {
\t\tsetWindowText(hwnd, "Ctrl+Shift+1: CVNSS4.0  •  Ctrl+Shift+2: VNI/Telex  •  Ctrl+Shift+0: ứng viên")
\t}
\tif hwnd := a.controls[idToggle]; hwnd != 0 {
\t\tif cfg.Enabled {
\t\t\tsetWindowText(hwnd, "TẮT")
\t\t} else {
\t\t\tsetWindowText(hwnd, "BẬT")
\t\t}
\t}
\ta.syncMethodCombo()
\ta.updateTrayIcon()
}'''
new_refresh = '''func (a *App) refreshMainStatus() {
\tcfg := a.runtimeConfig()
\tstate := "ĐANG BẬT"
\tif !cfg.Enabled {
\t\tstate = "ĐANG TẮT"
\t}
\tmethodRole := "CVNSS4.0 · LÕI"
\tif cfg.InputMethod == core.MethodVNITelex {
\t\tmethodRole = "VNI/TELEX · TỰ NHẬN DẠNG"
\t}
\tif hwnd := a.controls[idStatus]; hwnd != 0 {
\t\tsetWindowText(hwnd, fmt.Sprintf("   ✓  %s   •   %s   •   Unicode", state, methodRole))
\t}
\tif hwnd := a.controls[idCapsStatus]; hwnd != 0 {
\t\tsetWindowText(hwnd, "   ⌨  Ctrl+Shift+1: CVNSS4.0   •   Ctrl+Shift+2: VNI/Telex   •   Ctrl+Shift+0: ứng viên")
\t}
\tif hwnd := a.controls[idToggle]; hwnd != 0 {
\t\tsetWindowText(hwnd, "BẬT/TẮT")
\t}
\ta.syncMethodCombo()
\ta.updateTrayIcon()
}'''
replace_exact("internal/win/app_windows.go", old_refresh, new_refresh)

start = read("internal/win/app_windows.go")
old_about_start = start.index("func (a *App) showAbout() {")
old_about_end = start.index("\nfunc (a *App) showSettingsInfo()", old_about_start)
new_about = r'''func (a *App) showAbout() {
	text := "BilaKey PC " + settings.AppVersion + "\r\n" +
		"Kiểu gõ dấu chữ tiếng Việt\r\n\r\n" +
		"CVNSS4.0\r\n" +
		"• Kiểu gõ chuyên dụng lấy CVNSS4.0 làm lõi trung tâm.\r\n" +
		"• Có candidate graph, resolver và kiểm toán quy tắc.\r\n\r\n" +
		"VNI + Telex hợp nhất\r\n" +
		"• Gõ theo VNI hoặc Telex trong cùng một chế độ.\r\n" +
		"• Hệ thống tự nhận dạng và xuất chữ tiếng Việt Unicode.\r\n\r\n" +
		"Tác giả: Long Ngo phát triển\r\n" +
		"Dự án: CVNSS4.0\r\n" +
		"Hỗ trợ CVNSS4.0: NNC Trần Tư Bình và cộng đồng.\r\n" +
		"Giấy phép: MIT\r\n\r\n" +
		"Riêng tư: không telemetry, không network runtime."
	messageBox(a.mainHwnd, "THÔNG TIN · BilaKey PC", text, MB_OK|MB_ICONINFORMATION)
}
'''
write("internal/win/app_windows.go", start[:old_about_start] + new_about + start[old_about_end + 1:])

readme = r'''<div align="center">
  <img src="assets/brand/bilakey-logo.svg" alt="Logo BilaKey — chữ B trắng trên nền xanh" width="168" />

# 🌊 BilaKey PC 2.5.5

### Bộ gõ Windows với **CVNSS4.0 Core** và **VNI/Telex hợp nhất**

<a href="https://github.com/xulytiengviet/BilaKey_PC/releases/download/2.5.5/BilaKey-PC-2.5.5-CVNSS-Core-x64.exe">
  <img src="https://img.shields.io/badge/T%E1%BA%A2I_NGAY-Windows_x64-0756d8?style=for-the-badge&logo=windows&logoColor=white" alt="Tải BilaKey PC 2.5.5 Windows x64" />
</a>

**Tải `.exe`, mở và dùng ngay — portable, không cần cài đặt.**

[![ARM64](https://img.shields.io/badge/T%E1%BA%A3i-ARM64%20%7C%20Snapdragon-173ea5?style=flat-square&logo=windows)](https://github.com/xulytiengviet/BilaKey_PC/releases/download/2.5.5/BilaKey-PC-2.5.5-CVNSS-Core-arm64.exe)
[![x86](https://img.shields.io/badge/T%E1%BA%A3i-Windows%2032--bit-3158a8?style=flat-square&logo=windows)](https://github.com/xulytiengviet/BilaKey_PC/releases/download/2.5.5/BilaKey-PC-2.5.5-CVNSS-Core-x86.exe)
[![ZIP](https://img.shields.io/badge/T%E1%BA%A3i-G%C3%B3i%20Windows%20%C4%91%E1%BA%A7y%20%C4%91%E1%BB%A7-5427c7?style=flat-square&logo=github)](https://github.com/xulytiengviet/BilaKey_PC/releases/download/2.5.5/BilaKey-PC-2.5.5-Windows.zip)
[![SHA](https://img.shields.io/badge/Ki%E1%BB%83m_tra-SHA--256-6b7280?style=flat-square&logo=github)](https://github.com/xulytiengviet/BilaKey_PC/releases/download/2.5.5/SHA256SUMS-2.5.5.txt)
[![Release](https://img.shields.io/badge/Xem-Release%202.5.5-00a86b?style=flat-square&logo=github)](https://github.com/xulytiengviet/BilaKey_PC/releases/tag/2.5.5)

[![Version](https://img.shields.io/badge/version-2.5.5-0756d8?style=for-the-badge)](VERSION)
[![License](https://img.shields.io/badge/license-MIT-00a86b?style=for-the-badge)](LICENSE)
[![Platform](https://img.shields.io/badge/Windows-x86%20%7C%20x64%20%7C%20ARM64-0078d4?style=for-the-badge&logo=windows)](BUILD.md)
[![Privacy](https://img.shields.io/badge/telemetry-none-5427c7?style=for-the-badge)](SECURITY.md)

**Nhanh · Nhẹ · Unicode · Offline · Riêng tư · Kiểm toán được · Mã nguồn mở MIT**
</div>

---

## ✨ Giao diện mới trong 2.5.5

BilaKey PC 2.5.5 nâng cấp toàn diện lớp UI/UX Win32:

- cửa sổ chính **1000 × 650 px**, rộng và dễ đọc hơn, tương đương khoảng một phần tư màn hình Full HD;
- logo chữ **B**, tên sản phẩm và định vị “Kiểu gõ dấu chữ tiếng Việt” được đặt ở vùng nhận diện riêng;
- hai thẻ kiểu gõ lớn, nhấn trực tiếp: **CVNSS4.0 · LÕI** và **VNI/TELEX · TỰ ĐỘNG**;
- thanh trạng thái riêng cho bật/tắt, kiểu gõ hiện hành và Unicode;
- thanh phím tắt riêng, không còn chen vào trạng thái;
- khu vực **THÔNG TIN** giải thích hai kiểu gõ và ghi rõ tác giả, dự án, giấy phép;
- năm nút thao tác lớn: **BẬT/TẮT · MỞ RỘNG · THÔNG TIN · HƯỚNG DẪN · THOÁT**;
- bảng màu xanh–trắng, chữ Segoe UI lớn hơn và độ tương phản tốt hơn.

Dòng nhận diện trong ứng dụng:

> **CVNSS4.0 Core, Unicode, riêng tư, kiểm toán được**  
> **Kiểu gõ dấu chữ tiếng Việt**

## ⌨️ Hai kiểu gõ

| Kiểu gõ | Giải thích |
|---|---|
| 🧠 **CVNSS4.0** | Kiểu gõ chuyên dụng lấy CVNSS4.0 làm lõi trung tâm, có candidate graph và resolver kiểm toán được. |
| 🔁 **VNI/Telex hợp nhất** | Gõ theo VNI hoặc Telex trong cùng một chế độ; hệ thống tự nhận dạng và xuất Unicode tiếng Việt. |

```text
Telex       tieengs   → tiếng
VNI         tieng61   → tiếng
Kết hợp     vieet5    → việt
Kết hợp     d9oongf   → đồng
CVNSS4.0    qyl       → quỳ
```

Cấu hình cũ `Telex`, `VNI` hoặc `Telex/VNI` tiếp tục được tự động chuyển sang `VNI/Telex`.

## 📦 Tải trực tiếp

| Thiết bị | Liên kết |
|---|---|
| **Windows x64 — khuyến nghị** | [BilaKey-PC-2.5.5-CVNSS-Core-x64.exe](https://github.com/xulytiengviet/BilaKey_PC/releases/download/2.5.5/BilaKey-PC-2.5.5-CVNSS-Core-x64.exe) |
| **Windows ARM64 / Snapdragon** | [BilaKey-PC-2.5.5-CVNSS-Core-arm64.exe](https://github.com/xulytiengviet/BilaKey_PC/releases/download/2.5.5/BilaKey-PC-2.5.5-CVNSS-Core-arm64.exe) |
| **Windows x86 32-bit** | [BilaKey-PC-2.5.5-CVNSS-Core-x86.exe](https://github.com/xulytiengviet/BilaKey_PC/releases/download/2.5.5/BilaKey-PC-2.5.5-CVNSS-Core-x86.exe) |
| **Gói đầy đủ Windows ZIP** | [BilaKey-PC-2.5.5-Windows.zip](https://github.com/xulytiengviet/BilaKey_PC/releases/download/2.5.5/BilaKey-PC-2.5.5-Windows.zip) |
| **Bảng kiểm tra SHA-256** | [SHA256SUMS-2.5.5.txt](https://github.com/xulytiengviet/BilaKey_PC/releases/download/2.5.5/SHA256SUMS-2.5.5.txt) |

Workflow chỉ kết thúc khi xác minh đủ **3 EXE + ZIP + SHA-256**.

> Bản portable chưa ký Authenticode; Windows SmartScreen có thể cảnh báo ở lần chạy đầu. Chỉ tải từ Release chính thức và kiểm tra SHA-256.

## 🎛️ Phím tắt

| Phím | Tác vụ |
|---|---|
| `Ctrl+Shift+Space` | Bật/tắt BilaKey |
| `Ctrl+Shift+1` | Chọn CVNSS4.0 |
| `Ctrl+Shift+2` | Chọn VNI/Telex |
| `Ctrl+Shift+0` | Đổi ứng viên CVNSS4.0 |
| `Shift` một lần | Viết hoa một từ |
| `Shift` hai lần | Bật/tắt BilaCaps |

## 🔨 Build từ mã nguồn

```bash
git clone https://github.com/xulytiengviet/BilaKey_PC.git
cd BilaKey_PC
go test ./...
GO_BIN=go scripts/build_release.sh
```

Yêu cầu: Go 1.23+, Node.js 22+, Python 3.12+, `g++` và `xz`.

## 🔐 Quyền riêng tư

BilaKey không có telemetry, quảng cáo, tài khoản người dùng hoặc network runtime. Quy tắc được nhúng tĩnh; cấu hình được ghi atomic trong thư mục người dùng.

## 🤝 Ghi nhận

- **Tác giả và phát triển:** **Long Ngo**.
- **Dự án nền tảng:** **CVNSS4.0**.
- **Hỗ trợ CVNSS4.0:** **NNC Trần Tư Bình** và cộng đồng CVNSS4.0/BilaKey.
- **Cộng đồng:** [CVNSS4.0 và Bộ gõ BilaKey](https://www.facebook.com/groups/251479779599477).
- **Giấy phép:** **MIT**.

```text
Copyright (c) 2026 Long Ngo
```
'''
write("README.md", readme)

changelog = read("CHANGELOG.md")
if "## 2.5.5" not in changelog:
    marker = "# 📝 Changelog\n\n"
    entry = '''## 2.5.5 — UI/UX Edition · 2026-07-31

- Mở rộng cửa sổ chính lên 1000 × 650 px và tổ chức lại toàn bộ bố cục.
- Tạo vùng nhận diện logo B, BilaKey PC, thông điệp CVNSS4.0 Core/Unicode/riêng tư/kiểm toán.
- Tách hai kiểu gõ thành hai thẻ lớn: CVNSS4.0 và VNI/Telex tự động.
- Thêm thanh trạng thái, thanh phím tắt và khu vực THÔNG TIN riêng.
- Ghi rõ Long Ngo phát triển, Dự án CVNSS4.0 và giấy phép MIT trong GUI.
- Thêm nút THÔNG TIN và nâng độ tương phản, kích thước chữ, khoảng cách điều khiển.
- Duy trì engine hợp nhất VNI/Telex và CVNSS4.0 Core từ 2.5.0.

'''
    if not changelog.startswith(marker):
        raise RuntimeError("CHANGELOG.md: unexpected header")
    write("CHANGELOG.md", marker + entry + changelog[len(marker):])

write(
    "docs/RELEASE_NOTES_2.5.5.md",
    r'''# BilaKey PC 2.5.5 — UI/UX Edition

BilaKey PC 2.5.5 tập trung nâng trải nghiệm giao diện nhưng giữ nguyên hai lõi nhập liệu đã ổn định: **CVNSS4.0** và **VNI/Telex hợp nhất**.

## Điểm mới

- Cửa sổ chính 1000 × 650 px, bố cục thoáng và dễ đọc hơn.
- Header logo chữ B, tên BilaKey PC và mô tả “Kiểu gõ dấu chữ tiếng Việt”.
- Hai thẻ kiểu gõ lớn: CVNSS4.0 · LÕI và VNI/TELEX · TỰ ĐỘNG.
- Thanh trạng thái, phím tắt và khu vực THÔNG TIN được tách riêng.
- Hiển thị rõ: Long Ngo phát triển, Dự án CVNSS4.0, giấy phép MIT.
- Nút THÔNG TIN mới giải thích trực tiếp hai kiểu gõ.
- Bảng màu xanh–trắng, Segoe UI lớn hơn và tương phản tốt hơn.

## Tải đúng bản

- `BilaKey-PC-2.5.5-CVNSS-Core-x64.exe`: Intel/AMD Windows 64-bit, khuyến nghị.
- `BilaKey-PC-2.5.5-CVNSS-Core-arm64.exe`: Windows on ARM, Snapdragon.
- `BilaKey-PC-2.5.5-CVNSS-Core-x86.exe`: Windows 32-bit cũ.
- `BilaKey-PC-2.5.5-Windows.zip`: gói đầy đủ.
- `SHA256SUMS-2.5.5.txt`: bảng kiểm tra toàn vẹn.

Bản portable chưa ký Authenticode. Chỉ tải từ Release chính thức và kiểm tra SHA-256.
''',
)

write(
    "docs/UI_UX_2.5.5.md",
    r'''# UI/UX BilaKey PC 2.5.5

## Mục tiêu

Giao diện chính phải cho người dùng hiểu trong vài giây: đây là bộ gõ dấu chữ tiếng Việt, có hai kiểu gõ và đang chạy ở trạng thái nào.

## Phân cấp

1. Nhận diện: logo B, BilaKey PC, CVNSS4.0 Core/Unicode/riêng tư/kiểm toán.
2. Lựa chọn: CVNSS4.0 hoặc VNI/Telex hợp nhất.
3. Trạng thái: bật/tắt, kiểu hiện hành, Unicode.
4. Phím tắt: hiển thị trên một hàng riêng.
5. Thông tin: giải thích cách gõ, tác giả, dự án và giấy phép.
6. Hành động: bật/tắt, mở rộng, thông tin, hướng dẫn, thoát.

## Kích thước

Cửa sổ chính: 1000 × 650 px. Các nút kiểu gõ cao 84 px; nút thao tác cao 50 px. Font Segoe UI từ 14 đến 46 px theo vai trò.
''',
)

print("BilaKey PC 2.5.5 UI/UX migration prepared")
