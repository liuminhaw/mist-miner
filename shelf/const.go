package shelf

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	SHELF_DIR = ".miner"
)

// shelfDir returns the directory path of the shelf
// "executableDir/.miner
func shelfDir() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("shelf dir: get executable: %w", err)
	}

	return filepath.Join(filepath.Dir(execPath), SHELF_DIR), nil
}
