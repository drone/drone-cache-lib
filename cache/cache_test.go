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

		g.Describe("Rebuild", func() {
			g.It("Should rebuild with no errors", func() {
				s, err := dummy.New(dummyOpts)
				g.Assert(err == nil).IsTrue("failed to create storage")

				c := NewDefault(s)
				g.Assert(err == nil).IsTrue("failed to create cache")

				os.Chdir("/tmp/fixtures/mounts")
				err = c.Rebuild([]string{"test.txt", "subdir"}, "fixtures/tarfiles/file.tar")
				if err != nil {
					fmt.Printf("'Should rebuild with no errors' received unexpected error: %s\n", err)
				}
				g.Assert(err == nil).IsTrue("failed to rebuild the cache")
			})

			g.It("Should return error from channel", func() {
				s, err := dummy.New(dummyOpts)
				g.Assert(err == nil).IsTrue("failed to create storage")

				c := NewDefault(s)
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

				c := NewDefault(s)
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

				c := NewDefault(s)
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

				c := NewDefault(s)
				g.Assert(err == nil).IsTrue("failed to create cache")

				err = c.Restore("fixtures/test2.tar", "")
				g.Assert(err == nil).IsTrue("should not have returned error on missing file")
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
}

func cleanFixtures() {
	os.RemoveAll("/tmp/fixtures/")
}

func createDirectories() {
	directories := []string{
		"/tmp/fixtures/tarfiles",
		"/tmp/fixtures/mounts/subdir",
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
)
