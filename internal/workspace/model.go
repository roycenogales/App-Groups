package workspace

import "time"

// WindowSnapshot stores the information needed to restore a window.
type WindowSnapshot struct {
	ProcessName string `json:"process_name"`
	ProcessPath string `json:"process_path"`
	Title       string `json:"title"`
	X           int32  `json:"x"`
	Y           int32  `json:"y"`
	Width       int32  `json:"width"`
	Height      int32  `json:"height"`
	State       string `json:"state"`
}

// Workspace represents a saved group of application windows.
type Workspace struct {
	ID        string           `json:"id"`
	Name      string           `json:"name"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
	Windows   []WindowSnapshot `json:"windows"`
}
