package storage

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type FileStorage struct {
	basePath string
}

func NewFileStorage(basePath string) (*FileStorage, error) {
	absBasePath, err := filepath.Abs(basePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed determining absolute directory")
	}
	_, err = os.Stat(absBasePath)
	if err != nil {
		err := os.MkdirAll(absBasePath, os.ModePerm)
		if err != nil {
			return nil, errors.Wrap(err, "failed creating directory")
		}
	}

	return &FileStorage{absBasePath}, nil
}

func (s *FileStorage) ListFolders(path string) ([]string, error) {
	fullPath, err := s.securePath(path)
	if err != nil {
		return nil, errors.Wrap(err, "could not determine secure path")
	}
	stat, err := os.Stat(fullPath)
	if err != nil {
		return []string{}, nil
	}

	if !stat.IsDir() {
		return nil, errors.Wrap(err, "path is not a directory")
	}

	files, err := ioutil.ReadDir(fullPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read directory")
	}

	directories := []string{}
	for _, f := range files {
		if f.IsDir() {
			directories = append(directories, f.Name())
		}
	}

	return directories, nil
}

func (s *FileStorage) File(path string) (io.Reader, error) {
	return nil, nil
}

func (s *FileStorage) Exists(path string) bool {
	securePath, err := s.securePath(path)
	if err != nil {
		return false
	}

	_, err = os.Stat(securePath)
	if os.IsNotExist(err) {
		return false
	}

	return true
}

func (s *FileStorage) MkDir(path string) error {
	fullPath, err := s.securePath(path)
	if err != nil {
		return errors.Wrap(err, "could not determine secure path")
	}

	err = os.MkdirAll(fullPath, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "failed creating directory")
	}
	return nil
}

func (s *FileStorage) CreateFile(path string, file io.Reader) error {
	fullPath, err := s.securePath(path)
	if err != nil {
		return errors.Wrap(err, "could not determine secure path")
	}

	f, err := os.Create(fullPath)
	if err != nil {
		return errors.Wrap(err, "failed creating file")
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	r := bufio.NewReader(file)
	defer w.Flush()

	buf := make([]byte, 1024)
	for {
		// read a chunk
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return errors.Wrap(err, "failed to read from stream")
		}
		if n == 0 {
			break
		}

		// write a chunk
		if _, err := w.Write(buf[:n]); err != nil {
			return errors.Wrap(err, "failed to write to file")
		}
	}

	return nil
}

func (s *FileStorage) ReadFile(path string) (io.ReadCloser, error) {
	fullPath, err := s.securePath(path)
	if err != nil {
		return nil, errors.Wrap(err, "could not determine secure path")
	}

	_, err = os.Stat(fullPath)
	if os.IsNotExist(err) {
		return nil, errors.New(fmt.Sprintf("'%s' does not exist", path))
	}

	f, err := os.Open(fullPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open file")
	}

	return f, nil
}

func (s *FileStorage) securePath(path string) (string, error) {
	unsafePath := s.basePath + path
	absPath, err := filepath.Abs(unsafePath)
	if err != nil {
		return "", err
	}

	_, err = filepath.Rel(s.basePath, absPath)
	if err != nil {
		return "", err
	}

	return absPath, nil
}
