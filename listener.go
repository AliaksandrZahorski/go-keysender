package main

import (
	"bytes"
	"fmt"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"

	"./consts"
	"./key"
)

var (
	mod = windows.NewLazyDLL("user32.dll")
	procKeyBd = mod.NewProc("keybd_event")
	hwnd uintptr
	// procSendMessage = mod.NewProc("SendMessageW")
)

const (
	ModAlt = 1 << iota
	ModCtrl
	ModShift
	ModWin
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
			hwnd = getWindow("GetForegroundWindow")
			if hwnd != 0 {
				fmt.Println("# hwnd:", hwnd)
				time.Sleep(time.Millisecond * 1000)
				// key.PressKey(hwnd, consts.VK_F1)
			}

			if id == 3 { // CTRL+ALT+X = Exit
				fmt.Println("CTRL+ALT+X pressed, going to cycle...")
				out = true
			}
		}

		time.Sleep(time.Millisecond * 50)
	}

	uptimeTicker := time.NewTicker(3 * time.Second)
	// dateTicker := time.NewTicker(3 * time.Second)

	for {
		select {
			case <-uptimeTicker.C:
				key.PressKey(hwnd, consts.VK_F1)
			// case <-dateTicker.C:
			// 		periodic2()
		}
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

// func SendMessage(hwnd uintptr, msg uint32, wParam, lParam uintptr) uintptr {
// 	ret, _, _ := procSendMessage.Call(
// 		uintptr(hwnd),
// 		uintptr(msg),
// 		wParam,
// 		lParam)

// 	return ret
// }