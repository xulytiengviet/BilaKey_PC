//go:build windows

package win

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

const (
	hkeyCurrentUser = uintptr(0x80000001)
	keySetValue     = 0x0002
	regSZ           = 1
)

func setStartupWithWindows(enable bool) error {
	const subkey = `Software\Microsoft\Windows\CurrentVersion\Run`
	const valueName = "BilaKeyPC"
	var key uintptr
	var disposition uint32
	r, _, err := procRegCreateKeyExW.Call(
		hkeyCurrentUser,
		uintptr(unsafe.Pointer(utf16Ptr(subkey))),
		0, 0, 0, keySetValue, 0,
		uintptr(unsafe.Pointer(&key)),
		uintptr(unsafe.Pointer(&disposition)),
	)
	if r != 0 {
		return fmt.Errorf("RegCreateKeyExW lỗi %d: %v", r, err)
	}
	defer procRegCloseKey.Call(key)

	name := utf16Ptr(valueName)
	if !enable {
		r, _, _ := procRegDeleteValueW.Call(key, uintptr(unsafe.Pointer(name)))
		if r == 0 || r == 2 {
			return nil
		}
		return fmt.Errorf("RegDeleteValueW lỗi %d", r)
	}

	exe, err := os.Executable()
	if err != nil {
		return err
	}
	data := syscall.StringToUTF16(`"` + exe + `"`)
	r, _, err = procRegSetValueExW.Call(
		key,
		uintptr(unsafe.Pointer(name)),
		0,
		regSZ,
		uintptr(unsafe.Pointer(&data[0])),
		uintptr(len(data)*2),
	)
	if r != 0 {
		return fmt.Errorf("RegSetValueExW lỗi %d: %v", r, err)
	}
	return nil
}
