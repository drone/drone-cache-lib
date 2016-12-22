package timeflusher

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
		g.It("Should return subset of matched files", func() {
			f := New(nineDays)

			matched, err := f.Find(fileEntries)
			g.Assert(err == nil).IsTrue("failed to Find files")

			// Make sure we found 1 file
			g.Assert(len(matched)).Equal(1)

			// Make sure it's the right file
			g.Assert(matched[0]).Equal(fileEntries[2])
		})
	})
}

var (
	nineDays = time.Duration(9*24)*time.Hour

	fileEntries = []storage.FileEntry{
		{Path: "proj1/master/archive.txt", Info: fakefileinfo.New("archive.txt", int64(123), os.ModeType, time.Now().AddDate(0, 0, -1), false, nil)},
		{Path: "proj1/newtest/archive.txt", Info: fakefileinfo.New("archive.txt", int64(123), os.ModeType, time.Now().AddDate(0, 0, -3), false, nil)},
		{Path: "proj1/oldtest/archive.txt", Info: fakefileinfo.New("archive.txt", int64(123), os.ModeType, time.Now().AddDate(0, 0, -10), false, nil)},
	}
)
