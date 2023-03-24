package main

import (
	"syscall"
	"unsafe"

	"github.com/lxn/win"
)

// https://github.com/KaiCoder/router/blob/712947901108692a53d37f0098cea3cd66de24c6/gopath/src/github.com/shirou/w32/dwmapi.go
var (
	moddwmapi = syscall.NewLazyDLL("dwmapi.dll")

	procDwmGetWindowAttribute = moddwmapi.NewProc("DwmGetWindowAttribute")
	procDwmSetWindowAttribute = moddwmapi.NewProc("DwmSetWindowAttribute")
)

// C:\Windows\System32\rundll32.exe dwmapi.dll,DwmEnableMMCSS 0

// https://docs.microsoft.com/en-us/windows/win32/api/dwmapi/nf-dwmapi-dwmenablecomposition#parameters

type DWMWINDOWATTRIBUTE int32

const (
	DWMWA_NCRENDERING_ENABLED         DWMWINDOWATTRIBUTE = 1
	DWMWA_NCRENDERING_POLICY          DWMWINDOWATTRIBUTE = 2
	DWMWA_TRANSITIONS_FORCEDISABLED   DWMWINDOWATTRIBUTE = 3
	DWMWA_ALLOW_NCPAINT               DWMWINDOWATTRIBUTE = 4
	DWMWA_CAPTION_BUTTON_BOUNDS       DWMWINDOWATTRIBUTE = 5
	DWMWA_NONCLIENT_RTL_LAYOUT        DWMWINDOWATTRIBUTE = 6
	DWMWA_FORCE_ICONIC_REPRESENTATION DWMWINDOWATTRIBUTE = 7
	DWMWA_FLIP3D_POLICY               DWMWINDOWATTRIBUTE = 8
	DWMWA_EXTENDED_FRAME_BOUNDS       DWMWINDOWATTRIBUTE = 9
	DWMWA_HAS_ICONIC_BITMAP           DWMWINDOWATTRIBUTE = 10
	DWMWA_DISALLOW_PEEK               DWMWINDOWATTRIBUTE = 11
	DWMWA_EXCLUDED_FROM_PEEK          DWMWINDOWATTRIBUTE = 12
	DWMWA_CLOAK                       DWMWINDOWATTRIBUTE = 13
	DWMWA_CLOAKED                     DWMWINDOWATTRIBUTE = 14
	DWMWA_FREEZE_REPRESENTATION       DWMWINDOWATTRIBUTE = 15
	DWMWA_PASSIVE_UPDATE_MODE         DWMWINDOWATTRIBUTE = 16
	DWMWA_LAST                        DWMWINDOWATTRIBUTE = 17
)

type DWMNCRENDERINGPOLICY int32

const (
	DWMNCRP_USEWINDOWSTYLE DWMNCRENDERINGPOLICY = 0
	DWMNCRP_DISABLED       DWMNCRENDERINGPOLICY = 1
	DWMNCRP_ENABLED        DWMNCRENDERINGPOLICY = 2
	DWMNCRP_LAST           DWMNCRENDERINGPOLICY = 3
)

type (
	// HANDLE  uintptr
	// HWND    HANDLE
	LPCVOID unsafe.Pointer
	HRESULT int32
)

func DwmSetWindowAttribute(hWnd win.HWND, dwAttribute DWMWINDOWATTRIBUTE, pvAttribute *DWMNCRENDERINGPOLICY, cbAttribute uint32) int32 {
	ret, _, _ := procDwmSetWindowAttribute.Call(
		uintptr(hWnd),
		uintptr(dwAttribute),
		uintptr(unsafe.Pointer(pvAttribute)),
		uintptr(cbAttribute))
	return int32(ret)
}

func DwmGetWindowAttribute(hWnd win.HWND, dwAttribute DWMWINDOWATTRIBUTE, pvAttribute *DWMNCRENDERINGPOLICY, cbAttribute uint32) int32 {
	ret, _, _ := procDwmGetWindowAttribute.Call(
		uintptr(hWnd),
		uintptr(dwAttribute),
		uintptr(unsafe.Pointer(pvAttribute)),
		uintptr(cbAttribute))
	return int32(ret)
}
