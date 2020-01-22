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
	"fmt"

	"github.com/codenotary/objects/pkg/object"
	"github.com/codenotary/objects/pkg/uri"
)

var extractors = map[string]Extractor{}

// Extractor extract an object.Object referenced by the given uri.URI.
type Extractor func(uri.URI, ...Option) (*object.Object, error)

// Register the Extractor e for the given scheme
func Register(scheme string, e Extractor) {
	extractors[scheme] = e
}

// Schemes returns the list of registered schemes.
func Schemes() []string {
	schemes := make([]string, len(extractors))
	i := 0
	for scheme := range extractors {
		schemes[i] = scheme
		i++
	}
	return schemes
}

// Extract returns an object.Object for the given rawURI.
func Extract(rawURI string, options ...Option) (*object.Object, error) {
	u, err := uri.Parse(rawURI)
	if err != nil {
		return nil, err
	}

	if e, ok := extractors[u.Scheme]; ok {
		return e(*u, options...)
	}
	return nil, fmt.Errorf("%s scheme not yet supported", u.Scheme)
}
