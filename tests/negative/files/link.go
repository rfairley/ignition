// Copyright 2018 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package files

import (
	"github.com/coreos/ignition/tests/register"
	"github.com/coreos/ignition/tests/types"
)

func init() {
	register.Register(register.NegativeTest, WriteOverBrokenSymlink())
	register.Register(register.NegativeTest, WriteThroughRelativeSymlink())
	register.Register(register.NegativeTest, WriteThroughRelativeSymlinkBeyondRoot())
	register.Register(register.NegativeTest, WriteThroughAbsoluteSymlink())
}

func WriteOverBrokenSymlink() types.Test {
	name := "Write Over Broken Symlink at end of path"
	in := types.GetBaseDisk()
	out := types.GetBaseDisk()
	config := `{
	  "ignition": { "version": "$version" },
	  "storage": {
	    "files": [{
	      "path": "/etc/file",
	      "overwrite": false,
	      "mode": 420
	    }]
	  }
	}`
	in[0].Partitions.AddLinks("ROOT", []types.Link{
		{
			Node: types.Node{
				Name:      "file",
				Directory: "etc",
			},
			Target: "/usr/rofile",
		},
	})
	configMinVersion := "3.0.0-experimental"

	return types.Test{
		Name:             name,
		In:               in,
		Out:              out,
		Config:           config,
		ConfigMinVersion: configMinVersion,
	}
}

func WriteThroughRelativeSymlink() types.Test {
	name := "Write Through Relative Symlink on the Root Filesystem"
	in := types.GetBaseDisk()
	out := types.GetBaseDisk()
	// note this abuses the order in which ignition writes links and will break with 3.0.0
	// Also tests that Ignition does not try to resolve symlink targets
	config := `{
	  "ignition": { "version": "$version" },
	  "storage": {
	    "links": [{
	      "path": "/foo/bar",
	      "target": "../etc"
	    },
	    {
	      "path": "/foo/bar/baz",
	      "target": "somewhere/over/the/rainbow"
	    }]
	  }
	}`
	out[0].Partitions.AddLinks("ROOT", []types.Link{
		{
			Node: types.Node{
				Name:      "bar",
				Directory: "foo",
			},
			Target: "../etc",
			Hard:   false,
		},
		{
			Node: types.Node{
				Name:      "baz",
				Directory: "etc",
			},
			Target: "somewhere/over/the/rainbow",
			Hard:   false,
		},
	})
	configMinVersion := "3.0.0-experimental"

	return types.Test{
		Name:             name,
		In:               in,
		Out:              out,
		Config:           config,
		ConfigMinVersion: configMinVersion,
	}
}

func WriteThroughRelativeSymlinkBeyondRoot() types.Test {
	name := "Write Through Relative Symlink beyond the root on the Root Filesystem"
	in := types.GetBaseDisk()
	out := types.GetBaseDisk()
	// note this abuses the order in which ignition writes links and will break with 3.0.0
	// Also tests that Ignition does not try to resolve symlink targets
	config := `{
	  "ignition": { "version": "$version" },
	  "storage": {
	    "links": [{
	      "path": "/foo/bar",
	      "target": "../../etc"
	    },
	    {
	      "path": "/foo/bar/baz",
	      "target": "somewhere/over/the/rainbow"
	    }]
	  }
	}`
	out[0].Partitions.AddLinks("ROOT", []types.Link{
		{
			Node: types.Node{
				Name:      "bar",
				Directory: "foo",
			},
			Target: "../../etc",
			Hard:   false,
		},
		{
			Node: types.Node{
				Name:      "baz",
				Directory: "etc",
			},
			Target: "somewhere/over/the/rainbow",
			Hard:   false,
		},
	})
	configMinVersion := "3.0.0-experimental"

	return types.Test{
		Name:             name,
		In:               in,
		Out:              out,
		Config:           config,
		ConfigMinVersion: configMinVersion,
	}
}

func WriteThroughAbsoluteSymlink() types.Test {
	name := "Write Through Absolute Symlink on the Root Filesystem"
	in := types.GetBaseDisk()
	out := types.GetBaseDisk()
	// note this abuses the order in which ignition writes links and will break with 3.0.0
	// Also tests that Ignition does not try to resolve symlink targets
	config := `{
	  "ignition": { "version": "$version" },
	  "storage": {
	    "links": [{
	      "path": "/foo/bar",
	      "target": "/etc"
	    },
	    {
	      "path": "/foo/bar/baz",
	      "target": "somewhere/over/the/rainbow"
	    }]
	  }
	}`
	out[0].Partitions.AddLinks("ROOT", []types.Link{
		{
			Node: types.Node{
				Name:      "bar",
				Directory: "foo",
			},
			Target: "/etc",
			Hard:   false,
		},
		{
			Node: types.Node{
				Name:      "baz",
				Directory: "etc",
			},
			Target: "somewhere/over/the/rainbow",
			Hard:   false,
		},
	})
	configMinVersion := "3.0.0-experimental"

	return types.Test{
		Name:             name,
		In:               in,
		Out:              out,
		Config:           config,
		ConfigMinVersion: configMinVersion,
	}
}
