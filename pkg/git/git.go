package git

import (
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

// BasicAuth is a re-export of go-git BasicAuth
type BasicAuth = http.BasicAuth

// TokenAuth is a re-export of go-git TokenAuth
type TokenAuth = http.TokenAuth

// Clone clones the given repository at the given URL into an in-memory file system, using the given authentication mechanism
func Clone(url string, auth transport.AuthMethod) (*git.Repository, error) {
	fs := memfs.New()
	storer := memory.NewStorage()
	return git.Clone(storer, fs, &git.CloneOptions{
		Auth: auth,
		URL:  url,
	})
}
