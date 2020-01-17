package git

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

// BasicAuth is a re-export of go-git BasicAuth
type BasicAuth = http.BasicAuth

// TokenAuth is a re-export of go-git TokenAuth
type TokenAuth = http.TokenAuth

// // Clone clones the given repository at the given URL into an in-memory file system, using the given authentication mechanism
// func Clone(url string, auth transport.AuthMethod) (*git.Repository, error) {
// 	fs := memfs.New()
// 	storer := memory.NewStorage()
// 	return git.Clone(storer, fs, &git.CloneOptions{
// 		Auth: auth,
// 		URL:  url,
// 	})
// }

// Open opens a repository on-disk at the given path. All operations from here-on are performed using the on-disk filesystem.
func Open(path string) (*git.Repository, error) {
	return git.PlainOpenWithOptions(path, &git.PlainOpenOptions{
		DetectDotGit: false,
	})
}
