package shelf

import (
	"bytes"
	"compress/zlib"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

type IdentifierHashMap struct {
	Identifier string
	Hash       string
}

type IdentifierHashMaps struct {
	Hash     string
	Module   string
	Identity string
	Maps     []IdentifierHashMap
	buffer   bytes.Buffer
}

// Sort sorts the IdentifierHashMaps by the hash field.
func (ihm *IdentifierHashMaps) Sort() {
	slices.SortStableFunc(ihm.Maps, func(a, b IdentifierHashMap) int {
		return strings.Compare(a.Identifier, b.Identifier)
	})
}

func (lhm *IdentifierHashMaps) Write() error {
	err := lhm.calcHash()
	if err != nil {
		return fmt.Errorf("identifier hash maps write: calc hash: %w", err)
	}

	mapDir, err := lhm.dir()
	if err != nil {
		return fmt.Errorf("identifier hash maps write: %w", err)
	}

	err = os.MkdirAll(mapDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("identifier hash maps write: mkdir: %w", err)
	}

	mapFile := fmt.Sprintf("%s/%s", mapDir, lhm.Hash[2:])
	if _, err := os.Stat(mapFile); !errors.Is(err, os.ErrNotExist) {
		fmt.Printf("Identifier hash maps file already exists: %s\n", mapFile)
		return nil
	}

	f, err := os.Create(mapFile)
	if err != nil {
		return fmt.Errorf("identifier hash maps write: create file: %w", err)
	}
	defer f.Close()

	w := zlib.NewWriter(f)
	_, err = w.Write(lhm.buffer.Bytes())
	if err != nil {
		return fmt.Errorf("identifier hash maps write: write file: %w", err)
	}
	defer w.Close()

	fmt.Printf("Identifier hash maps file written: %s\n", mapFile)
	return nil
}

// dir returns the directory path of IdentifierHashMaps.
func (lhm *IdentifierHashMaps) dir() (string, error) {
	sd, err := shelfDir()
	if err != nil {
		return "", fmt.Errorf("identifier hash maps dir: %w", err)
	}

	return fmt.Sprintf("%s/labels/%s/%s/maps/%s", sd, lhm.Module, lhm.Identity, lhm.Hash[:2]), nil
}

// calcHash calculates the hash of Maps in IdentifierHashMaps.
// Maps is first write to a buffer with content "hash identifier"
// and then the buffer is hashed with sha256 to get the hash value.
func (lhm *IdentifierHashMaps) calcHash() error {
	for _, m := range lhm.Maps {
		fmt.Fprintf(&lhm.buffer, "%s %s\n", m.Hash, m.Identifier)
	}

	h := sha256.New()
	_, err := h.Write(lhm.buffer.Bytes())
	if err != nil {
		return fmt.Errorf("calc hash: %w", err)
	}
	lhm.Hash = fmt.Sprintf("%x", h.Sum(nil))

	return nil
}

type headMark struct {
	module    string
	identity  string
	reference []byte
}

// file returns the path of the HEAD file.
func (hm *headMark) file() (string, error) {
	sd, err := shelfDir()
	if err != nil {
		return "", fmt.Errorf("head mark file: %w", err)
	}

	return fmt.Sprintf("%s/labels/%s/%s/marks/HEAD", sd, hm.module, hm.identity), nil
}

// write writes the reference to the HEAD file.
// reference is the hash of the latest label mark.
func (hm *headMark) write() error {
	headFile, err := hm.file()
	if err != nil {
		return fmt.Errorf("head mark write: %w", err)
	}

	err = os.MkdirAll(filepath.Dir(headFile), os.ModePerm)
	if err != nil {
		return fmt.Errorf("head mark write: mkdir: %w", err)
	}

	f, err := os.Create(headFile)
	if err != nil {
		return fmt.Errorf("head mark write: create file: %w", err)
	}
	defer f.Close()

	_, err = f.Write(hm.reference)
	if err != nil {
		return fmt.Errorf("head mark write: write file: %w", err)
	}

	return nil
}

// currentRef returns the reference of headMark
// which is the hash of the latest label mark.
func (hm *headMark) currentRef() ([]byte, error) {
	headFile, err := hm.file()
	if err != nil {
		return nil, fmt.Errorf("head mark current ref: %w", err)
	}

	if _, err := os.Stat(headFile); errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("head mark read: file not found: %w", err)
	}

	f, err := os.Open(headFile)
	if err != nil {
		return nil, fmt.Errorf("head mark read: open file: %w", err)
	}
	defer f.Close()

	b := make([]byte, sha256.BlockSize)
	_, err = f.Read(b)
	if err != nil {
		return nil, fmt.Errorf("head mark read: read file: %w", err)
	}

	return b, nil
}

type LabelMark struct {
	Module       string
	Identity     string
	Hash         string
	TimeStamp    time.Time
	Parent       string
	LabelMapHash string
	buffer       bytes.Buffer
}

// NewMark creates a new label mark with the given plugin name, plugin id and label map hash.
func NewMark(plugName, plugId string, mapHash string) (*LabelMark, error) {
	mark := LabelMark{
		Module:       plugName,
		Identity:     plugId,
		TimeStamp:    time.Now(),
		LabelMapHash: mapHash,
	}

	head := headMark{
		module:   plugName,
		identity: plugId,
	}
	headFile, err := head.file()
	if err != nil {
		return nil, fmt.Errorf("new label mark: %w", err)
	}
	if _, err := os.Stat(headFile); !errors.Is(err, os.ErrNotExist) {
		parent, err := head.currentRef()
		if err != nil {
			return nil, fmt.Errorf("new label mark: head ref: %w", err)
		}
		mark.Parent = string(parent)
	}

	return &mark, nil
}

// Update writes the label mark to a file in format:
// timestamp
// parent
// label map hash
//
// And also updates the HEAD reference to the hash of the latest label mark.
func (lm *LabelMark) Update() error {
	// Write the label mark to a file.
	err := lm.calcHash()
	if err != nil {
		return fmt.Errorf("label mark update: %w", err)
	}

	markDir, err := lm.dir()
	if err != nil {
		return fmt.Errorf("label mark update: %w", err)
	}

	err = os.MkdirAll(markDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("label mark update: mkdir: %w", err)
	}

	fmt.Printf("Label buffer: %s\n", lm.buffer.String())

	markFile := fmt.Sprintf("%s/%s", markDir, lm.Hash[2:])
	if _, err := os.Stat(markFile); !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("Label mark file already exists, there maybe a collision: %s\n", markFile)
	}

	f, err := os.Create(markFile)
	if err != nil {
		return fmt.Errorf("label mark update: create file: %w", err)
	}
	defer f.Close()

	w := zlib.NewWriter(f)
	_, err = w.Write(lm.buffer.Bytes())
	if err != nil {
		return fmt.Errorf("label mark update: write file: %w", err)
	}
	defer w.Close()
	fmt.Printf("Label mark file written: %s\n", markFile)

	// Update the HEAD reference.
	head := headMark{
		module:    lm.Module,
		identity:  lm.Identity,
		reference: []byte(lm.Hash),
	}
	if err := head.write(); err != nil {
		return fmt.Errorf("label mark update: head write: %w", err)
	}

	return nil
}

// calcHash calculates the hash of the label mark file content.
func (lm *LabelMark) calcHash() error {
	var parent string
	if lm.Parent == "" {
		parent = "nil"
	} else {
		parent = lm.Parent
	}

	fmt.Printf("Parent: %s\n", parent)
	fmt.Fprintf(&lm.buffer, "%v\n", lm.TimeStamp)
	fmt.Fprintf(&lm.buffer, "%s\n", parent)
	fmt.Fprintf(&lm.buffer, "%s\n", lm.LabelMapHash)

	h := sha256.New()
	_, err := h.Write(lm.buffer.Bytes())
	if err != nil {
		return fmt.Errorf("calc hash: %w", err)
	}
	lm.Hash = fmt.Sprintf("%x", h.Sum(nil))

	return nil
}

// dir returns the directory path of LabelMark.
func (lm *LabelMark) dir() (string, error) {
	sd, err := shelfDir()
	if err != nil {
		return "", fmt.Errorf("label mark dir: %w", err)
	}

	return fmt.Sprintf("%s/labels/%s/%s/marks/%s", sd, lm.Module, lm.Identity, lm.Hash[:2]), nil
}
