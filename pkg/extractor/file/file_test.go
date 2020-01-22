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

package file

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/stretchr/testify/assert"

	"os"
	"testing"

	"github.com/codenotary/objects/pkg/object"
	"github.com/codenotary/objects/pkg/uri"
)

func TestFile(t *testing.T) {
	file, err := ioutil.TempFile("", "codenotary-test-scheme-file")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())
	err = ioutil.WriteFile(file.Name(), []byte("123\n"), 0644)
	if err != nil {
		log.Fatal(err)
	}
	u, _ := uri.Parse("file://" + file.Name())

	o, err := Extract(*u)
	assert.NoError(t, err)
	assert.NotNil(t, o)
	assert.Equal(t, *u, o.URI)
	assert.Equal(t, filepath.Base(file.Name()), o.Metadata["name"])
	assert.Equal(t, "sha256:181210f8f9c779c26da1d9b2075bde0127302ee0e3fca38c9a83f5b1dd8e5d3b", o.Digest.String())

	u, _ = uri.Parse(file.Name())
	o, err = Extract(*u)
	assert.NoError(t, err)
	assert.NotNil(t, o)
	u.Scheme = "file"
	assert.Equal(t, *u, o.URI)
	assert.Equal(t, filepath.Base(file.Name()), o.Metadata["name"])
	assert.Equal(t, "sha256:181210f8f9c779c26da1d9b2075bde0127302ee0e3fca38c9a83f5b1dd8e5d3b", o.Digest.String())

	u, _ = uri.Parse("file_test.go")
	o, err = Extract(*u)
	assert.NoError(t, err)
	assert.NotNil(t, o)
	assert.Equal(t, uri.URI{Scheme: "file", Opaque: "file_test.go"}, o.URI)
	assert.Equal(t, o.ContentType, "text/plain; charset=utf-8")
	assert.Equal(t, o.Metadata, object.Metadata{
		"name": "file_test.go",
	})
}
