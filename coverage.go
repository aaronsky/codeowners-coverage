// Package coverage will generate a coverage report for CODEOWNERS in the repository.
package coverage

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aaronsky/codeowners-coverage/internal/codeowners"
	"github.com/aaronsky/codeowners-coverage/internal/git"
	"gopkg.in/src-d/go-billy.v4"
)

// Report contains information on the codeowner coverage of files in a repository
type Report struct {
	RemoteURL         string  `json:"remote_url"`
	SHA               string  `json:"sha"`
	CoveredFilesCount int     `json:"covered_files_count"`
	TotalFilesCount   int     `json:"total_files_count"`
	CoverageRatio     float64 `json:"coverage_ratio"`
}

// NewCoverageReport produces a coverage report from the given repository
// Modifies state of the given repository by performing a git-clean.
func NewCoverageReport(path string) (*Report, error) {
	repository, err := git.Open(path)
	if err != nil {
		return nil, err
	}

	remote, err := repository.Remote("origin")
	if err != nil {
		return nil, err
	}
	remoteURL := remote.Config().URLs[0]
	headSHA, err := repository.Head()
	if err != nil {
		return nil, err
	}

	worktree, err := repository.Worktree()
	if err != nil {
		return nil, err
	}
	status, err := worktree.Status()
	if err != nil {
		return nil, err
	}
	git.CleanWorktree(worktree)

	fs := worktree.Filesystem
	owners, err := codeowners.LoadFromFilesystem(fs)
	if err != nil {
		return nil, err
	}

	report := &Report{RemoteURL: remoteURL, SHA: headSHA.Hash().String()}
	err = report.setCoverage(status, fs, owners)
	if err != nil {
		return nil, err
	}
	return report, nil
}

// setCoverage mutates the Report object to store information on covered files and the ratio of coverage
func (r *Report) setCoverage(status git.Status, fs billy.Filesystem, owners codeowners.Codeowners) error {
	var totalFilesCount int
	var filesToCheckCoverage []string
	var coveredFilesCount int

	err := git.WalkTree(fs, func(path string, info os.FileInfo, err error) error {
		if !info.Mode().IsRegular() {
			// not file
			return nil
		} else if codeowners.PathIsCodeowners(path, fs) {
			// skip codeowners
			return nil
		} else if status.IsUntracked(path) {
			// path is untracked
			return nil
		}

		totalFilesCount++
		filesToCheckCoverage = append(filesToCheckCoverage, path)

		return nil
	})
	if err != nil {
		return err
	}

	for _, path := range filesToCheckCoverage {
		ownersForPath := owners.Owners(path)
		if len(ownersForPath) > 0 {
			coveredFilesCount++
		}
	}

	r.CoveredFilesCount = coveredFilesCount
	r.TotalFilesCount = totalFilesCount
	if totalFilesCount > 0 {
		r.CoverageRatio = float64(coveredFilesCount) / float64(totalFilesCount)
	}

	return nil
}

type reportFormat string

const (
	// ReportFormatJSON is a constant representing the JSON format for a Report object
	ReportFormatJSON reportFormat = "json"
)

// ToFormat converts the report to a string in the given format.
// Currently only "json" is supported.
func (r *Report) ToFormat(format reportFormat) (string, error) {
	switch format {
	case ReportFormatJSON:
		bytes, err := json.Marshal(r)
		if err != nil {
			return "", err
		}
		return string(bytes), nil
	default:
		return "", fmt.Errorf("unsupported reportFormat")
	}
}
