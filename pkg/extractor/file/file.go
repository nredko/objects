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

package file

import (
	"crypto/sha256"
	"io"
	"os"
	"strings"

	"github.com/codenotary/objects/pkg/extractor"
	"github.com/codenotary/objects/pkg/object"
	"github.com/codenotary/objects/pkg/uri"

	digest "github.com/opencontainers/go-digest"
)

// Scheme for file
const Scheme = "file"

var _ extractor.Extractor = Extract

// Extract returns a file *object.Object from a given u
func Extract(u uri.URI, options ...extractor.Option) (*object.Object, error) {

	if u.Scheme != "" && u.Scheme != Scheme {
		return nil, nil
	}

	path := strings.TrimPrefix(u.Opaque, "//")

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Hash
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return nil, err
	}

	// Name and Size
	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	// ContentType
	ct, err := contentType(f)
	if err != nil {
		return nil, err
	}

	// Metadata container
	m := object.Metadata{
		"name": stat.Name(),
	}

	// Infer version from filename
	if version := inferVer(stat.Name()); version != "" {
		m["version"] = version
	}

	// Sniff executable info, if any
	if ok, data, _ := xInfo(f, &ct); ok {
		m.SetValues(data)
	}

	u.Scheme = Scheme
	return &object.Object{
		Descriptor: object.Descriptor{
			Digest:      digest.NewDigest(digest.SHA256, h),
			Size:        uint64(stat.Size()),
			ContentType: ct,
		},
		Metadata: m,
		URI:      u,
	}, nil
}
