package main

import (
	"log"
	"syscall"
	"unsafe"

	"github.com/lxn/win"
)

func main() {
	listener = &Listener{}
	go startListenerMessageLoop()

	list, err := GetAllWindows()
	if err != nil {
		log.Fatal(err)
	}

	for _, h := range list {
		hwnd := win.HWND(h)
		if getDWMactive(hwnd) {
			// https://docs.microsoft.com/en-us/windows/win32/api/dwmapi/nf-dwmapi-dwmsetwindowattribute
			var policyParameter = DWMNCRP_DISABLED
			// Use with DwmSetWindowAttribute. Sets the non-client rendering policy. The pvAttribute parameter points to a value from the DWMNCRENDERINGPOLICY enumeration.
			DwmSetWindowAttribute(hwnd, DWMWA_NCRENDERING_POLICY, &policyParameter, 8) // unsafe.Sizeof(int(0))

			policyParameter = DWMNCRP_DISABLED
			// Use with DwmSetWindowAttribute. Forces the window to display an iconic thumbnail or peek representation (a static bitmap), even if a live or snapshot representation of the window is available. This value is normally set during a window's creation, and not changed throughout the window's lifetime. Some scenarios, however, might require the value to change over time. The pvAttribute parameter points to a value of type BOOL. TRUE to require a iconic thumbnail or peek representation; otherwise, FALSE.
			DwmSetWindowAttribute(hwnd, DWMWA_FORCE_ICONIC_REPRESENTATION, &policyParameter, 8) // unsafe.Sizeof(int(0))
		}
	}

	select {}
}

func GetAllWindows() ([]syscall.Handle, error) {
	var hwndList []syscall.Handle
	cb := syscall.NewCallback(func(h syscall.Handle, p uintptr) uintptr {
		hwndList = append(hwndList, h)
		return 1 // continue enumeration
	})
	EnumWindows(cb, 0)

	return hwndList, nil
}

func getDWMactive(hWnd win.HWND) bool {
	var isNCRenderingEnabled = DWMNCRP_USEWINDOWSTYLE
	if DwmGetWindowAttribute(hWnd, DWMWA_NCRENDERING_ENABLED, &isNCRenderingEnabled, uint32(unsafe.Sizeof(isNCRenderingEnabled))) != 0 { // unsafe.Sizeof(int(0))
		log.Println("getDWMactive Err", hWnd)
	}
	return isNCRenderingEnabled == 1
}
