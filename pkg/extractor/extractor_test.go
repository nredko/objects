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

package extractor

import (
	"testing"

	"github.com/codenotary/objects/pkg/object"
	"github.com/codenotary/objects/pkg/uri"
	"github.com/stretchr/testify/assert"
)

func TestExtract(t *testing.T) {

	dummyExtractor := func(uri.URI, ...Option) (*object.Object, error) {
		return &object.Object{}, nil
	}

	Register("foo", dummyExtractor)

	o, err := Extract("foo://bar")
	assert.NotNil(t, o)
	assert.NoError(t, err)

	o, err = Extract("missing://something")
	assert.Nil(t, o)
	assert.Error(t, err)

	o, err = Extract("noschemenofallback")
	assert.Nil(t, o)
	assert.Error(t, err)

	SetFallbackScheme("foo")

	o, err = Extract("noschemewithfallback")
	assert.NotNil(t, o)
	assert.NoError(t, err)
}
