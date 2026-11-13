package workspace

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// storagePath returns the location of the workspace JSON file.
func storagePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("get user configuration directory: %w", err)
	}

	return filepath.Join(configDir, "App-Groups", "workspaces.json"), nil
}

// Load reads all saved workspaces.
func Load() ([]Workspace, error) {
	path, err := storagePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return []Workspace{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read workspaces: %w", err)
	}

	var workspaces []Workspace

	if err := json.Unmarshal(data, &workspaces); err != nil {
		return nil, fmt.Errorf("parse workspace file: %w", err)
	}

	if workspaces == nil {
		workspaces = []Workspace{}
	}

	return workspaces, nil
}

// Save writes all workspaces to disk using a temporary file.
func Save(workspaces []Workspace) error {
	path, err := storagePath()
	if err != nil {
		return err
	}

	directory := filepath.Dir(path)

	if err := os.MkdirAll(directory, 0700); err != nil {
		return fmt.Errorf("create workspace directory: %w", err)
	}

	data, err := json.MarshalIndent(workspaces, "", "    ")
	if err != nil {
		return fmt.Errorf("encode workspaces: %w", err)
	}

	tempFile, err := os.CreateTemp(directory, "workspaces-*.tmp")
	if err != nil {
		return fmt.Errorf("create temporary workspace file: %w", err)
	}

	tempPath := tempFile.Name()
	defer os.Remove(tempPath)

	if _, err := tempFile.Write(data); err != nil {
		tempFile.Close()
		return fmt.Errorf("write temporary workspace file: %w", err)
	}

	if err := tempFile.Sync(); err != nil {
		tempFile.Close()
		return fmt.Errorf("sync temporary workspace file: %w", err)
	}

	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("close temporary workspace file: %w", err)
	}

	if err := os.Rename(tempPath, path); err != nil {
		return fmt.Errorf("replace workspace file: %w", err)
	}

	return nil
}
