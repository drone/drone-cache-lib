package cache

import (
	"io"
	"time"

	"github.com/drone/drone-cache-lib/archive"
	"github.com/drone/drone-cache-lib/archive/tar"
	"github.com/drone/drone-cache-lib/flusher"
	"github.com/drone/drone-cache-lib/flusher/timeflusher"
	"github.com/drone/drone-cache-lib/storage"

	log "github.com/Sirupsen/logrus"
)

type Cache struct {
	s storage.Storage
	a archive.Archive
	f flusher.Flusher
}

func New(s storage.Storage, a archive.Archive, f flusher.Flusher) Cache {
	return Cache{ s: s, a: a, f: f }
}

func NewDefault(s storage.Storage) Cache {
	// Return default Cache that uses tar and flushes items after 7 days
	return New(s, tar.New(), timeflusher.New(time.Duration(7*24)*time.Hour))
}

func (c Cache) Rebuild(srcs []string, dst string) error {
	return rebuildCache(srcs, dst, c.s, c.a)
}

func (c Cache) Restore(src string, fallback string) error {
	err := restoreCache(src, c.s, c.a)

	if err != nil && fallback != "" && fallback != src {
		log.Warnf("Failed to retrieve %s, trying %s", src, fallback)
		err = restoreCache(fallback, c.s, c.a)
	}

	// Cache plugin should print an error but it should not return it
	// this is so the build continues even if the cache cant be restored
	if err != nil {
		log.Warnf("Cache could not be restored %s", err)
	}

	return nil
}

func (c Cache) Cleanup(src string) error {
	log.Infof("Cleaning files from %s", src)

	files, err := getList(src, c.s)
	if err != nil {
		return err
	}

	files, err = c.f.Find(files)
	if err != nil {
		return err
	}

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
