package shelf

import (
	"compress/zlib"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/liuminhaw/mist-miner/shared"
)

type Stuff struct {
	Hash     string
	Module   string
	Identity string
	Resource []byte
}

// NewStuff creates a new Stuff from a plugin name and a MinerResource
func NewStuff(plugName, plugId string, resource shared.MinerResource) (*Stuff, error) {
	b, err := json.Marshal(resource)
	if err != nil {
		return nil, fmt.Errorf("new blob: marshal: %w", err)
	}

	h := sha256.New()
	h.Write(b)

	return &Stuff{
		Hash:     fmt.Sprintf("%x", h.Sum(nil)),
		Module:   plugName,
		Identity: plugId,
		Resource: b,
	}, nil
}

// Write writes the Stuff resource content to a file
func (s *Stuff) Write() error {
	blobDir, err := s.dir()
	if err != nil {
		return fmt.Errorf("blob write: %w", err)
	}

	err = os.MkdirAll(blobDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("blob write: mkdir: %w", err)
	}

	blobFile := fmt.Sprintf("%s/%s", blobDir, s.Hash[2:])
	if _, err := os.Stat(blobFile); !errors.Is(err, os.ErrNotExist) {
		fmt.Printf("Stuff file already exists: %s\n", blobFile)
		return nil
	}

	f, err := os.Create(blobFile)
	if err != nil {
		return fmt.Errorf("blob write: create file: %w", err)
	}
	defer f.Close()

	w := zlib.NewWriter(f)
	_, err = w.Write(s.Resource)
	if err != nil {
		return fmt.Errorf("blob write: write file: %w", err)
	}
	defer w.Close()

	fmt.Printf("Stuff file written: %s\n", blobFile)
	return nil
}

// Read reads the Stuff resource content from a file and stores it in the Stuff struct
func (s *Stuff) Read() error {
	blobDir, err := s.dir()
	if err != nil {
		return fmt.Errorf("blob read: %w", err)
	}

	blobFile := fmt.Sprintf("%s/%s", blobDir, s.Hash[2:])

	if _, err := os.Stat(blobFile); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("blob read: stuff not found: %s", blobFile)
	}

	f, err := os.Open(blobFile)
	if err != nil {
		return fmt.Errorf("blob read: open file: %w", err)
	}
	defer f.Close()

	r, err := zlib.NewReader(f)
	if err != nil {
		return fmt.Errorf("blob read: create reader: %w", err)
	}
	defer r.Close()

	b, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("blob read: read all: %w", err)
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

// dir returns the directory path of Stuff
// If absDir is empty, the relative path is returned
func (s *Stuff) dir() (string, error) {
	sd, err := shelfDir()
	if err != nil {
		return "", fmt.Errorf("stuff dir: %w", err)
	}

	return fmt.Sprintf("%s/stuffs/%s/%s/%s", sd, s.Module, s.Identity, s.Hash[:2]), nil
}
