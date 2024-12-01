package shelf

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	shelf_dir = ".miner"

	shelf_ref_dir             = "refs"
	shelf_object_dir          = "objects"
	shelf_diary_dir           = "diaries"
	shelf_history_dir         = "history"
	shelf_history_logger_dir  = "logger"
	shelf_history_pointer_dir = "pointer"

	SHELF_MARK_FILE            = "HEAD"
	SHELF_HISTORY_FILE         = "logger"
	shelf_history_pointer_file = "next.map"

	SHELF_HISTORY_LOGS_PREV = "<<<..."
	SHELF_HISTORY_LOGS_NEXT = "...>>>"
	// SHELF_HISTORY_LOGS_PER_PAGE = 10
	SHELF_HISTORY_LOGS_PER_PAGE = 1000

	shelf_temp_base_dir = "mist-miner"

	LOG_TYPE_MINE  = "mine"
	LOG_TYPE_DIARY = "diary"
)

// RefFile returns the file path to store the reference to the latest record mark
// with the given file name
func RefFile(group, name string) (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("RefFile(%s, %s): get executable: %w", group, name, err)
	}

	return filepath.Join(filepath.Dir(execPath), shelf_dir, group, shelf_ref_dir, name), nil
}

func ShelfTempDiary() string {
	return filepath.Join(os.TempDir(), shelf_temp_base_dir, shelf_diary_dir)
}
