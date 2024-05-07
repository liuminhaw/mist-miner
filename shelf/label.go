package shelf

import (
	"bufio"
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

	mapFile, err := ObjectFile(lhm.Group, lhm.Hash)
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

type RefMark struct {
	Name      string
	Group     string
	Reference []byte
}

// write writes the reference to the HEAD file.
// reference is the hash of the latest label mark.
func (m *RefMark) write() error {
	RefFile, err := RefFile(m.Group, m.Name)
	if err != nil {
		return fmt.Errorf("ref mark write: %w", err)
	}

	err = os.MkdirAll(filepath.Dir(RefFile), os.ModePerm)
	if err != nil {
		return fmt.Errorf("ref mark write: mkdir: %w", err)
	}

	f, err := os.Create(RefFile)
	if err != nil {
		return fmt.Errorf("ref mark write: create file: %w", err)
	}
	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("%s\n", m.Reference))
	if err != nil {
		return fmt.Errorf("ref mark write: write file: %w", err)
	}

	return nil
}

// CurrentRef returns the reference of headMark
// which is the hash of the latest label mark.
func (m *RefMark) CurrentRef() ([]byte, error) {
	RefFile, err := RefFile(m.Group, m.Name)
	if err != nil {
		return nil, fmt.Errorf("ref mark current ref: %w", err)
	}

	if _, err := os.Stat(RefFile); errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("ref mark read: file not found: %w", err)
	}

	f, err := os.Open(RefFile)
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

// NewMark creates a new label mark with the given plugin name, plugin id (group) and label map hash.
func NewMark(plugName, group, gmapHash string) (*LabelMark, error) {
	mark := LabelMark{
		Group:     group,
		TimeStamp: time.Now(),
		Mappings:  []MarkMapping{},
	}

	head := RefMark{
		Name:  "HEAD",
		Group: group,
	}
	headFile, err := RefFile(head.Group, head.Name)
	if err != nil {
		return nil, fmt.Errorf("new label mark: %w", err)
	}
	if _, err := os.Stat(headFile); !errors.Is(err, os.ErrNotExist) {
		parent, err := head.CurrentRef()
		if err != nil {
			return nil, fmt.Errorf("new label mark: head ref: %w", err)
		}
		mark.Parent = string(parent)
	}

	return &mark, nil
}

// ReadMark reads the label mark from the file (hash) using given group and map hash.
func ReadMark(group, mapHash string) (*LabelMark, error) {
	mark := LabelMark{
		Hash:     mapHash,
		Group:    group,
		Mappings: []MarkMapping{},
	}
	hashFile, err := ObjectFile(mark.Group, mark.Hash)
	if err != nil {
		return nil, fmt.Errorf("read label mark: %w", err)
	}

	f, err := os.Open(hashFile)
	if err != nil {
		return nil, fmt.Errorf("read label mark: open file: %w", err)
	}
	defer f.Close()

	r, err := zlib.NewReader(f)
	if err != nil {
		return nil, fmt.Errorf("read label mark: create zlib reader: %w", err)
	}
	defer r.Close()

	scanner := bufio.NewScanner(r)
	// Scan the timestamp.
	if ok := scanner.Scan(); ok {
		mark.TimeStamp, err = time.Parse(time.RFC3339, scanner.Text())
		if err != nil {
			return nil, fmt.Errorf("read label mark: parse time: %w", err)
		}
	}
	// Scan the parent hash.
	if ok := scanner.Scan(); ok {
		mark.Parent = scanner.Text()
	}
	// Scan the mappings.
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, fmt.Errorf("read label mark: invalid mapping: %s", line)
		}
		mark.Mappings = append(mark.Mappings, MarkMapping{
			Hash:   fields[0],
			Module: fields[1],
		})
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

	markFile, err := ObjectFile(lm.Group, lm.Hash)
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
	head := RefMark{
		Name:      "HEAD",
		Group:     lm.Group,
		Reference: []byte(lm.Hash),
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
	// fmt.Fprintf(&lm.buffer, "%v\n", lm.TimeStamp)
	fmt.Fprintf(&lm.buffer, "%v\n", lm.TimeStamp.Format(time.RFC3339))
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
