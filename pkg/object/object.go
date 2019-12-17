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
	"github.com/codenotary/objects/pkg/uri"

	// See https://github.com/opencontainers/go-digest#usage
	_ "crypto/sha256"
	_ "crypto/sha512"

	digest "github.com/opencontainers/go-digest"
)

// Descriptor describes the disposition of targeted content.
type Descriptor struct {
	// Digest is the digest of the targeted content.
	Digest digest.Digest `json:"digest"`

	// Size specifies the size in bytes of the targeted content.
	Size uint64 `json:"size"`

	// ContentType specifies the media type string of the targeted content.
	ContentType string `json:"contentType"`
}

// Object represents the set of all relevant information gathered from a digital object referenced by URI.
type Object struct {
	Descriptor
	Metadata `json:"metadata"`
	uri.URI  `json:"uri"`
}
