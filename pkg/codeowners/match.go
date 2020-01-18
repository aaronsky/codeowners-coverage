package codeowners

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	leadingHashOrBangPattern         = regexp.MustCompile(`^(\#|\!)`)
	prependSlashForFileGlobPattern   = regexp.MustCompile(`([^\/+])/.*\*\.`)
	escapingDotSlashPattern          = regexp.MustCompile(`\.`)
	innerDoubleAsteriskPattern       = regexp.MustCompile(`/\*\*/`)
	leadingDoubleAsteriskPattern     = regexp.MustCompile(`\*\*/`)
	trailingDoubleAsteriskPattern    = regexp.MustCompile(`/\*\*`)
	veryEscapedSingleAsteriskPattern = regexp.MustCompile(`\\\*`)
	escapedSingleAsteriskPattern     = regexp.MustCompile(`\*`)
)

// Pattern aliases string to add some additional documentation
type Pattern struct {
	pattern *regexp.Regexp
	negate  bool
}

// NewPattern creates a new Pattern object behind a pointer
func NewPattern(regex *regexp.Regexp, negate bool) *Pattern {
	return &Pattern{pattern: regex, negate: negate}
}

func (p *Pattern) String() string {
	negatedString := ""
	if p.negate {
		negatedString += "not "
	}
	return fmt.Sprintf("%s%s", negatedString, p.pattern)
}

// Matches returns whether or not a given path matches the Codeowners path pattern
func (p *Pattern) Matches(path string) bool {
	path = strings.Replace(path, string(os.PathSeparator), "/", -1)
	if p.negate {
		return !p.pattern.MatchString(path)
	}
	return p.pattern.MatchString(path)
}

// CompilePattern takes a given ignore pattern and attempts to create a Pattern object from it
func CompilePattern(pattern string) (*Pattern, error) {
	// Trim OS-specific carriage returns.
	pattern = strings.TrimRight(pattern, "\r")

	// A line starting with # serves as a comment
	if strings.HasPrefix(pattern, `#`) {
		return nil, fmt.Errorf("intentionally not compiling comment pattern")
	}

	// Trailing spaces are ignored unless they are quoted with backslash ("\").
	pattern = strings.Trim(pattern, " ")

	// A blank line matches no files, so it can serve as a separator for readability.
	if pattern == "" {
		return nil, fmt.Errorf("intentionally not compiling empty pattern")
	}

	// An optional prefix "!" which negates the pattern; any matching file excluded by a previous
	// pattern will become included again. It is not possible to re-include a file if a parent
	// directory of that file is excluded. Git doesn’t list excluded directories for performance
	// reasons, so any patterns on contained files have no effect, no matter where they are defined.
	// Put a backslash ("\") in front of the first "!" for patterns that begin with a literal "!",
	// for example, "\!important!.txt".
	negatePattern := false
	if pattern[0] == '!' {
		negatePattern = true
		pattern = pattern[1:]
	}

	// Trailing spaces are ignored unless they are quoted with backslash ("\").
	// Put a backslash ("\") in front of the first hash for patterns that begin with a hash.
	if leadingHashOrBangPattern.MatchString(pattern) {
		pattern = pattern[1:]
	}

	// The slash / is used as the directory separator. Separators may occur at the beginning,
	// middle or end of the .gitignore search pattern.
	if prependSlashForFileGlobPattern.MatchString(pattern) && pattern[0] != '/' {
		pattern = "/" + pattern
	}

	// Handle escaping the "." char
	pattern = escapingDotSlashPattern.ReplaceAllString(pattern, `\.`)

	magicStar := "#$~"
	pattern = handleConsecutiveAsterisks(pattern, magicStar)

	// Handle escaping the "?" char
	pattern = strings.Replace(pattern, "?", `\?`, -1)
	// An asterisk "*" matches anything except a slash. The character "?" matches any one character
	// except "/". The range notation, e.g. [a-zA-Z], can be used to match one of the characters in
	// a range. See fnmatch(3) and the FNM_PATHNAME flag for a more detailed description.
	pattern = strings.Replace(pattern, magicStar, "*", -1)

	var expr = ""
	// If there is a separator at the end of the pattern then the pattern will only match
	// directories, otherwise the pattern can match both files and directories.
	if strings.HasSuffix(pattern, "/") {
		expr = pattern + "(|.*)$"
	} else {
		expr = pattern + "(|/.*)$"
	}
	// If there is a separator at the beginning or middle (or both) of the pattern, then the
	// pattern is relative to the directory level of the particular .gitignore file itself.
	// Otherwise the pattern may also match at any level below the .gitignore level.
	if strings.HasPrefix(expr, "/") {
		expr = "^(|/)" + expr[1:]
	} else {
		expr = "^(|.*/)" + expr
	}

	regex, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}

	return NewPattern(regex, negatePattern), nil
}

func handleConsecutiveAsterisks(pattern, magicStar string) string {
	// Two consecutive asterisks ("**") in patterns matched against full pathname may have special meaning

	// Handle "/**/" usage
	if strings.HasPrefix(pattern, "/**/") {
		pattern = pattern[1:]
	}
	// A slash followed by two consecutive asterisks then a slash matches zero or more directories.
	// For example, "a/**\/b" matches "a/b", "a/x/b", "a/x/y/b" and so on.
	pattern = innerDoubleAsteriskPattern.ReplaceAllString(pattern, `(/|/.+/)`)
	// A leading "**" followed by a slash means match in all directories. For example, "**\/foo"
	// matches file or directory "foo" anywhere, the same as pattern "foo". "**\/foo/bar" matches
	// file or directory "bar" anywhere that is directly under directory "foo".
	pattern = leadingDoubleAsteriskPattern.ReplaceAllString(pattern, `(|.`+magicStar+`/)`)
	// A trailing "/**" matches everything inside. For example, "abc/**" matches all files inside directory "abc", relative to the location of the .gitignore file, with infinite depth.
	pattern = trailingDoubleAsteriskPattern.ReplaceAllString(pattern, `(|/.`+magicStar+`)`)
	// Other consecutive asterisks are considered regular asterisks and will match according to the previous rules.

	// Handle escaping the "*" char
	pattern = veryEscapedSingleAsteriskPattern.ReplaceAllString(pattern, `\`+magicStar)
	pattern = escapedSingleAsteriskPattern.ReplaceAllString(pattern, `([^/]*)`)

	return pattern
}
