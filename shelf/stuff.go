package shelf

import (
	"compress/zlib"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/liuminhaw/mist-miner/shared"
)

type Stuff struct {
	Hash string
	// Module   string
	Identity string
	Resource []byte
}

// NewStuff creates a new Stuff from a plugin name and a MinerResource
func NewStuff(plugId string, resource shared.MinerResource) (*Stuff, error) {
	b, err := json.Marshal(resource)
	if err != nil {
		return nil, fmt.Errorf("new blob: marshal: %w", err)
	}

	h := sha256.New()
	h.Write(b)

	return &Stuff{
		Hash:     fmt.Sprintf("%x", h.Sum(nil)),
		Identity: plugId,
		Resource: b,
	}, nil
}

// Write writes the Stuff resource content to a file
func (s *Stuff) Write() error {
	stuffFile, err := objectFile(s.Identity, s.Hash)
	if err != nil {
		return fmt.Errorf("stuff write: %w", err)
	}
	if _, err := os.Stat(stuffFile); !errors.Is(err, os.ErrNotExist) {
		fmt.Printf("Stuff file already exists: %s\n", stuffFile)
		return nil
	}

	err = os.MkdirAll(filepath.Dir(stuffFile), os.ModePerm)
	if err != nil {
		return fmt.Errorf("stuff write: mkdir: %w", err)
	}

	f, err := os.Create(stuffFile)
	if err != nil {
		return fmt.Errorf("stuff write: create file: %w", err)
	}
	defer f.Close()

	w := zlib.NewWriter(f)
	_, err = w.Write(s.Resource)
	if err != nil {
		return fmt.Errorf("stuff write: write file: %w", err)
	}
	defer w.Close()

	fmt.Printf("Stuff file written: %s\n", stuffFile)
	return nil
}

// Read reads the Stuff resource content from a file and stores it in the Stuff struct
func (s *Stuff) Read() error {
	stuffFile, err := objectFile(s.Identity, s.Hash)
	if err != nil {
		return fmt.Errorf("stuff read: %w", err)
	}
	if _, err := os.Stat(stuffFile); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("stuff read: stuff not found: %s", stuffFile)
	}

	f, err := os.Open(stuffFile)
	if err != nil {
		return fmt.Errorf("stuff read: open file: %w", err)
	}
	defer f.Close()

	r, err := zlib.NewReader(f)
	if err != nil {
		return fmt.Errorf("stuff read: create reader: %w", err)
	}
	defer r.Close()

	b, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("stuff read: read all: %w", err)
	}
	s.Resource = b

	return nil
}

// ResourceIdentifier extract and return Identifier value from stored resource
func (s *Stuff) ResourceIdentifier() (string, error) {
	resource := shared.MinerResource{}

	err := json.Unmarshal(s.Resource, &resource)
	if err != nil {
		return "", fmt.Errorf("resource identifier: unmarshal: %w", err)
	}

	return resource.Identifier, nil
}
