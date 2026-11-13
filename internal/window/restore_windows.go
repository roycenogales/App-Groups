//go:build windows

package window

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	setWindowPos = user32.NewProc("SetWindowPos")
	showWindow   = user32.NewProc("ShowWindow")
)

const (
	swRestore  = 9
	swMinimize = 6
	swMaximize = 3

	swpNoZOrder   = 0x0004
	swpNoActivate = 0x0010
)

// RestoreWindow finds an existing application window or launches its
// executable, then restores its saved position and state.
func RestoreWindow(
	processPath string,
	title string,
	x, y, width, height int32,
	state string,
) error {
	if strings.TrimSpace(processPath) == "" {
		return fmt.Errorf("no executable path was saved for window %q", title)
	}

	// First, look for an existing window belonging to this application.
	hwnd, err := findWindow(processPath, title)

	if err != nil {
		return err
	}

	// Launch the application if its saved window is not currently present.
	if hwnd == 0 {
		command := exec.Command(processPath)

		if err := command.Start(); err != nil {
			return fmt.Errorf("launch %s: %w", processPath, err)
		}

		// Release the process handle without waiting for the application.
		if err := command.Process.Release(); err != nil {
			return fmt.Errorf("release process handle: %w", err)
		}

		hwnd, err = waitForWindow(processPath, title, 15*time.Second)
		if err != nil {
			return err
		}
	}

	// Restore the window before changing its geometry.
	showWindow.Call(hwnd, swRestore)

	result, _, callErr := setWindowPos.Call(
		hwnd,
		0,
		uintptr(int64(x)),
		uintptr(int64(y)),
		uintptr(width),
		uintptr(height),
		swpNoZOrder|swpNoActivate,
	)

	if result == 0 {
		return fmt.Errorf("position window %q: %w", title, callErr)
	}

	// Apply the saved window state.
	switch state {
	case "minimized":
		showWindow.Call(hwnd, swMinimize)
	case "maximized":
		showWindow.Call(hwnd, swMaximize)
	default:
		showWindow.Call(hwnd, swRestore)
	}

	return nil
}

// findWindow searches for a window matching the saved executable and title.
func findWindow(processPath, title string) (uintptr, error) {
	windows, err := List()
	if err != nil {
		return 0, err
	}

	for _, win := range windows {
		if strings.EqualFold(win.ProcessPath, processPath) &&
			win.Title == title {
			handle, err := strconv.ParseUint(
				strings.TrimPrefix(win.HWND, "0x"),
				16,
				64,
			)
			if err != nil {
				return 0, fmt.Errorf("parse window handle: %w", err)
			}

			return uintptr(handle), nil
		}
	}

	return 0, nil
}

// waitForWindow waits for the application to create its saved window.
func waitForWindow(
	processPath string,
	title string,
	timeout time.Duration,
) (uintptr, error) {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		hwnd, err := findWindow(processPath, title)
		if err != nil {
			return 0, err
		}

		if hwnd != 0 {
			return hwnd, nil
		}

		time.Sleep(250 * time.Millisecond)
	}

	return 0, fmt.Errorf(
		"timed out waiting for window %q from %s",
		title,
		processPath,
	)
}

// Keep unsafe imported for Windows API pointer compatibility.
var _ = unsafe.Pointer(nil)

// Ensure the Windows package remains available for this file's API calls.
var _ windows.Handle
