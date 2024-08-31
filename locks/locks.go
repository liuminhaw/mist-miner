package locks

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gofrs/flock"
)

type Lock struct {
	Flock *flock.Flock
	Path  string
}

// NewLock creates a new Lock object with the given group and filename.
// the lock file will be created in the appropriate directory based on the OS.
// Example of lock file path will be like: OS_DIR_PATH/GROUP/FILENAME
// If group is empty, the lock file will be created in the root directory.
// Example of lock file path will be like: OS_DIR_PATH/FILENAME
func NewLock(group, filename string) (Lock, error) {
	filepath, err := lockPath(group, filename)
	if err != nil {
		return Lock{}, fmt.Errorf("NewLock(%s): %w", filename, err)
	}

	f := flock.New(filepath)

	return Lock{Flock: f, Path: filepath}, nil
}

// TryLock acquires an exclusive lock on the file. Is a wrapper for flock.Lock().
func (l *Lock) TryLock() (bool, error) {
	return l.Flock.TryLock()
}

// TryRLock acquires a shared lock on the file. Is a wrapper for flock.RLock().
func (l *Lock) TryRLock() (bool, error) {
	return l.Flock.TryRLock()
}

// Unlock releases the lock on the file. Is a wrapper for flock.Unlock().
func (l *Lock) Unlock() error {
	return l.Flock.Unlock()
}

// lockPath returns the full path of a lock file based on the OS and given filename.
func lockPath(group, filename string) (string, error) {
	osType := runtime.GOOS

	var path string
	switch osType {
	case "linux":
		if group == "" {
			path = filepath.Join(LINUX_DIR_PATH, filename)
		} else {
			path = filepath.Join(LINUX_DIR_PATH, group, filename)
		}
	default:
		return "", fmt.Errorf("lockPath(%s, %s): unsupported OS: %s", group, filename, osType)
	}

	os.MkdirAll(filepath.Dir(path), os.ModePerm)

	return path, nil
}
