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

package codeowners

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

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
type Codeowners []Entry

// Entry contains owners for a given pattern
type Entry struct {
	LineNo  uint64
	Pattern Pattern
	Owners  []string
}

func (e Entry) String() string {
	return fmt.Sprintf("line %d: %s\t%v", e.LineNo, e.Pattern.String(), strings.Join(e.Owners, ", "))
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

	return nil, fmt.Errorf("No CODEOWNERS found in the root, docs/, or .github/ directory of the repository")
}

func parseCodeowners(r io.Reader) ([]Entry, error) {
	var e []Entry
	s := bufio.NewScanner(r)
	no := uint64(0)
	for s.Scan() {
		no++
		fields := strings.Fields(s.Text())

		if len(fields) == 0 { // empty
			continue
		}

		if strings.HasPrefix(fields[0], "#") { // comment
			continue
		}
		pattern, err := CompilePattern(fields[0])
		if err != nil {
			return nil, err
		}
		owners := fields[1:]

		e = append(e, Entry{
			Pattern: *pattern,
			Owners:  owners,
			LineNo:  no,
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
