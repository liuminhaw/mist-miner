package shelf

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func HistoryDir(group string) (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("HistoryDir(%s): %w", group, err)
	}

	return filepath.Join(
		filepath.Dir(execPath),
		SHELF_DIR,
		group,
		SHELF_REF_DIR,
		SHELF_HISTORY_DIR,
	), nil
}

// NewHistoryFile creates a new history file with the given group and index
// and returns the file handler.
func NewHistoryFile(group string, idx int) (io.ReadWriteCloser, error) {
	dir, err := HistoryDir(group)
	if err != nil {
		return nil, fmt.Errorf("NewHistoryFile(%s, %d): %w", group, idx, err)
	}

	file, err := os.Create(filepath.Join(dir, fmt.Sprintf("%s.%d", SHELF_HISTORY_FILE, idx)))
	if err != nil {
		return nil, fmt.Errorf("NewHistoryFile(%s, %d): %w", group, idx, err)
	}

	return file, nil
}
