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
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/codenotary/bundle"
	"github.com/codenotary/objects/pkg/extractor"
	"github.com/codenotary/objects/pkg/object"
	"github.com/codenotary/objects/pkg/uri"
)

// Scheme for dir
const Scheme = "dir"

// ManifestKey is the metadata's key for storing the manifest
const ManifestKey = "manifest"

// PathKey is the metadata's key for the directory path
const PathKey = "path"

type opts struct {
	initIgnoreFile bool
}

var _ extractor.Extractor = Extract

// Extract returns a file *api.Extract from a given u
func Extract(u uri.URI, options ...extractor.Option) (*object.Object, error) {

	if u.Scheme != "" && u.Scheme != Scheme {
		return nil, nil
	}

	opts := &opts{}
	if err := extractor.Options(options).Apply(opts); err != nil {
		return nil, err
	}

	path := strings.TrimPrefix(u.Opaque, "//")
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// get file info and check if is a directory
	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if !stat.IsDir() {
		return nil, fmt.Errorf("read %s: is not a directory", path)
	}

	if opts.initIgnoreFile {
		if err := initIgnoreFile(path); err != nil {
			return nil, err
		}
	}

	files, err := walk(path)
	if err != nil {
		return nil, err
	}

	manifest := bundle.NewManifest(files...)
	digest, err := manifest.Digest()
	if err != nil {
		return nil, err
	}

	return &object.Object{
		Descriptor: object.Descriptor{
			Digest: digest,
			Size:   uint64(stat.Size()),
		},
		URI:      uri.URI{Scheme: Scheme, Opaque: u.Opaque},
		Metadata: object.Metadata{ManifestKey: files, PathKey: path},
	}, nil
}

// WithIgnoreFileInit returns a functional option to instruct the dir's extractor to create the defualt ignore file
// when not yet present into the targeted directory.
func WithIgnoreFileInit() extractor.Option {
	return func(o interface{}) error {
		if o, ok := o.(*opts); ok {
			o.initIgnoreFile = true
		}
		return nil
	}
}
