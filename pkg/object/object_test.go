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

package object

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/codenotary/objects/pkg/uri"

	digest "github.com/opencontainers/go-digest"
)

func TestObjectMarshalling(t *testing.T) {

	testData := []byte(
		`{"descriptor":{"digest":"sha256:6a6c10b3cb05a670d3ff7af93f45df97c051af60c230939db98be2a977aaddd1","size":6,"contentType":"text/plain"},"metadata":{"key1":"value1","key2":"value2"},"uri":"file://info.txt"}`,
	)

	o := Object{
		Descriptor: Descriptor{
			Digest:      digest.FromString("sha256:b691886c974b911d9f6ccb8ee3c330ccaf4b1317b109212a3d37a49e11248e4d"),
			Size:        6,
			ContentType: "text/plain",
		},
		Metadata: Metadata(map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
		}),
		URI: uri.URI{
			Scheme: "file",
			Opaque: "//info.txt",
		},
	}

	b, err := json.Marshal(o)

	assert.NoError(t, err)
	assert.Equal(t, testData, b)

}
