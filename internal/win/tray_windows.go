//go:build windows

package win

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/xulytiengviet/BilaKey_PC/internal/settings"
)

const (
	trayCallbackMessage = WM_APP + 1
	methodCycleMessage  = WM_APP + 2
	hotkeyActionMessage = WM_APP + 3
	trayIconID          = 1

	trayCmdToggle = 9001
	trayCmdShow   = 9002
	trayCmdCVNSS  = 9003
	trayCmdTelex  = 9004
	trayCmdVNI    = 9005
	trayCmdExit   = 9006
)

func (a *App) initIcons() {
	if a.iconSmall == 0 {
		a.iconSmall = createBIcon(16)
	}
	if a.iconBig == 0 {
		a.iconBig = createBIcon(32)
	}
}

func (a *App) destroyIcons() {
	if a.iconSmall != 0 {
		procDestroyIcon.Call(a.iconSmall)
		a.iconSmall = 0
	}
	if a.iconBig != 0 {
		procDestroyIcon.Call(a.iconBig)
		a.iconBig = 0
	}
}

func createBIcon(size int) uintptr {
	if size < 8 {
		size = 8
	}
	rowBytes := size * 4
	xorBits := make([]byte, rowBytes*size)
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			setIconPixel(xorBits, size, x, y, 0xAB, 0x69, 0x0C, 0xFF)
		}
	}
	pattern := []string{"11110", "10001", "10001", "11110", "10001", "10001", "11110"}
	scale := size / 10
	if scale < 1 {
		scale = 1
	}
	glyphW := 5 * scale
	glyphH := 7 * scale
	x0 := (size - glyphW) / 2
	y0 := (size - glyphH) / 2
	for gy, row := range pattern {
		for gx, bit := range row {
			if bit != '1' {
				continue
			}
			for sy := 0; sy < scale; sy++ {
				for sx := 0; sx < scale; sx++ {
					setIconPixel(xorBits, size, x0+gx*scale+sx, y0+gy*scale+sy, 0xFF, 0xFF, 0xFF, 0xFF)
				}
			}
		}
	}
	maskStride := ((size + 31) / 32) * 4
	andBits := make([]byte, maskStride*size)
	hinst, _, _ := procGetModuleHandleW.Call(0)
	h, _, _ := procCreateIcon.Call(hinst, uintptr(size), uintptr(size), 1, 32, uintptr(unsafe.Pointer(&andBits[0])), uintptr(unsafe.Pointer(&xorBits[0])))
	return h
}

func setIconPixel(buf []byte, size, x, y int, b, g, r, a byte) {
	if x < 0 || y < 0 || x >= size || y >= size {
		return
	}
	row := size - 1 - y
	i := (row*size + x) * 4
	buf[i+0] = b
	buf[i+1] = g
	buf[i+2] = r
	buf[i+3] = a
}

func (a *App) installTrayIcon() error {
	if a.mainHwnd == 0 || a.iconSmall == 0 {
		return fmt.Errorf("main window/icon chưa sẵn sàng")
	}
	nid := a.makeTrayData()
	r, _, err := procShellNotifyIconW.Call(NIM_ADD, uintptr(unsafe.Pointer(&nid)))
	if r == 0 {
		return fmt.Errorf("Shell_NotifyIcon(NIM_ADD): %v", err)
	}
	a.trayAdded = true
	return nil
}

func (a *App) removeTrayIcon() {
	if !a.trayAdded || a.mainHwnd == 0 {
		return
	}
	nid := a.makeTrayData()
	procShellNotifyIconW.Call(NIM_DELETE, uintptr(unsafe.Pointer(&nid)))
	a.trayAdded = false
}

func (a *App) updateTrayIcon() {
	if !a.trayAdded || a.mainHwnd == 0 {
		return
	}
	nid := a.makeTrayData()
	procShellNotifyIconW.Call(NIM_MODIFY, uintptr(unsafe.Pointer(&nid)))
}

func (a *App) makeTrayData() notifyIconData {
	cfg := a.runtimeConfig()
	state := "ON"
	if !cfg.Enabled {
		state = "OFF"
	}
	tip := fmt.Sprintf("BilaKey PC %s · %s · %s · %s", settings.AppVersion, state, cfg.InputMethod, a.capsModeLabel())
	nid := notifyIconData{HWnd: a.mainHwnd, UID: trayIconID, UFlags: NIF_MESSAGE | NIF_ICON | NIF_TIP, UCallbackMessage: trayCallbackMessage, HIcon: a.iconSmall}
	nid.CbSize = uint32(unsafe.Sizeof(nid))
	copyUTF16Fixed(nid.SzTip[:], tip)
	return nid
}

func copyUTF16Fixed(dst []uint16, s string) {
	if len(dst) == 0 {
		return
	}
	u := syscall.StringToUTF16(s)
	if len(u) > len(dst) {
		u = u[:len(dst)]
		u[len(u)-1] = 0
	}
	copy(dst, u)
}

func (a *App) handleTrayMessage(lParam uintptr) {
	switch uint32(lParam & 0xffff) {
	case WM_LBUTTONUP, WM_LBUTTONDBLCLK:
		a.showMainWindow()
	case WM_RBUTTONUP:
		a.showTrayMenu()
	}
}

func (a *App) showTrayMenu() {
	menu, _, _ := procCreatePopupMenu.Call()
	if menu == 0 {
		return
	}
	defer procDestroyMenu.Call(menu)

	cfg := a.runtimeConfig()
	toggleText := "Tắt BilaKey"
	if !cfg.Enabled {
		toggleText = "Bật BilaKey"
	}
	appendMenuText(menu, MF_STRING, trayCmdToggle, toggleText)
	appendMenuText(menu, MF_STRING, trayCmdShow, "Mở BilaKey PC")
	procAppendMenuW.Call(menu, MF_SEPARATOR, 0, 0)
	appendMenuChecked(menu, trayCmdCVNSS, "CVNSS4.0 · Lõi", cfg.InputMethod == "CVNSS4.0")
	appendMenuChecked(menu, trayCmdTelex, "Telex · tương thích", cfg.InputMethod == "Telex")
	appendMenuChecked(menu, trayCmdVNI, "VNI · tương thích", cfg.InputMethod == "VNI")
	procAppendMenuW.Call(menu, MF_SEPARATOR, 0, 0)
	appendMenuText(menu, MF_STRING, trayCmdExit, "Thoát")

	var p point
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&p)))
	procSetForegroundWindow.Call(a.mainHwnd)
	cmd, _, _ := procTrackPopupMenu.Call(menu, TPM_RIGHTBUTTON|TPM_RETURNCMD, uintptr(p.X), uintptr(p.Y), 0, a.mainHwnd, 0)
	procPostMessageW.Call(a.mainHwnd, 0, 0, 0)

	switch int(cmd) {
	case trayCmdToggle:
		a.updateConfig(func(c *settings.Config) { c.Enabled = !c.Enabled })
	case trayCmdShow:
		a.showMainWindow()
	case trayCmdCVNSS:
		a.updateConfig(func(c *settings.Config) { c.InputMethod = "CVNSS4.0" })
		a.syncMethodCombo()
	case trayCmdTelex:
		a.updateConfig(func(c *settings.Config) { c.InputMethod = "Telex" })
		a.syncMethodCombo()
	case trayCmdVNI:
		a.updateConfig(func(c *settings.Config) { c.InputMethod = "VNI" })
		a.syncMethodCombo()
	case trayCmdExit:
		a.exitApp()
	}
}

func appendMenuText(menu uintptr, flags uintptr, id int, text string) {
	procAppendMenuW.Call(menu, flags, uintptr(id), uintptr(unsafe.Pointer(utf16Ptr(text))))
}

func appendMenuChecked(menu uintptr, id int, text string, checked bool) {
	flags := uintptr(MF_STRING)
	if checked {
		flags |= MF_CHECKED
	}
	appendMenuText(menu, flags, id, text)
}
