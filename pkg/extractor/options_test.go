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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type dummyOpts struct {
	flag bool
}

func withTestOption() Option {
	return func(o interface{}) error {
		if oo, ok := o.(*dummyOpts); ok {
			oo.flag = true
		}
		return nil
	}
}

func withTestOptionWithError() Option {
	return func(o interface{}) error {
		return errors.New("some error")
	}
}

func TestApply(t *testing.T) {

	opts := &dummyOpts{}

	err := Options([]Option{
		withTestOption(),
	}).Apply(opts)
	assert.NoError(t, err)
	assert.True(t, opts.flag)

	err = Options([]Option{
		withTestOptionWithError(),
	}).Apply(opts)
	assert.Error(t, err)
}
