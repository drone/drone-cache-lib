package tgz

// special thanks to this medium article:
// https://medium.com/@skdomino/taring-untaring-files-in-go-6b07cf56bc07

import (
	"io"
	"compress/gzip"

	. "github.com/drone/drone-cache-lib/archive"
	"github.com/drone/drone-cache-lib/archive/tar"
)

type tgzArchive struct{}

// NewTarArchive creates an Archive that uses the .tar file format.
func New() Archive {
	return &tgzArchive{}
}

func (a *tgzArchive) Pack(srcs []string, w io.Writer) error {
	gw := gzip.NewWriter(w)
	defer gw.Close()

	tar_wr := tar.New()

	err := tar_wr.Pack(srcs, gw)

	return err
}

func (a *tgzArchive) Unpack(dst string, r io.Reader) error {
	gr, err := gzip.NewReader(r)

	if err != nil	{
		return err
	}

	tar_r := tar.New()

	fw_err := tar_r.Unpack(dst, gr)

	return fw_err
}
