// Copyright 2019 Red Hat, Inc.
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

// The storage stage is responsible for partitioning disks, creating RAID
// arrays, formatting partitions, writing files, writing systemd units, and
// writing network units.

package disks

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"syscall"

	"github.com/coreos/ignition/internal/config/types"
	"github.com/coreos/ignition/internal/exec/stages"
	"github.com/coreos/ignition/internal/exec/util"
	"github.com/coreos/ignition/internal/log"
	"github.com/coreos/ignition/internal/resource"
)

const (
	name = "mount"
)

func init() {
	stages.Register(creator{})
}

type creator struct{}

func (creator) Create(logger *log.Logger, root string, f resource.Fetcher) stages.Stage {
	return &stage{
		Util: util.Util{
			DestDir: root,
			Logger:  logger,
		},
	}
}

func (creator) Name() string {
	return name
}

type stage struct {
	util.Util
}

func (stage) Name() string {
	return name
}

func (s stage) Run(config types.Config) error {
	fss := config.Storage.Filesystems
	sort.Slice(fss, func(i, j int) bool { return util.Depth(fss[i].Path) < util.Depth(fss[j].Path) })
	for _, fs := range fss {
		if err := s.mountFs(fs); err != nil {
			return err
		}
	}
	return nil
}

// checkForNonDirectories returns an error if any element of path is not a directory
func checkForNonDirectories(path string) error {
	p := "/"
	for _, component := range util.SplitPath(path) {
		p = filepath.Join(p, component)
		st, err := os.Lstat(p)
		if err != nil && os.IsNotExist(err) {
			return nil // nonexistent is ok
		} else if err != nil {
			return err
		}
		if !st.Mode().IsDir() {
			return fmt.Errorf("Mount path %q contains non-directory component %q", path, p)
		}
	}
	return nil
}

func (s stage) mountFs(fs types.Filesystem) error {
	if fs.Format == "swap" {
		return nil
	}

	// mount paths shouldn't include symlinks or other non-directories so we can use filepath.Join()
	// instead of s.JoinPath(). Check that the resulting path is composed of only directories.
	path := filepath.Join(s.DestDir, fs.Path)
	if err := checkForNonDirectories(path); err != nil {
		return err
	}

	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	if err := s.Logger.LogOp(func() error { return syscall.Mount(fs.Device, path, fs.Format, 0, "") },
		"mounting %q at %q with type %q", fs.Device, path, fs.Format,
	); err != nil {
		return err
	}
	return nil
}
