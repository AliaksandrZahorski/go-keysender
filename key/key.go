package key

import (
	"time"
	"golang.org/x/sys/windows"

	"../consts"
)

var (
	mod = windows.NewLazyDLL("user32.dll")
	procSendMessage = mod.NewProc("SendMessageW")
)


func PressKey(hwnd uintptr, key uintptr) uintptr {
	SendMessage(hwnd, consts.WM_KEYDOWN, key, 0)
	time.Sleep(time.Millisecond * 50)
	ret := SendMessage(hwnd, consts.WM_KEYUP, key, 0)
	return ret
}

func SendMessage(hwnd uintptr, msg uint32, wParam, lParam uintptr) uintptr {
	ret, _, _ := procSendMessage.Call(
		uintptr(hwnd),
		uintptr(msg),
		wParam,
		lParam)

	return ret
}
