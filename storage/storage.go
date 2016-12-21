package storage

import (
	"io"
	"os"
)

type FileEntry struct {
	Path string
	Info os.FileInfo
}

// Storage is a place that files can be written to and read from.
type Storage interface {
	Get(p string, dst io.Writer) error
	Put(p string, src io.Reader) error
	List(p string) ([]FileEntry, error)
	Delete(p string) error
}
