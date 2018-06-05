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
	SendMessage(hwnd, consts.WM_KEYDOWN, key, lParamDown(key))
	time.Sleep(time.Millisecond * 100)
	ret := SendMessage(hwnd, consts.WM_KEYUP, key, lParamUp(key))
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

func lParamDown(key uintptr) uintptr {
	var repeatCount = 0;//                         always 1 for a WM_KEYUP
	// var scanCode = 1;// 0x70 for VK_F1
	var extended = 0;//Alt code
	var context = 0;//always 0 for a WM_KEYDOWN           0 for a WM_KEYUP
	var previousState = 0;//                              1 for a WM_KEYUP
	var transition = 0;//always 0 for a WM_KEYDOWN        1 for a WM_KEYUP

	lParam := repeatCount | (int(key) << 16) | (extended << 24) | (context << 29) | (previousState << 30) | (transition << 31);
	return uintptr(lParam)
}

func lParamUp(key uintptr) uintptr {
	var repeatCount = 1;//                         always 1 for a WM_KEYUP
	// var scanCode = 0x30;// 0x70 for VK_F1
	var extended = 0;//Alt code
	var context = 0;//always 0 for a WM_KEYDOWN           0 for a WM_KEYUP
	var previousState = 1;//                              1 for a WM_KEYUP
	var transition = 1;//always 0 for a WM_KEYDOWN        1 for a WM_KEYUP

	lParam := repeatCount | (int(key) << 16) | (extended << 24) | (context << 29) | (previousState << 30) | (transition << 31);
	return uintptr(lParam)
}
