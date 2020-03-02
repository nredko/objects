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

package stdin

import (
	"bufio"
	"crypto/sha256"
	"io"
	"os"

	"github.com/codenotary/objects/pkg/extractor"
	"github.com/codenotary/objects/pkg/object"
	"github.com/codenotary/objects/pkg/uri"

	digest "github.com/opencontainers/go-digest"
)

// Scheme for stdin
const Scheme = "stdin"

var _ extractor.Extractor = Extract

// Extract returns a file *object.Object from a given u
func Extract(u uri.URI, options ...extractor.Option) (*object.Object, error) {

	if u.Scheme != Scheme {
		return nil, nil
	}

	reader := bufio.NewReader(os.Stdin)

	// Hash
	h := sha256.New()

	if _, err := io.Copy(h, reader); err != nil {
		return nil, err
	}

	// Name and Size
	stat, err := os.Stdin.Stat()
	if err != nil {
		return nil, err
	}

	return &object.Object{
		Descriptor: object.Descriptor{
			Digest:      digest.NewDigest(digest.SHA256, h),
			Size:        uint64(stat.Size()),
			ContentType: "text/octet-stream",
		},
		Metadata: object.Metadata{},
		URI:      uri.URI{Scheme: u.Scheme},
	}, nil
}
