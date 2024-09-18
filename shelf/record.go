package shelf

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type ShelfRecord struct {
	Group string
	Type  string
	Hash  string
}

func NewObjectRecord(group, hash string) ShelfRecord {
	return newShelfRecord(group, shelf_object_dir, hash)
}

func NewDiaryRecord(group, hash string) ShelfRecord {
	return newShelfRecord(group, shelf_diary_dir, hash)
}

func newShelfRecord(group, recordType, hash string) ShelfRecord {
	return ShelfRecord{
		Group: group,
		Type:  recordType,
		Hash:  hash,
	}
}

// RecordFile returns the file path of the shelf record
func (r ShelfRecord) RecordFile() (string, error) {
	dir, err := r.recordDir()
	if err != nil {
		return "", fmt.Errorf("RecordFile(): %w", err)
	}

	return filepath.Join(dir, r.Hash[2:]), nil
}

// RecordRead returns the entire content of the shelf record, and will try to
// prettify the content if it is a JSON
func (sr ShelfRecord) RecordRead() (string, error) {
	r, err := sr.RecordReadCloser()
	if err != nil {
		return "", fmt.Errorf("RecordRead(): %w", err)
	}
	defer r.Close()

	b, err := io.ReadAll(r)
	if err != nil {
		return "", fmt.Errorf("RecordRead(): %w", err)
	}

	var prettyJson bytes.Buffer
	if err := json.Indent(&prettyJson, b, "", "  "); err != nil {
		return string(b), nil
	} else {
		return prettyJson.String(), nil
	}
}

// RecordReadCloser returns a io.ReaderCloser of the shelf record for more control
// on the reading process
func (sr ShelfRecord) RecordReadCloser() (io.ReadCloser, error) {
	path, err := sr.RecordFile()
	if err != nil {
		return nil, fmt.Errorf("recordReader(): %w", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("recordReader(): %w", err)
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("recordReader(): %w", err)
	}
	defer f.Close()

	r, err := zlib.NewReader(f)
	if err != nil {
		return nil, fmt.Errorf("recordReader(): %w", err)
	}
	defer r.Close()

	return r, nil
}

// recordDir returns the directory path of the shelf record
func (r ShelfRecord) recordDir() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("RecordDir(): get executable: %w", err)
	}

	return filepath.Join(
		filepath.Dir(execPath),
		shelf_dir,
		r.Group,
		r.Type,
		r.Hash[:2],
	), nil
}
