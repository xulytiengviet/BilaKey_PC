//go:build windows

package main

import (
	"fmt"
	"runtime"

	"github.com/xulytiengviet/BilaKey_PC/internal/settings"
	bilawin "github.com/xulytiengviet/BilaKey_PC/internal/win"
)

func main() {
	runtime.LockOSThread()
	store, err := settings.Open()
	if err != nil {
		panic(fmt.Errorf("mở cấu hình: %w", err))
	}
	app, err := bilawin.NewApp(store)
	if err != nil {
		panic(err)
	}
	if err := app.Run(); err != nil {
		panic(err)
	}
}
