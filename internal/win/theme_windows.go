//go:build windows

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
