package shelf

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type DiaryMeta struct {
	Group      string
	Plugin     string
	Identifier string
	Alias      string
}

type Diary struct {
	Hash    string // current diary's record hash value
	Content string
	Meta    DiaryMeta
}

func NewDiary(group, plugin, identifier, alias, hash string) Diary {
	return Diary{
		Meta: DiaryMeta{
			Group:      group,
			Plugin:     plugin,
			Identifier: identifier,
			Alias:      alias,
		},
		Hash: hash,
	}
}

// Exist checks if the diary record (object) exists
func (d *Diary) Exist() bool {
	if d.Hash == "" {
		return false
	}

	record := NewObjectRecord(d.Meta.Group, d.Hash)
	return record.Exist()
}

// NewTempFile creates a temporary file for the diary record if it does not exist
// and returns the file path with error if any
func (d *Diary) NewTempFile() (DiaryTempFile, error) {
	path, err := d.tempFile()
	if err != nil {
		return DiaryTempFile{}, fmt.Errorf("diary NewTempFile: %w", err)
	}

	return DiaryTempFile{
		Path:       path,
		StaticPath: d.staticTempFile(),
		Meta:       d.Meta,
	}, nil
}

// tempDir returns the temporary directory to generate the diary record
func (d *Diary) tempDir() string {
	return filepath.Join(os.TempDir(), shelf_temp_base_dir, d.Meta.Group, d.Meta.Plugin)
}

// tempFile returns the temporary file to generate the diary record
// the file will be interpret in markdown format
// filename format: <identifierBase64Encode>.<random>.md
func (d *Diary) tempFile() (string, error) {
	identifierEncode := base64.RawURLEncoding.EncodeToString([]byte(d.Meta.Identifier))

	randBytes, err := randomHex(4)
	if err != nil {
		return "", fmt.Errorf("diary temp file: %w", err)
	}

	filename := fmt.Sprintf("%s.%s.md", identifierEncode, randBytes)
	return filepath.Join(d.tempDir(), filename), nil
}

// staticTempFile returns the static temporary file to store the edited diary record
// filename format: <identifierBase64Encode>.md
func (d *Diary) staticTempFile() string {
	identifierEncode := base64.RawURLEncoding.EncodeToString([]byte(d.Meta.Identifier))
	filename := fmt.Sprintf("%s.md", identifierEncode)

	return filepath.Join(d.tempDir(), filename)
}

type DiaryTempFile struct {
	File       *os.File
	Path       string
	StaticPath string
	Meta       DiaryMeta
}

func (d *DiaryTempFile) Init() error {
	err := d.create()
	if err != nil {
		return fmt.Errorf("diary temp file init: %w", err)
	}

	var initMsg string
	if d.Meta.Alias != "" {
		initMsg = fmt.Sprintf("# %s\n\n", d.Meta.Alias)
	} else {
		initMsg = fmt.Sprintf("# %s\n\n", d.Meta.Identifier)
	}
	initMsg += fmt.Sprint("## Basic Info\n\n")
	initMsg += fmt.Sprintf("- **Group:** %s\n", d.Meta.Group)
	initMsg += fmt.Sprintf("- **Plugin:** %s\n", d.Meta.Plugin)
	if d.Meta.Alias != "" {
		initMsg += fmt.Sprintf("- **Identifier:** %s\n", d.Meta.Identifier)
	}

	if _, err := d.File.WriteString(initMsg); err != nil {
		return fmt.Errorf("diary temp file init: %w", err)
	}

	return nil
}

func (d *DiaryTempFile) Close() error {
	if d.File != nil {
		if err := d.File.Close(); err != nil {
			return fmt.Errorf("diary temp file close: %w", err)
		}
	}

	return nil
}

func (d *DiaryTempFile) Exist() bool {
	if _, err := os.Stat(d.StaticPath); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}

// ToStaticTemp renames the temporary file to a static temporary file
// static temporary file format: <identifierBase64Encode>.md
func (d *DiaryTempFile) ToStaticTemp() error {
	if err := os.Rename(d.Path, d.StaticPath); err != nil {
		return fmt.Errorf("diaryTempFile ToStaticTemp: %w", err)
	}

	return nil
}

// CopyFromStatic copies the static temporary file to the temporary file for editing
func (d *DiaryTempFile) CopyFromStatic() error {
	sourceFile, err := os.Open(d.StaticPath)
	if err != nil {
		return fmt.Errorf("diaryTempFile CopyFromStatic: open static: %w", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(d.Path)
	if err != nil {
		return fmt.Errorf("diaryTempFile CopyFromStatic: create temp: %w", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return fmt.Errorf("diaryTempFile CopyFromStatic: copy: %w", err)
	}
	if err := destFile.Sync(); err != nil {
		return fmt.Errorf("diaryTempFile CopyFromStatic: sync: %w", err)
	}

	return nil
}

func (d *DiaryTempFile) create() error {
	if _, err := os.Stat(d.Path); errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(filepath.Dir(d.Path), 0755); err != nil {
			return fmt.Errorf("diary temp file create: %w", err)
		}

		d.File, err = os.Create(d.Path)
		if err != nil {
			return fmt.Errorf("diary temp file create: %w", err)
		}
	}

	return nil
}
