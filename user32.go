package main

import (
	"syscall"
)

var (
	user32          = syscall.NewLazyDLL("user32.dll")
	procEnumWindows = user32.NewProc("EnumWindows")
)

func EnumWindows(enumFunc uintptr, lparam uintptr) (err error) {
	r1, _, e1 := syscall.Syscall(procEnumWindows.Addr(), 2, uintptr(enumFunc), uintptr(lparam), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}
