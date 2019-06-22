package storage

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

type FileStorage struct {
	basePath string
}

func NewFileStorage(basePath string) (*FileStorage, error) {
	_, err := os.Stat(basePath)
	if err != nil {
		return nil, errors.Wrap(err, "file storage directory does not exist")
	}

	return &FileStorage{basePath}, nil
}

func (s *FileStorage) ListFolders(path string) ([]string, error) {
	files, err := ioutil.ReadDir(s.basePath + path)
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
	_, err := os.Stat(s.basePath + path)
	if os.IsNotExist(err) {
		return false
	}

	return true
}

func (s *FileStorage) MkDir(path string) error {
	err := os.MkdirAll(s.basePath+path, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "failed creating directory")
	}
	return nil
}

func (s *FileStorage) CreateFile(path string, file io.Reader) error {
	f, err := os.Create(s.basePath + path)
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
	_, err := os.Stat(s.basePath + path)
	if os.IsNotExist(err) {
		return nil, errors.New(fmt.Sprintf("'%s' does not exist", path))
	}

	f, err := os.Open(s.basePath + path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open file")
	}

	return f, nil
}
