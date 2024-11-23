package shelf

import (
	"compress/zlib"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type DiaryMeta struct {
	Group      string
	Plugin     string
	Identifier string
	Alias      string
}

// tempDir returns the temporary directory basetargeted for generating the diary record
func (md DiaryMeta) tempDir() string {
	return filepath.Join(
		os.TempDir(),
		shelf_temp_base_dir,
		shelf_diary_dir,
		md.Group,
		md.Plugin,
	)
}

// tempFile returns the temporary file to generate the diary record
// the file will be interpret in markdown format
// filename format: <identifierBase64Encode>.<random>.md
func (md DiaryMeta) tempFile() (string, error) {
	identifierEncode := base64.RawURLEncoding.EncodeToString([]byte(md.Identifier))

	randBytes, err := randomHex(4)
	if err != nil {
		return "", fmt.Errorf("diary temp file: %w", err)
	}

	filename := fmt.Sprintf("%s.%s.md", identifierEncode, randBytes)
	return filepath.Join(md.tempDir(), "temp", filename), nil
}

// staticTempFile returns the static temporary file to store the edited diary record
// filename format: <identifierBase64Encode>.md
func (md DiaryMeta) staticTempFile() string {
	identifierEncode := base64.RawURLEncoding.EncodeToString([]byte(md.Identifier))
	filename := fmt.Sprintf("%s.md", identifierEncode)

	return filepath.Join(md.tempDir(), "static", filename)
}

func metaFromStaticTempFile(path string) (DiaryMeta, error) {
	filename := filepath.Base(path)
	identifierEncode := strings.TrimSuffix(filename, ".md")
	identifier, err := base64.RawURLEncoding.DecodeString(identifierEncode)
	if err != nil {
		return DiaryMeta{}, fmt.Errorf("metaFromStaticTempFile: filename decode: %w", err)
	}

	return DiaryMeta{
		Group:      filepath.Base(filepath.Dir(filepath.Dir(filepath.Dir(path)))),
		Plugin:     filepath.Base(filepath.Dir(filepath.Dir(path))),
		Identifier: string(identifier),
	}, nil
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

// Write writes the diary record to file using hash value as directory and filename
// func (d *Diary) Write() error {
// 	diaryFile, err := NewObjectRecord(d.Meta.Group, d.Hash).RecordFile()
// 	if err != nil {
// 		return fmt.Errorf("diary write: %w", err)
// 	}
// 	if _, err := os.Stat(diaryFile); !errors.Is(err, os.ErrNotExist) {
// 		fmt.Printf("diary file already exists: %s\n", diaryFile)
// 		return nil
// 	}
//
// 	err = os.MkdirAll(filepath.Dir(diaryFile), os.ModePerm)
// 	if err != nil {
// 		return fmt.Errorf("diary write: mkdir: %w", err)
// 	}
//
// 	f, err := os.Create(diaryFile)
// 	if err != nil {
// 		return fmt.Errorf("diary write: create file: %w", err)
// 	}
// 	defer f.Close()
//
// 	w := zlib.NewWriter(f)
// 	_, err = w.Write([]byte(d.Content))
// 	if err != nil {
// 		return fmt.Errorf("diary write: write file: %w", err)
// 	}
// 	defer w.Close()
//
// 	return nil
// }

// NewTempFile creates a temporary file for the diary record if it does not exist
// and returns the file path with error if any
func (d *Diary) NewTempFile() (DiaryTempFile, error) {
	path, err := d.Meta.tempFile()
	if err != nil {
		return DiaryTempFile{}, fmt.Errorf("diary NewTempFile: %w", err)
	}

	return DiaryTempFile{
		Path: path,
		Meta: d.Meta,
	}, nil
}

type DiaryTempFile struct {
	File *os.File
	Path string
	Meta DiaryMeta
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

// StaticExist checks if the static temporary file exists
func (d *DiaryTempFile) StaticExist() bool {
	if _, err := os.Stat(d.Meta.staticTempFile()); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}

// TempExist checks if the temporary file exists
func (d *DiaryTempFile) TempExist() bool {
	if _, err := os.Stat(d.Path); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}

// ToStaticTemp renames the temporary file to a static temporary file
// static temporary file format: <identifierBase64Encode>.md
func (d *DiaryTempFile) ToStaticTemp() error {
	if !d.TempExist() {
		return fmt.Errorf("diaryTempFile ToStaticTemp: temp file does not exist")
	}

	if err := os.MkdirAll(filepath.Dir(d.Meta.staticTempFile()), 0755); err != nil {
		return fmt.Errorf("diaryTempFile ToStaticTemp: mkdir static temp dir: %w", err)
	}

	fileInfo, _ := os.Stat(d.Path)
	tempModtime := fileInfo.ModTime()

	var staticModtime time.Time
	if d.StaticExist() {
		fileInfo, _ := os.Stat(d.Meta.staticTempFile())
		staticModtime = fileInfo.ModTime()
	}

	if tempModtime.After(staticModtime) {
		if err := os.Rename(d.Path, d.Meta.staticTempFile()); err != nil {
			return fmt.Errorf("diaryTempFile ToStaticTemp: %w", err)
		}
	} else {
		if err := os.Remove(d.Path); err != nil {
			return fmt.Errorf("diaryTempFile ToStaticTemp: remove temp: %w", err)
		}
	}

	return nil
}

// CopyFromStatic copies the static temporary file to the temporary file for editing
func (d *DiaryTempFile) CopyFromStatic() error {
	sourceFile, err := os.Open(d.Meta.staticTempFile())
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

// DiaryStaticTempFile is a wrapper for Static temp file over DiaryTempFile
type DiaryStaticTempFile struct {
	Path string
	Meta DiaryMeta
}

func NewDiaryStaticTempFile(staticPath string) (DiaryStaticTempFile, error) {
	meta, err := metaFromStaticTempFile(staticPath)
	if err != nil {
		return DiaryStaticTempFile{}, fmt.Errorf("NewDiaryStaticTempFile: %w", err)
	}

	return DiaryStaticTempFile{
		Path: staticPath,
		Meta: meta,
	}, nil
}

func (d *DiaryStaticTempFile) Read() (string, error) {
	if _, err := os.Stat(d.Path); errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("DiaryStaticTempFile stat file: %w", err)
	}

	content, err := os.ReadFile(d.Path)
	if err != nil {
		return "", fmt.Errorf("DiaryStaticTempFile read file: %w", err)
	}

	return string(content), nil
}

func (d *DiaryStaticTempFile) CalcHash() (string, error) {
	content, err := d.Read()
	if err != nil {
		return "", fmt.Errorf("DiaryStaticTempFile CalcHash: %w", err)
	}

	h := sha256.New()
	_, err = h.Write([]byte(content))
	if err != nil {
		return "", fmt.Errorf("DiaryStaticTempFile CalcHash: %w", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func (d *DiaryStaticTempFile) WriteDiary() (Diary, error) {
	content, err := d.Read()
	if err != nil {
		return Diary{}, fmt.Errorf("diaryStaticTempFile WriteDiary: %w", err)
	}
	hash, err := d.CalcHash()
	if err != nil {
		return Diary{}, fmt.Errorf("diaryStaticTempFile WriteDiary: %w", err)
	}
	diaryFile, err := NewObjectRecord(d.Meta.Group, hash).RecordFile()
	if err != nil {
		return Diary{}, fmt.Errorf("diaryStaticTempFile WriteDiary: %w", err)
	}

	if _, err := os.Stat(diaryFile); !errors.Is(err, os.ErrNotExist) {
		return Diary{}, fmt.Errorf(
			"diaryStaticTempFile WriteDiary: diary file already exists: %s",
			diaryFile,
		)
	}

	err = os.MkdirAll(filepath.Dir(diaryFile), os.ModePerm)
	if err != nil {
		return Diary{}, fmt.Errorf("diaryStaticTempFile WriteDiary: mkdir: %w", err)
	}

	f, err := os.Create(diaryFile)
	if err != nil {
		return Diary{}, fmt.Errorf("diaryStaticTempFile WriteDiary: create file: %w", err)
	}
	defer f.Close()

	w := zlib.NewWriter(f)
	_, err = w.Write([]byte(content))
	if err != nil {
		return Diary{}, fmt.Errorf("diaryStaticTempFile WriteDiary: write file: %w", err)
	}
	defer w.Close()

	return NewDiary(d.Meta.Group, d.Meta.Plugin, d.Meta.Identifier, d.Meta.Alias, hash), nil
}
