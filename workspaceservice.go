package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"time"

	"changeme/internal/window"
	"changeme/internal/workspace"
)

// WorkspaceService exposes workspace operations to the frontend.
type WorkspaceService struct {
	mu sync.Mutex
}

// ListWorkspaces returns all saved workspaces.
func (s *WorkspaceService) ListWorkspaces() ([]workspace.Workspace, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return workspace.Load()
}

// SaveWorkspace saves the selected windows under a workspace name.
func (s *WorkspaceService) SaveWorkspace(
	name string,
	windows []window.Window,
) (workspace.Workspace, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	name = strings.TrimSpace(name)
	if name == "" {
		return workspace.Workspace{}, fmt.Errorf("workspace name cannot be empty")
	}

	// Load existing workspaces before adding the new one.
	workspaces, err := workspace.Load()
	if err != nil {
		return workspace.Workspace{}, err
	}

	// Generate a unique ID.
	idBytes := make([]byte, 16)
	if _, err := rand.Read(idBytes); err != nil {
		return workspace.Workspace{}, fmt.Errorf("generate workspace ID: %w", err)
	}

	now := time.Now()

	saved := workspace.Workspace{
		ID:        hex.EncodeToString(idBytes),
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
		Windows:   make([]workspace.WindowSnapshot, 0, len(windows)),
	}

	// Copy only the information needed to restore each window.
	for _, win := range windows {
		saved.Windows = append(saved.Windows, workspace.WindowSnapshot{
			ProcessName: win.ProcessName,
			ProcessPath: win.ProcessPath,
			Title:       win.Title,
			X:           win.X,
			Y:           win.Y,
			Width:       win.Width,
			Height:      win.Height,
			State:       win.State,
		})
	}

	workspaces = append(workspaces, saved)

	if err := workspace.Save(workspaces); err != nil {
		return workspace.Workspace{}, err
	}

	return saved, nil
}

// DeleteWorkspace removes a workspace by its ID.
func (s *WorkspaceService) DeleteWorkspace(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	workspaces, err := workspace.Load()
	if err != nil {
		return err
	}

	updated := make([]workspace.Workspace, 0, len(workspaces))
	found := false

	for _, item := range workspaces {
		if item.ID == id {
			found = true
			continue
		}

		updated = append(updated, item)
	}

	if !found {
		return fmt.Errorf("workspace not found: %s", id)
	}

	return workspace.Save(updated)
}

// RestoreWorkspace launches applications and restores their saved windows.
func (s *WorkspaceService) RestoreWorkspace(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	workspaces, err := workspace.Load()
	if err != nil {
		return err
	}

	var selected *workspace.Workspace

	for i := range workspaces {
		if workspaces[i].ID == id {
			selected = &workspaces[i]
			break
		}
	}

	if selected == nil {
		return fmt.Errorf("workspace not found: %s", id)
	}

	for _, saved := range selected.Windows {
		err := window.RestoreWindow(
			saved.ProcessPath,
			saved.Title,
			saved.X,
			saved.Y,
			saved.Width,
			saved.Height,
			saved.State,
		)

		if err != nil {
			return fmt.Errorf(
				"restore window %q: %w",
				saved.Title,
				err,
			)
		}
	}

	return nil
}
