package git

import (
	"gopkg.in/src-d/go-git.v4"
)

// Open opens a repository on-disk at the given path. All operations from here-on are performed using the on-disk filesystem.
func Open(path string) (*git.Repository, error) {
	return git.PlainOpenWithOptions(path, &git.PlainOpenOptions{
		DetectDotGit: false,
	})
}
