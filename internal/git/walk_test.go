package git

import (
	"os"
	"testing"

	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
)

func TestWalkTree(t *testing.T) {
	mockFs, countFiles := setupPopulatedFilesystem()
	var discoveredFiles int
	WalkTree(mockFs, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			discoveredFiles++
		}
		return nil
	})
	if discoveredFiles != countFiles {
		t.Errorf("discovered files count %d does not equal expected files count %d", discoveredFiles, countFiles)
	}
}

func TestWalkTreeEmptyTree(t *testing.T) {
	mockFs := memfs.New()
	err := WalkTree(mockFs, func(path string, info os.FileInfo, err error) error {
		return err
	})
	if err == nil {
		t.Error(err)
	}
}

func setupPopulatedFilesystem() (fs billy.Filesystem, countFiles int) {
	fs = memfs.New()

	fs.MkdirAll(".git/hooks", os.ModeDir)
	fs.MkdirAll("dirA/dirB/dirC", os.ModeDir)
	fs.MkdirAll("dirA/dirE/dirF", os.ModeDir)
	fs.MkdirAll("dirA/dirE/dirG", os.ModeDir)
	fs.MkdirAll("dirH/dirI", os.ModeDir)
	fs.MkdirAll("dirJ", os.ModeDir)

	fs.Create("dog.txt")
	countFiles++
	fs.Create("dirA/dog.txt")
	countFiles++
	fs.Create("dirA/dirB/dog.txt")
	countFiles++
	fs.Create("dirA/dirB/dirC/dog.txt")
	countFiles++
	fs.Create("dirA/dirB/dirC/cat.txt")
	countFiles++
	fs.Create("dirA/dirE/dog.txt")
	countFiles++
	fs.Create("dirA/dirE/dirF/dog.txt")
	countFiles++
	fs.Create("dirA/dirE/dirG/dog.txt")
	countFiles++
	fs.Create("dirH/dog.txt")
	countFiles++
	fs.Create("dirH/dirI/dog.txt")
	countFiles++

	return fs, countFiles
}
