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

func NewLock(filename string) (Lock, error) {
	filepath, err := lockPath(filename)
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
func lockPath(filename string) (string, error) {
	osType := runtime.GOOS

	var path string
	switch osType {
	case "linux":
		path = filepath.Join(LINUX_DIR_PATH, filename)
	default:
		return "", fmt.Errorf("FilePath(%s): unsupported OS: %s", filename, osType)
	}

	os.Mkdir(filepath.Dir(path), os.ModePerm)

	return path, nil
}
