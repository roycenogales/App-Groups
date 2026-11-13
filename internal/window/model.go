package window

// Window describes a visible desktop window and its current geometry.
type Window struct {
	HWND        string `json:"hwnd"`
	PID         uint32 `json:"pid"`
	ProcessName string `json:"process_name"`
	ProcessPath string `json:"process_path"`
	Title       string `json:"title"`
	X           int32  `json:"x"`
	Y           int32  `json:"y"`
	Width       int32  `json:"width"`
	Height      int32  `json:"height"`
	State       string `json:"state"`
}
