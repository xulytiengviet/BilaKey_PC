//go:build windows

package win

import (
	"syscall"
	"unsafe"
)

var (
	user32   = syscall.NewLazyDLL("user32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	gdi32    = syscall.NewLazyDLL("gdi32.dll")
	shell32  = syscall.NewLazyDLL("shell32.dll")
	comdlg32 = syscall.NewLazyDLL("comdlg32.dll")
	advapi32 = syscall.NewLazyDLL("advapi32.dll")

	procRegisterClassExW    = user32.NewProc("RegisterClassExW")
	procCreateWindowExW     = user32.NewProc("CreateWindowExW")
	procDefWindowProcW      = user32.NewProc("DefWindowProcW")
	procShowWindow          = user32.NewProc("ShowWindow")
	procUpdateWindow        = user32.NewProc("UpdateWindow")
	procGetMessageW         = user32.NewProc("GetMessageW")
	procTranslateMessage    = user32.NewProc("TranslateMessage")
	procDispatchMessageW    = user32.NewProc("DispatchMessageW")
	procPostQuitMessage     = user32.NewProc("PostQuitMessage")
	procDestroyWindow       = user32.NewProc("DestroyWindow")
	procSendMessageW        = user32.NewProc("SendMessageW")
	procSetWindowTextW      = user32.NewProc("SetWindowTextW")
	procGetWindowTextW      = user32.NewProc("GetWindowTextW")
	procGetWindowTextLen    = user32.NewProc("GetWindowTextLengthW")
	procMessageBoxW         = user32.NewProc("MessageBoxW")
	procEnableWindow        = user32.NewProc("EnableWindow")
	procSetFocus            = user32.NewProc("SetFocus")
	procGetKeyState         = user32.NewProc("GetKeyState")
	procGetAsyncKeyState    = user32.NewProc("GetAsyncKeyState")
	procSetWindowsHookEx    = user32.NewProc("SetWindowsHookExW")
	procUnhookWindowsHook   = user32.NewProc("UnhookWindowsHookEx")
	procCallNextHookEx      = user32.NewProc("CallNextHookEx")
	procSendInput           = user32.NewProc("SendInput")
	procOpenClipboard       = user32.NewProc("OpenClipboard")
	procCloseClipboard      = user32.NewProc("CloseClipboard")
	procEmptyClipboard      = user32.NewProc("EmptyClipboard")
	procSetClipboardData    = user32.NewProc("SetClipboardData")
	procGetOpenFileNameW    = comdlg32.NewProc("GetOpenFileNameW")
	procCreateIcon          = user32.NewProc("CreateIcon")
	procDestroyIcon         = user32.NewProc("DestroyIcon")
	procSetForegroundWindow = user32.NewProc("SetForegroundWindow")
	procGetForegroundWindow = user32.NewProc("GetForegroundWindow")
	procGetCursorPos        = user32.NewProc("GetCursorPos")
	procCreatePopupMenu     = user32.NewProc("CreatePopupMenu")
	procAppendMenuW         = user32.NewProc("AppendMenuW")
	procTrackPopupMenu      = user32.NewProc("TrackPopupMenu")
	procDestroyMenu         = user32.NewProc("DestroyMenu")
	procPostMessageW        = user32.NewProc("PostMessageW")
	procGetDlgCtrlID        = user32.NewProc("GetDlgCtrlID")
	procInvalidateRect      = user32.NewProc("InvalidateRect")
	procFillRect            = user32.NewProc("FillRect")
	procFrameRect           = user32.NewProc("FrameRect")
	procDrawTextW           = user32.NewProc("DrawTextW")
	procGetGUIThreadInfo    = user32.NewProc("GetGUIThreadInfo")
	procGetWindowLongW      = user32.NewProc("GetWindowLongW")
	procSetProcessDPIAware  = user32.NewProc("SetProcessDPIAware")

	procCreateSolidBrush = gdi32.NewProc("CreateSolidBrush")
	procDeleteObject     = gdi32.NewProc("DeleteObject")
	procSetTextColor     = gdi32.NewProc("SetTextColor")
	procSetBkColor       = gdi32.NewProc("SetBkColor")
	procSetBkMode        = gdi32.NewProc("SetBkMode")
	procSelectObject     = gdi32.NewProc("SelectObject")
	procCreateFontW      = gdi32.NewProc("CreateFontW")

	procShellNotifyIconW = shell32.NewProc("Shell_NotifyIconW")

	procGlobalAlloc      = kernel32.NewProc("GlobalAlloc")
	procGlobalLock       = kernel32.NewProc("GlobalLock")
	procGlobalUnlock     = kernel32.NewProc("GlobalUnlock")
	procGlobalFree       = kernel32.NewProc("GlobalFree")
	procRtlMoveMemory    = kernel32.NewProc("RtlMoveMemory")
	procGetModuleHandleW = kernel32.NewProc("GetModuleHandleW")
	procCreateMutexW     = kernel32.NewProc("CreateMutexW")
	procCloseHandle      = kernel32.NewProc("CloseHandle")

	procRegCreateKeyExW = advapi32.NewProc("RegCreateKeyExW")
	procRegSetValueExW  = advapi32.NewProc("RegSetValueExW")
	procRegDeleteValueW = advapi32.NewProc("RegDeleteValueW")
	procRegCloseKey     = advapi32.NewProc("RegCloseKey")
)

const (
	CS_HREDRAW = 0x0002
	CS_VREDRAW = 0x0001

	WS_OVERLAPPED       = 0x00000000
	WS_CAPTION          = 0x00C00000
	WS_SYSMENU          = 0x00080000
	WS_MINIMIZEBOX      = 0x00020000
	WS_VISIBLE          = 0x10000000
	WS_CHILD            = 0x40000000
	WS_TABSTOP          = 0x00010000
	WS_BORDER           = 0x00800000
	WS_VSCROLL          = 0x00200000
	WS_OVERLAPPEDWINDOW = 0x00CF0000

	BS_PUSHBUTTON    = 0x00000000
	BS_DEFPUSHBUTTON = 0x00000001
	BS_CHECKBOX      = 0x00000002
	BS_AUTOCHECKBOX  = 0x00000003
	BS_GROUPBOX      = 0x00000007
	BS_OWNERDRAW     = 0x0000000B
	BS_FLAT          = 0x00008000

	ES_AUTOHSCROLL = 0x0080
	ES_PASSWORD    = 0x0020

	CBS_DROPDOWNLIST = 0x0003
	LBS_NOTIFY       = 0x0001

	SW_SHOW = 5
	SW_HIDE = 0

	WM_CREATE          = 0x0001
	WM_DESTROY         = 0x0002
	WM_CLOSE           = 0x0010
	WM_DRAWITEM        = 0x002B
	WM_SETFONT         = 0x0030
	WM_SETICON         = 0x0080
	WM_COMMAND         = 0x0111
	WM_KEYDOWN         = 0x0100
	WM_KEYUP           = 0x0101
	WM_SYSKEYDOWN      = 0x0104
	WM_SYSKEYUP        = 0x0105
	WM_CTLCOLORBTN     = 0x0135
	WM_CTLCOLOREDIT    = 0x0133
	WM_CTLCOLORLISTBOX = 0x0134
	WM_CTLCOLORSTATIC  = 0x0138
	WM_LBUTTONUP       = 0x0202
	WM_LBUTTONDBLCLK   = 0x0203
	WM_RBUTTONUP       = 0x0205
	WM_APP             = 0x8000

	BN_CLICKED    = 0
	CBN_SELCHANGE = 1
	LBN_SELCHANGE = 1

	BM_GETCHECK   = 0x00F0
	BM_SETCHECK   = 0x00F1
	BST_UNCHECKED = 0
	BST_CHECKED   = 1

	CB_ADDSTRING    = 0x0143
	CB_GETCURSEL    = 0x0147
	CB_SETCURSEL    = 0x014E
	LB_ADDSTRING    = 0x0180
	LB_RESETCONTENT = 0x0184
	LB_GETCURSEL    = 0x0188

	MB_OK              = 0x00000000
	MB_ICONINFORMATION = 0x00000040
	MB_ICONWARNING     = 0x00000030

	WH_KEYBOARD_LL = 13
	HC_ACTION      = 0
	LLKHF_INJECTED = 0x00000010

	VK_BACK    = 0x08
	VK_TAB     = 0x09
	VK_RETURN  = 0x0D
	VK_SHIFT   = 0x10
	VK_LSHIFT  = 0xA0
	VK_RSHIFT  = 0xA1
	VK_CONTROL = 0x11
	VK_MENU    = 0x12
	VK_CAPITAL = 0x14
	VK_ESCAPE  = 0x1B
	VK_SPACE   = 0x20
	VK_LWIN    = 0x5B
	VK_RWIN    = 0x5C

	ICON_SMALL = 0
	ICON_BIG   = 1

	ODS_SELECTED = 0x0001
	ODS_DISABLED = 0x0004
	ODS_FOCUS    = 0x0010

	DT_CENTER     = 0x00000001
	DT_VCENTER    = 0x00000004
	DT_SINGLELINE = 0x00000020
	TRANSPARENT   = 1

	MF_STRING       = 0x00000000
	MF_SEPARATOR    = 0x00000800
	MF_CHECKED      = 0x00000008
	TPM_RIGHTBUTTON = 0x0002
	TPM_RETURNCMD   = 0x0100

	NIM_ADD     = 0x00000000
	NIM_MODIFY  = 0x00000001
	NIM_DELETE  = 0x00000002
	NIF_MESSAGE = 0x00000001
	NIF_ICON    = 0x00000002
	NIF_TIP     = 0x00000004

	FW_NORMAL           = 400
	FW_BOLD             = 700
	DEFAULT_CHARSET     = 1
	OUT_DEFAULT_PRECIS  = 0
	CLIP_DEFAULT_PRECIS = 0
	DEFAULT_QUALITY     = 0
	FIXED_PITCH         = 1
	FF_MODERN           = 48
	DEFAULT_PITCH       = 0
	FF_SWISS            = 32

	KEYEVENTF_KEYUP   = 0x0002
	KEYEVENTF_UNICODE = 0x0004
	INPUT_KEYBOARD    = 1

	CF_UNICODETEXT = 13
	GMEM_MOVEABLE  = 0x0002

	OFN_PATHMUSTEXIST = 0x00000800
	OFN_FILEMUSTEXIST = 0x00001000
	OFN_EXPLORER      = 0x00080000
	GWL_STYLE         = ^uint32(15)
)

type point struct{ X, Y int32 }
type rect struct{ Left, Top, Right, Bottom int32 }

type guiThreadInfo struct {
	CbSize        uint32
	Flags         uint32
	HwndActive    uintptr
	HwndFocus     uintptr
	HwndCapture   uintptr
	HwndMenuOwner uintptr
	HwndMoveSize  uintptr
	HwndCaret     uintptr
	RcCaret       rect
}

type drawItemStruct struct {
	CtlType    uint32
	CtlID      uint32
	ItemID     uint32
	ItemAction uint32
	ItemState  uint32
	HwndItem   uintptr
	HDC        uintptr
	RcItem     rect
	ItemData   uintptr
}

type notifyIconData struct {
	CbSize           uint32
	HWnd             uintptr
	UID              uint32
	UFlags           uint32
	UCallbackMessage uint32
	HIcon            uintptr
	SzTip            [128]uint16
	DwState          uint32
	DwStateMask      uint32
	SzInfo           [256]uint16
	UVersion         uint32
	SzInfoTitle      [64]uint16
	DwInfoFlags      uint32
	GuidItem         [16]byte
	HBalloonIcon     uintptr
}

type msg struct {
	Hwnd     uintptr
	Message  uint32
	WParam   uintptr
	LParam   uintptr
	Time     uint32
	Pt       point
	LPrivate uint32
}

type wndClassEx struct {
	CbSize        uint32
	Style         uint32
	LpfnWndProc   uintptr
	CbClsExtra    int32
	CbWndExtra    int32
	HInstance     uintptr
	HIcon         uintptr
	HCursor       uintptr
	HbrBackground uintptr
	LpszMenuName  *uint16
	LpszClassName *uint16
	HIconSm       uintptr
}

type kbdLLHookStruct struct {
	VkCode      uint32
	ScanCode    uint32
	Flags       uint32
	Time        uint32
	DwExtraInfo uintptr
}

type keybdInput struct {
	WVKey       uint16
	WScan       uint16
	DwFlags     uint32
	Time        uint32
	DwExtraInfo uintptr
}

type input struct {
	Type uint32
	Pad  uint32
	Ki   keybdInput
	Tail [8]byte
}

type openFileName struct {
	LStructSize       uint32
	HwndOwner         uintptr
	HInstance         uintptr
	LpstrFilter       *uint16
	LpstrCustomFilter *uint16
	NMaxCustFilter    uint32
	NFilterIndex      uint32
	LpstrFile         *uint16
	NMaxFile          uint32
	LpstrFileTitle    *uint16
	NMaxFileTitle     uint32
	LpstrInitialDir   *uint16
	LpstrTitle        *uint16
	Flags             uint32
	NFileOffset       uint16
	NFileExtension    uint16
	LpstrDefExt       *uint16
	LCustData         uintptr
	LpfnHook          uintptr
	LpTemplateName    *uint16
	PvReserved        uintptr
	DwReserved        uint32
	FlagsEx           uint32
}

func utf16Ptr(s string) *uint16 {
	p, _ := syscall.UTF16PtrFromString(s)
	return p
}

func loword(v uintptr) uint16 { return uint16(v & 0xffff) }
func hiword(v uintptr) uint16 { return uint16((v >> 16) & 0xffff) }

func setWindowText(hwnd uintptr, s string) {
	procSetWindowTextW.Call(hwnd, uintptr(unsafe.Pointer(utf16Ptr(s))))
}

func getWindowText(hwnd uintptr) string {
	ln, _, _ := procGetWindowTextLen.Call(hwnd)
	buf := make([]uint16, ln+1)
	procGetWindowTextW.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	return syscall.UTF16ToString(buf)
}

func sendMessage(hwnd uintptr, msg uint32, wparam, lparam uintptr) uintptr {
	r, _, _ := procSendMessageW.Call(hwnd, uintptr(msg), wparam, lparam)
	return r
}

func messageBox(hwnd uintptr, title, text string, flags uintptr) {
	procMessageBoxW.Call(hwnd, uintptr(unsafe.Pointer(utf16Ptr(text))), uintptr(unsafe.Pointer(utf16Ptr(title))), flags)
}
