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

package uri

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURI(t *testing.T) {
	u, err := Parse("scheme://opaque")
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, URI{Scheme: "scheme", Opaque: "//opaque"}, *u)
	assert.Equal(t, "scheme://opaque", u.String())

	u, err = Parse("file.txt")
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, URI{Scheme: "", Opaque: "file.txt"}, *u)
	assert.Equal(t, "file.txt", u.String())
}
