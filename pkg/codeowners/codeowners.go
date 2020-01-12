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

	"gopkg.in/src-d/go-git.v4"
)

// Codeowners is the deserialized form of a given CODEOWNERS file
type Codeowners []Entry

// Entry contains owners for a given pattern
type Entry struct {
	LineNo  uint64
	Pattern string
	Owners  []string
}

func (e Entry) String() string {
	return fmt.Sprintf("line %d: %s\t%v", e.LineNo, e.Pattern, strings.Join(e.Owners, ", "))
}

// NewFromTree loads and deserializes a CODEOWNERS file from the given repository, if one exists
func NewFromTree(worktree *git.Worktree) (Codeowners, error) {
	r, err := openCodeownersFile(worktree)
	if err != nil {
		return nil, err
	}

	entries := parseCodeowners(r)
	return entries, nil
}

// openCodeownersFile finds a CODEOWNERS file and returns content.
// see: https://help.github.com/articles/about-code-owners/#codeowners-file-location
func openCodeownersFile(worktree *git.Worktree) (io.Reader, error) {
	fs := worktree.Filesystem
	for _, p := range []string{".", "docs", ".github"} {
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

func parseCodeowners(r io.Reader) []Entry {
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

		e = append(e, Entry{
			Pattern: fields[0],
			Owners:  fields[1:],
			LineNo:  no,
		})
	}

	return e
}
