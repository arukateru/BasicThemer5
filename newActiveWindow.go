package main

import (
	"log"
	"sync"

	"github.com/lxn/win"
)

var listener *Listener

type Listener struct {
	AllPIDs                 map[uint32]struct{}
	EVENT_OBJECT_CREATE     win.HWINEVENTHOOK
	EVENT_SYSTEM_FOREGROUND win.HWINEVENTHOOK
	EVENT_OBJECT_DESTROY    win.HWINEVENTHOOK
	mutex                   sync.Mutex
}

const (
	OBJID_WINDOW            = 0
	CHILDID_SELF            = 0
	WM_APPEXIT              = 0x0400 + 1
	ProcessQueryInformation = 0x0400 // https://docs.microsoft.com/en-us/windows/win32/procthread/process-security-and-access-rights
)

// newActiveWindowCallback is passed to Windows to be called whenever the active window changes.
// When it is called it will attempt to find the process of an associated handle, then get the executable associated with that.
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nc-winuser-wineventproc
func (l *Listener) newActiveWindowCallback(
	hWinEventHook win.HWINEVENTHOOK, // Handle to an event hook function. This value is returned by SetWinEventHook when the hook function is installed and is specific to each instance of the hook function.
	event uint32, // Specifies the event that occurred. This value is one of the event constants.
	hwnd win.HWND, // Handle to the window that generates the event, or NULL if no window is associated with the event. For example, the mouse pointer is not associated with a window.
	idObject int32, // Identifies the object associated with the event. This is one of the object identifiers or a custom object ID.
	idChild int32, // Identifies whether the event was triggered by an object or a child element of the object. If this value is CHILDID_SELF, the event was triggered by the object; otherwise, this value is the child ID of the element that triggered the event.
	idEventThread uint32,
	dwmsEventTime uint32, // Specifies the time, in milliseconds, that the event was generated.
) (ret uintptr) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if idObject != OBJID_WINDOW || idChild != CHILDID_SELF {
		return
	}

	if event == win.EVENT_OBJECT_CREATE {
		return
	}

	if hwnd == 0 {
		return
	}

	if getDWMactive(hwnd) {
		// https://docs.microsoft.com/en-us/windows/win32/api/dwmapi/nf-dwmapi-dwmsetwindowattribute
		var policyParameter = DWMNCRP_DISABLED
		DwmSetWindowAttribute(hwnd, DWMWA_NCRENDERING_POLICY, &policyParameter, 8)          // unsafe.Sizeof(int(0))
		DwmSetWindowAttribute(hwnd, DWMWA_FORCE_ICONIC_REPRESENTATION, &policyParameter, 8) // unsafe.Sizeof(int(0))
	}

	return 0
}

func startListenerMessageLoop() {
	var err error
	EVENT_OBJECT_CREATE, err := setActiveWindowWinEventHook(listener.newActiveWindowCallback, win.EVENT_OBJECT_CREATE)
	if err != nil {
		_ = err
		panic(err)
	}
	defer win.UnhookWinEvent(EVENT_OBJECT_CREATE)

	EVENT_SYSTEM_FOREGROUND, err := setActiveWindowWinEventHook(listener.newActiveWindowCallback, win.EVENT_SYSTEM_FOREGROUND)
	if err != nil {
		panic(err)
	}
	defer win.UnhookWinEvent(EVENT_SYSTEM_FOREGROUND)

	var msg win.MSG
	for win.GetMessage(&msg, 0, 0, 0) != 0 {
		log.Println(msg)
		if msg.Message == WM_APPEXIT {
			break
		}
		win.TranslateMessage(&msg)
		win.DispatchMessage(&msg)
	}
}

// setActiveWindowWinEventHook is for informing windows which function should be called whenever a
// foreground window has changed. https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-setwineventhook
func setActiveWindowWinEventHook(callbackFunction win.WINEVENTPROC, event uint32) (win.HWINEVENTHOOK, error) {
	ret, err := win.SetWinEventHook(
		event,            // Specifies the event constant for the lowest event value in the range of events that are handled by the hook function. This parameter can be set to EVENT_MIN to indicate the lowest possible event value.
		event,            // Specifies the event constant for the highest event value in the range of events that are handled by the hook function. This parameter can be set to EVENT_MAX to indicate the highest possible event value.
		0,                // Handle to the DLL that contains the hook function at lpfnWinEventProc, if the WINEVENT_INCONTEXT flag is specified in the dwFlags parameter. If the hook function is not located in a DLL, or if the WINEVENT_OUTOFCONTEXT flag is specified, this parameter is NULL.
		callbackFunction, // Pointer to the event hook function. For more information about this function, see WinEventProc.
		0,                // Specifies the ID of the process from which the hook function receives events. Specify zero (0) to receive events from all processes on the current desktop.
		0,                // Specifies the ID of the thread from which the hook function receives events. If this parameter is zero, the hook function is associated with all existing threads on the current desktop.
		win.WINEVENT_OUTOFCONTEXT|win.WINEVENT_SKIPOWNPROCESS, // Flag values that specify the location of the hook function and of the events to be skipped. The following flags are valid:
	)
	if ret == 0 {
		return 0, err
	}

	return ret, nil
}
