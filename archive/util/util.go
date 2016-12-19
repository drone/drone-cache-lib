package util

import (
	"fmt"
	"strings"

	. "github.com/drone/drone-cache-lib/archive"
	"github.com/drone/drone-cache-lib/archive/tar"
)

// FromFilename determines the archive format to use based on the name.
func FromFilename(name string) (Archive, error) {
	if strings.HasSuffix(name, ".tar") {
		return tar.New(), nil
	}

	return nil, fmt.Errorf("Unknown file format for archive %s", name)
}
