/*
Copyright 2019-2020 vChain, Inc.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
	http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package dir

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/src-d/go-git.v4/plumbing/format/gitignore"
)

// IgnoreFilename is the name of the ignore file used by this package when processing directories.
//
// The ignore file supports all pattern formats as specified in the gitignore specification:
//  - https://git-scm.com/docs/gitignore
//  - https://github.com/src-d/go-git/blob/master/plumbing/format/gitignore/doc.go
//
// However, this package implementation:
//  - only uses the ignore file within the root directory (nested ignore files will be treated as normal files)
//  - always ignores the manifest file (it cannot be excluded by the ignore file)
//
const IgnoreFilename = ".vcnignore"

// DefaultIgnoreFileContent is the content of ignore file with default patterns.
const DefaultIgnoreFileContent = `# Windows thumbnail cache files
Thumbs.db
Thumbs.db:encryptable
ehthumbs.db
ehthumbs_vista.db

# Windows folder config file
[Dd]esktop.ini

# Windows Recycle Bin used on file shares
$RECYCLE.BIN/

# macOS
.DS_Store
.AppleDouble
.LSOverride

# macOS Thumbnails
._*

# macOS files that might appear in the root of a volume
.DocumentRevisions-V100
.fseventsd
.Spotlight-V100
.TemporaryItems
.Trashes
.VolumeIcon.icns
.com.apple.timemachine.donotpresent

# Directories potentially created on remote AFP share
.AppleDB
.AppleDesktop
Network Trash Folder
Temporary Items
.apdisk

# temporary files which can be created if a process still has a handle open of a deleted file
.fuse_hidden*

# KDE directory preferences
.directory

# Linux trash folder which might appear on any partition or disk
.Trash-*

# .nfs files are created when an open file is removed but is still being accessed
.nfs*
`

const (
	ignorefileCommentPrefix = "#"
	ignorefileEOL           = "\n"
)

type nullMatcher struct{}

func (n *nullMatcher) Match(path []string, isDir bool) bool {
	return false
}

// newIgnoreFileMatcher reads and parses the ignore file in path and return a gitignore.Matcher.
// If the ignore file does not exists, a matcher that will never match is returned instead.
func newIgnoreFileMatcher(path string) (m gitignore.Matcher, err error) {
	f, err := os.Open(filepath.Join(path, IgnoreFilename))
	if err != nil {
		if os.IsNotExist(err) {
			return &nullMatcher{}, nil
		}
		return
	}
	defer f.Close()
	ps := []gitignore.Pattern{}
	if data, err := ioutil.ReadAll(f); err == nil {
		for _, s := range strings.Split(string(data), ignorefileEOL) {
			if !strings.HasPrefix(s, ignorefileCommentPrefix) && len(strings.TrimSpace(s)) > 0 {
				ps = append(ps, gitignore.ParsePattern(s, nil))
			}
		}
	}
	m = gitignore.NewMatcher(ps)
	return
}

// initIgnoreFile writes the default ignore file if it does not exist.
func initIgnoreFile(root string) error {
	filename := filepath.Join(root, IgnoreFilename)

	// create and open the file if not exists
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		if os.IsExist(err) {
			return nil // file exists already
		}
		return err
	}

	// otherwise, write the default content
	_, err = f.WriteString(DefaultIgnoreFileContent)
	return err
}
