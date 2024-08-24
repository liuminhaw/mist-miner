package locks

import (
	"fmt"
	"path/filepath"
	"runtime"
)

const (
	LINUX_DIR_PATH = "/var/lock"

	HISTORY_LOCK_FILE = "mist-miner-history.lock"
)

// FilePath returns the full path of a lock file based on the OS and given filename.
func FilePath(filename string) (string, error) {
	osType := runtime.GOOS

	var path string
	switch osType {
	case "linux":
		path = filepath.Join(LINUX_DIR_PATH, filename)
	default:
		return "", fmt.Errorf("FilePath(%s): unsupported OS: %s", filename, osType)
	}

	return path, nil
}
