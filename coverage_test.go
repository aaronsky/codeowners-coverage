package coverage

import (
	"os"
	"testing"

	"github.com/aaronsky/codeowners-coverage/internal/codeowners"
	"github.com/aaronsky/codeowners-coverage/internal/git"
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
	go_git "gopkg.in/src-d/go-git.v4"
)

func TestSetCoverage(t *testing.T) {
	report := Report{
		RemoteURL: "https://github.com/gitignore/gitignore",
		SHA:       "0b18ca88",
	}

	mockStatus, mockFs, _ := setupPopulatedFilesystem()
	owners, err := codeowners.LoadFromFilesystem(mockFs)
	if err != nil {
		t.Error(err)
	}

	err = report.setCoverage(mockStatus, mockFs, owners)
	if err != nil {
		t.Error(err)
	}
	if report.TotalFilesCount != 5 {
		t.Errorf("expected total file count to be 5, but it was %d", report.TotalFilesCount)
	}
	if report.CoveredFilesCount != 2 {
		t.Errorf("expected covered file count to be 2, but it was %d", report.CoveredFilesCount)
	}
	if report.CoverageRatio != 0.4 {
		t.Errorf("expected coverage ratio to be 0.4, but it was %f", report.CoverageRatio)
	}
}

func TestToFormatWithData(t *testing.T) {
	report := Report{
		RemoteURL:         "https://github.com/gitignore/gitignore",
		SHA:               "0b18ca88",
		CoveredFilesCount: 30,
		TotalFilesCount:   50,
		CoverageRatio:     0.6,
	}
	expectedJSON := `{"remote_url":"https://github.com/gitignore/gitignore","sha":"0b18ca88","covered_files_count":30,"total_files_count":50,"coverage_ratio":0.6}`
	jsonString, err := report.ToFormat(ReportFormatJSON)
	if err != nil {
		t.Error(err)
	}
	if jsonString != expectedJSON {
		t.Error("json did not match expected")
	}
}

func TestToFormatWithoutData(t *testing.T) {
	report := Report{}
	expectedJSON := `{"remote_url":"","sha":"","covered_files_count":0,"total_files_count":0,"coverage_ratio":0}`
	jsonString, err := report.ToFormat(ReportFormatJSON)
	if err != nil {
		t.Error(err)
	}
	if jsonString != expectedJSON {
		t.Error("json did not match expected")
	}
}

func TestToFormatWithInvalidFormat(t *testing.T) {
	report := Report{}
	jsonString, err := report.ToFormat(reportFormat("dogs"))
	if err == nil {
		t.Error("Expected error that format is unsupported")
	}
	if jsonString != "" {
		t.Error("json did not match expected")
	}
}

func setupPopulatedFilesystem() (status git.Status, fs billy.Filesystem, countFiles int) {
	status = git.Status{}
	fs = memfs.New()

	fs.MkdirAll("src", os.ModeDir)

	createFile("CODEOWNERS", []byte(`*.js		@org/team_reviewers`), fs, &status)
	countFiles++

	createFile("README.md", nil, fs, &status)
	countFiles++

	createFile("index.js", nil, fs, &status)
	countFiles++

	createFile("src/app.js", nil, fs, &status)
	countFiles++

	createFile("src/index.html", nil, fs, &status)
	countFiles++

	createFile("src/index.css", nil, fs, &status)
	countFiles++

	return status, fs, countFiles
}

func createFile(path string, content []byte, fs billy.Filesystem, status *git.Status) {
	file, _ := fs.Create(path)
	if len(content) > 0 {
		file.Write(content)
	}
	fileStatus := status.File(path)
	fileStatus.Staging = go_git.Unmodified
	fileStatus.Worktree = go_git.Unmodified
}
