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

const (
	SHELF_DIR = ".miner"

	SHELF_REF_DIR     = "refs"
	SHELF_OBJECT_DIR  = "objects"
	SHELF_HISTORY_DIR = "history"

	SHELF_MARK_FILE    = "HEAD"
	SHELF_HISTORY_FILE = "logger"

	SHELF_HISTORY_LOGS_PREV = "<<<..."
	SHELF_HISTORY_LOGS_NEXT = "...>>>"
	// SHELF_HISTORY_LOGS_PER_PAGE = 10
	SHELF_HISTORY_LOGS_PER_PAGE = 1000
)

// ObjectDir returns the directory path to store all object type of records
// - type 1: marks
//   - Record of each mine execution, in the format of
//     TIMESTAMP
//     PARENT MARK HASH
//     ACCORDING MAP HASH 1
//     ACCORDING MAP HASH 2
//     ...
//
// - type 2: maps
//   - A list of stuff records, each module will have it own stuff record, in the format of
//     ACCORDING STUFF HASH 1 <SPACE> STUFF IDENTIFIER
//     ACCORDING STUFF HASH 2 <SPACE> STUFF IDENTIFIER
//     ...
//
// - type 3: stuff
//   - Record of fetched stuff information,
//     in json format of struct shared.MinerResource
func ObjectDir(group, prefixBytes string) (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("object dir: get executable: %w", err)
	}

	return filepath.Join(
		filepath.Dir(execPath),
		SHELF_DIR,
		group,
		SHELF_OBJECT_DIR,
		prefixBytes,
	), nil
}

// ObjectFile returns the file path to store the object record with the given file name
func ObjectFile(group, objectHash string) (string, error) {
	dir, err := ObjectDir(group, objectHash[:2])
	if err != nil {
		return "", fmt.Errorf("object file: %w", err)
	}

	return filepath.Join(dir, objectHash[2:]), nil
}

// ObjectRead returns the content of the object with the given hash value
func ObjectRead(group, objectHash string) (string, error) {
	r, err := objectReader(group, objectHash)
	if err != nil {
		return "", fmt.Errorf("ObjectRead(%s, %s): %w", group, objectHash, err)
	}

	b, err := io.ReadAll(r)
	if err != nil {
		return "", fmt.Errorf("ObjectRead(%s, %s): %w", group, objectHash, err)
	}

	var prettyJson bytes.Buffer
	if err := json.Indent(&prettyJson, b, "", "  "); err != nil {
		return string(b), nil
	} else {
		return prettyJson.String(), nil
	}
}

// RefFile returns the file path to store the reference to the latest record mark
// with the given file name
func RefFile(group, name string) (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("RefFile(%s, %s): get executable: %w", group, name, err)
	}

	return filepath.Join(filepath.Dir(execPath), SHELF_DIR, group, SHELF_REF_DIR, name), nil
}

// objectReader returns a io ReaderCloser to the object with the given group and hash value
func objectReader(group, objectHash string) (io.ReadCloser, error) {
	objectFile, err := ObjectFile(group, objectHash)
	if err != nil {
		return nil, fmt.Errorf("objectReader(%s, %s): %w", group, objectHash, err)
	}

	if _, err := os.Stat(objectFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("objectReader(%s, %s): %w", group, objectHash, err)
	}

	f, err := os.Open(objectFile)
	if err != nil {
		return nil, fmt.Errorf("objectReader(%s, %s): %w", group, objectHash, err)
	}
	defer f.Close()

	r, err := zlib.NewReader(f)
	if err != nil {
		return nil, fmt.Errorf("objectReader(%s, %s): %w", group, objectHash, err)
	}
	defer r.Close()

	return r, nil
}
