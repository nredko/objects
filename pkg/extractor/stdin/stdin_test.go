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
	"io/ioutil"
	"log"

	"github.com/stretchr/testify/assert"

	"os"
	"testing"

	"github.com/codenotary/objects/pkg/uri"
)

func TestStdin(t *testing.T) {
	file, err := ioutil.TempFile("", "codenotary-test-scheme-stdin")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())
	err = ioutil.WriteFile(file.Name(), []byte("123\n"), 0644)
	if err != nil {
		log.Fatal(err)
	}

	u, _ := uri.Parse("stdin://")

	myStdin := os.Stdin
	defer func() { os.Stdin = myStdin }()
	os.Stdin = file

	o, err := Extract(*u)

	stat, _ := file.Stat()
	size := stat.Size()

	if err := file.Close(); err != nil {
		log.Fatal(err)
	}

	assert.NoError(t, err)
	assert.NotNil(t, o)
	assert.Equal(t, *&u.Scheme, o.URI.Scheme)
	assert.Equal(t, "sha256:181210f8f9c779c26da1d9b2075bde0127302ee0e3fca38c9a83f5b1dd8e5d3b", o.Digest.String())
	assert.Equal(t, uri.URI{Scheme: "stdin", Opaque: ""}, o.URI)
	assert.Equal(t, "text/octet-stream", o.ContentType)
	assert.Equal(t, size, int64(o.Size))
}
