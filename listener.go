package main

import (
	"bytes"
	"fmt"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	mod = windows.NewLazyDLL("user32.dll")
	procKeyBd = mod.NewProc("keybd_event")
	procSendMessage = mod.NewProc("SendMessageW")
)

const (
	ModAlt = 1 << iota
	ModCtrl
	ModShift
	ModWin

	_KEYEVENTF_KEYUP       = 0x0002

	WM_KEYUP               = 0x0101
	WM_KEYDOWN             = 0x0100
	WM_CHAR                = 0x0102
	
	VK_F1                  = 0x70
	VK_F2                  = 0x71
	VK_F3                  = 0x72
	VK_F4                  = 0x73
	VK_F5                  = 0x74
	VK_F6                  = 0x75
	VK_F7                  = 0x76
	VK_F8                  = 0x77
	VK_F9                  = 0x78
	VK_F10                 = 0x79
	VK_F11                 = 0x7A
	VK_F12                 = 0x7B

	VK_A                   = 0x41
)

type Hotkey struct {
	Id        int // Unique id
	Modifiers int // Mask of modifiers
	KeyCode   int // Key code, e.g. 'A'
}

// String returns a human-friendly display name of the hotkey
// such as "Hotkey[Id: 1, Alt+Ctrl+O]"
func (h *Hotkey) String() string {
	mod := &bytes.Buffer{}
	if h.Modifiers&ModAlt != 0 {
		mod.WriteString("Alt+")
	}
	if h.Modifiers&ModCtrl != 0 {
		mod.WriteString("Ctrl+")
	}
	if h.Modifiers&ModShift != 0 {
		mod.WriteString("Shift+")
	}
	if h.Modifiers&ModWin != 0 {
		mod.WriteString("Win+")
	}
	return fmt.Sprintf("Hotkey[Id: %d, %s%c]", h.Id, mod, h.KeyCode)
}

func main() {
	user32 := syscall.MustLoadDLL("user32")
	defer user32.Release()

	reghotkey := user32.MustFindProc("RegisterHotKey")

	// Hotkeys to listen to:
	keys := map[int16]*Hotkey{
		1: &Hotkey{1, ModAlt + ModCtrl, 'O'},  // ALT+CTRL+O
		2: &Hotkey{2, ModAlt + ModShift, 'M'}, // ALT+SHIFT+M
		3: &Hotkey{3, ModAlt + ModCtrl, 'X'},  // ALT+CTRL+X
	}

	// Register hotkeys:
	for _, v := range keys {
		r1, _, err := reghotkey.Call(
			0, uintptr(v.Id), uintptr(v.Modifiers), uintptr(v.KeyCode))
		if r1 == 1 {
			fmt.Println("Registered", v)
		} else {
			fmt.Println("Failed to register", v, ", error:", err)
		}
	}

	peekmsg := user32.MustFindProc("PeekMessageW")

	out := false
	for out == false {
		var msg = &MSG{}
		peekmsg.Call(uintptr(unsafe.Pointer(msg)), 0, 0, 0, 1)

		// Registered id is in the WPARAM field:
		if id := msg.WPARAM; id != 0 {
			fmt.Println("Hotkey pressed:", keys[id])

			// hwnd of current window
			if hwnd := getWindow("GetForegroundWindow"); hwnd != 0 {
				fmt.Println("# hwnd:", hwnd)
				time.Sleep(time.Millisecond * 1000)
				fmt.Println("sending:", SendMessage(hwnd, VK_A, 0, 0))
				SendMessage(hwnd, WM_KEYDOWN, VK_F1, 0)
				time.Sleep(time.Millisecond * 50)
				SendMessage(hwnd, WM_KEYUP, VK_F1, 0)
			}

			if id == 3 { // CTRL+ALT+X = Exit
				fmt.Println("CTRL+ALT+X pressed, goodbye...")
				out = true
				// return
			}
		}

		time.Sleep(time.Millisecond * 50)
	}

	for i := 0; i < 3; i++ {
		fmt.Println("cycle")
		time.Sleep(time.Millisecond * 1000)
	}
}

type MSG struct {
	HWND   uintptr
	UINT   uintptr
	WPARAM int16
	LPARAM int64
	DWORD  int32
	POINT  struct{ X, Y int64 }
}

func getWindow(funcName string) uintptr {
	proc := mod.NewProc(funcName)
	hwnd, _, _ := proc.Call()
	return hwnd
}

func SendMessage(hwnd uintptr, msg uint32, wParam, lParam uintptr) uintptr {
	ret, _, _ := procSendMessage.Call(
		uintptr(hwnd),
		uintptr(msg),
		wParam,
		lParam)

	return ret
}