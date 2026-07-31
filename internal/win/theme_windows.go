//go:build windows

package win

import "unsafe"

const (
	oceanBG     = uint32(0x00AB690C)
	oceanDeep   = uint32(0x007C4706)
	oceanPanel  = uint32(0x00C47B0B)
	oceanAccent = uint32(0x00E09B20)
	oceanWhite  = uint32(0x00FFFFFF)
	oceanPale   = uint32(0x00FFF2DC)
	oceanMuted  = uint32(0x00E8D2B6)
)

func (a *App) initTheme() {
	if a.theme.font != 0 {
		return
	}
	a.theme.bgBrush = createSolidBrush(oceanBG)
	a.theme.panelBrush = createSolidBrush(oceanPanel)
	a.theme.editBrush = createSolidBrush(oceanDeep)
	a.theme.cyanBrush = createSolidBrush(oceanAccent)
	a.theme.grayBrush = createSolidBrush(oceanMuted)
	a.theme.font = createUIFont(FW_NORMAL, 16)
	a.theme.fontBold = createUIFont(FW_BOLD, 18)
}

func (a *App) destroyTheme() {
	for _, h := range []uintptr{
		a.theme.font,
		a.theme.fontBold,
		a.theme.bgBrush,
		a.theme.panelBrush,
		a.theme.editBrush,
		a.theme.cyanBrush,
		a.theme.grayBrush,
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

func (a *App) applyControlTheme(hwnd uintptr, bold bool) {
	if hwnd == 0 {
		return
	}
	font := a.theme.font
	if bold {
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
	procSetTextColor.Call(hdc, uintptr(oceanWhite))

	switch msg {
	case WM_CTLCOLOREDIT, WM_CTLCOLORLISTBOX:
		procSetBkMode.Call(hdc, 2)
		procSetBkColor.Call(hdc, uintptr(oceanDeep))
		procSetTextColor.Call(hdc, uintptr(oceanWhite))
		return a.theme.editBrush
	case WM_CTLCOLORSTATIC:
		id, _, _ := procGetDlgCtrlID.Call(child)
		if int(id) == idSubtle {
			procSetTextColor.Call(hdc, uintptr(oceanPale))
		}
		return a.theme.bgBrush
	case WM_CTLCOLORBTN:
		procSetTextColor.Call(hdc, uintptr(oceanWhite))
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

	fill := a.theme.panelBrush
	border := a.theme.grayBrush
	textColor := oceanWhite
	if a.isMethodTab(dis.CtlID) && a.methodTabActive(dis.CtlID) {
		fill = a.theme.cyanBrush
		border = a.theme.cyanBrush
	}
	if dis.ItemState&ODS_SELECTED != 0 {
		fill = a.theme.cyanBrush
		border = a.theme.cyanBrush
	}
	if dis.ItemState&ODS_DISABLED != 0 {
		textColor = oceanMuted
	}

	procFillRect.Call(dis.HDC, uintptr(unsafe.Pointer(&dis.RcItem)), fill)
	procFrameRect.Call(dis.HDC, uintptr(unsafe.Pointer(&dis.RcItem)), border)
	procSetBkMode.Call(dis.HDC, TRANSPARENT)
	procSetTextColor.Call(dis.HDC, uintptr(textColor))
	if a.theme.fontBold != 0 {
		procSelectObject.Call(dis.HDC, a.theme.fontBold)
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
