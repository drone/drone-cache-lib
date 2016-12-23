package noopflusher

import (
	"os"
	"testing"
  "time"

	. "github.com/franela/goblin"
	"github.com/drone/drone-cache-lib/storage"
	"github.com/nowk/go-fakefileinfo"
)

func TestTimeFlusher(t *testing.T) {
	g := Goblin(t)

	g.Describe("Find", func() {
		g.It("Should always return no matched files", func() {
			f := New()

			matched, err := f.Find(fileEntries)
			g.Assert(err == nil).IsTrue("failed to Find files")

			// Make sure we found no files
			g.Assert(len(matched)).Equal(0)
		})
	})
}

var (
	fileEntries = []storage.FileEntry{
		{Path: "proj1/master/archive.txt", Info: fakefileinfo.New("archive.txt", int64(123), os.ModeType, time.Now().AddDate(0, 0, -1), false, nil)},
		{Path: "proj1/newtest/archive.txt", Info: fakefileinfo.New("archive.txt", int64(123), os.ModeType, time.Now().AddDate(0, 0, -3), false, nil)},
		{Path: "proj1/oldtest/archive.txt", Info: fakefileinfo.New("archive.txt", int64(123), os.ModeType, time.Now().AddDate(0, 0, -10), false, nil)},
	}
)
