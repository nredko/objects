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
	"os"
	"path/filepath"
	"strings"

	"github.com/codenotary/bundle"
)

func walk(root string) (files []bundle.Descriptor, err error) {
	files = make([]bundle.Descriptor, 0)
	ignore, err := newIgnoreFileMatcher(root)
	if err != nil {
		return
	}
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		// skip irregular files (e.g. dir, symlink, pipe, socket, device...)
		if !info.Mode().IsRegular() {
			return nil
		}

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		// descriptor's path must be OS agnostic
		relPath = filepath.ToSlash(relPath)

		// skip files matching the ignore patterns
		if ignore.Match(strings.Split(relPath, "/"), false) {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		d, err := bundle.NewDescriptor(relPath, file)
		file.Close()
		if err != nil {
			return err
		}
		files = append(files, *d)

		return nil
	})
	return
}
