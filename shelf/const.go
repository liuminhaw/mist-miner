package shelf

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	SHELF_DIR = ".miner"

	SHELF_REF_DIR     = "refs"
	shelf_object_dir  = "objects"
	shelf_diary_dir   = "diaries"
	SHELF_HISTORY_DIR = "history"

	SHELF_MARK_FILE    = "HEAD"
	SHELF_HISTORY_FILE = "logger"

	SHELF_HISTORY_LOGS_PREV = "<<<..."
	SHELF_HISTORY_LOGS_NEXT = "...>>>"
	// SHELF_HISTORY_LOGS_PER_PAGE = 10
	SHELF_HISTORY_LOGS_PER_PAGE = 1000
)

// RefFile returns the file path to store the reference to the latest record mark
// with the given file name
func RefFile(group, name string) (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("RefFile(%s, %s): get executable: %w", group, name, err)
	}

	return filepath.Join(filepath.Dir(execPath), SHELF_DIR, group, SHELF_REF_DIR, name), nil
}
