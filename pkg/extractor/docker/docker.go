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
	"context"
	"fmt"
	"strings"

	"github.com/codenotary/objects/pkg/extractor"
	"github.com/codenotary/objects/pkg/object"
	"github.com/codenotary/objects/pkg/uri"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	digest "github.com/opencontainers/go-digest"
)

// Scheme for docker
const Scheme = "docker"

// Extract returns a file *object.Object from a given u
func Extract(u uri.URI, options ...extractor.Option) (*object.Object, error) {

	if u.Scheme != Scheme {
		return nil, nil
	}

	dockerClient, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}
	imageID := strings.TrimPrefix(u.Opaque, "//")
	imgInspect, _, err := dockerClient.ImageInspectWithRaw(context.Background(), imageID)

	if err != nil {
		return nil, fmt.Errorf("failed to inspect docker image: %s", err)
	}

	m := object.Metadata{
		"name":         getName(imgInspect),
		"architecture": imgInspect.Architecture,
		"platform":     imgInspect.Os,
		Scheme:         imgInspect,
	}

	if version := inferVer(imgInspect); version != "" {
		m["version"] = version
	}

	return &object.Object{
		Descriptor: object.Descriptor{
			Digest: digest.FromString(imgInspect.ID),
			Size:   uint64(imgInspect.Size),
		},
		Metadata: m,
		URI:      u,
	}, nil
}

func getName(i types.ImageInspect) string {
	if len(i.RepoTags) > 0 {
		return i.RepoTags[0]
	}
	return strings.TrimSpace(i.ID)
}

func inferVer(i types.ImageInspect) string {
	if len(i.RepoTags) > 0 {
		parts := strings.SplitN(i.RepoTags[0], ":", 2)
		if len(parts) > 1 && parts[1] != "latest" {
			return parts[1]
		}
	}

	return ""
}
