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

package object

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetadata(t *testing.T) {
	m := Metadata{}

	// Set/Get
	m.Set("key", "value")
	assert.Equal(t, "value", m.Get("key", nil))
	assert.Equal(t, "default", m.Get("nonExistingKey", "default"))

	// Multiple values
	m.SetValues(map[string]interface{}{"key": "newValue", "a": "one", "b": 2})
	assert.Equal(t, "newValue", m.Get("key", nil))
	assert.Equal(t, "one", m.Get("a", nil))
	assert.Equal(t, 2, m.Get("b", nil))
}
