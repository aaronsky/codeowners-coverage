package git

import (
	"os"
	"path/filepath"

	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-git.v4"
)

// Worktree is a re-export of go-git Worktree
type Worktree = git.Worktree

// SkipDir is a re-export of the filepath.SkipDir error
var SkipDir = filepath.SkipDir

// WalkFunc is the callback function signature used by WalkTree
type WalkFunc func(path string, info os.FileInfo, err error) error

// WalkTree will recursively crawl a go-git Worktree and invoke a callback for each file and directory
func WalkTree(fs billy.Filesystem, cb WalkFunc) error {
	root := fs.Join(".")
	info, err := fs.Stat(root)
	if err != nil {
		err = cb(root, nil, err)
	} else {
		err = walk(fs, root, info, cb)
	}
	if err == SkipDir {
		return nil
	}
	return err
}

func walk(fs billy.Filesystem, path string, info os.FileInfo, cb WalkFunc) error {
	if !info.IsDir() {
		return cb(path, info, nil)
	} else if info.Name() == ".git" {
		return nil
	}

	infos, err := fs.ReadDir(path)
	err1 := cb(path, info, err)

	if err != nil || err1 != nil {
		return err1
	}

	for _, info := range infos {
		filename := fs.Join(path, info.Name())
		fileInfo, err := fs.Stat(filename)
		if err != nil {
			if err := cb(filename, fileInfo, err); err != nil && err != SkipDir {
				return err
			}
			continue
		}
		err = walk(fs, filename, fileInfo, cb)
		if err != nil && (!fileInfo.IsDir() || err != SkipDir) {
			return err
		}
	}

	return nil
}
