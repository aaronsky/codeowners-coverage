package git

import (
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-git.v4"
)

// Open opens a repository on-disk at the given path. All operations from here-on are performed using the on-disk filesystem.
func Open(path string) (*git.Repository, error) {
	return git.PlainOpenWithOptions(path, &git.PlainOpenOptions{
		DetectDotGit: false,
	})
}

// IsPathTracked returns whether or not a path is tracked under the current repository
func IsPathTracked(path string, fs billy.Filesystem) (bool, error) {
	// FIXME: currently a no-op with dangerous implications
	return true, nil
}
