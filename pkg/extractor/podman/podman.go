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

package podman

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/codenotary/objects/pkg/extractor"
	"github.com/codenotary/objects/pkg/object"
	"github.com/codenotary/objects/pkg/uri"

	digest "github.com/opencontainers/go-digest"
	"github.com/varlink/go/varlink"
)

// Scheme for podman
const Scheme = "podman"

var service *varlink.Service
var bridge string

type imageData struct {
	Id           string `json:"Id"`
	Digest       string
	RepoTags     []string
	Architecture string
	Os           string
	Size         int
}

// Extract returns a file *object.Object from a given u
func Extract(u uri.URI, options ...extractor.Option) (*object.Object, error) {

	if u.Scheme != Scheme {
		return nil, nil
	}
	imageID := strings.TrimPrefix(u.Opaque, "//")

	var meth string
	if strings.Contains(imageID, "/") {
		meth = "InspectImage"
	} else {
		meth = "InspectContainer"
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var args []string = []string{"unix://run/podman/io.podman/io.podman." + meth, "{\"name\":\"" + imageID + "\"}"}

	resp := varlink_call(ctx, args)

	var imageResp map[string]interface{}
	rdr := bytes.NewReader(resp)
	err := json.NewDecoder(rdr).Decode(&imageResp)
	if err != nil {
		return nil, err
	}

	var r interface{}
	for key := range imageResp {
		r = imageResp[key]
	}

	str := fmt.Sprintf("%v", r)

	var res map[string]interface{}
	rd := strings.NewReader(str)
	err = json.NewDecoder(rd).Decode(&res)
	if err != nil {
		return nil, err
	}
	jsonbody, _ := json.Marshal(res)

	image := imageData{}
	if err := json.Unmarshal(jsonbody, &image); err != nil {
		return nil, err
	}

	m := object.Metadata{
		"name":         getName(image),
		"architecture": image.Architecture,
		"platform":     image.Os,
		Scheme:         res,
	}

	if version := inferVer(image); version != "" {
		m["version"] = version
	}

	return &object.Object{
		Descriptor: object.Descriptor{
			Digest: digest.FromString(image.Id),
			Size:   uint64(image.Size),
		},
		Metadata: m,
		URI:      u,
	}, nil
}

func getName(i imageData) string {
	if len(i.RepoTags) > 0 {
		return i.RepoTags[0]
	}
	return strings.TrimSpace(i.Id)
}

func inferVer(i imageData) string {
	if len(i.RepoTags) > 0 {
		parts := strings.SplitN(i.RepoTags[0], ":", 2)
		if len(parts) > 1 && parts[1] != "latest" {
			return parts[1]
		}
	}

	return ""
}

func ErrPrintf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s ", "")
	fmt.Fprintf(os.Stderr, format, a...)
}

func varlink_call(ctx context.Context, args []string) []byte {
	var err error
	var oneway bool

	callFlags := flag.NewFlagSet("help", flag.ExitOnError)
	callFlags.BoolVar(&oneway, "-oneway", false, "Use bridge for connection")
	var help bool
	callFlags.BoolVar(&help, "help", false, "Prints help information")
	var usage = func() { print_usage(callFlags, "<[ADDRESS/]INTERFACE.METHOD> [ARGUMENTS]") }
	callFlags.Usage = usage

	_ = callFlags.Parse(args)

	if help {
		usage()
	}

	var con *varlink.Connection
	var address string
	var methodName string

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if len(bridge) != 0 {
		con, err = varlink.NewBridge(bridge)

		if err != nil {
			ErrPrintf("Cannot connect with bridge '%s': %v\n", bridge, err)
			os.Exit(2)
		}
		address = "bridge:" + bridge
		methodName = callFlags.Arg(0)
	} else {
		uri := callFlags.Arg(0)
		if uri == "" {
			usage()
		}

		li := strings.LastIndex(uri, "/")

		if li == -1 {
			ErrPrintf("Invalid address '%s'\n", uri)
			os.Exit(2)
		}

		address = uri[:li]
		methodName = uri[li+1:]

		con, err = varlink.NewConnection(ctx, address)

		if err != nil {
			ErrPrintf("Cannot connect to '%s': %v\n", address, err)
			os.Exit(2)
		}
	}
	var parameters string
	var params json.RawMessage

	parameters = callFlags.Arg(1)
	if parameters == "" {
		params = nil
	} else {
		json.Unmarshal([]byte(parameters), &params)
	}

	var flags uint64
	flags = 0
	if oneway {
		flags |= varlink.Oneway
	}
	recv, err := con.Send(ctx, methodName, params, flags)

	var retval map[string]interface{}

	// FIXME: Use cont
	_, err = recv(ctx, &retval)

	if err != nil {
		if e, ok := err.(*varlink.Error); ok {
			ErrPrintf("Call failed with error: %v\n", e.Name)
			errorRawParameters := e.Parameters.(*json.RawMessage)
			if errorRawParameters != nil {
				var param map[string]interface{}
				_ = json.Unmarshal(*errorRawParameters, &param)
				c, _ := json.Marshal(param)
				fmt.Fprintf(os.Stderr, "%v\n", string(c))
			}
			os.Exit(2)
		}
		ErrPrintf("Error calling '%s': %v\n", methodName, err)
		os.Exit(2)
	}
	c, _ := json.Marshal(retval)
	return c
}

func print_usage(set *flag.FlagSet, arg_help string) {
	if set == nil {
		fmt.Fprintf(os.Stderr, "Usage: %s [GLOBAL OPTIONS] COMMAND ...\n", os.Args[0])
	} else {
		fmt.Fprintf(os.Stderr, "Usage: %s [GLOBAL OPTIONS] %s [OPTIONS] %s\n", os.Args[0], set.Name(), arg_help)
	}

	fmt.Fprintln(os.Stderr, "\nGlobal Options:")
	flag.PrintDefaults()

	if set == nil {
		fmt.Fprintln(os.Stderr, "\nCommands:")
		fmt.Fprintln(os.Stderr, "  info\tPrint information about a service")
		fmt.Fprintln(os.Stderr, "  help\tPrint interface description or service information")
		fmt.Fprintln(os.Stderr, "  call\tCall a method")
	} else {
		fmt.Fprintln(os.Stderr, "\nOptions:")
		set.PrintDefaults()
	}
	os.Exit(1)
}
