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
		shelf_dir,
		group,
		shelf_ref_dir,
		shelf_history_dir,
		shelf_history_logger_dir,
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

	// Use flock to prevent other processes from writing to the files
	objFileLock, err := locks.NewLock("", locks.OBJECTS_LOCKFILE)
	if err != nil {
		return fmt.Errorf("GenerateHistoryRecords(%s): %w", group, err)
	}
	locked, err := objFileLock.TryRLock()
	if err != nil {
		return fmt.Errorf("GenerateHistoryRecords(%s): %w", group, err)
	}
	defer objFileLock.Unlock()

	histFileLock, err := locks.NewLock(group, locks.HISTORY_LOCKFILE)
	if err != nil {
		return fmt.Errorf("GenerateHistoryRecords(%s): %w", group, err)
	}
	locked, err = histFileLock.TryLock()
	if err != nil {
		return fmt.Errorf("GenerateHistoryRecords(%s): %w", group, err)
	}
	defer histFileLock.Unlock()

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
			prev := fmt.Sprintf("%s%s %s %v\n", SHELF_HISTORY_LOGS_PREV, prevMark.Hash[:8], prevMark.LogType, prevMark.TimeStamp.Format(time.RFC3339))
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

			_, err = w.Write([]byte(fmt.Sprintf("%s %s %v\n", mark.Hash, mark.LogType, mark.TimeStamp.Format(time.RFC3339))))
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
				next := fmt.Sprintf("%s%s %s %v\n", tmpMark.Hash[:8], SHELF_HISTORY_LOGS_NEXT, tmpMark.LogType, tmpMark.TimeStamp.Format(time.RFC3339))
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

// GenerateHistoryPointers generates history pointers files for the given group.
// history pointer file is stored in the format of: <sha hash> <next log sha hash>
// which is for getting the next history hash from the first column sha hash.
func GenerateHistoryPointers(group string) error {
	head, err := NewRefMark(SHELF_MARK_FILE, group)
	if err != nil {
		return fmt.Errorf("GenerateHistoryPointers(%s): %w", group, err)
	}

	currentSha := string(head.Reference)
	for {
		mark, err := ReadMark(group, currentSha)
		if err != nil {
			return fmt.Errorf("GenerateHistoryPointers(%s): %w", group, err)
		}

		parentSha := mark.Parent
		if mark.Parent == "nil" {
			break
		}

		// next.map format: <parent sha> <current sha>
		if err := writeNextMap(group, parentSha, currentSha); err != nil {
			return fmt.Errorf("GenerateHistoryPointers(%s): %w", group, err)
		}
		currentSha = parentSha
	}

	return nil
}

// writeNextMap writes the parent sha and current sha to the history pointer file
// in the format of: <parent sha> <current sha>.
func writeNextMap(group, parentSha, currentSha string) error {
	file, err := pointerFile(group, parentSha)
	if err != nil {
		return fmt.Errorf("writeNextMap(%s, %s, %s): %w", group, parentSha, currentSha, err)
	}
	if err := os.MkdirAll(filepath.Dir(file), os.ModePerm); err != nil {
		return fmt.Errorf("writeNextMap(%s, %s, %s): %w", group, parentSha, currentSha, err)
	}

	// Read pointer file content
	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf(
			"writeNextMap(%s, %s, %s): %w", group, parentSha, currentSha, err,
		)
	}
	defer f.Close()

	var content []byte
	if _, err := os.Stat(file); os.IsExist(err) {
		r, err := zlib.NewReader(f)
		if err != nil {
			return fmt.Errorf(
				"writeNextMap(%s, %s, %s): %w", group, parentSha, currentSha, err,
			)
		}
		defer r.Close()

		// Append new
		content, err = io.ReadAll(r)
		if err != nil {
			return fmt.Errorf(
				"writeNextMap(%s, %s, %s): read all: %w",
				group,
				parentSha,
				currentSha,
				err,
			)
		}
	}
	content = append(content, []byte(fmt.Sprintf("%s %s\n", parentSha, currentSha))...)

	w := zlib.NewWriter(f)
	_, err = w.Write(content)
	if err != nil {
		return fmt.Errorf("writeNextMap(%s, %s, %s): %w", group, parentSha, currentSha, err)
	}
	defer w.Close()

	return nil
}

// pointerDir returns the directory path to store the history pointers for the given group and sha
// using the first two characters of the sha as the subdirectory
func pointerFile(group, sha string) (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("PointerDir(%s, %s): get executable: %w", group, sha, err)
	}

	return filepath.Join(
		filepath.Dir(execPath),
		shelf_dir,
		group,
		shelf_ref_dir,
		shelf_history_dir,
		shelf_history_pointer_dir,
		sha[:2],
		shelf_history_pointer_file,
	), nil
}
