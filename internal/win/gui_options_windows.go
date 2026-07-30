//go:build windows

package win

import "github.com/xulytiengviet/BilaKey_PC/internal/settings"

func (a *App) showOptions() {
	if a.optionsHwnd != 0 {
		procShowWindow.Call(a.optionsHwnd, SW_SHOW)
		procSetFocus.Call(a.optionsHwnd)
		return
	}
	hwnd := createWindow(classOptions, "Mở rộng · BilaKey PC "+settings.AppVersion, WS_OVERLAPPED|WS_CAPTION|WS_SYSMENU, 250, 80, 920, 680, a.mainHwnd, 0)
	if hwnd == 0 {
		messageBox(a.mainHwnd, "Lỗi", "Không thể mở cửa sổ tùy chọn.", MB_OK|MB_ICONWARNING)
		return
	}
	a.optionsHwnd = hwnd
	if a.iconSmall != 0 {
		sendMessage(hwnd, WM_SETICON, ICON_SMALL, a.iconSmall)
	}
	procShowWindow.Call(hwnd, SW_SHOW)
	procUpdateWindow.Call(hwnd)
}

func checkHandle(parent uintptr, text string, x, y, w int32, id int, checked bool) uintptr {
	h := createControl(parent, "BUTTON", text, BS_AUTOCHECKBOX|WS_TABSTOP, x, y, w, 25, id)
	if checked {
		sendMessage(h, BM_SETCHECK, BST_CHECKED, 0)
	}
	return h
}

func isChecked(hwnd uintptr) bool { return sendMessage(hwnd, BM_GETCHECK, 0, 0) == BST_CHECKED }

func optionsWndProc(hwnd uintptr, msg uint32, wParam, lParam uintptr) uintptr {
	a := currentApp
	if a == nil {
		r, _, _ := procDefWindowProcW.Call(hwnd, uintptr(msg), wParam, lParam)
		return r
	}
	switch msg {
	case WM_CREATE:
		cfg := a.runtimeConfig()
		createControl(hwnd, "STATIC", "Tùy chọn mở rộng · BilaKey PC "+settings.AppVersion, 0, 20, 12, 850, 24, idHeader)

		createControl(hwnd, "BUTTON", "Tùy chọn khác", BS_GROUPBOX, 20, 42, 850, 165, 0)
		a.controls[idOptFreeTone] = checkHandle(hwnd, "Cho phép gõ tự do", 40, 72, 350, idOptFreeTone, cfg.FreeToneMarking)
		a.controls[idOptOldTone] = checkHandle(hwnd, "Đặt dấu oà, uý (thay vì òa, úy)", 40, 106, 380, idOptOldTone, cfg.OldToneStyle)
		a.controls[idOptClipboard] = checkHandle(hwnd, "Tương thích clipboard (ghi đè nội dung hiện tại)", 40, 140, 405, idOptClipboard, cfg.AlwaysUseClipboardUnicode)
		a.controls[idOptSpell] = checkHandle(hwnd, "Bật kiểm tra chính tả", 470, 72, 340, idOptSpell, cfg.SpellCheck)
		a.controls[idOptRestore] = checkHandle(hwnd, "Tự động khôi phục phím với từ sai", 470, 106, 360, idOptRestore, cfg.AutoRestoreWrongKey)
		a.controls[idOptFeedback] = checkHandle(hwnd, "Hiện thông báo phản hồi", 470, 140, 340, idOptFeedback, cfg.ShowFeedback)

		createControl(hwnd, "BUTTON", "Viết hoa thông minh", BS_GROUPBOX, 20, 216, 850, 145, 0)
		a.controls[idOptAutoCapInitial] = checkHandle(hwnd, "Tự động viết hoa chữ đầu tiên", 40, 246, 370, idOptAutoCapInitial, cfg.AutoCapInitial)
		a.controls[idOptAutoCapSentence] = checkHandle(hwnd, "Tự động viết hoa sau . ! ? và xuống dòng", 40, 282, 390, idOptAutoCapSentence, cfg.AutoCapSentence)
		a.controls[idOptDoubleShiftCaps] = checkHandle(hwnd, "SHIFT ×2: khóa viết hoa liên tục (BilaCaps)", 470, 246, 360, idOptDoubleShiftCaps, cfg.DoubleShiftCaps)
		a.controls[idOptRestoreDelimiter] = checkHandle(hwnd, "Backspace sau dấu cách: quay lại từ để sửa", 470, 282, 370, idOptRestoreDelimiter, cfg.RestoreAfterDelimiter)
		createControl(hwnd, "STATIC", "SHIFT ×1 luôn là one-shot: viết hoa một chữ rồi tự trở lại chữ thường.", 0, 40, 322, 780, 22, 0)

		createControl(hwnd, "BUTTON", "Tùy chọn gõ tắt", BS_GROUPBOX, 20, 375, 410, 175, 0)
		a.controls[idOptMacro] = checkHandle(hwnd, "Cho phép gõ tắt", 40, 408, 330, idOptMacro, cfg.MacroEnabled)
		a.controls[idOptMacroOff] = checkHandle(hwnd, "Cho phép gõ tắt cả khi tắt tiếng Việt", 40, 444, 350, idOptMacroOff, cfg.MacroWhileOff)
		createDOSButton(hwnd, "BẢNG GÕ TẮT...", 45, 490, 190, 38, idOptMacroTable)

		createControl(hwnd, "BUTTON", "Hệ thống", BS_GROUPBOX, 450, 375, 420, 175, 0)
		a.controls[idOptShowStartup] = checkHandle(hwnd, "Bật hội thoại này khi khởi động", 470, 408, 350, idOptShowStartup, cfg.ShowDialogAtStartup)
		a.controls[idOptStartWindows] = checkHandle(hwnd, "Khởi động cùng Windows", 470, 444, 330, idOptStartWindows, cfg.StartWithWindows)
		a.controls[idOptVietnamese] = checkHandle(hwnd, "Giao diện tiếng Việt", 470, 480, 330, idOptVietnamese, cfg.VietnameseInterface)
		procEnableWindow.Call(a.controls[idOptVietnamese], 0)
		a.controls[idOptPausePassword] = checkHandle(hwnd, "Tạm dừng trong ô mật khẩu Win32", 470, 516, 350, idOptPausePassword, cfg.PauseInPasswordFields)

		createDOSButton(hwnd, "LƯU", 655, 575, 100, 40, idOptSave)
		createDOSButton(hwnd, "ĐÓNG", 770, 575, 100, 40, idOptClose)
		return 0

	case WM_DRAWITEM:
		if a.drawDOSButton(lParam) {
			return 1
		}
	case WM_CTLCOLORSTATIC, WM_CTLCOLORBTN, WM_CTLCOLOREDIT, WM_CTLCOLORLISTBOX:
		return a.themeControlColor(msg, wParam, lParam)

	case WM_COMMAND:
		id := int(loword(wParam))
		if hiword(wParam) == BN_CLICKED {
			switch id {
			case idOptMacroTable:
				a.showMacroEditor()
			case idOptClose:
				procDestroyWindow.Call(hwnd)
			case idOptSave:
				cfg := a.runtimeConfig()
				cfg.FreeToneMarking = isChecked(a.controls[idOptFreeTone])
				cfg.OldToneStyle = isChecked(a.controls[idOptOldTone])
				cfg.AlwaysUseClipboardUnicode = isChecked(a.controls[idOptClipboard])
				cfg.SpellCheck = isChecked(a.controls[idOptSpell])
				cfg.AutoRestoreWrongKey = isChecked(a.controls[idOptRestore])
				cfg.ShowFeedback = isChecked(a.controls[idOptFeedback])
				cfg.MacroEnabled = isChecked(a.controls[idOptMacro])
				cfg.MacroWhileOff = isChecked(a.controls[idOptMacroOff])
				cfg.ShowDialogAtStartup = isChecked(a.controls[idOptShowStartup])
				cfg.StartWithWindows = isChecked(a.controls[idOptStartWindows])
				cfg.AutoCapInitial = isChecked(a.controls[idOptAutoCapInitial])
				cfg.AutoCapSentence = isChecked(a.controls[idOptAutoCapSentence])
				cfg.DoubleShiftCaps = isChecked(a.controls[idOptDoubleShiftCaps])
				cfg.RestoreAfterDelimiter = isChecked(a.controls[idOptRestoreDelimiter])
				cfg.PauseInPasswordFields = isChecked(a.controls[idOptPausePassword])
				if err := a.store.Replace(cfg); err != nil {
					messageBox(hwnd, "Lỗi", err.Error(), MB_OK|MB_ICONWARNING)
					return 0
				}
				a.rebuildEngine()
				a.resetTypingState()
				if err := setStartupWithWindows(cfg.StartWithWindows); err != nil {
					messageBox(hwnd, "Khởi động cùng Windows", err.Error(), MB_OK|MB_ICONWARNING)
				} else {
					messageBox(hwnd, "BilaKey PC", "Đã lưu cấu hình và reset state machine viết hoa.", MB_OK|MB_ICONINFORMATION)
				}
				a.refreshMainStatus()
			}
			return 0
		}
	case WM_CLOSE:
		procDestroyWindow.Call(hwnd)
		return 0
	case WM_DESTROY:
		a.optionsHwnd = 0
		return 0
	}
	r, _, _ := procDefWindowProcW.Call(hwnd, uintptr(msg), wParam, lParam)
	return r
}
