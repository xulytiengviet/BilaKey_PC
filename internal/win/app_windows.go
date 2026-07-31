//go:build windows

package win

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"syscall"
	"unsafe"

	"github.com/xulytiengviet/BilaKey_PC/internal/core"
	"github.com/xulytiengviet/BilaKey_PC/internal/macro"
	"github.com/xulytiengviet/BilaKey_PC/internal/settings"
	"github.com/xulytiengviet/BilaKey_PC/internal/typingstate"
)

type commitSnapshot struct {
	valid      bool
	raw        string
	rendered   string
	wordCase   typingstate.WordCase
	capsBefore typingstate.Snapshot
	delimiter  string
}

type themeHandles struct {
	font       uintptr
	fontBold   uintptr
	bgBrush    uintptr
	panelBrush uintptr
	editBrush  uintptr
	cyanBrush  uintptr
	grayBrush  uintptr
}

type App struct {
	mu             sync.Mutex
	store          *settings.Store
	macros         *macro.Table
	engine         atomic.Pointer[core.Engine]
	runtimeCfg     atomic.Value // settings.Config
	mainHwnd       uintptr
	optionsHwnd    uintptr
	macroHwnd      uintptr
	hook           uintptr
	raw            string
	rendered       string
	wordCase       typingstate.WordCase
	wordCaseActive bool
	lastCommit     commitSnapshot
	caps           *typingstate.Capitalizer
	controls       map[int]uintptr
	macroControls  map[int]uintptr
	theme          themeHandles
	iconSmall      uintptr
	iconBig        uintptr
	trayAdded      bool
	exiting        bool
}

var currentApp *App

func NewApp(store *settings.Store) (*App, error) {
	cfg := store.Get()
	macros := macro.New(cfg.MacroFile)
	if err := macros.Load(); err != nil {
		return nil, fmt.Errorf("đọc bảng gõ tắt: %w", err)
	}
	app := &App{
		store:         store,
		macros:        macros,
		caps:          typingstate.New(cfg.AutoCapInitial, cfg.AutoCapSentence, cfg.DoubleShiftCaps),
		controls:      make(map[int]uintptr),
		macroControls: make(map[int]uintptr),
	}
	app.runtimeCfg.Store(cfg)
	app.rebuildEngine()
	currentApp = app
	return app, nil
}

func (a *App) rebuildEngine() {
	cfg := a.store.Get()
	a.runtimeCfg.Store(cfg)
	a.engine.Store(core.New(cfg.InputMethod, core.Options{
		OldToneStyle:        cfg.OldToneStyle,
		FreeToneMarking:     cfg.FreeToneMarking,
		SpellCheck:          cfg.SpellCheck,
		AutoRestoreWrongKey: cfg.AutoRestoreWrongKey,
	}))
	a.mu.Lock()
	a.caps.Configure(cfg.AutoCapInitial, cfg.AutoCapSentence, cfg.DoubleShiftCaps)
	a.mu.Unlock()
}

func (a *App) runtimeConfig() settings.Config {
	v := a.runtimeCfg.Load()
	if v == nil {
		return a.store.Get()
	}
	return v.(settings.Config)
}

func (a *App) Run() error {
	mutex, alreadyRunning, err := acquireSingleInstance()
	if err != nil {
		return err
	}
	if alreadyRunning {
		messageBox(0, settings.AppName, "BilaKey PC đã chạy trong khay hệ thống.", MB_OK|MB_ICONINFORMATION)
		return nil
	}
	defer procCloseHandle.Call(mutex)
	procSetProcessDPIAware.Call()

	a.initTheme()
	defer a.destroyTheme()
	a.initIcons()
	defer a.destroyIcons()

	if err := a.registerClasses(); err != nil {
		return err
	}
	if err := a.createMainWindow(); err != nil {
		return err
	}
	if err := a.installTrayIcon(); err != nil {
		messageBox(a.mainHwnd, "BilaKey PC", "Không thể tạo biểu tượng khay hệ thống: "+err.Error(), MB_OK|MB_ICONWARNING)
	}
	defer a.removeTrayIcon()

	if err := a.installKeyboardHook(); err != nil {
		messageBox(a.mainHwnd, "BilaKey PC", "Không thể cài keyboard hook toàn cục: "+err.Error(), MB_OK|MB_ICONWARNING)
	}
	defer a.uninstallKeyboardHook()

	cfg := a.runtimeConfig()
	if cfg.ShowDialogAtStartup {
		procShowWindow.Call(a.mainHwnd, SW_SHOW)
	} else {
		procShowWindow.Call(a.mainHwnd, SW_HIDE)
	}
	procUpdateWindow.Call(a.mainHwnd)
	a.refreshMainStatus()

	var m msg
	for {
		r, _, err := procGetMessageW.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0)
		if int32(r) == -1 {
			return err
		}
		if r == 0 {
			break
		}
		procTranslateMessage.Call(uintptr(unsafe.Pointer(&m)))
		procDispatchMessageW.Call(uintptr(unsafe.Pointer(&m)))
	}
	return nil
}

func acquireSingleInstance() (uintptr, bool, error) {
	handle, _, callErr := procCreateMutexW.Call(0, 0, uintptr(unsafe.Pointer(utf16Ptr("Local\\BilaKeyPC.LongNgo.1"))))
	if handle == 0 {
		return 0, false, fmt.Errorf("CreateMutexW: %w", callErr)
	}
	alreadyRunning := callErr == syscall.Errno(183)
	return handle, alreadyRunning, nil
}

func (a *App) updateConfig(fn func(*settings.Config)) {
	if err := a.store.Update(fn); err != nil {
		messageBox(a.mainHwnd, "Lỗi cấu hình", err.Error(), MB_OK|MB_ICONWARNING)
		return
	}
	a.rebuildEngine()
	if a.runtimeConfig().Enabled {
		a.resetComposition()
	} else {
		a.resetTypingState()
	}
	a.refreshMainStatus()
}

func (a *App) resetComposition() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.raw = ""
	a.rendered = ""
	a.wordCase = typingstate.WordCase{}
	a.wordCaseActive = false
	a.lastCommit = commitSnapshot{}
}

func (a *App) resetTypingState() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.raw = ""
	a.rendered = ""
	a.wordCase = typingstate.WordCase{}
	a.wordCaseActive = false
	a.lastCommit = commitSnapshot{}
	a.caps.Reset()
}

func (a *App) capsModeLabel() string {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.caps.ModeLabel()
}

func (a *App) cycleCVNSSCandidate() {
	if a.runtimeConfig().InputMethod != core.MethodCVNSS {
		a.feedback("Đổi ứng viên chỉ dùng cho CVNSS4.0")
		return
	}
	a.mu.Lock()
	raw := a.raw
	oldRendered := a.rendered
	candidates := core.CVNSSCandidates(raw)
	if raw == "" || len(candidates) < 2 {
		a.mu.Unlock()
		a.feedback("Từ hiện tại không có nhập nhằng CVNSS4.0")
		return
	}
	next := 0
	for i, candidate := range candidates {
		if candidate == oldRendered {
			next = (i + 1) % len(candidates)
			break
		}
	}
	a.rendered = candidates[next]
	newRendered := a.rendered
	a.mu.Unlock()
	if err := replaceRendered(oldRendered, newRendered, a.runtimeConfig().AlwaysUseClipboardUnicode, a.mainHwnd); err != nil {
		a.feedback("Lỗi đổi ứng viên: " + err.Error())
		return
	}
	a.feedback(fmt.Sprintf("Ứng viên CVNSS %d/%d: %s", next+1, len(candidates), newRendered))
}

func (a *App) refreshMainStatus() {
	cfg := a.runtimeConfig()
	state := "ĐANG BẬT"
	if !cfg.Enabled {
		state = "ĐANG TẮT"
	}
	methodRole := "CVNSS4.0 CORE"
	if cfg.InputMethod == core.MethodVNITelex {
		methodRole = "VNI/TELEX · TỰ NHẬN DẠNG"
	}
	if hwnd := a.controls[idStatus]; hwnd != 0 {
		setWindowText(hwnd, fmt.Sprintf("  %s  •  %s  •  Unicode  •  Ctrl+Shift+Space: bật/tắt", state, methodRole))
	}
	if hwnd := a.controls[idCapsStatus]; hwnd != 0 {
		setWindowText(hwnd, "Ctrl+Shift+1: CVNSS4.0  •  Ctrl+Shift+2: VNI/Telex  •  Ctrl+Shift+0: ứng viên")
	}
	if hwnd := a.controls[idToggle]; hwnd != 0 {
		if cfg.Enabled {
			setWindowText(hwnd, "TẮT")
		} else {
			setWindowText(hwnd, "BẬT")
		}
	}
	a.syncMethodCombo()
	a.updateTrayIcon()
}

func (a *App) configFolder() string {
	p := a.store.Path()
	return filepath.Dir(p)
}

func (a *App) showHelp() {
	text := "BilaKey PC " + settings.AppVersion + " · CVNSS Core Edition\r\n\r\n" +
		"1. BilaKey chỉ có hai kiểu gõ: CVNSS4.0 và VNI/Telex.\r\n" +
		"2. VNI/Telex tự nhận cả phím chữ Telex và phím số VNI; có thể đổi cách gõ giữa từng từ.\r\n" +
		"3. Ctrl+Shift+Space bật/tắt; Ctrl+Shift+1 chọn CVNSS4.0; Ctrl+Shift+2 chọn VNI/Telex.\r\n" +
		"4. Ctrl+Shift+0 đổi ứng viên khi mã CVNSS có nhiều cách giải.\r\n" +
		"5. Đầu ra luôn là Unicode; chữ đầu tiên và chữ sau . ! ? tự động viết hoa.\r\n" +
		"6. Shift×1 viết hoa một chữ; Shift×2 khóa viết hoa liên tục.\r\n" +
		"7. Backspace ngay sau dấu cách/dấu câu quay lại từ vừa gõ để sửa tiếp.\r\n" +
		"8. Mặc định BilaKey tạm dừng trong ô mật khẩu Win32 và chỉ cho chạy một phiên bản.\r\n\r\n" +
		"Ví dụ: toiy → tôi; iwy → yêu; vidf → việt; qyl → quỳ; ses → sẽ."
	messageBox(a.mainHwnd, "Hướng dẫn BilaKey PC", text, MB_OK|MB_ICONINFORMATION)
}

func (a *App) showAbout() {
	text := "BilaKey PC " + settings.AppVersion + "\r\n" +
		"CVNSS Core Edition\r\n\r\n" +
		"Lõi nhập liệu: CVNSS4.0\r\n" +
		"Kiểu gõ hợp nhất: VNI/Telex tự nhận dạng\r\n" +
		"Rule oracle: CVNSS4.0 " + core.CVNSSRuleVersion + "\r\n" +
		"Phát triển: Long Ngo, 2026\r\n" +
		"Hỗ trợ CVNSS4.0: NNC Trần Tư Bình và cộng đồng.\r\n\r\n" +
		"Runtime: Go + Win32 native.\r\n" +
		"Toolchain audit: Python + C/C++ + Rust/Bamboo Core.\r\n" +
		"Icon B nền xanh/chữ trắng · không telemetry · không network runtime."
	messageBox(a.mainHwnd, "Thông tin", text, MB_OK|MB_ICONINFORMATION)
}

func (a *App) showSettingsInfo() {
	exe, _ := os.Executable()
	text := "Thư mục cấu hình:\r\n" + a.configFolder() + "\r\n\r\n" +
		"Tệp thực thi:\r\n" + exe + "\r\n\r\n" +
		"Cấu hình được lưu tự động. Tùy chọn 'Khởi động cùng Windows' dùng HKCU, không cần quyền Administrator.\r\n" +
		"Đóng cửa sổ bằng nút X chỉ ẩn BilaKey xuống khay hệ thống; dùng THOÁT hoặc menu tray để kết thúc tiến trình."
	messageBox(a.mainHwnd, "Cài đặt", text, MB_OK|MB_ICONINFORMATION)
}

func (a *App) showMainWindow() {
	if a.mainHwnd == 0 {
		return
	}
	procShowWindow.Call(a.mainHwnd, SW_SHOW)
	procSetForegroundWindow.Call(a.mainHwnd)
	procUpdateWindow.Call(a.mainHwnd)
}

func (a *App) exitApp() {
	a.exiting = true
	if a.mainHwnd != 0 {
		procDestroyWindow.Call(a.mainHwnd)
	}
}

func registerClass(name string, proc uintptr) error {
	hinst, _, _ := procGetModuleHandleW.Call(0)
	bg := uintptr(6)
	var iconBig, iconSmall uintptr
	if currentApp != nil {
		if currentApp.theme.bgBrush != 0 {
			bg = currentApp.theme.bgBrush
		}
		iconBig = currentApp.iconBig
		iconSmall = currentApp.iconSmall
	}
	wc := wndClassEx{
		CbSize:        uint32(unsafe.Sizeof(wndClassEx{})),
		Style:         CS_HREDRAW | CS_VREDRAW,
		LpfnWndProc:   proc,
		HInstance:     hinst,
		HIcon:         iconBig,
		HbrBackground: bg,
		LpszClassName: utf16Ptr(name),
		HIconSm:       iconSmall,
	}
	r, _, err := procRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc)))
	if r == 0 && err != syscall.Errno(1410) {
		return err
	}
	return nil
}
