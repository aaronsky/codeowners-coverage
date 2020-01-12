package git

import (
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

func Clone(url string) (*git.Repository, error) {
	fs := memfs.New()
	storer := memory.NewStorage()
	return git.Clone(storer, fs, &git.CloneOptions{
		URL: url,
	})
}
