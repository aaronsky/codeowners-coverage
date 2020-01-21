// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package codeowners contains logic for loading and parsing patterns in CODEOWNERS files
package codeowners

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aaronsky/codeowners-coverage/internal/git"
	"gopkg.in/src-d/go-billy.v4"
)

var codeownersDirectories = []string{".", "docs", ".github"}

// PathIsCodeowners returns whether or not the provided path is for a valid CODEOWNERS file
// see: https://help.github.com/articles/about-code-owners/#codeowners-file-location
func PathIsCodeowners(path string, fs billy.Filesystem) bool {
	for _, dir := range codeownersDirectories {
		if path == fs.Join(dir, "CODEOWNERS") {
			return true
		}
	}
	return false
}

// Codeowners is the deserialized form of a given CODEOWNERS file
type Codeowners []OwnerEntry

// OwnerEntry contains owners for a given pattern
type OwnerEntry struct {
	lineNumber uint64
	Pattern    git.IgnorePattern
	Owners     []string
}

func (e OwnerEntry) String() string {
	return fmt.Sprintf("line %d: %s\t%v", e.lineNumber, e.Pattern.String(), strings.Join(e.Owners, ", "))
}

// LoadFromFilesystem loads and deserializes a CODEOWNERS file from the given repository, if one exists
func LoadFromFilesystem(fs billy.Filesystem) (Codeowners, error) {
	r, err := openCodeownersFile(fs)
	if err != nil {
		return nil, err
	}

	return parseCodeowners(r)
}

// openCodeownersFile finds a CODEOWNERS file and returns content.
// see: https://help.github.com/articles/about-code-owners/#codeowners-file-location
func openCodeownersFile(fs billy.Filesystem) (io.Reader, error) {
	for _, p := range codeownersDirectories {
		path := fs.Join(p)
		if _, err := fs.Stat(path); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}

		file := fs.Join(path, "CODEOWNERS")
		_, err := fs.Stat(file)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}

		return fs.Open(file)
	}

	return nil, fmt.Errorf("no CODEOWNERS found in the root, docs/, or .github/ directory of the repository")
}

func parseCodeowners(r io.Reader) ([]OwnerEntry, error) {
	var e []OwnerEntry
	s := bufio.NewScanner(r)
	var lineNumber uint64
	for s.Scan() {
		lineNumber++
		fields := strings.Fields(s.Text())

		if len(fields) == 0 { // empty
			continue
		}

		if strings.HasPrefix(fields[0], "#") { // comment
			continue
		}
		pattern, err := git.CompileIgnorePattern(fields[0])
		if err != nil {
			return nil, err
		}
		owners := fields[1:]

		e = append(e, OwnerEntry{
			lineNumber: lineNumber,
			Pattern:    *pattern,
			Owners:     owners,
		})
	}

	return e, nil
}

// Owners returns the list of owners for a given path, in the event of a match
func (o *Codeowners) Owners(path string) []string {
	owners := []string{}
	if o == nil {
		return owners
	}
	for _, entry := range *o {
		if entry.Pattern.Matches(path) {
			owners = entry.Owners
		}
	}
	return owners
}
