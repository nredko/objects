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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInferVer(t *testing.T) {
	testCases := map[string]string{
		// Supported
		"cn-v0.4.0-darwin-10.6-amd64":     "0.4.0",
		"cn-v0.4.0-linux-amd64":           "0.4.0",
		"cn-v0.4.0-windows-4.0-amd64.exe": "0.4.0",

		// Unsupported
		"codenotary_cn_0.4.0_setup.exe": "",
	}

	for filename, ver := range testCases {
		assert.Equal(t, ver, inferVer(filename), "wrong version for %s", filename)
	}

}
