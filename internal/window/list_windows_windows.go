//go:build windows

package window

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32 = syscall.NewLazyDLL("user32.dll")

	getWindowTextW = user32.NewProc("GetWindowTextW")
	getWindowRect  = user32.NewProc("GetWindowRect")
	isIconic       = user32.NewProc("IsIconic")
	isZoomed       = user32.NewProc("IsZoomed")
)

// List returns visible desktop windows, excluding this application itself.
func List() ([]Window, error) {
	var results []Window

	callback := syscall.NewCallback(func(hwnd uintptr, _ uintptr) uintptr {
		handle := windows.HWND(hwnd)

		// Ignore invalid or invisible windows.
		if hwnd == 0 || !windows.IsWindowVisible(handle) {
			return 1
		}

		// Read the window title.
		titleBuffer := make([]uint16, 1024)

		titleLength, _, _ := getWindowTextW.Call(
			hwnd,
			uintptr(unsafe.Pointer(&titleBuffer[0])),
			uintptr(len(titleBuffer)),
		)

		if titleLength == 0 {
			return 1
		}

		title := syscall.UTF16ToString(titleBuffer[:titleLength])

		// Get the window's screen coordinates.
		var rect windows.Rect

		rectResult, _, _ := getWindowRect.Call(
			hwnd,
			uintptr(unsafe.Pointer(&rect)),
		)

		if rectResult == 0 {
			return 1
		}

		// Identify the process owning this window.
		var pid uint32

		_, err := windows.GetWindowThreadProcessId(handle, &pid)
		if err != nil || pid == 0 {
			return 1
		}

		// Do not include App-Groups itself in its own workspace.
		if pid == uint32(os.Getpid()) {
			return 1
		}

		processPath := ""
		processName := ""

		// Retrieve the executable path when Windows permits access.
		processHandle, err := windows.OpenProcess(
			windows.PROCESS_QUERY_LIMITED_INFORMATION,
			false,
			pid,
		)

		if err == nil {
			defer windows.CloseHandle(processHandle)

			pathBuffer := make([]uint16, 32768)
			pathLength := uint32(len(pathBuffer))

			err = windows.QueryFullProcessImageName(
				processHandle,
				0,
				&pathBuffer[0],
				&pathLength,
			)

			if err == nil {
				processPath = syscall.UTF16ToString(
					pathBuffer[:pathLength],
				)
				processName = filepath.Base(processPath)
			}
		}

		if processName == "" {
			processName = fmt.Sprintf("PID %d", pid)
		}

		// Determine the window's current state.
		state := "normal"

		if minimized, _, _ := isIconic.Call(hwnd); minimized != 0 {
			state = "minimized"
		} else if maximized, _, _ := isZoomed.Call(hwnd); maximized != 0 {
			state = "maximized"
		}

		results = append(results, Window{
			HWND:        fmt.Sprintf("0x%X", hwnd),
			PID:         pid,
			ProcessName: processName,
			ProcessPath: processPath,
			Title:       title,
			X:           rect.Left,
			Y:           rect.Top,
			Width:       rect.Right - rect.Left,
			Height:      rect.Bottom - rect.Top,
			State:       state,
		})

		return 1
	})

	// Ask Windows to enumerate top-level desktop windows.
	if err := windows.EnumWindows(callback, nil); err != nil {
		return nil, fmt.Errorf("enumerate desktop windows: %w", err)
	}

	return results, nil
}
