package cache

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
	. "github.com/franela/goblin"

	"github.com/drone/drone-cache-lib/storage"
	"github.com/drone/drone-cache-lib/storage/dummy"
)

func TestFlusher(t *testing.T) {
	g := Goblin(t)
	wd, _ := os.Getwd()

	g.Describe("flusher package", func() {

		g.Before(func() {
			os.Chdir("/tmp")
			createFlusherFixtures()
		})

		g.BeforeEach(func() {
			os.Chdir("/tmp")
		})

		g.After(func() {
			os.Chdir(wd)
			cleanFixtures()
		})

		g.Describe("Cleanup", func() {

			g.BeforeEach(func() {
				createCleanupContent()
			})

			g.It("Should find no files to cleanup", func() {
				s, err := dummy.New(dummyOpts)
				g.Assert(err == nil).IsTrue("failed to create storage")

				f := NewFlusher(s, noFind)

				f.Flush("fixtures/cleanup/proj1")
				g.Assert(err == nil).IsTrue("failed to cleanup nothing")

				// Check expected files still exist
				checkFileExists("/tmp/fixtures/cleanup/proj1/master/archive.txt", g)
				checkFileExists("/tmp/fixtures/cleanup/proj1/oldtest/archive.txt", g)
				checkFileExists("/tmp/fixtures/cleanup/proj1/newtest/archive.txt", g)
			})

			g.It("Should find some files to cleanup", func() {
				s, err := dummy.New(dummyOpts)
				g.Assert(err == nil).IsTrue("failed to create storage")

				f := NewDefaultFlusher(s)

				// Perform Cleanup
				f.Flush("fixtures/cleanup/proj1")
				g.Assert(err == nil).IsTrue("failed to cleanup nothing")

				// Check expected files no longer exist
				checkFileRemoved("/tmp/fixtures/cleanup/proj1/oldtest/archive.txt", g)

				// Check expected files still exist
				checkFileExists("/tmp/fixtures/cleanup/proj1/master/archive.txt", g)
				checkFileExists("/tmp/fixtures/cleanup/proj1/newtest/archive.txt", g)
			})
		})
	})
}

func createFlusherFixtures() {
	createDirectories(flusherFixtureDirectories)
	createCleanupContent()
}

func createCleanupContent() {
	var name string
	var err error
	for _, element := range cleanupFiles {
		name = "/tmp/fixtures/cleanup/" + element.Path
		err = ioutil.WriteFile(name, []byte(element.Content), 0644)
		if err != nil {
			log.Fatalln(err)
		}
		err = os.Chtimes(name, element.Time, element.Time)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

var (
	cleanupFiles = []testFile{
		{Path: "proj1/master/archive.txt", Content: "hello\ngo\n", Time: time.Now()},
		{Path: "proj1/newtest/archive.txt", Content: "hello2\ngo\n", Time: time.Now().AddDate(0, 0, -1)},
		{Path: "proj1/oldtest/archive.txt", Content: "hello\ngo\n", Time: time.Now().AddDate(0, 0, -40)},
	}

	noFind = func(file storage.FileEntry) bool {
		if file.Info.IsDir() {
			return false
		}

		if file.Info.ModTime().Before(time.Now().AddDate(0, 0, -60)) {
			return true
		}

		return false
	}

	flusherFixtureDirectories = []string{
		"/tmp/fixtures/cleanup/proj1/master",
		"/tmp/fixtures/cleanup/proj1/newtest",
		"/tmp/fixtures/cleanup/proj1/oldtest",
		"/tmp/fixtures/cleanup/proj2/master",
		"/tmp/fixtures/cleanup/proj2/newtest",
		"/tmp/fixtures/cleanup/proj2/oldtest",
	}
)
