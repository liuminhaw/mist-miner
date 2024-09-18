package shelf

import (
	"compress/zlib"
	"errors"
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
	if err := objFileLock.TryRLock(); err != nil {
		if errors.Is(err, locks.ErrIsLocked) {
			return err
		}
		return fmt.Errorf("GenerateHistoryRecords(%s): %w", group, err)
	}
	defer objFileLock.Unlock()

	histFileLock, err := locks.NewLock(group, locks.HISTORY_LOCKFILE)
	if err != nil {
		return fmt.Errorf("GenerateHistoryRecords(%s): %w", group, err)
	}
	if err := histFileLock.TryLock(); err != nil {
		if errors.Is(err, locks.ErrIsLocked) {
			return err
		}
		return fmt.Errorf("GenerateHistoryRecords(%s): %w", group, err)
	}
	defer histFileLock.Unlock()

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

	// flock to prevent writing simultaneously
	objFileLock, err := locks.NewLock("", locks.OBJECTS_LOCKFILE)
	if err != nil {
		return fmt.Errorf("GenerateHistoryPointers(%s): %w", group, err)
	}
	if err := objFileLock.TryRLock(); err != nil {
		if errors.Is(err, locks.ErrIsLocked) {
			return err
		}
		return fmt.Errorf("GenerateHistoryPointers(%s): %w", group, err)
	}
	defer objFileLock.Unlock()

	ptrFileLock, err := locks.NewLock(group, locks.HISTORY_POINTER_LOCKFILE)
	if err != nil {
		return fmt.Errorf("GenerateHistoryPointers(%s): %w", group, err)
	}
	if err := ptrFileLock.TryLock(); err != nil {
		if errors.Is(err, locks.ErrIsLocked) {
			return err
		}
		return fmt.Errorf("GenerateHistoryPointers(%s): %w", group, err)
	}
	defer ptrFileLock.Unlock()

	ptrDir, err := pointerDir(group)
	if err != nil {
		return fmt.Errorf("GenerateHistoryPointers(%s): %w", group, err)
	}
	if err := os.RemoveAll(ptrDir); err != nil {
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
		pointer := NewHistoryPointer(group, parentSha, currentSha)
		if err := pointer.WriteNextMap(); err != nil {
			return fmt.Errorf("GenerateHistoryPointers(%s): %w", group, err)
		}
		currentSha = parentSha
	}

	return nil
}

type HistoryPointer struct {
	Group       string
	ParentHash  string
	CurrentHash string
}

func NewHistoryPointer(group, parentHash, currentHash string) HistoryPointer {
	return HistoryPointer{Group: group, ParentHash: parentHash, CurrentHash: currentHash}
}

// WriteNextMap writes the parent sha and current sha to the history pointer file
// in the format of: <parent sha> <current sha>.
func (hp HistoryPointer) WriteNextMap() error {
	file, err := hp.filepath()
	if err != nil {
		return fmt.Errorf("writeNextMap: get pointer file: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(file), os.ModePerm); err != nil {
		return fmt.Errorf("writeNextMap: create dir: %w", err)
	}

	var content []byte
	if _, err := os.Stat(file); !os.IsNotExist(err) {
		content, err = hp.readFile()
		if err != nil {
			return fmt.Errorf("writeNextMap: read pointer content: %w", err)
		}
	}
	content = append(
		content,
		[]byte(fmt.Sprintf("%s %s\n", hp.ParentHash, hp.CurrentHash))...)

	f, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("writeNextMap: create pointer file: %w", err)
	}
	defer f.Close()

	w := zlib.NewWriter(f)
	_, err = w.Write(content)
	if err != nil {
		return fmt.Errorf("writeNextMap: zlib write: %w", err)
	}
	defer w.Close()

	return nil
}

// readFile reads the history pointer file and returns the content as a byte slice
func (hp HistoryPointer) readFile() ([]byte, error) {
	file, err := hp.filepath()
	if err != nil {
		return nil, fmt.Errorf("hp.readPointerFile: get pointer file: %w", err)
	}

	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("hp.readPointerFile: open pointer file: %w", err)
	}
	defer f.Close()

	r, err := zlib.NewReader(f)
	if err != nil {
		return nil, fmt.Errorf("hp.readPointerFile: zlib read: %w", err)
	}
	defer r.Close()

	content, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("hp.readPointerFile: read pointer content: %w", err)
	}

	return content, nil
}

// filepath returns the file path to store the history pointers of HistoryPointer struct
func (hp HistoryPointer) filepath() (string, error) {
	dir, err := pointerDir(hp.Group)
	if err != nil {
		return "", fmt.Errorf("hp.pointerFile: %w", err)
	}

	return filepath.Join(
		dir,
		hp.ParentHash[:2],
		shelf_history_pointer_file,
	), nil
}

// pointerDir returns the directory path to store the history pointers for the given group and sha
// using the first two characters of the sha as the subdirectory
func pointerDir(group string) (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("PointerDir: get executable: %w", err)
	}

	return filepath.Join(
		filepath.Dir(execPath),
		shelf_dir,
		group,
		shelf_ref_dir,
		shelf_history_dir,
		shelf_history_pointer_dir,
	), nil
}
