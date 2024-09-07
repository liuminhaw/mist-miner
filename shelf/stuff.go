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
	"strings"
)

func ReadStuffOutline(group, hash string) (*StuffOutline, error) {
	r, err := objectReader(group, hash)
	if err != nil {
		return nil, fmt.Errorf("read stuff outline: %w", err)
	}

	content, err := io.ReadAll(r)
	fields := strings.Fields(string(content))
	if len(fields) != 2 {
		return nil, fmt.Errorf("read stuff outline: invalid content: %s", content)
	}

	return &StuffOutline{
		Hash:         hash,
		Group:        group,
		ResourceHash: fields[0],
		DiaryHash:    fields[1],
	}, nil
}

type StuffOutline struct {
	Hash         string
	Group        string
	ResourceHash string
	DiaryHash    string
	Content      []byte
}

func NewStuffOutline(group, resourceHash, diaryHash string) StuffOutline {
	content := fmt.Sprintf("%s %s", resourceHash, diaryHash)

	h := sha256.New()
	h.Write([]byte(content))

	return StuffOutline{
		Hash:         fmt.Sprintf("%x", h.Sum(nil)),
		Group:        group,
		ResourceHash: resourceHash,
		DiaryHash:    diaryHash,
		Content:      []byte(content),
	}
}

func (s *StuffOutline) Write() error {
	outlineFile, err := ObjectFile(s.Group, s.Hash)
	if err != nil {
		return fmt.Errorf("stuff outline write: %w", err)
	}
	if _, err := os.Stat(outlineFile); !errors.Is(err, os.ErrNotExist) {
		fmt.Printf("Stuff outline file already exists: %s\n", outlineFile)
		return nil
	}

	err = os.MkdirAll(filepath.Dir(outlineFile), os.ModePerm)
	if err != nil {
		return fmt.Errorf("stuff outline write: mkdir: %w", err)
	}

	f, err := os.Create(outlineFile)
	if err != nil {
		return fmt.Errorf("stuff outline write: create file: %w", err)
	}
	defer f.Close()

	w := zlib.NewWriter(f)
	_, err = w.Write(s.Content)
	if err != nil {
		return fmt.Errorf("stuff outline write: write file: %w", err)
	}
	defer w.Close()

	fmt.Printf("Stuff outline file written: %s\n", outlineFile)
	return nil
}

type Stuff struct {
	Hash     string
	Group    string
	Resource []byte
}

// NewStuff creates a new Stuff from a plugin name and a MinerResource
func NewStuff(group string, a any) (*Stuff, error) {
	b, err := json.Marshal(a)
	if err != nil {
		return nil, fmt.Errorf("new blob: marshal: %w", err)
	}

	h := sha256.New()
	h.Write(b)

	return &Stuff{
		Hash:     fmt.Sprintf("%x", h.Sum(nil)),
		Group:    group,
		Resource: b,
	}, nil
}

// Write writes the Stuff resource content to a file
func (s *Stuff) Write() error {
	stuffFile, err := ObjectFile(s.Group, s.Hash)
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

