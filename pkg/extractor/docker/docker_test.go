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

package docker

import (
	"os/exec"
	"testing"

	"github.com/codenotary/objects/pkg/uri"

	"github.com/stretchr/testify/assert"
)

func TestDocker(t *testing.T) {
	_, err := exec.Command("docker", "pull", "hello-world").Output()
	if err != nil {
		t.Skip("docker not available")
	}

	u, _ := uri.Parse("docker://hello-world")
	o, err := Extract(*u)
	assert.NoError(t, err)
	assert.NotNil(t, o)
	assert.Equal(t, "hello-world:latest", o.Metadata["name"])
	assert.NoError(t, o.Digest.Validate())
	assert.NotZero(t, o.Size)
}
