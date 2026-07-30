//go:build windows

package win

import (
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	"github.com/xulytiengviet/BilaKey_PC/internal/settings"
)

func (a *App) showMacroEditor() {
	if a.macroHwnd != 0 {
		procShowWindow.Call(a.macroHwnd, SW_SHOW)
		procSetFocus.Call(a.macroHwnd)
		return
	}
	hwnd := createWindow(classMacro, "Bảng gõ tắt · BilaKey PC "+settings.AppVersion, WS_OVERLAPPED|WS_CAPTION|WS_SYSMENU, 300, 100, 820, 610, a.optionsHwnd, 0)
	if hwnd == 0 {
		messageBox(a.mainHwnd, "Lỗi", "Không thể mở bảng gõ tắt.", MB_OK|MB_ICONWARNING)
		return
	}
	a.macroHwnd = hwnd
	if a.iconSmall != 0 {
		sendMessage(hwnd, WM_SETICON, ICON_SMALL, a.iconSmall)
	}
	procShowWindow.Call(hwnd, SW_SHOW)
	procUpdateWindow.Call(hwnd)
}

func macroWndProc(hwnd uintptr, msg uint32, wParam, lParam uintptr) uintptr {
	a := currentApp
	if a == nil {
		r, _, _ := procDefWindowProcW.Call(hwnd, uintptr(msg), wParam, lParam)
		return r
	}
	switch msg {
	case WM_CREATE:
		createControl(hwnd, "STATIC", "[ USER MACRO TABLE · TEEN / BÙI HIỀN / CUSTOM ]", 0, 24, 15, 730, 24, idHeader)
		createControl(hwnd, "STATIC", "Thay thế:", 0, 24, 50, 140, 24, 0)
		createControl(hwnd, "STATIC", "Bởi:", 0, 190, 50, 80, 24, 0)
		a.macroControls[idMacroTrigger] = createControl(hwnd, "EDIT", "", WS_BORDER|ES_AUTOHSCROLL|WS_TABSTOP, 24, 78, 145, 31, idMacroTrigger)
		a.macroControls[idMacroReplacement] = createControl(hwnd, "EDIT", "", WS_BORDER|ES_AUTOHSCROLL|WS_TABSTOP, 190, 78, 400, 31, idMacroReplacement)
		a.macroControls[idMacroList] = createControl(hwnd, "LISTBOX", "", WS_BORDER|WS_VSCROLL|LBS_NOTIFY|WS_TABSTOP, 24, 130, 566, 285, idMacroList)
		createDOSButton(hwnd, "LƯU", 615, 130, 140, 40, idMacroSave)
		createDOSButton(hwnd, "+ THÊM", 615, 185, 140, 40, idMacroAdd)
		createDOSButton(hwnd, "- XÓA", 615, 240, 140, 40, idMacroDelete)
		createDOSButton(hwnd, "ĐÓNG", 615, 375, 140, 40, idMacroClose)
		createControl(hwnd, "STATIC", "File gõ tắt:", 0, 24, 440, 110, 24, 0)
		pathH := createControl(hwnd, "STATIC", a.macros.Path(), WS_BORDER, 140, 436, 450, 28, idMacroPath)
		a.macroControls[idMacroPath] = pathH
		createDOSButton(hwnd, "CHỌN FILE...", 24, 485, 160, 38, idMacroChoose)
		createDOSButton(hwnd, "FILE MẶC ĐỊNH", 198, 485, 175, 38, idMacroDefault)
		createControl(hwnd, "STATIC", "Định dạng: trigger<TAB>replacement · dữ liệu local · không network.", 0, 395, 492, 370, 35, 0)
		a.refreshMacroList()
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
		if id == idMacroList && notify == LBN_SELCHANGE {
			sel := int(sendMessage(a.macroControls[idMacroList], LB_GETCURSEL, 0, 0))
			entries := a.macros.Entries()
			if sel >= 0 && sel < len(entries) {
				setWindowText(a.macroControls[idMacroTrigger], entries[sel].Trigger)
				setWindowText(a.macroControls[idMacroReplacement], entries[sel].Replacement)
			}
			return 0
		}
		if notify == BN_CLICKED {
			switch id {
			case idMacroAdd:
				trigger := strings.TrimSpace(getWindowText(a.macroControls[idMacroTrigger]))
				repl := getWindowText(a.macroControls[idMacroReplacement])
				if trigger == "" {
					messageBox(hwnd, "Bảng gõ tắt", "Vui lòng nhập chuỗi 'Thay thế'.", MB_OK|MB_ICONWARNING)
					return 0
				}
				a.macros.Upsert(trigger, repl)
				a.refreshMacroList()
			case idMacroDelete:
				trigger := strings.TrimSpace(getWindowText(a.macroControls[idMacroTrigger]))
				if trigger != "" {
					a.macros.Delete(trigger)
					a.refreshMacroList()
				}
			case idMacroSave:
				trigger := strings.TrimSpace(getWindowText(a.macroControls[idMacroTrigger]))
				if trigger != "" {
					a.macros.Upsert(trigger, getWindowText(a.macroControls[idMacroReplacement]))
				}
				if err := a.macros.Save(); err != nil {
					messageBox(hwnd, "Lỗi lưu bảng gõ tắt", err.Error(), MB_OK|MB_ICONWARNING)
				} else {
					_ = a.store.Update(func(c *settings.Config) { c.MacroFile = a.macros.Path() })
					messageBox(hwnd, "Bảng gõ tắt", "Đã lưu bảng gõ tắt.", MB_OK|MB_ICONINFORMATION)
				}
				a.refreshMacroList()
			case idMacroChoose:
				if path, ok := chooseMacroFile(hwnd, a.macros.Path()); ok {
					a.macros.SetPath(path)
					if err := a.macros.Load(); err != nil {
						messageBox(hwnd, "Lỗi đọc file", err.Error(), MB_OK|MB_ICONWARNING)
					} else {
						_ = a.store.Update(func(c *settings.Config) { c.MacroFile = path })
						setWindowText(a.macroControls[idMacroPath], path)
						a.refreshMacroList()
					}
				}
			case idMacroDefault:
				path, err := settings.DefaultMacroPath()
				if err == nil {
					a.macros.SetPath(path)
					_ = a.macros.Load()
					_ = a.store.Update(func(c *settings.Config) { c.MacroFile = path })
					setWindowText(a.macroControls[idMacroPath], path)
					a.refreshMacroList()
				}
			case idMacroClose:
				procDestroyWindow.Call(hwnd)
			}
			return 0
		}
	case WM_CLOSE:
		procDestroyWindow.Call(hwnd)
		return 0
	case WM_DESTROY:
		a.macroHwnd = 0
		return 0
	}
	r, _, _ := procDefWindowProcW.Call(hwnd, uintptr(msg), wParam, lParam)
	return r
}

func (a *App) refreshMacroList() {
	list := a.macroControls[idMacroList]
	if list == 0 {
		return
	}
	sendMessage(list, LB_RESETCONTENT, 0, 0)
	for _, e := range a.macros.Entries() {
		text := e.Trigger + "   →   " + e.Replacement
		sendMessage(list, LB_ADDSTRING, 0, uintptr(unsafe.Pointer(utf16Ptr(text))))
	}
	if p := a.macroControls[idMacroPath]; p != 0 {
		setWindowText(p, a.macros.Path())
	}
}

func chooseMacroFile(owner uintptr, current string) (string, bool) {
	buf := make([]uint16, 4096)
	if current != "" {
		copy(buf, syscall.StringToUTF16(current))
	}
	filter := syscall.StringToUTF16("Bảng gõ tắt TSV/TXT\x00*.tsv;*.txt\x00Tất cả tệp\x00*.*\x00\x00")
	title := utf16Ptr("Chọn bảng gõ tắt BilaKey PC")
	of := openFileName{
		LStructSize: uint32(unsafe.Sizeof(openFileName{})),
		HwndOwner:   owner,
		LpstrFilter: &filter[0],
		LpstrFile:   &buf[0],
		NMaxFile:    uint32(len(buf)),
		LpstrTitle:  title,
		Flags:       OFN_PATHMUSTEXIST | OFN_FILEMUSTEXIST | OFN_EXPLORER,
	}
	r, _, _ := procGetOpenFileNameW.Call(uintptr(unsafe.Pointer(&of)))
	if r == 0 {
		return "", false
	}
	path := syscall.UTF16ToString(buf)
	if path == "" {
		return "", false
	}
	return filepath.Clean(path), true
}
