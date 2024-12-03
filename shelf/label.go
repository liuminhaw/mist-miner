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

	"github.com/liuminhaw/mist-miner/locks"
)

type IdentifierHashMap struct {
	Identifier string
	Alias      string
	// Hash of pointed stuff outline
	Hash string
}

type IdentifierHashMaps struct {
	Hash string
	// Module   string
	Group  string
	Maps   []IdentifierHashMap
	buffer bytes.Buffer
}

// ReadIdentifierHashMaps reads the identifier hash maps from the file (hash) using the given group and map hash
// then convert the content to IdentifierHashMaps and return it.
func ReadIdentifierHashMaps(group, hash string) (*IdentifierHashMaps, error) {
	idHashMaps := IdentifierHashMaps{
		Hash:  hash,
		Group: group,
		Maps:  []IdentifierHashMap{},
	}

	r, err := NewObjectRecord(idHashMaps.Group, idHashMaps.Hash).RecordReadCloser()
	if err != nil {
		return nil, fmt.Errorf("read identifier hash maps: %w", err)
	}

	// Scan in hash identifier pairs
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		switch len(fields) {
		case 2:
			idHashMaps.Maps = append(idHashMaps.Maps, IdentifierHashMap{
				Hash:       fields[0],
				Identifier: fields[1],
			})
		case 3:
			idHashMaps.Maps = append(idHashMaps.Maps, IdentifierHashMap{
				Hash:       fields[0],
				Identifier: fields[1],
				Alias:      fields[2],
			})
		default:
			return nil, fmt.Errorf("read identifier hash maps: invalid mapping: %s", line)
		}
	}

	return &idHashMaps, nil
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

	mapFile, err := NewObjectRecord(lhm.Group, lhm.Hash).RecordFile()
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
		if m.Alias != "" {
			fmt.Fprintf(&lhm.buffer, "%s %s %s\n", m.Hash, m.Identifier, m.Alias)
		} else {
			fmt.Fprintf(&lhm.buffer, "%s %s\n", m.Hash, m.Identifier)
		}
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

func NewRefMark(name, group string) (RefMark, error) {
	mark := RefMark{Name: name, Group: group}
	if err := mark.getReference(); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return RefMark{}, ErrRefHeadNotFound
		}
		return RefMark{}, fmt.Errorf("NewRefMark(%s, %s): %w", name, group, err)
	}

	return mark, nil
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

// getReference reads the reference from the HEAD file and sets it to the Reference field.
// Will use flock to prevent writing while reading,
// return locks.ErrIsLocked if file lock is not acquired.
func (m *RefMark) getReference() error {
	file, err := RefFile(m.Group, m.Name)
	if err != nil {
		return fmt.Errorf("getReference(): %w", err)
	}

	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("getReference(): %w", err)
	}

	fileLock, err := locks.NewLock(m.Group, locks.REF_MARK_LOCKFILE)
	if err != nil {
		return fmt.Errorf("getReference(): %w", err)
	}
	if err := fileLock.TryRLock(); err != nil {
		if errors.Is(err, locks.ErrIsLocked) {
			return err
		}
		return fmt.Errorf("getReference(): %w", err)
	}
	defer fileLock.Unlock()

	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("getReference(): %w", err)
	}
	defer f.Close()

	b := make([]byte, sha256.BlockSize)
	_, err = f.Read(b)
	if err != nil {
		return fmt.Errorf("gerReference: %w", err)
	}

	m.Reference = b
	return nil
}

type MarkMapping struct {
	Module string
	Hash   string
}

type LabelMark struct {
	// Module       string
	Hash      string
	TimeStamp time.Time
	LogType   string
	Parent    string
	Mappings  []MarkMapping
	Group     string
	// LabelMapHash string
	buffer bytes.Buffer
}

// NewMark creates a new label mark with the given plugin name, plugin id (group) and label map hash.
func NewMark(group, logType string) (*LabelMark, error) {
	mark := LabelMark{
		Group:     group,
		LogType:   logType,
		TimeStamp: time.Now(),
		Mappings:  []MarkMapping{},
	}

	head, err := NewRefMark(SHELF_MARK_FILE, group)
	if errors.Is(err, ErrRefHeadNotFound) {
		mark.Parent = "nil"
	} else if err != nil {
		return nil, fmt.Errorf("new label mark: %w", err)
	} else {
		mark.Parent = string(head.Reference)
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

	r, err := NewObjectRecord(mark.Group, mark.Hash).RecordReadCloser()
	if err != nil {
		return nil, fmt.Errorf("read label mark: %w", err)
	}

	scanner := bufio.NewScanner(r)
	// Scan the timestamp.
	if ok := scanner.Scan(); ok {
		mark.TimeStamp, err = time.Parse(time.RFC3339, scanner.Text())
		if err != nil {
			return nil, fmt.Errorf("read label mark: parse time: %w", err)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read label mark timestamp: %w", err)
	}
	// Scan the log type.
	if ok := scanner.Scan(); ok {
		mark.LogType = scanner.Text()
		if mark.LogType != "mine" && mark.LogType != "diary" {
			return nil, fmt.Errorf("read label mark: invalid log type: %s", mark.LogType)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read label mark log type: %w", err)
	}
	// Scan the parent hash.
	if ok := scanner.Scan(); ok {
		mark.Parent = scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read label mark parent hash: %w", err)
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
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read label mark: %w", err)
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

	markFile, err := NewObjectRecord(lm.Group, lm.Hash).RecordFile()
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
	if lm.Parent == "" || lm.Parent == "nil" {
		parent = "nil"
	} else {
		parent = lm.Parent
	}

	// fmt.Printf("Parent: %s\n", parent)
	// fmt.Fprintf(&lm.buffer, "%v\n", lm.TimeStamp)
	fmt.Fprintf(&lm.buffer, "%v\n", lm.TimeStamp.Format(time.RFC3339))
	fmt.Fprintf(&lm.buffer, "%s\n", lm.LogType)
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
