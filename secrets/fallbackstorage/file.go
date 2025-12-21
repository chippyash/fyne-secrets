package fallbackstorage

import (
	"fmt"
	"io"
	"strings"

	"fyne.io/fyne/v2"
	fs "fyne.io/fyne/v2/storage"
	"github.com/pkg/errors"
)

// FileStorage - a simple wrapper around Fyne storage
type FileStorage struct {
	root fyne.URI
}

func NewFileStorage() Storing {
	return &FileStorage{
		root: fyne.CurrentApp().Storage().RootURI(),
	}
}

func (s *FileStorage) Store(pathAndName string, data []byte) error {
	if strings.Contains(pathAndName, "/") {
		if err := s.makePath(pathAndName); err != nil {
			return err
		}
	}
	uri, err := fs.Child(s.root, pathAndName)
	if err != nil {
		return err
	}
	w, err := fs.Writer(uri)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to open file: %s", pathAndName))
	}
	defer func(w fyne.URIWriteCloser) {
		_ = w.Close()
	}(w)
	_, err = w.Write(data)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to write to file: %s", pathAndName))
	}
	return nil
}

func (s *FileStorage) Exists(pathAndName string) (bool, error) {
	uri, err := fs.Child(s.root, pathAndName)
	if err != nil {
		return false, err
	}
	return fs.Exists(uri)
}

func (s *FileStorage) Load(pathAndName string) ([]byte, error) {
	uri, err := fs.Child(s.root, pathAndName)
	if err != nil {
		return nil, err
	}
	r, err := fs.Reader(uri)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to open file: %s", pathAndName))
	}
	defer func(r fyne.URIReadCloser) {
		_ = r.Close()
	}(r)

	return io.ReadAll(r)
}

func (s *FileStorage) Delete(pathAndName string) error {
	if flag, _ := s.Exists(pathAndName); !flag {
		return nil
	}

	uri, err := fs.Child(s.root, pathAndName)
	if err != nil {
		return err
	}
	return fs.Delete(uri)
}

func (s *FileStorage) makePath(pathAndName string) error {
	path, err := fs.Child(s.root, pathAndName)
	if err != nil {
		return err
	}
	uri, err := fs.Parent(path)
	if err != nil {
		return err
	}
	err = fs.CreateListable(uri)
	if err != nil {
		if !strings.Contains(err.Error(), "file exists") {
			return errors.Wrap(err, fmt.Sprintf("failed to create %s directory", path))
		}
	}
	return nil
}
