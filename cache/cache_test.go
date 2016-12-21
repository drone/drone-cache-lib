package cache

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
	. "github.com/franela/goblin"

	"github.com/drone/drone-cache-lib/storage/dummy"
)

func TestCache(t *testing.T) {
	g := Goblin(t)
	wd, _ := os.Getwd()

	g.Describe("cache package", func() {

		g.Before(func() {
			os.Chdir("/tmp")
			createFixtures()
		})

		g.BeforeEach(func() {
			os.Chdir("/tmp")
		})

		g.After(func() {
			os.Chdir(wd)
			cleanFixtures()
		})

		g.Describe("New", func() {
			g.It("Should create new Cache", func() {
				s, err := dummy.New(dummyOpts)
				g.Assert(err == nil).IsTrue("failed to create storage")

				_, err = New(s)
				g.Assert(err == nil).IsTrue("failed to create cache")
			})
		})

		g.Describe("Rebuild", func() {
			g.It("Should rebuild with no errors", func() {
				s, err := dummy.New(dummyOpts)
				g.Assert(err == nil).IsTrue("failed to create storage")

				c, err := New(s)
				g.Assert(err == nil).IsTrue("failed to create cache")

				os.Chdir("/tmp/fixtures/mounts")
				err = c.Rebuild([]string{"test.txt", "subdir"}, "fixtures/tarfiles/file.tar")
				if err != nil {
					fmt.Printf("'Should rebuild with no errors' received unexpected error: %s\n", err)
				}
				g.Assert(err == nil).IsTrue("failed to rebuild the cache")
			})

			g.It("Should return error on failure", func() {
				s, err := dummy.New(dummyOpts)
				g.Assert(err == nil).IsTrue("failed to create storage")

				c, err := New(s)
				g.Assert(err == nil).IsTrue("failed to create cache")

				err = c.Rebuild([]string{"mount1", "mount2"}, "file.ttt")
				g.Assert(err != nil).IsTrue("failed to return error")
				g.Assert(err.Error()).Equal("Unknown file format for archive file.ttt")
			})

			g.It("Should return error from channel", func() {
				s, err := dummy.New(dummyOpts)
				g.Assert(err == nil).IsTrue("failed to create storage")

				c, err := New(s)
				g.Assert(err == nil).IsTrue("failed to create cache")

				err = c.Rebuild([]string{"mount1", "mount2"}, "file.tar")
				g.Assert(err != nil).IsTrue("failed to return error")
				g.Assert(err.Error()).Equal("stat mount1: no such file or directory")
			})
		})

		g.Describe("Restore", func() {
			g.It("Should restore with no errors", func() {
				s, err := dummy.New(dummyOpts)
				g.Assert(err == nil).IsTrue("failed to create storage")

				c, err := New(s)
				g.Assert(err == nil).IsTrue("failed to create cache")

				err = c.Restore("fixtures/test.tar", "")
				if err != nil {
					fmt.Printf("Received unexpected error: %s\n", err)
				}
				g.Assert(err == nil).IsTrue("failed to rebuild the cache")
			})

			g.It("Should restore from fallback if path does not exist", func() {
				s, err := dummy.New(dummyOpts)
				g.Assert(err == nil).IsTrue("failed to create storage")

				c, err := New(s)
				g.Assert(err == nil).IsTrue("failed to create cache")

				err = c.Restore("fixtures/test2.tar", "fixtures/test.tar")
				if err != nil {
					fmt.Printf("Received unexpected error: %s\n", err)
				}
				g.Assert(err == nil).IsTrue("failed to rebuild the cache")
			})

			g.It("Should not return error on missing file", func() {
				s, err := dummy.New(dummyOpts)
				g.Assert(err == nil).IsTrue("failed to create storage")

				c, err := New(s)
				g.Assert(err == nil).IsTrue("failed to create cache")

				err = c.Restore("fixtures/test2.tar", "")
				g.Assert(err == nil).IsTrue("should not have returned error on missing file")
			})

			g.It("Should return error on unknown file format", func() {
				s, err := dummy.New(dummyOpts)
				g.Assert(err == nil).IsTrue("failed to create storage")

				c, err := New(s)
				g.Assert(err == nil).IsTrue("failed to create cache")

				err = c.Restore("fixtures/test2.ttt", "")
				g.Assert(err != nil).IsTrue("failed to return filetype error")
			})
		})

		g.Describe("Cleanup", func() {

			g.BeforeEach(func() {
				createCleanupContent()
			})

			g.It("Should find no files to cleanup", func() {
				s, err := dummy.New(dummyOpts)
				g.Assert(err == nil).IsTrue("failed to create storage")

				c, err := New(s)
				g.Assert(err == nil).IsTrue("failed to create cache")

				c.Cleanup("fixtures/cleanup/proj1", time.Duration(20*24)*time.Hour)
				g.Assert(err == nil).IsTrue("failed to cleanup nothing")

				// Check expected files still exist
				checkFileExists("/tmp/fixtures/cleanup/proj1/master/archive.txt", g)
				checkFileExists("/tmp/fixtures/cleanup/proj1/oldtest/archive.txt", g)
				checkFileExists("/tmp/fixtures/cleanup/proj1/newtest/archive.txt", g)
			})

			g.It("Should find some files to cleanup", func() {
				s, err := dummy.New(dummyOpts)
				g.Assert(err == nil).IsTrue("failed to create storage")

				c, err := New(s)
				g.Assert(err == nil).IsTrue("failed to create cache")

				// Perform Cleanup
				c.Cleanup("fixtures/cleanup/proj1", time.Duration(9*24)*time.Hour)
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

func checkFileExists(fileName string, g *G) {
	_, err := os.Stat(fileName)
	g.Assert(err == nil).IsTrue(fileName + " should still exist")
}

func checkFileRemoved(fileName string, g *G) {
	_, err := os.Stat(fileName)
	g.Assert(err != nil).IsTrue("Failed to clean " + fileName)
}

func createFixtures() {
	createDirectories()
	createMountContent()
	createCleanupContent()
}

func cleanFixtures() {
	os.RemoveAll("/tmp/fixtures/")
	// os.RemoveAll("/tmp/extracted/")
}

func createDirectories() {
	directories := []string{
		"/tmp/fixtures/tarfiles",
		"/tmp/fixtures/mounts/subdir",
		"/tmp/fixtures/cleanup/proj1/master",
		"/tmp/fixtures/cleanup/proj1/newtest",
		"/tmp/fixtures/cleanup/proj1/oldtest",
		"/tmp/fixtures/cleanup/proj2/master",
		"/tmp/fixtures/cleanup/proj2/newtest",
		"/tmp/fixtures/cleanup/proj2/oldtest",
	}

	for _, directory := range directories {
		if _, err := os.Stat(directory); os.IsNotExist(err) {
			os.MkdirAll(directory, os.FileMode(int(0755)))
		}
	}
}

func createMountContent() {
	var err error
	for _, element := range mountFiles {
		err = ioutil.WriteFile("/tmp/fixtures/mounts/" + element.Path, []byte(element.Content), 0644)
		if err != nil {
			log.Fatalln(err)
		}
	}
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

type testFile struct {
	Path string
	Content string
	Time time.Time
}

var (
	dummyOpts = &dummy.Options{
		Server:   "myserver.com",
		Username: "johndoe",
		Password: "supersecret",
	}

	mountFiles = []testFile{
		{Path: "test.txt", Content: "hello\ngo\n"},
		{Path: "subdir/test2.txt", Content: "hello2\ngo\n"},
	}

	cleanupFiles = []testFile{
		{Path: "proj1/master/archive.txt", Content: "hello\ngo\n", Time: time.Now()},
		{Path: "proj1/newtest/archive.txt", Content: "hello2\ngo\n", Time: time.Now().AddDate(0, 0, -1)},
		{Path: "proj1/oldtest/archive.txt", Content: "hello\ngo\n", Time: time.Now().AddDate(0, 0, -10)},
	}
)
