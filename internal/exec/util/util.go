// Copyright 2015 CoreOS, Inc.
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

package util

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/coreos/ignition/internal/log"
	"github.com/coreos/ignition/internal/resource"
)

var (
	errEscapedMountpoint = errors.New("Symlink traversal resulted in path outside of filesystem")
)

// Util encapsulates logging and destdir indirection for the util methods.
type Util struct {
	DestDir string // directory prefix to use in applying fs paths.
	Fetcher resource.Fetcher
	*log.Logger
}

// SplitPath splits /a/b/c/d into [a, b, c, d]
// golang-- for making me write this

func SplitPath(p string) []string {
	dir, file := filepath.Split(p)
	if dir == "" || dir == "/" {
		return []string{file}
	}
	dir = filepath.Clean(dir)
	return append(SplitPath(dir), file)
}

func wantsToEscape(p string) bool {
	return filepath.Join("/", p) == filepath.Join("/a", p)
}

// JoinPath returns a path into the context ala filepath.Join(d, args)
// It resolves symlinks as if they were rooted at u.DestDir. This means
// that the resulting path will always be under u.DestDir.
// The last element of the path is never followed.
func (u Util) JoinPath(path ...string) (string, error) {
	components := []string{}
	for _, tmp := range path {
		components = append(components, SplitPath(tmp)...)
	}
	last := components[len(components)-1]
	components = components[:len(components)-1]

	realpath := "/"
	for _, component := range components {
		tmp := filepath.Join(realpath, component)
		s, err := os.Lstat(filepath.Join(u.DestDir, tmp))
		if os.IsNotExist(err) {
			realpath = tmp
			continue
		} else if err != nil {
			return "", err
		}

		if s.Mode()&os.ModeSymlink == 0 {
			realpath = tmp
			continue
		}

		symlinkPath, err := os.Readlink(filepath.Join(u.DestDir, tmp))
		if err != nil {
			return "", err
		}
		if filepath.IsAbs(symlinkPath) {
			realpath = "/"
		}
		realpath = filepath.Join(realpath, symlinkPath)
	}

	return filepath.Join(u.DestDir, realpath, last), nil
}
