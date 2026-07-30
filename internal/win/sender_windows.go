//go:build windows

package win

import (
	"fmt"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

func sendInputs(inputs []input) error {
	if len(inputs) == 0 {
		return nil
	}
	r, _, err := procSendInput.Call(
		uintptr(len(inputs)),
		uintptr(unsafe.Pointer(&inputs[0])),
		unsafe.Sizeof(input{}),
	)
	if int(r) != len(inputs) {
		return fmt.Errorf("SendInput gửi %d/%d sự kiện: %v", r, len(inputs), err)
	}
	return nil
}

func keyboardInput(vk uint16, scan uint16, flags uint32) input {
	return input{Type: INPUT_KEYBOARD, Ki: keybdInput{WVKey: vk, WScan: scan, DwFlags: flags}}
}

func sendBackspaces(n int) error {
	if n <= 0 {
		return nil
	}
	inputs := make([]input, 0, n*2)
	for i := 0; i < n; i++ {
		inputs = append(inputs,
			keyboardInput(VK_BACK, 0, 0),
			keyboardInput(VK_BACK, 0, KEYEVENTF_KEYUP),
		)
	}
	return sendInputs(inputs)
}

func sendUnicodeText(s string) error {
	units := utf16.Encode([]rune(s))
	inputs := make([]input, 0, len(units)*2)
	for _, u := range units {
		inputs = append(inputs,
			keyboardInput(0, u, KEYEVENTF_UNICODE),
			keyboardInput(0, u, KEYEVENTF_UNICODE|KEYEVENTF_KEYUP),
		)
	}
	return sendInputs(inputs)
}

func sendLiteralKey(vk uint16) error {
	return sendInputs([]input{
		keyboardInput(vk, 0, 0),
		keyboardInput(vk, 0, KEYEVENTF_KEYUP),
	})
}

func sendPaste() error {
	return sendInputs([]input{
		keyboardInput(VK_CONTROL, 0, 0),
		keyboardInput('V', 0, 0),
		keyboardInput('V', 0, KEYEVENTF_KEYUP),
		keyboardInput(VK_CONTROL, 0, KEYEVENTF_KEYUP),
	})
}

func writeClipboardUnicode(text string, owner uintptr) error {
	r, _, err := procOpenClipboard.Call(owner)
	if r == 0 {
		return fmt.Errorf("OpenClipboard: %v", err)
	}
	defer procCloseClipboard.Call()
	if r, _, err := procEmptyClipboard.Call(); r == 0 {
		return fmt.Errorf("EmptyClipboard: %v", err)
	}

	units := syscall.StringToUTF16(text)
	size := uintptr(len(units) * 2)
	h, _, err := procGlobalAlloc.Call(GMEM_MOVEABLE, size)
	if h == 0 {
		return fmt.Errorf("GlobalAlloc: %v", err)
	}
	p, _, err := procGlobalLock.Call(h)
	if p == 0 {
		procGlobalFree.Call(h)
		return fmt.Errorf("GlobalLock: %v", err)
	}
	procRtlMoveMemory.Call(p, uintptr(unsafe.Pointer(&units[0])), size)
	procGlobalUnlock.Call(h)

	if r, _, err := procSetClipboardData.Call(CF_UNICODETEXT, h); r == 0 {
		procGlobalFree.Call(h)
		return fmt.Errorf("SetClipboardData: %v", err)
	}
	return nil
}

func replaceRendered(oldText, newText string, useClipboard bool, owner uintptr) error {
	if oldText == newText {
		return nil
	}
	if err := sendBackspaces(len([]rune(oldText))); err != nil {
		return err
	}
	if newText == "" {
		return nil
	}
	if useClipboard {
		if err := writeClipboardUnicode(newText, owner); err != nil {
			return err
		}
		return sendPaste()
	}
	return sendUnicodeText(newText)
}
