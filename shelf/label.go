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
	Hash string
	// Module   string
	Group  string
	Maps   []IdentifierHashMap
	buffer bytes.Buffer
}

// Sort sorts the IdentifierHashMaps by the hash field.
func (ihm *IdentifierHashMaps) Sort() {
	slices.SortStableFunc(ihm.Maps, func(a, b IdentifierHashMap) int {
		if a.Identifier == b.Identifier {
			return strings.Compare(a.Hash, b.Hash)
		}
		return strings.Compare(a.Identifier, b.Identifier)
	})
}

func (lhm *IdentifierHashMaps) Write() error {
	err := lhm.calcHash()
	if err != nil {
		return fmt.Errorf("identifier hash maps write: calc hash: %w", err)
	}

	mapFile, err := objectFile(lhm.Group, lhm.Hash)
	if err != nil {
		return fmt.Errorf("identifier hash maps write: %w", err)
	}
	if _, err := os.Stat(mapFile); !errors.Is(err, os.ErrNotExist) {
		fmt.Printf("Identifier hash maps file already exists: %s\n", mapFile)
		return nil
	}

	err = os.MkdirAll(filepath.Dir(mapFile), os.ModePerm)
	if err != nil {
		return fmt.Errorf("identifier hash maps write: mkdir: %w", err)
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

type refMark struct {
	name      string
	group     string
	reference []byte
}

// write writes the reference to the HEAD file.
// reference is the hash of the latest label mark.
func (m *refMark) write() error {
	refFile, err := refFile(m.group, m.name)
	if err != nil {
		return fmt.Errorf("ref mark write: %w", err)
	}

	err = os.MkdirAll(filepath.Dir(refFile), os.ModePerm)
	if err != nil {
		return fmt.Errorf("ref mark write: mkdir: %w", err)
	}

	f, err := os.Create(refFile)
	if err != nil {
		return fmt.Errorf("ref mark write: create file: %w", err)
	}
	defer f.Close()

	_, err = f.Write(m.reference)
	if err != nil {
		return fmt.Errorf("ref mark write: write file: %w", err)
	}

	return nil
}

// currentRef returns the reference of headMark
// which is the hash of the latest label mark.
func (m *refMark) currentRef() ([]byte, error) {
	refFile, err := refFile(m.group, m.name)
	if err != nil {
		return nil, fmt.Errorf("ref mark current ref: %w", err)
	}

	if _, err := os.Stat(refFile); errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("ref mark read: file not found: %w", err)
	}

	f, err := os.Open(refFile)
	if err != nil {
		return nil, fmt.Errorf("ref mark read: open file: %w", err)
	}
	defer f.Close()

	b := make([]byte, sha256.BlockSize)
	_, err = f.Read(b)
	if err != nil {
		return nil, fmt.Errorf("ref mark read: read file: %w", err)
	}

	return b, nil
}

type MarkMapping struct {
	Module string
	Hash   string
}

type LabelMark struct {
	// Module       string
	Hash      string
	TimeStamp time.Time
	Parent    string
	Mappings  []MarkMapping
	Group     string
	// LabelMapHash string
	buffer bytes.Buffer
}

// NewMark creates a new label mark with the given plugin name, plugin id and label map hash.
func NewMark(plugName, plugId string, mapHash string) (*LabelMark, error) {
	mark := LabelMark{
		Group:     plugId,
		TimeStamp: time.Now(),
		Mappings:  []MarkMapping{},
	}

	head := refMark{
		name:  "HEAD",
		group: plugId,
	}
	headFile, err := refFile(head.group, head.name)
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

// AddMapping adds a new mark mapping to the label mark.
func (lm *LabelMark) AddMapping(module, hash string) {
	lm.Mappings = append(lm.Mappings, MarkMapping{
		Module: module,
		Hash:   hash,
	})
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

	markFile, err := objectFile(lm.Group, lm.Hash)
	if err != nil {
		return fmt.Errorf("label mark update: %w", err)
	}
	if _, err := os.Stat(markFile); !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("Label mark file already exists, there maybe a collision: %s\n", markFile)
	}

	markDir := filepath.Dir(markFile)
	err = os.MkdirAll(markDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("label mark update: mkdir: %w", err)
	}

	fmt.Printf("Label mark buffer: %s\n", lm.buffer.String())

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
	head := refMark{
		name:      "HEAD",
		group:     lm.Group,
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

	// fmt.Printf("Parent: %s\n", parent)
	fmt.Fprintf(&lm.buffer, "%v\n", lm.TimeStamp)
	fmt.Fprintf(&lm.buffer, "%s\n", parent)

	lm.sort()
	for _, m := range lm.Mappings {
		fmt.Fprintf(&lm.buffer, "%s %s\n", m.Hash, m.Module)
	}

	h := sha256.New()
	_, err := h.Write(lm.buffer.Bytes())
	if err != nil {
		return fmt.Errorf("calc hash: %w", err)
	}
	lm.Hash = fmt.Sprintf("%x", h.Sum(nil))

	return nil
}

// sort sorts the MarkMappings slice in LabelMark
// first by the module field then by the hash field.
func (lm *LabelMark) sort() {
	slices.SortStableFunc(lm.Mappings, func(a, b MarkMapping) int {
		if a.Module == b.Module {
			return strings.Compare(a.Hash, b.Hash)
		}
		return strings.Compare(a.Module, b.Module)
	})
}

