package internal

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/aaronsky/codeowners-coverage/pkg/codeowners"
	"github.com/aaronsky/codeowners-coverage/pkg/git"
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
func NewCoverageReport(remote, token string) (*Report, error) {
	// It doesn't matter what we provide here for a username so we just pass the token twice
	repository, err := git.Clone(remote, &git.BasicAuth{
		Password: token,
	})
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
	owners, err := codeowners.NewFromTree(worktree)
	if err != nil {
		return nil, err
	}

	report := &Report{RemoteURL: remote, SHA: headSHA.Hash().String()}
	err = report.setCoverage(worktree, owners)
	if err != nil {
		return nil, err
	}
	return report, nil
}

// setCoverage mutates the Report object to store information on covered files and the ratio of coverage
func (r *Report) setCoverage(worktree *git.Worktree, owners codeowners.Codeowners) error {
	var totalFilesCount int
	var filesToCheckCoverage []string
	var coveredFilesCount int

	err := git.WalkTree(worktree, func(path string, info os.FileInfo, err error) error {
		if info.Mode().IsRegular() && !codeowners.PathIsCodeowners(path, worktree) {
			totalFilesCount++
			filesToCheckCoverage = append(filesToCheckCoverage, path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	for _, path := range filesToCheckCoverage {
		ownersForPath, err := owners.Owners(path, worktree)
		if err != nil {
			continue
		}
		if len(ownersForPath) > 0 {
			coveredFilesCount++
		}
	}

	r.CoveredFilesCount = coveredFilesCount
	r.TotalFilesCount = totalFilesCount
	r.CoverageRatio = float64(coveredFilesCount) / float64(totalFilesCount)

	return nil
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
