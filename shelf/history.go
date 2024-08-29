package shelf

import (
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/liuminhaw/mist-miner/locks"
)

type HistoryRecord struct {
	Dir   string
	File  io.ReadWriteCloser
	Index int
}

func NewHistoryRecord(group string, index int) (HistoryRecord, error) {
	dir, err := historyDir(group)
	if err != nil {
		return HistoryRecord{}, fmt.Errorf("NewHistoryRecord(%s): %w", group, err)
	}

	return HistoryRecord{Dir: dir, Index: index}, nil
}

// NewFile creates a new history file with the given group and index number
// and assigns it to the HistoryRecord struct File field.
func (hr *HistoryRecord) NewFile() error {
	file, err := os.Create(
		filepath.Join(hr.Dir, fmt.Sprintf("%s.%d", SHELF_HISTORY_FILE, hr.Index)),
	)
	if err != nil {
		return fmt.Errorf("NewHistoryFile(%d): %w", hr.Index, err)
	}
	hr.File = file

	return nil
}

// ReadFile reads the file from the group and index set in the HistoryRecord struct,
// extract with zlib and return the result as a io.ReadCloser
func (hr *HistoryRecord) Read() (io.ReadCloser, error) {
	if err := hr.openFile(); err != nil {
		return nil, fmt.Errorf("Read(): %w", err)
	}

	r, err := zlib.NewReader(hr.File)
	if err != nil {
		return nil, fmt.Errorf("Read(): %w", err)
	}
	defer r.Close()

	return r, nil
}

// CloseFile closes the file from the HistoryRecord struct
func (hr *HistoryRecord) CloseFile() error {
	if hr.File != nil {
		return hr.File.Close()
	}

	return nil
}

// openFile opens the history file from the group and index set in the HistoryRecord struct
func (hr *HistoryRecord) openFile() error {
	file, err := os.Open(
		filepath.Join(hr.Dir, fmt.Sprintf("%s.%d", SHELF_HISTORY_FILE, hr.Index)),
	)
	if err != nil {
		return fmt.Errorf("OpenHistoryFile(%d): %w", hr.Index, err)
	}
	hr.File = file

	return nil
}

// historyDir returns the directory path to store the history records for the given group
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

	record, err := NewHistoryRecord(group, 0)
	if err != nil {
		return fmt.Errorf("GenerateHistoryRecords(%s): %w", group, err)
	}

	if err := os.RemoveAll(record.Dir); err != nil {
		return fmt.Errorf("GenerateHistoryRecords(%s): %w", group, err)
	}
	if err := os.MkdirAll(record.Dir, os.ModePerm); err != nil {
		return fmt.Errorf("GenerateHistoryRecords(%s): %w", group, err)
	}

	reference := string(head.Reference)
	prevMark := LabelMark{}
filesLoop:
	for {
		record.NewFile()
		defer record.File.Close()

		w := zlib.NewWriter(record.File)
		defer w.Close()

		if record.Index != 0 {
			prev := fmt.Sprintf("%s%s %v\n", SHELF_HISTORY_LOGS_PREV, prevMark.Hash[:8], prevMark.TimeStamp.Format(time.RFC3339))
			_, err := w.Write([]byte(prev))
			if err != nil {
				return fmt.Errorf("GenerateHistoryRecords(%s): %w", group, err)
			}
		}
		for i := 0; i < recordsPerPage; i++ {
			mark, err := ReadMark(group, reference)
			if err != nil {
				return fmt.Errorf("GenerateHistoryRecords(%s): %w", group, err)
			}

			_, err = w.Write([]byte(fmt.Sprintf("%s %v\n", mark.Hash, mark.TimeStamp.Format(time.RFC3339))))
			if err != nil {
				return fmt.Errorf("GenerateHistoryRecords(%s): %w", group, err)
			}

			if mark.Parent == "nil" {
				break filesLoop
			} else if i == recordsPerPage-1 {
				prevMark = *mark
				reference = mark.Parent
				tmpMark, err := ReadMark(group, reference)
				if err != nil {
					return fmt.Errorf("GenerateHistoryRecords(%s): %w", group, err)
				}
				next := fmt.Sprintf("%s%s %v\n", tmpMark.Hash[:8], SHELF_HISTORY_LOGS_NEXT, tmpMark.TimeStamp.Format(time.RFC3339))
				_, err = w.Write([]byte(next))
				if err != nil {
					return fmt.Errorf("GenerateHistoryRecords(%s): %w", group, err)
				}
			} else {
				reference = mark.Parent
			}
		}
		record.Index++
	}

	return nil
}
