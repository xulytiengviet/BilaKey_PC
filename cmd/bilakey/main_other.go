//go:build !windows

package main

import "fmt"

import "github.com/xulytiengviet/BilaKey_PC/internal/settings"

func main() {
	fmt.Printf("BilaKey PC %s là ứng dụng Windows. Hãy build bằng GOOS=windows GOARCH=amd64.\n", settings.AppVersion)
}
