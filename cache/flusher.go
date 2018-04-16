package cache

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/drone/drone-cache-lib/storage"
)

// DirtyFunc defines when an cache item is outdated.
type DirtyFunc func(storage.FileEntry) bool

// Flusher defines an object to clear the cache.
type Flusher struct {
	store storage.Storage
	dirty func(storage.FileEntry) bool
}

// NewFlusher creates a new cache flusher.
func NewFlusher(s storage.Storage, fn DirtyFunc) Flusher {
	return Flusher{store: s, dirty: fn}
}

// NewDefaultFlusher creates a new cache flusher with default expire.
func NewDefaultFlusher(s storage.Storage) Flusher {
	return Flusher{store: s, dirty: IsExpired}
}

// Flush cleans the cache if it's expired.
func (f *Flusher) Flush(src string) error {
	log.Infof("Cleaning files from %s", src)

	files, err := f.store.List(src)
	if err != nil {
		return err
	}

	for _, file := range files {
		if f.dirty(file) {
			err := f.store.Delete(file.Path)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// IsExpired checks if the cache is expired.
func IsExpired(file storage.FileEntry) bool {
	// Check if older then 30 days
	return file.LastModified.Before(time.Now().AddDate(0, 0, -30))
}
