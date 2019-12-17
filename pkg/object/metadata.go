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

// Metadata holds custom artifact attributes
type Metadata map[string]interface{}

func (m *Metadata) init() {
	if (*m) == nil {
		(*m) = Metadata{}
	}
}

// SetValues sets given values into this Metadata instance
func (m *Metadata) SetValues(values map[string]interface{}) {
	m.init()
	for k, v := range values {
		(*m)[k] = v
	}
}

// Set sets the value for given key
func (m *Metadata) Set(key string, value interface{}) {
	m.init()
	(*m)[key] = value
}

// Get returns the value for the given key, if any, otherwise returns defaultValue
func (m Metadata) Get(key string, defaultValue interface{}) interface{} {
	if m == nil {
		return defaultValue
	}
	if v, ok := m[key]; ok {
		return v
	}
	return defaultValue
}
