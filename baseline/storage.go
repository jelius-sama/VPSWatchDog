package baseline

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	DefaultBaselinePath = "./data/baselines"
	CPUBaselineFile     = "cpu_baseline.json"
)

// Storage handles filesystem operations for baselines
type Storage struct {
	BasePath string
}

// NewStorage creates a new storage handler
func NewStorage(basePath string) *Storage {
	if basePath == "" {
		basePath = DefaultBaselinePath
	}
	return &Storage{BasePath: basePath}
}

// EnsureDirectory creates the baseline directory if it doesn't exist
func (s *Storage) EnsureDirectory() error {
	return os.MkdirAll(s.BasePath, 0755)
}

// SaveCPUBaseline saves CPU baseline to filesystem
func (s *Storage) SaveCPUBaseline(baseline *CPUBaseline) error {
	if err := s.EnsureDirectory(); err != nil {
		return fmt.Errorf("failed to create baseline directory: %w", err)
	}

	data, err := baseline.Marshal()
	if err != nil {
		return fmt.Errorf("failed to marshal baseline: %w", err)
	}

	filePath := filepath.Join(s.BasePath, CPUBaselineFile)

	// Write to temp file first, then rename (atomic operation)
	tempPath := filePath + ".tmp"
	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write baseline file: %w", err)
	}

	if err := os.Rename(tempPath, filePath); err != nil {
		os.Remove(tempPath) // Clean up temp file
		return fmt.Errorf("failed to rename baseline file: %w", err)
	}

	return nil
}

// LoadCPUBaseline loads CPU baseline from filesystem
func (s *Storage) LoadCPUBaseline() (*CPUBaseline, error) {
	filePath := filepath.Join(s.BasePath, CPUBaselineFile)

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // No baseline exists yet
		}
		return nil, fmt.Errorf("failed to read baseline file: %w", err)
	}

	baseline := &CPUBaseline{}
	if err := baseline.Unmarshal(data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal baseline: %w", err)
	}

	return baseline, nil
}

// BaselineExists checks if a baseline file exists
func (s *Storage) BaselineExists() bool {
	filePath := filepath.Join(s.BasePath, CPUBaselineFile)
	_, err := os.Stat(filePath)
	return err == nil
}

// DeleteCPUBaseline removes the baseline file (useful for reset)
func (s *Storage) DeleteCPUBaseline() error {
	filePath := filepath.Join(s.BasePath, CPUBaselineFile)
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete baseline file: %w", err)
	}
	return nil
}
