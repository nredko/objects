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

package git

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

const commitBody = `tree 7483d954ea8be04d8fc4d67692a5ce08705f5c40
parent d3977204a7d9ed24ba0c33ba67310a73d8aa59fd
author moshix <moshix@gmail.com> 1546245813 -0600
committer GitHub <noreply@github.com> 1546245813 -0600
gpgsig -----BEGIN PGP SIGNATURE-----
 
 wsBcBAABCAAQBQJcKda1CRBK7hj4Ov3rIwAAdHIIAFvt8sZMLGXq5fhVtr7nUDIP
 X2bQ0eBSC0nwYtnHTbDbK1gAX39G1JCKiuPZhpQlxpuemPVip9zYTvOL0VKvVATP
 GVTfyxiHcsDmv+7EOSMScJqN0aMJ8QSt8PKM4TYlwuQ41PG1SrNH3NYAiIJKTDUT
 jYI7taBBqW6LRe5t3nIIN+0y7nzZ0lJVlUs1NgQLsGUokqAMJliFhdmOmWQpWymC
 6MkFaSqc4eZTCrzfZWKdzhvbsmzack7IpFsNtYEProMwc/B2mCAyTIjxHlvYFY/e
 EnP0gAl92Fb7mMakx3GxVPHWzv6p0/rAXY6P/ryT1Pp+WKzrRQHoR0fFxweU3p8=
 =Pofq
 -----END PGP SIGNATURE-----
 

Update vcn.go 

Made copyright 2018-2019`

const commitSHA256 = "sha256:2592cad3ca67e7e51c75c5321287e5a11b1709a17ba5c6d2a732e01825c7e151"

func TestDigestCommit(t *testing.T) {
	o := &plumbing.MemoryObject{}
	o.SetType(plumbing.CommitObject)
	_, err := o.Write([]byte(commitBody))
	if err != nil {
		t.Fatal(err)
	}
	c := &object.Commit{}
	err = c.Decode(o)
	if err != nil {
		t.Fatal(err)
	}

	digest, size, err := digestCommit(*c)
	assert.NoError(t, err)
	assert.Equal(t, uint64(len(commitBody)), size)
	assert.Equal(t, commitSHA256, digest.String())
}
