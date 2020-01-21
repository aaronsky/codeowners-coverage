package codeowners

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/src-d/go-billy.v4/memfs"
)

func TestWorktreeContainsCodeowners(t *testing.T) {
	mockFs := memfs.New()
	if !PathIsCodeowners(filepath.Join(".", "CODEOWNERS"), mockFs) {
		t.Error("expected root to be valid for containing CODEOWNERS")
	}
	if !PathIsCodeowners(filepath.Join("docs", "CODEOWNERS"), mockFs) {
		t.Error("expected docs/ to be valid for containing CODEOWNERS")
	}
	if !PathIsCodeowners(filepath.Join(".github", "CODEOWNERS"), mockFs) {
		t.Error("expected .github/ to be valid for containing CODEOWNERS")
	}
	if PathIsCodeowners(filepath.Join("src", "CODEOWNERS"), mockFs) {
		t.Error("expected src to be invalid for containing CODEOWNERS")
	}
	if PathIsCodeowners(filepath.Join("github", "OWNERS"), mockFs) {
		t.Error("expected github to be invalid for containing CODEOWNERS")
	}
}

func TestLoadFromFilesystem(t *testing.T) {
	mockFs := memfs.New()
	file, _ := mockFs.Create("CODEOWNERS")
	file.Write([]byte(`# a fun comment!
	*.js		@org/team_reviewers
	!boat/*		@org/team_reviewers @jeffery`))
	owners, err := LoadFromFilesystem(mockFs)
	if err != nil {
		t.Error(err)
	}
	if owners == nil {
		t.Error("expected codeowners object to be loaded")
	}
	names := owners.Owners("dog.js")
	if len(names) == 0 {
		t.Error("expected owners to be found for 'dog.js'")
	}
}

func TestLoadFromFilesystemFails(t *testing.T) {
	mockFs := memfs.New()
	err := mockFs.MkdirAll(".github", os.ModeDir)
	if err != nil {
		t.Error(err)
	}
	err = mockFs.MkdirAll("docs", os.ModeDir)
	if err != nil {
		t.Error(err)
	}
	_, err = LoadFromFilesystem(mockFs)
	if err == nil {
		t.Error("expected to fail loading from filesystem")
	}
}

func TestOwnersFromNilReturnsEmpty(t *testing.T) {
	var o Codeowners
	owners := o.Owners("jeff")
	if len(owners) > 0 {
		t.Error("expected no owners to be returned for nil Codeowners object")
	}
}
