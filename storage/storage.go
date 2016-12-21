package storage

import (
	"io"
	"os"
)

type File struct {
	Path string
	FileInfo os.FileInfo
}

// Storage is a place that files can be written to and read from.
type Storage interface {
	Get(p string, dst io.Writer) error
	Put(p string, src io.Reader) error
	List(p string) ([]File, error)
	Delete(p string) error
}
