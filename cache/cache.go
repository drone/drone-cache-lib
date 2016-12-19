package cache

import (
	"io"

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
