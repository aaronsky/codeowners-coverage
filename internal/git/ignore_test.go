package git

import "testing"

func TestCommentLinePattern(t *testing.T) {
	pattern, err := CompileIgnorePattern("# This is a comment.")
	if err == nil {
		t.Error("expected string to not compile")
	} else if pattern != nil {
		t.Errorf("expected pattern to be nil (%s)", pattern)
	}
}

func TestEmptyLinePattern(t *testing.T) {
	pattern, err := CompileIgnorePattern("                    ")
	if err == nil {
		t.Error("expected string to not compile")
	} else if pattern != nil {
		t.Errorf("expected pattern to be nil (%s)", pattern)
	}
}

func TestSingleAsteriskPattern(t *testing.T) {
	pattern, err := CompileIgnorePattern("*")
	if err != nil {
		t.Error(err)
	} else if pattern == nil {
		t.Error("expected pattern not to be nil")
	}
	if pattern.negate {
		t.Error("expected pattern should not be negated")
	}
	if !pattern.Matches("it doesn't matter what I put here because it should match") {
		t.Error("expected string to match pattern")
	}
}

func TestFileGlobPattern(t *testing.T) {
	pattern, err := CompileIgnorePattern("*.js")
	if err != nil {
		t.Error(err)
	} else if pattern == nil {
		t.Error("expected pattern not to be nil")
	}
	if pattern.negate {
		t.Error("expected pattern should not be negated")
	}
	if pattern.Matches("it doesn't matter what I put here because it shouldn't match") {
		t.Error("expected string not to match pattern")
	}
	if !pattern.Matches("horse/cat/dog.js") {
		t.Error("expected string to match pattern")
	}
}

func TestFilePathGlobPattern(t *testing.T) {
	pattern, err := CompileIgnorePattern("damn/*.js")
	if err != nil {
		t.Error(err)
	} else if pattern == nil {
		t.Error("expected pattern not to be nil")
	}
	if pattern.negate {
		t.Error("expected pattern should not be negated")
	}
	if pattern.Matches("it doesn't matter what I put here because it shouldn't match") {
		t.Error("expected string not to match pattern")
	}
	if pattern.Matches("dang/dog.js") {
		t.Error("expected string not to match pattern")
	}
	if pattern.Matches("horse/damn/dog.js") {
		t.Error("expected string not to match pattern")
	}
	if !pattern.Matches("damn/dog.js") {
		t.Error("expected string to match pattern")
	}
}

func TestTrailingGlobPattern(t *testing.T) {
	pattern, err := CompileIgnorePattern("docs/*")
	if err != nil {
		t.Error(err)
	} else if pattern == nil {
		t.Error("expected pattern not to be nil")
	}
	if pattern.negate {
		t.Error("expected pattern should not be negated")
	}
	if pattern.Matches("it doesn't matter what I put here because it shouldn't match") {
		t.Error("expected string not to match pattern")
	}
	if pattern.Matches("rock/docs") {
		t.Error("expected string not to match pattern")
	}
	if !pattern.Matches("docs/dog.js") {
		t.Error("expected string to match pattern")
	}
	if !pattern.Matches("rock/docs/dog.js") {
		t.Error("expected string to match pattern")
	}
}

func TestLeadingSlashPattern(t *testing.T) {
	pattern, err := CompileIgnorePattern("/build/logs/")
	if err != nil {
		t.Error(err)
	} else if pattern == nil {
		t.Error("expected pattern not to be nil")
	}
	if pattern.negate {
		t.Error("expected pattern should not be negated")
	}
	if pattern.Matches("it doesn't matter what I put here because it shouldn't match") {
		t.Error("expected string not to match pattern")
	}
	if pattern.Matches("horse/build/logs/dog.js") {
		t.Error("expected string not to match pattern")
	}
	if !pattern.Matches("build/logs/dog.js") {
		t.Error("expected string to match pattern")
	}
	if !pattern.Matches("/build/logs/dog.js") {
		t.Error("expected string to match pattern")
	}
}

func TestNegationPattern(t *testing.T) {
	pattern, err := CompileIgnorePattern("!*.js")
	if err != nil {
		t.Error(err)
	} else if pattern == nil {
		t.Error("expected pattern not to be nil")
	}
	if !pattern.negate {
		t.Error("expected pattern should be negated")
	}
	if pattern.Matches("build/logs/dog.js") {
		t.Error("expected string to not match pattern")
	}
	if !pattern.Matches("build/logs/dog.css") {
		t.Error("expected string to match pattern")
	}
}
