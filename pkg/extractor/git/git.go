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

package git

import (
	"path/filepath"
	"strings"

	git "gopkg.in/src-d/go-git.v4"

	"github.com/codenotary/objects/pkg/extractor"
	"github.com/codenotary/objects/pkg/object"
	"github.com/codenotary/objects/pkg/uri"
)

// Scheme for git
const Scheme = "git"

var _ extractor.Extractor = Extract

// Extract returns a git *object.Object from a given u
func Extract(u uri.URI, options ...extractor.Option) (*object.Object, error) {

	if u.Scheme != Scheme {
		return nil, nil
	}

	path := strings.TrimPrefix(u.Opaque, "//")
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	commit, err := lastCommit(repo)
	if err != nil {
		return nil, err
	}

	d, size, err := digestCommit(*commit)
	if err != nil {
		return nil, err
	}

	// Metadata container
	m := object.Metadata{
		Scheme: map[string]interface{}{
			"Commit": commit.Hash.String(),
			"Tree":   commit.TreeHash.String(),
			"Parents": func() []string {
				res := make([]string, len(commit.ParentHashes))
				for i, h := range commit.ParentHashes {
					res[i] = h.String()
				}
				return res
			}(),
			"Author":       commit.Author,
			"Committer":    commit.Committer,
			"Message":      commit.Message,
			"PGPSignature": commit.PGPSignature,
		},
	}

	name := filepath.Base(path)
	if remotes, err := repo.Remotes(); err == nil && len(remotes) > 0 {
		urls := remotes[0].Config().URLs
		if len(urls) > 0 {
			name = urls[0]
		}
	}
	name += "@" + commit.Hash.String()

	m["name"] = name

	return &object.Object{
		Descriptor: object.Descriptor{
			Digest: d,
			Size:   size,
		},
		Metadata: m,
		URI:      u,
	}, nil
}
