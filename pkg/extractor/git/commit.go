/*
Copyright 2019 vChain, Inc.
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
	"crypto/sha256"
	"io"

	dgst "github.com/opencontainers/go-digest"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func lastCommit(repo *git.Repository) (*object.Commit, error) {
	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}

	return repo.CommitObject(ref.Hash())
}

func digestCommit(c object.Commit) (digest dgst.Digest, size uint64, err error) {
	o := &plumbing.MemoryObject{}
	c.Encode(o)

	reader, err := o.Reader()
	if err != nil {
		return
	}
	defer reader.Close()

	h := sha256.New()
	n, err := io.Copy(h, reader)
	if err != nil {
		return
	}
	return dgst.NewDigest(dgst.SHA256, h), uint64(n), nil
}
