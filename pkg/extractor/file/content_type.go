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
	"io"
	"net/http"
	"os"

	"github.com/h2non/filetype"
)

func contentType(file *os.File) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	file.Seek(0, 0)
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil {
		if n == 0 && err == io.EOF {
			// empty file, no content type
			return "", nil
		}
		return "", err
	}

	kind, err := filetype.Match(buf)
	if err == nil && kind != filetype.Unknown {
		return kind.MIME.Value, nil
	}

	// As fallback, use the net/http package's handy DectectContentType function.
	// Always returns a valid content-type by returning "application/octet-stream"
	// if no others seemed to match.
	contentType := http.DetectContentType(buf)

	return contentType, nil
}
