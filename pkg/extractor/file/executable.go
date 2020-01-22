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

package file

import (
	"os"
	"strings"

	"github.com/codenotary/objects/pkg/extractor/file/internal/sniff"
	"github.com/codenotary/objects/pkg/object"
)

func xInfo(file *os.File, contentType *string) (bool, object.Metadata, error) {
	if strings.HasPrefix(*contentType, "application/") {
		d, err := sniff.File(file)
		if err != nil {
			return false, nil, err
		}
		*contentType = d.ContentType()
		return true, object.Metadata{
			"architecture": strings.ToLower(d.Arch),
			"platform":     d.Platform,
			"file":         d,
		}, nil
	}
	return false, nil, nil
}
