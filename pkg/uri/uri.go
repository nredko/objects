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
	"encoding/json"
	"fmt"
	"strings"
)

// URI represents the canonical identification of a digital object
type URI struct {
	Scheme string // Scheme identifies a kind of objects
	Opaque string // Rest of encoded data
}

// String implements the Stringer interface
func (u *URI) String() string {
	if u.Scheme != "" {
		return fmt.Sprintf("%s:%s", u.Scheme, u.Opaque)
	}
	return u.Opaque
}

// Parse converts a rawURI string into an URI structure
func Parse(rawURI string) (*URI, error) {
	parts := strings.Split(rawURI, "://")
	l := len(parts)
	if l == 1 {
		return &URI{
			Scheme: "",
			Opaque: rawURI,
		}, nil

	}
	if l == 2 {
		return &URI{
			Scheme: parts[0],
			Opaque: "//" + parts[1],
		}, nil
	}
	return nil, fmt.Errorf("invalid URI: %s", rawURI)
}

// MarshalJSON implements the json.Marshaller interface
func (u URI) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

// UnmarshalJSON parses URI
func (u *URI) UnmarshalJSON(input []byte) error {
	var rawURI string
	if err := json.Unmarshal(input, &rawURI); err != nil {
		return err
	}
	pu, err := Parse(rawURI)
	if err != nil {
		return err
	}
	u.Scheme = pu.Scheme
	u.Opaque = pu.Opaque
	return nil
}
