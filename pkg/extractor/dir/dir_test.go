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
	"path/filepath"
	"testing"

	"github.com/codenotary/objects/pkg/uri"
	"github.com/stretchr/testify/assert"
)

func TestArtifact(t *testing.T) {

	tmpDir, err := ioutil.TempDir("", "/TempDir")
	if err != nil {
		t.Fatal(err)
	}

	tmpFile := filepath.Join(tmpDir, "file")
	err = ioutil.WriteFile(tmpFile, nil, 0644)
	if err != nil {
		t.Fatal(err)
	}
	// dir - OK
	u, _ := uri.Parse("dir://" + tmpDir)
	a, err := Extract(*u)
	assert.NoError(t, err)
	assert.NotNil(t, a)
	assert.Equal(t, "dir", a.URI.Scheme)
	assert.Equal(t, filepath.Base(tmpDir), filepath.Base(a.URI.Opaque))
	assert.NotEmpty(t, a.Descriptor.Digest)

	// dir (no schema) - OK
	u, _ = uri.Parse(tmpDir)
	a, err = Extract(*u)
	assert.NoError(t, err)
	assert.NotNil(t, a)
	assert.Equal(t, "dir", a.URI.Scheme)
	assert.Equal(t, filepath.Base(tmpDir), filepath.Base(a.URI.Opaque))
	assert.NotEmpty(t, a.Descriptor.Digest)

	// wrong schema - SKIP (no error)
	u, _ = uri.Parse("file://" + tmpDir)
	a, err = Extract(*u)
	assert.NoError(t, err)
	assert.Nil(t, a)

	// not a dir - ERROR
	u, _ = uri.Parse("dir://" + tmpFile)
	a, err = Extract(*u)
	assert.Error(t, err)
	assert.Nil(t, a)

	// not existing dir - ERROR
	u, _ = uri.Parse("dir://" + tmpDir + "/not-existing")
	a, err = Extract(*u)
	assert.Error(t, err)
	assert.Nil(t, a)
}
