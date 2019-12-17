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

package sniff

import (
	"errors"
	"os"
)

type Data struct {
	Format   string `json:"format"`
	Type     string `json:"type"`
	Platform string `json:"platform"`
	Arch     string `json:"arch"`
	X64      bool   `json:"x64"`
}

func (d Data) ContentType() string {
	switch true {
	case d.Platform == Platform_MachO:
		return "application/x-mach-binary"
	case d.Platform == Platform_PE:
		return "application/x-dosexec"
	case d.Format == "ELF":
		return "application/x-executable"
	}
	return "application/octet-stream"
}

var sniffers = []func(*os.File) (*Data, error){
	ELF,
	PE,
	MachO,
}

func File(file *os.File) (*Data, error) {

	for _, sniffer := range sniffers {
		if d, e := sniffer(file); e == nil {
			return d, nil
		}
	}

	return nil, errors.New("Nothing found")
}
