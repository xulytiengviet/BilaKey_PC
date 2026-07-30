//go:build windows

package win

import (
	"fmt"
	"strings"
	"syscall"
	"unicode"
	"unicode/utf8"
	"unsafe"

	"github.com/xulytiengviet/BilaKey_PC/internal/core"
	"github.com/xulytiengviet/BilaKey_PC/internal/hotkey"
	"github.com/xulytiengviet/BilaKey_PC/internal/typingstate"
)

var keyboardProcPtr = syscall.NewCallback(keyboardProc)

func (a *App) installKeyboardHook() error {
	hinst, _, _ := procGetModuleHandleW.Call(0)
	hook, _, err := procSetWindowsHookEx.Call(WH_KEYBOARD_LL, keyboardProcPtr, hinst, 0)
	if hook == 0 {
		return err
	}
	a.hook = hook
	return nil
}

func (a *App) uninstallKeyboardHook() {
	if a.hook != 0 {
		procUnhookWindowsHook.Call(a.hook)
		a.hook = 0
	}
}

func keyboardProc(nCode int, wParam uintptr, kbd *kbdLLHookStruct) uintptr {
	a := currentApp
	if a == nil || nCode < HC_ACTION || kbd == nil || kbd.Flags&LLKHF_INJECTED != 0 {
		return callNextKeyboardHook(nCode, wParam, kbd)
	}

	switch wParam {
	case WM_KEYDOWN, WM_SYSKEYDOWN:
		if isShiftVK(kbd.VkCode) {
			if a.runtimeConfig().Enabled {
				a.handleShiftDown()
			}
			return callNextKeyboardHook(nCode, wParam, kbd)
		}
		if a.runtimeConfig().Enabled {
			a.markShiftUsed()
		}
		if a.handleKeydown(kbd.VkCode) {
			return 1
		}
	case WM_KEYUP, WM_SYSKEYUP:
		if isShiftVK(kbd.VkCode) && a.runtimeConfig().Enabled {
			if a.handleShiftUp(int64(kbd.Time)) {
				a.refreshMainStatus()
			}
		}
	}
	return callNextKeyboardHook(nCode, wParam, kbd)
}

func callNextKeyboardHook(nCode int, wParam uintptr, kbd *kbdLLHookStruct) uintptr {
	var p uintptr
	if kbd != nil {
		p = uintptr(unsafe.Pointer(kbd))
	}
	r, _, _ := procCallNextHookEx.Call(0, uintptr(nCode), wParam, p)
	return r
}

func isShiftVK(vk uint32) bool {
	return vk == VK_SHIFT || vk == VK_LSHIFT || vk == VK_RSHIFT
}

func (a *App) handleShiftDown() {
	a.mu.Lock()
	a.caps.ShiftDown()
	a.mu.Unlock()
}

func (a *App) markShiftUsed() {
	a.mu.Lock()
	a.caps.MarkShiftUsed()
	a.mu.Unlock()
}

func (a *App) handleShiftUp(nowMS int64) bool {
	a.mu.Lock()
	changed := a.caps.ShiftUp(nowMS)
	if changed {
		a.lastCommit = commitSnapshot{}
	}
	a.mu.Unlock()
	return changed
}

func (a *App) handleKeydown(vk uint32) bool {
	cfg := a.runtimeConfig()

	// Global shortcuts are resolved without I/O. The UI thread performs state
	// changes, which keeps WH_KEYBOARD_LL responsive and avoids stealing Ctrl+Tab.
	action := hotkey.Resolve(
		vk,
		isModifierDown(VK_CONTROL),
		isModifierDown(VK_SHIFT),
		isModifierDown(VK_MENU),
		isModifierDown(VK_LWIN) || isModifierDown(VK_RWIN),
	)
	if action != hotkey.None {
		if action != hotkey.CycleCandidate {
			a.resetComposition()
		}
		if a.mainHwnd != 0 {
			procPostMessageW.Call(a.mainHwnd, hotkeyActionMessage, uintptr(action), 0)
		}
		return true
	}

	if isModifierDown(VK_CONTROL) || isModifierDown(VK_MENU) || isModifierDown(VK_LWIN) || isModifierDown(VK_RWIN) {
		a.resetComposition()
		return false
	}
	if vk == VK_ESCAPE {
		a.resetComposition()
		return false
	}
	if cfg.PauseInPasswordFields && isNativePasswordField() {
		a.resetComposition()
		return false
	}

	if !cfg.Enabled {
		if cfg.MacroEnabled && cfg.MacroWhileOff {
			return a.handleMacroWhileOff(vk)
		}
		a.resetComposition()
		return false
	}

	if vk == VK_BACK {
		return a.handleBackspace(cfg.AlwaysUseClipboardUnicode, cfg.RestoreAfterDelimiter)
	}

	if delim, ok := delimiterFromVK(vk); ok {
		return a.commitDelimiter(delim, cfg.AlwaysUseClipboardUnicode, cfg.MacroEnabled, cfg.RestoreAfterDelimiter)
	}

	ch, ok := keyRune(vk)
	if !ok {
		a.resetComposition()
		return false
	}
	if cfg.InputMethod != core.MethodVNI && unicode.IsDigit(ch) {
		// Numbers are literal delimiters for Telex/CVNSS; do not consume them.
		a.resetComposition()
		return false
	}

	engine := a.engine.Load()
	if engine == nil {
		return false
	}

	a.mu.Lock()
	oldRendered := a.rendered
	startingWord := a.raw == "" && unicode.IsLetter(ch)
	modeChanged := false
	if startingWord && !a.wordCaseActive {
		beforeMode := a.caps.ModeLabel()
		a.wordCase = a.caps.BeginWord()
		a.wordCaseActive = true
		modeChanged = beforeMode != a.caps.ModeLabel()
	}
	if unicode.IsLetter(ch) {
		ch = applyCompositionCase(ch, utf8.RuneCountInString(a.raw), a.wordCase)
	}
	a.lastCommit = commitSnapshot{}
	newRaw := a.raw + string(ch)
	newRendered := engine.Transform(newRaw)
	a.raw = newRaw
	a.rendered = newRendered
	a.mu.Unlock()

	if err := replaceRendered(oldRendered, newRendered, cfg.AlwaysUseClipboardUnicode, a.mainHwnd); err != nil {
		a.feedback("Lỗi Unicode: " + err.Error())
		return false
	}
	if modeChanged {
		a.refreshMainStatus()
	}
	if cfg.ShowFeedback {
		a.feedback(fmt.Sprintf("%s → %s", newRaw, newRendered))
	}
	return true
}

func isNativePasswordField() bool {
	info := guiThreadInfo{CbSize: uint32(unsafe.Sizeof(guiThreadInfo{}))}
	r, _, _ := procGetGUIThreadInfo.Call(0, uintptr(unsafe.Pointer(&info)))
	if r == 0 || info.HwndFocus == 0 {
		return false
	}
	style, _, _ := procGetWindowLongW.Call(info.HwndFocus, uintptr(GWL_STYLE))
	return uint32(style)&ES_PASSWORD != 0
}

func applyCompositionCase(ch rune, index int, wc typingstate.WordCase) rune {
	if wc.AllUpper {
		return unicode.ToUpper(ch)
	}
	if index == 0 && wc.UpperFirst {
		return unicode.ToUpper(ch)
	}
	return ch
}

func (a *App) handleBackspace(useClipboard bool, restoreAfterDelimiter bool) bool {
	a.mu.Lock()
	if a.raw == "" {
		if restoreAfterDelimiter && a.lastCommit.valid {
			s := a.lastCommit
			a.raw = s.raw
			a.rendered = s.rendered
			a.wordCase = s.wordCase
			a.wordCaseActive = true
			a.caps.Restore(s.capsBefore)
			a.lastCommit = commitSnapshot{}
			a.mu.Unlock()
			// Returning false lets Windows delete only the delimiter. The previous
			// word remains on screen while our raw/rendered buffers are restored,
			// so the next edit continues as real Quốc ngữ instead of plain ASCII.
			a.refreshMainStatus()
			return false
		}
		a.mu.Unlock()
		return false
	}

	r := []rune(a.raw)
	a.raw = string(r[:len(r)-1])
	oldRendered := a.rendered
	newRendered := ""
	if a.raw != "" {
		if engine := a.engine.Load(); engine != nil {
			newRendered = engine.Transform(a.raw)
		}
	}
	a.rendered = newRendered
	a.lastCommit = commitSnapshot{}
	a.mu.Unlock()

	if err := replaceRendered(oldRendered, newRendered, useClipboard, a.mainHwnd); err != nil {
		a.feedback("Lỗi Backspace: " + err.Error())
		return false
	}
	return true
}

func (a *App) commitDelimiter(delim string, useClipboard bool, macroEnabled bool, restoreAfterDelimiter bool) bool {
	a.mu.Lock()
	raw := a.raw
	rendered := a.rendered
	capsBefore := a.caps.Snapshot()
	wordCase := a.wordCase

	a.raw, a.rendered = "", ""
	a.wordCase = wordCase
	a.wordCaseActive = false
	if raw != "" && restoreAfterDelimiter {
		a.lastCommit = commitSnapshot{
			valid:      true,
			raw:        raw,
			rendered:   rendered,
			wordCase:   wordCase,
			capsBefore: capsBefore,
			delimiter:  delim,
		}
	} else {
		// A second delimiter (for example period + space) invalidates the
		// single-step rollback snapshot, preventing accidental restoration of
		// the wrong word.
		a.lastCommit = commitSnapshot{}
	}
	a.caps.ObserveDelimiter(delim)
	a.mu.Unlock()

	if strings.ContainsAny(delim, ".!?\r\n") {
		a.refreshMainStatus()
	}

	if macroEnabled && raw != "" {
		if replacement, ok := a.macros.Lookup(raw); ok {
			a.mu.Lock()
			a.lastCommit = commitSnapshot{}
			a.mu.Unlock()
			if err := replaceRendered(rendered, replacement, useClipboard, a.mainHwnd); err != nil {
				a.feedback("Lỗi gõ tắt: " + err.Error())
				return false
			}
			if err := sendUnicodeText(delim); err != nil {
				return false
			}
			return true
		}
	}
	return false
}

func (a *App) handleMacroWhileOff(vk uint32) bool {
	if vk == VK_BACK {
		a.mu.Lock()
		if a.raw != "" {
			r := []rune(a.raw)
			a.raw = string(r[:len(r)-1])
			a.rendered = a.raw
		}
		a.mu.Unlock()
		return false
	}
	if delim, ok := delimiterFromVK(vk); ok {
		a.mu.Lock()
		raw := a.raw
		a.raw, a.rendered = "", ""
		a.mu.Unlock()
		if raw != "" {
			if replacement, found := a.macros.Lookup(raw); found {
				if err := sendBackspaces(len([]rune(raw))); err == nil {
					_ = sendUnicodeText(replacement + delim)
					return true
				}
			}
		}
		return false
	}
	ch, ok := keyRune(vk)
	if !ok {
		a.resetComposition()
		return false
	}
	a.mu.Lock()
	a.raw += string(ch)
	a.rendered = a.raw
	a.mu.Unlock()
	return false
}

func (a *App) feedback(text string) {
	if h := a.controls[idStatus]; h != 0 {
		setWindowText(h, text)
	}
}

func isModifierDown(vk int) bool {
	r, _, _ := procGetAsyncKeyState.Call(uintptr(vk))
	return uint16(r)&0x8000 != 0
}

func keyRune(vk uint32) (rune, bool) {
	if vk >= 'A' && vk <= 'Z' {
		shift := keyStateActive(VK_SHIFT)
		caps := keyStateToggle(VK_CAPITAL)
		upper := shift != caps
		r := rune(vk)
		if !upper {
			r = unicode.ToLower(r)
		}
		return r, true
	}
	if vk >= '0' && vk <= '9' {
		return rune(vk), true
	}
	return 0, false
}

func keyStateActive(vk int) bool {
	r, _, _ := procGetKeyState.Call(uintptr(vk))
	return uint16(r)&0x8000 != 0
}

func keyStateToggle(vk int) bool {
	r, _, _ := procGetKeyState.Call(uintptr(vk))
	return uint16(r)&1 != 0
}

func delimiterFromVK(vk uint32) (string, bool) {
	shift := keyStateActive(VK_SHIFT)
	switch vk {
	case '1':
		if shift {
			return "!", true
		}
	case '2':
		if shift {
			return "@", true
		}
	case '3':
		if shift {
			return "#", true
		}
	case '4':
		if shift {
			return "$", true
		}
	case '5':
		if shift {
			return "%", true
		}
	case '6':
		if shift {
			return "^", true
		}
	case '7':
		if shift {
			return "&", true
		}
	case '8':
		if shift {
			return "*", true
		}
	case '9':
		if shift {
			return "(", true
		}
	case '0':
		if shift {
			return ")", true
		}
	case VK_SPACE:
		return " ", true
	case VK_RETURN:
		return "\r", true
	case VK_TAB:
		return "\t", true
	case 0xBA:
		if shift {
			return ":", true
		}
		return ";", true
	case 0xBB:
		if shift {
			return "+", true
		}
		return "=", true
	case 0xBC:
		if shift {
			return "<", true
		}
		return ",", true
	case 0xBD:
		if shift {
			return "_", true
		}
		return "-", true
	case 0xBE:
		if shift {
			return ">", true
		}
		return ".", true
	case 0xBF:
		if shift {
			return "?", true
		}
		return "/", true
	case 0xC0:
		if shift {
			return "~", true
		}
		return "`", true
	case 0xDB:
		if shift {
			return "{", true
		}
		return "[", true
	case 0xDC:
		if shift {
			return "|", true
		}
		return "\\", true
	case 0xDD:
		if shift {
			return "}", true
		}
		return "]", true
	case 0xDE:
		if shift {
			return "\"", true
		}
		return "'", true
	}
	return "", false
}

func normalizeMethodName(s string) string {
	switch strings.ToLower(s) {
	case "telex":
		return core.MethodTelex
	case "vni":
		return core.MethodVNI
	default:
		return core.MethodCVNSS
	}
}
