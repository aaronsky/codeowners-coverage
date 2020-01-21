// Package git wraps functionality in the go-git package to manipulate Git repositories
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

// Status is a re-export of go-git Status
type Status = git.Status

// CleanWorktree runs the equivalent of `git clean -xfd .` on the current worktree
func CleanWorktree(worktree *git.Worktree) error {
	return worktree.Clean(&git.CleanOptions{
		Dir: true,
	})
}
