package util

import (
	"testing"
	. "github.com/franela/goblin"
)

func TestArchive(t *testing.T) {
	g := Goblin(t)

	g.Describe("FromFilename", func() {
		g.It("Should return tarArchive for .tar", func() {
			_, err := FromFilename("filename.tar")
			g.Assert(err == nil).IsTrue("failed to determine .tar suffix")
		})

		g.It("Should return error for everything else", func() {
			_, err := FromFilename("filename.ttt")
			g.Assert(err != nil).IsTrue("failed to return error")
			g.Assert(err.Error()).Equal("Unknown file format for archive filename.ttt")
		})
	})
}
