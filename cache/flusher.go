package cache

import(
	"time"

	"github.com/drone/drone-cache-lib/storage"

	log "github.com/Sirupsen/logrus"
)

type dirtyFunc func(storage.FileEntry) bool

type Flusher struct {
	store storage.Storage
	dirty func(storage.FileEntry) bool
}

func NewFlusher(s storage.Storage, fn dirtyFunc) Flusher {
	return Flusher{ store: s, dirty: fn }
}

func NewDefaultFlusher(s storage.Storage) Flusher {
	return Flusher{ store: s, dirty: TimeOp() }
}

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

func TimeOp() dirtyFunc {
	return func(file storage.FileEntry) bool {
		// Do not delete directories
		if file.Info.IsDir() {
			return false
		}

		// Check if older then 30 days
		if file.Info.ModTime().Before(time.Now().AddDate(0, 0, -30)) {
			return true
		}

		// No match
		return false
	}
}
