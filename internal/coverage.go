package internal

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aaronsky/codeowners-coverage/pkg/codeowners"
	"github.com/aaronsky/codeowners-coverage/pkg/git"
	"github.com/aaronsky/codeowners-coverage/pkg/github"
)

// Report contains information on the codeowner coverage of files in a repository
type Report struct {
	RemoteURL         string `json:"remote_url"`
	SHA               string `json:"sha"`
	CoveredFilesCount int    `json:"covered_files_count"`
	TotalFilesCount   int    `json:"total_files_count"`
}

// NewCoverageReport produces a coverage report from the given repository
func NewCoverageReport(remote, token string) (*Report, error) {
	// get list of all files in the repository
	// check for a codeowners file
	// 		if none, produce output of 0% coverage and a message that no file was found
	// parse each of the patterns of the codeowners file
	// coverage % = sum of count of file matches for each pattern / count of all files
	// produce output

	hub, err := github.NewGithub(remote, token)
	if err != nil {
		return nil, err
	}
	repository, err := git.Clone(remote)
	if err != nil {
		return nil, err
	}
	headSHA, err := repository.Head()
	if err != nil {
		return nil, err
	}
	worktree, err := repository.Worktree()
	if err != nil {
		return nil, err
	}
	codeowners, err := codeowners.NewFromTree(worktree)
	if err != nil {
		return nil, err
	}

	// fs := worktree.Filesystem

	fmt.Println(hub)
	fmt.Println(codeowners)
	return &Report{
		RemoteURL:         remote,
		SHA:               headSHA.Hash().String(),
		TotalFilesCount:   0,
		CoveredFilesCount: 0,
	}, nil
}

// ToFormat converts the report to a string in the given format.
// Currently only "json" is supported.
func (r *Report) ToFormat(format string) (string, error) {
	if strings.ToLower(format) == "json" {
		bytes, err := json.Marshal(r)
		if err != nil {
			return "", nil
		}
		return string(bytes), nil
	}
	return "", nil
}
