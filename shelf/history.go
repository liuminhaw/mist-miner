package shelf

import (
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/liuminhaw/mist-miner/locks"
)

type HistoryRecord struct {
	Dir  string
	File io.ReadWriteCloser
}

func NewHistoryRecord(group string) (HistoryRecord, error) {
	dir, err := historyDir(group)
	if err != nil {
		return HistoryRecord{}, fmt.Errorf("NewHistoryRecord(%s): %w", group, err)
	}

	return HistoryRecord{Dir: dir}, nil
}

// NewFile creates a new history file with the given group and index number
// and assigns it to the HistoryRecord struct File field.
func (hr *HistoryRecord) NewFile(idx int) error {
	file, err := os.Create(filepath.Join(hr.Dir, fmt.Sprintf("%s.%d", SHELF_HISTORY_FILE, idx)))
	if err != nil {
		return fmt.Errorf("NewHistoryFile(%d): %w", idx, err)
	}
	hr.File = file

	return nil
}

func historyDir(group string) (string, error) {
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

// GenerateHistoryRecords generates history records files for the given group
// with the given number of records per log file.
// The file will be stored sequentially with the format of: SHELF_HISTORY_FILE.<index>
// Will use flock to prevent writing simultaneously. Return locks.ErrIsLocked if file lock is not acquired.
func GenerateHistoryRecords(group string, recordsPerPage int) error {
	head, err := NewRefMark(SHELF_MARK_FILE, group)
	if err != nil {
		return fmt.Errorf("GenerateHistoryRecords(%s): %w", group, err)
	}

	// Use flock to prevent other processes from writing to the file
	fileLock, err := locks.NewLock(locks.HISTORY_LOCKFILE)
	if err != nil {
		return fmt.Errorf("GenerateHistoryRecords(%s): %w", group, err)
	}
	locked, err := fileLock.TryLock()
	if err != nil {
		return fmt.Errorf("GenerateHistoryRecords(%s): %w", group, err)
	}
	defer fileLock.Unlock()

	if !locked {
		return locks.ErrIsLocked
	}

	record, err := NewHistoryRecord(group)
	if err != nil {
		return fmt.Errorf("GenerateHistoryRecords(%s): %w", group, err)
	}

	if err := os.RemoveAll(record.Dir); err != nil {
		return fmt.Errorf("GenerateHistoryRecords(%s): %w", group, err)
	}
	if err := os.MkdirAll(record.Dir, os.ModePerm); err != nil {
		return fmt.Errorf("GenerateHistoryRecords(%s): %w", group, err)
	}

	filesIdx := 0
	reference := string(head.Reference)
filesLoop:
	for {
		record.NewFile(filesIdx)
		defer record.File.Close()

		w := zlib.NewWriter(record.File)
		defer w.Close()

		for i := 0; i < recordsPerPage; i++ {
			mark, err := ReadMark(group, reference)
			if err != nil {
				return fmt.Errorf("GenerateHistoryRecords(%s): %w", group, err)
			}

			_, err = w.Write([]byte(fmt.Sprintf("%s\n", mark.Hash)))
			if err != nil {
				return fmt.Errorf("GenerateHistoryRecords(%s): %w", group, err)
			}

			if mark.Parent == "nil" {
				break filesLoop
			} else if i == recordsPerPage-1 {
				_, err = w.Write([]byte("more...\n"))
				if err != nil {
					return fmt.Errorf("GenerateHistoryRecords(%s): %w", group, err)
				}
				reference = mark.Parent
			} else {
				reference = mark.Parent
			}
		}
		filesIdx++
	}

	return nil
}
