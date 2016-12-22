package timeflusher

import (
	"time"

	log "github.com/Sirupsen/logrus"

	. "github.com/drone/drone-cache-lib/flusher"
	"github.com/drone/drone-cache-lib/storage"
)

type timeFlusher struct{
	Age time.Duration
}

// New creates an Flusher that ops on time.Duration
func New(age time.Duration) Flusher {
	return &timeFlusher{
		Age: age,
	}
}

func (f *timeFlusher) Find(files []storage.FileEntry) ([]storage.FileEntry, error) {
  var matchedFiles []storage.FileEntry
	for _, file := range files {
		// Match files (not dirs) older then age
		if !file.Info.IsDir() && file.Info.ModTime().Before(time.Now().Add(-1 * f.Age)) {
			matchedFiles = append(matchedFiles, file)
		}
	}

	log.Infof("Found %s files", len(matchedFiles))

	return matchedFiles, nil
}
