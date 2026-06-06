package main

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	kernel32 = windows.NewLazySystemDLL("kernel32.dll")
	user32   = windows.NewLazySystemDLL("user32.dll")

	EnumWindows        = user32.NewProc("EnumWindows")
	ShowWindow         = user32.NewProc("ShowWindow")
	IsWindowVisible    = user32.NewProc("IsWindowVisible")
	GetWindowTextW     = user32.NewProc("GetWindowTextW")
	AllocConsole       = kernel32.NewProc("AllocConsole")
	GetModuleFileNameW = kernel32.NewProc("GetModuleFileNameW")
)

func xe() string {
	var buf [windows.MAX_PATH]uint16
	r, _, _ := GetModuleFileNameW.Call(
		0,
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(len(buf)),
	)
	if r == 0 {
		return ""
	}
	path := windows.UTF16ToString(buf[:r])
	if i := strings.LastIndex(path, "\\"); i != -1 {
		return strings.ToLower(path[i+1:])
	}
	return strings.ToLower(path)
}

func hideBanner() error {
	capt := "Version"
	partialLower := strings.ToLower(capt)
	var found bool

	cb := syscall.NewCallback(func(hwnd windows.Handle, _ uintptr) uintptr {
		if found {
			return 1
		}

		var buf [256]uint16
		n, _, _ := GetWindowTextW.Call(
			uintptr(hwnd),
			uintptr(unsafe.Pointer(&buf[0])),
			uintptr(len(buf)),
		)
		if n > 0 {
			title := windows.UTF16ToString(buf[:n])
			if strings.Contains(strings.ToLower(title), partialLower) {

				visible, _, _ := IsWindowVisible.Call(uintptr(hwnd))
				if visible != 0 {
					ShowWindow.Call(uintptr(hwnd), 0)
				}
				found = true
				return 0
			}
		}
		return 1
	})

	_, _, err := EnumWindows.Call(cb, 0)
	if err != nil && err.Error() != "The operation completed successfully." {
		return err
	}

	if !found {
		return fmt.Errorf("[Nava::Error] no window name: '%s'", capt)
	}
	return nil
}
