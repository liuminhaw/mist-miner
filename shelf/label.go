package shelf

import (
	"bytes"
	"compress/zlib"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
)

// func mapDir() string {
//     return fmt.Sprintf(
// }

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

	mapDir := lhm.dir()
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
func (lhm *IdentifierHashMaps) dir() string {
    return fmt.Sprintf("%s/labels/%s/%s/maps/%s", SHELF_DIR, lhm.Module, lhm.Identity, lhm.Hash[:2])
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
