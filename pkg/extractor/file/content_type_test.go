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
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContenType(t *testing.T) {
	emptyFile, err := ioutil.TempFile("", "TestContenType")
	if err != nil {
		t.Fatal(err)
	}
	txtFile, err := ioutil.TempFile("", "TestContenType")
	if err != nil {
		t.Fatal(err)
	}
	txtFile.Write([]byte{99, 105, 97, 111})

	ct, err := contentType(emptyFile)
	assert.NoError(t, err)
	assert.Empty(t, ct)

	ct, err = contentType(txtFile)
	assert.NoError(t, err)
	assert.Equal(t, "application/octet-stream", ct)
}
