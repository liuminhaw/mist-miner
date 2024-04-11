package shelf

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/liuminhaw/mist-miner/shared"
)

const (
	SHELF_DIR = ".miner"
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
	blobDir := fmt.Sprintf("%s/stuffs/%s/%s/%s", SHELF_DIR, s.Module, s.Identity, s.Hash[:2])
	err := os.MkdirAll(blobDir, os.ModePerm)
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

	_, err = f.Write(s.Resource)
	if err != nil {
		return fmt.Errorf("blob write: write file: %w", err)
	}

	fmt.Printf("Stuff file written: %s\n", blobFile)
	return nil
}
