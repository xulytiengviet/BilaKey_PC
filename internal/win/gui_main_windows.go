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
	idToggle     = 101
	idExit       = 102
	idExpand     = 103
	idHelp       = 104
	idAbout      = 105
	idSettings   = 106
	idStatus     = 107
	idCapsStatus = 108
	idHeader     = 109
	idSubtle     = 110
	idTabCVNSS   = 111
	idTabTelex   = 112
	idTabVNI     = 113

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
		currentApp.applyControlTheme(hwnd, id == idHeader)
	}
	return hwnd
}

func createDOSButton(parent uintptr, text string, x, y, w, h int32, id int) uintptr {
	return createControl(parent, "BUTTON", text, BS_OWNERDRAW|WS_TABSTOP, x, y, w, h, id)
}

func (a *App) createMainWindow() error {
	hwnd := createWindow(classMain, "BilaKey PC "+settings.AppVersion, WS_OVERLAPPED|WS_CAPTION|WS_SYSMENU|WS_MINIMIZEBOX, 260, 150, 600, 360, 0, 0)
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
		createControl(hwnd, "STATIC", "B", WS_BORDER, 24, 20, 48, 48, idHeader)
		createControl(hwnd, "STATIC", "BilaKey PC", 0, 88, 18, 220, 28, idHeader)
		createControl(hwnd, "STATIC", "CVNSS4.0 Core · Unicode · riêng tư · kiểm toán được", 0, 88, 48, 440, 22, idSubtle)

		createControl(hwnd, "STATIC", "Lõi nhập liệu và adapter tương thích", 0, 24, 88, 310, 22, idSubtle)
		a.controls[idTabCVNSS] = createDOSButton(hwnd, "CVNSS4.0 · LÕI", 24, 114, 290, 42, idTabCVNSS)
		a.controls[idTabTelex] = createDOSButton(hwnd, "TELEX", 324, 114, 110, 42, idTabTelex)
		a.controls[idTabVNI] = createDOSButton(hwnd, "VNI", 444, 114, 110, 42, idTabVNI)

		status := createControl(hwnd, "STATIC", "", WS_BORDER, 24, 176, 530, 34, idStatus)
		a.controls[idStatus] = status
		caps := createControl(hwnd, "STATIC", "", 0, 24, 216, 530, 24, idCapsStatus)
		a.controls[idCapsStatus] = caps

		toggle := createDOSButton(hwnd, "TẮT", 24, 258, 110, 42, idToggle)
		a.controls[idToggle] = toggle
		createDOSButton(hwnd, "MỞ RỘNG", 144, 258, 130, 42, idExpand)
		createDOSButton(hwnd, "HƯỚNG DẪN", 284, 258, 130, 42, idHelp)
		createDOSButton(hwnd, "THOÁT", 424, 258, 130, 42, idExit)

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
			case idTabTelex:
				a.selectInputMethod(core.MethodTelex)
			case idTabVNI:
				a.selectInputMethod(core.MethodVNI)
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
		// Plain Tab cycles the three tabs only while the BilaKey window is active.
		// Global switching uses structured Ctrl+Shift shortcuts; Tab/Ctrl+Tab in
		// editors and forms is never stolen.
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
	method := a.runtimeConfig().InputMethod
	switch method {
	case core.MethodCVNSS:
		method = core.MethodTelex
	case core.MethodTelex:
		method = core.MethodVNI
	default:
		method = core.MethodCVNSS
	}
	a.selectInputMethod(method)
}

// Kept under the historical name to minimize call-site churn. The combo box
// was replaced by three direct, owner-drawn method tabs.
func (a *App) syncMethodCombo() {
	for _, id := range []int{idTabCVNSS, idTabTelex, idTabVNI} {
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
	case hotkey.SelectTelex:
		a.selectInputMethod(core.MethodTelex)
	case hotkey.SelectVNI:
		a.selectInputMethod(core.MethodVNI)
	case hotkey.CycleCandidate:
		a.cycleCVNSSCandidate()
	}
}
