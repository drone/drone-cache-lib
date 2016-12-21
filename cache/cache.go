package cache

import (
	"io"
	"time"

	"github.com/drone/drone-cache-lib/archive"
	"github.com/drone/drone-cache-lib/archive/util"
	"github.com/drone/drone-cache-lib/storage"

	log "github.com/Sirupsen/logrus"
)

type Cache struct {
	s storage.Storage
}

func New(s storage.Storage) (Cache, error) {
	return Cache{
		s: s,
	}, nil
}

func (c Cache) Rebuild(srcs []string, dst string) error {
	a, err := util.FromFilename(dst)

	if err != nil {
		return err
	}

	return rebuildCache(srcs, dst, c.s, a)
}

func (c Cache) Restore(src string, fallback string) error {
	a, err := util.FromFilename(src)

	if err != nil {
		return err
	}

	err = restoreCache(src, c.s, a)

	if err != nil && fallback != "" && fallback != src {
		log.Warnf("Failed to retrieve %s, trying %s", src, fallback)
		err = restoreCache(fallback, c.s, a)
	}

	// Cache plugin should print an error but it should not return it
	// this is so the build continues even if the cache cant be restored
	if err != nil {
		log.Warnf("Cache could not be restored %s", err)
	}

	return nil
}

func (c Cache) Cleanup(src string, age time.Duration) error {
	log.Infof("Cleaning files from %s older then %s", src, age)

	files, err := getList(src, c.s)
	if err != nil {
		return err
	}

	files = findFiles(files, age)

	err = deleteFiles(files, c.s)

	return err
}

func restoreCache(src string, s storage.Storage, a archive.Archive) error {
	reader, writer := io.Pipe()

	cw := make(chan error, 1)
	defer close(cw)

	go func() {
		defer writer.Close()

		cw <- s.Get(src, writer)
	}()

	err := a.Unpack("", reader)
	werr := <-cw

	if werr != nil {
		return werr
	}

	return err
}

func rebuildCache(srcs []string, dst string, s storage.Storage, a archive.Archive) error {
	log.Infof("Rebuilding cache at %s to %s", srcs, dst)

	reader, writer := io.Pipe()
	defer reader.Close()

	cw := make(chan error, 1)
	defer close(cw)

	go func() {
		defer writer.Close()

		cw <- a.Pack(srcs, writer)
	}()

	err := s.Put(dst, reader)
	werr := <-cw

	if werr != nil {
		return werr
	}

	return err
}


func getList(src string, s storage.Storage) ([]storage.FileEntry, error) {
	return s.List(src)
}

func findFiles(files []storage.FileEntry, age time.Duration) []storage.FileEntry {
	var matchedFiles []storage.FileEntry
	for _, file := range files {
		// Match files (not dirs) older then age
		if !file.Info.IsDir() && file.Info.ModTime().Before(time.Now().Add(-1 * age)) {
			matchedFiles = append(matchedFiles, file)
		}
	}

	return matchedFiles
}

func deleteFiles(files []storage.FileEntry, s storage.Storage) error {
	var err error
	for _, file := range files {
		err = s.Delete(file.Path)
		if err != nil {
			return err
		}
	}

	return nil
}
