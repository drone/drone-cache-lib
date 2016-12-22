package flusher

import "github.com/drone/drone-cache-lib/storage"

// Flusher is an interface for finding files to flush
type Flusher interface {
	// Find subset of files
	Find(files []storage.FileEntry) ([]storage.FileEntry, error)
}
