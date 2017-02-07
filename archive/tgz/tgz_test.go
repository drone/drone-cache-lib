package tgz

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"
	. "github.com/franela/goblin"

	"github.com/drone/drone-cache-lib/archive"
)

type mountFile struct {
	Path string
	Content string
}

func TestTgzArchive(t *testing.T) {
	g := Goblin(t)
	wd, _ := os.Getwd()

	g.Describe("tgz package", func() {
		g.Before(func() {
			// Create necessary fixtures
			createFixtures()
		})

		g.After(func() {
			// Remove fixtures
			cleanFixtures()
		})

		g.Describe("New", func() {
			g.It("Should return tgzArchive", func() {
				tga := New()
				g.Assert(tga != nil).IsTrue("failed to create tgzArchive")
			})
		})

		g.Describe("Pack", func() {
			g.It("Should return no error", func() {
				tga := New()
				g.Assert(tga != nil).IsTrue("failed to create tgzArchive")

				os.Chdir("/tmp/fixtures/mounts")
				err, werr := packIt(tga, validMount, "/tmp/fixtures/tarfiles/test.tar.gz")
				os.Chdir(wd)

				if err != nil {
					fmt.Printf("Received unexpected err: %s\n", err)
				}
				g.Assert(err == nil).IsTrue("Failed to read the stream")
				if werr != nil {
					fmt.Printf("Received unexpected werr: %s\n", werr)
				}
				g.Assert(werr == nil).IsTrue("Failed to pack")
			})

			g.It("Should return error if mount does not exist", func() {
				tga := New()
				g.Assert(tga != nil).IsTrue("failed to create tgzArchive")

				err, werr := packIt(tga, invalidMount, "/tmp/fixtures/tarfiles/invalidMount.tar.gz")

				g.Assert(err == nil).IsTrue("Failed to read the stream")
				g.Assert(werr != nil).IsTrue("Failed to properly stat 'mount'")
				g.Assert(werr.Error()).Equal("stat mount1: no such file or directory")
			})
		})

		g.Describe("Unpack", func() {
			g.It("Should return no error", func() {
				tga := New()
				g.Assert(tga != nil).IsTrue("failed to create tgzArchive")

				err := unpackIt(tga, validFile)

				if err != nil {
					fmt.Printf("Received unexpected err: %s\n", err)
				}
				g.Assert(err == nil).IsTrue("Failed to unpack")
			})

			g.It("Should create files in correct strucutre", func() {
				g.Assert(exists("/tmp/extracted/test.txt")).IsTrue("failed to create test.txt")
				g.Assert(exists("/tmp/extracted/subdir")).IsTrue("failed to create subdir")
				g.Assert(exists("/tmp/extracted/subdir/test2.txt")).IsTrue("failed to create subdir/test2.txt")
				g.Assert(exists("/tmp/extracted/subdir/linkto_test.txt")).IsTrue("failed to create subdir/linkto_test.txt")
			})

			g.It("Should create files with correct content", func() {
				var err error
				var content []byte
				for _, element := range mountFiles {
					content, err = ioutil.ReadFile("/tmp/extracted/" + element.Path)
					g.Assert(err == nil).IsTrue("failed to read" + element.Path)
					g.Assert(string(content)).Equal(element.Content)
				}

				content, err = ioutil.ReadFile("/tmp/extracted/subdir/linkto_test.txt")
				g.Assert(err == nil).IsTrue("failed to read /tmp/extracted/subdir/linkto_test.txt")
				g.Assert(string(content)).Equal("hello\ngo\n")
			})

			g.It("Should return error on invalid tarfile", func() {
				tga := New()
				g.Assert(tga != nil).IsTrue("failed to create tgzArchive")

				err := unpackIt(tga, invalidFile)

				g.Assert(err != nil).IsTrue("Failed to return error")
				g.Assert(err.Error()).Equal("unexpected EOF")
			})

			g.It("Should return error on missing file", func() {
				tga := New()
				g.Assert(tga != nil).IsTrue("failed to create tgzArchive")

				err := unpackIt(tga, missingFile)

				g.Assert(err != nil).IsTrue("Failed to return error")
				g.Assert(err.Error()).Equal("open /tmp/fixtures/tarfiles/test2.tar.gz: no such file or directory")
			})
		})
	})
}

func packIt(a archive.Archive, srcs []string, dst string) (error, error) {
	reader, writer := io.Pipe()
	defer reader.Close()

	cw := make(chan error, 1)
	defer close(cw)

	go func() {
		defer writer.Close()

		cw <- a.Pack(srcs, writer)
	}()

	bytes, err := ioutil.ReadAll(reader)
	ioutil.WriteFile(dst, bytes, 0644)

	werr := <-cw

	return err, werr
}

func unpackIt(a archive.Archive, src string) error {
	reader, writer := io.Pipe()

	cw := make(chan error, 1)
	defer close(cw)

	f, err := os.Open(src)

	if err != nil {
		return err
	}

	go func() {
		defer writer.Close()

		_, err = io.Copy(writer, f)

		if err != nil {
			cw <- err
			return
		}
	}()

	return a.Unpack("/tmp/extracted", reader)
}

func createBadTgzfile() {
	content := []byte("hello\ngo\n")
	err := ioutil.WriteFile("/tmp/fixtures/tarfiles/bad.tar.gz", content, 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func createMountContent() {
	// Write files and their content
	var err error
	for _, element := range mountFiles {
		err = ioutil.WriteFile("/tmp/fixtures/mounts/" + element.Path, []byte(element.Content), 0644)
		if err != nil {
			log.Fatalln(err)
		}
	}

	// Create a symlink
	os.Symlink("../test.txt", "/tmp/fixtures/mounts/subdir/linkto_test.txt")
}

func createFixtures() {
	createDirectories()
	createBadTgzfile()
	createMountContent()
}

func cleanFixtures() {
	os.RemoveAll("/tmp/fixtures/")
	os.RemoveAll("/tmp/extracted/")
}

func createDirectories() {
	directories := []string{
		"/tmp/fixtures/tarfiles",
		"/tmp/fixtures/mounts/subdir",
		"/tmp/extracted",
	}

	for _, directory := range directories {
		if _, err := os.Stat(directory); os.IsNotExist(err) {
			os.MkdirAll(directory, os.FileMode(int(0755)))
		}
	}
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil { return true }
	if os.IsNotExist(err) { return false }
	return true
}

var (
	invalidMount = []string{
		"mount1",
		"mount2",
	}

	mountFiles = []mountFile{
		{Path: "test.txt", Content: "hello\ngo\n"},
		{Path: "subdir/test2.txt", Content: "hello2\ngo\n"},
	}

	validMount = []string{
		"test.txt",
		"subdir",
	}

	validFile = "/tmp/fixtures/tarfiles/test.tar.gz"
	invalidFile = "/tmp/fixtures/tarfiles/bad.tar.gz"
	missingFile = "/tmp/fixtures/tarfiles/test2.tar.gz"
)
