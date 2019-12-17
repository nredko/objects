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
	"debug/macho"
	"os"
	"strings"
)

const Platform_MachO = "Mach"

func MachO(file *os.File) (*Data, error) {
	f, err := macho.NewFile(file)
	if err != nil {
		return nil, err
	}

	cpu := strings.TrimPrefix(f.Cpu.String(), "Cpu")

	d := &Data{
		Type:     f.Type.String(),
		Platform: Platform_MachO,
		Arch:     cpu,
		X64:      strings.HasSuffix(cpu, "64"),
	}
	return d, nil
}
