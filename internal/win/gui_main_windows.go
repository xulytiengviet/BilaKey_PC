//go:build windows

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
	idHeader       = 119 // shared title style for options and macro windows

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
