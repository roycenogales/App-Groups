package main

import "changeme/internal/window"

// WindowService exposes desktop window operations to the frontend.
type WindowService struct{}

// ListWindows returns the currently detected desktop windows.
func (s *WindowService) ListWindows() ([]window.Window, error) {
	return window.List()
}
