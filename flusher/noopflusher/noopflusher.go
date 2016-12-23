package noopflusher

import (
	. "github.com/drone/drone-cache-lib/flusher"
	"github.com/drone/drone-cache-lib/storage"
)

type noopFlusher struct{}

// New creates an Flusher that does nothing
func New() Flusher {
	return &noopFlusher{}
}

func (f *noopFlusher) Find(files []storage.FileEntry) ([]storage.FileEntry, error) {
	return []storage.FileEntry{}, nil
}
