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
	"io"
	"sort"

	// See https://github.com/opencontainers/go-digest#usage
	_ "crypto/sha256"
	_ "crypto/sha512"

	digest "github.com/opencontainers/go-digest"
)

// Descriptor describes the disposition of targeted content.
type Descriptor struct {
	// Digest is the digest of the targeted content.
	Digest digest.Digest `json:"digest"`

	// Size specifies the size in bytes of the targeted content.
	Size uint64 `json:"size"`

	// Paths specifies the relative locations of the targeted content.
	Paths []string `json:"paths"`
}

func (d *Descriptor) sortUnique() {
	tmp := make(map[string]bool, len(d.Paths))
	for _, p := range d.Paths {
		tmp[p] = true
	}
	d.Paths = make([]string, len(tmp))
	i := 0
	for p := range tmp {
		d.Paths[i] = p
		i++
	}
	sort.Strings(d.Paths)
}

// NewDescriptor returns a new *Descriptor for the provided path and src.
func NewDescriptor(path string, src io.Reader) (*Descriptor, error) {
	digester := digest.SHA256.Digester()
	size, err := io.Copy(digester.Hash(), src)
	if err != nil {
		return nil, err
	}

	return &Descriptor{
		Paths:  []string{path},
		Digest: digester.Digest(),
		Size:   uint64(size),
	}, nil
}
