package dummy

import (
	"io"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	. "github.com/drone/drone-cache-lib/storage"
)

// Options contains configuration for the S3 connection.
type Options struct {
	Server   string
	Username string
	Password string
}

type dummyStorage struct {
	opts   *Options
}

// New creates an implementation of Storage with Dummy as the backend.
func New(opts *Options) (Storage, error) {
	return &dummyStorage{
		opts:   opts,
	}, nil
}

func (s *dummyStorage) Get(p string, dst io.Writer) error {
	if _, err := os.Stat(p); err != nil {
		return err
	}

	return nil
}

func (s *dummyStorage) Put(p string, src io.Reader) error {
	log.Infof("Reading for %s", p)

	_, err := ioutil.ReadAll(src)

	if err != nil {
		log.Errorf("Failed to read for %s", p)
		return err
	}

	log.Infof("Finished reading for %s", p)

	return nil
}
