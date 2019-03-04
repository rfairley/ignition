// Copyright 2017 CoreOS, Inc.
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

package distro

import (
	"fmt"
	"os"
)

// Distro-specific settings that can be overridden at link time with e.g.
// -X github.com/coreos/ignition/internal/distro.mdadmCmd=/opt/bin/mdadm
var (
	// Device node directories and paths
	diskByIDDir       = "/dev/disk/by-id"
	diskByLabelDir    = "/dev/disk/by-label"
	diskByPartUUIDDir = "/dev/disk/by-partuuid"

	// File paths
	kernelCmdlinePath = "/proc/cmdline"
	// initramfs directory containing distro-provided base config
	systemConfigDir = "/usr/lib/ignition"

	// Helper programs
	chrootCmd   = "chroot"
	groupaddCmd = "groupadd"
	idCmd       = "id"
	mdadmCmd    = "mdadm"
	mountCmd    = "mount"
	sgdiskCmd   = "sgdisk"
	udevadmCmd  = "udevadm"
	usermodCmd  = "usermod"
	useraddCmd  = "useradd"

	// The restorecon tool is embedded inside of a systemd unit
	// and as such requires the absolute path
	restoreconCmd = "/usr/sbin/restorecon"

	// Filesystem tools
	btrfsMkfsCmd = "mkfs.btrfs"
	ext4MkfsCmd  = "mkfs.ext4"
	swapMkfsCmd  = "mkswap"
	vfatMkfsCmd  = "mkfs.vfat"
	xfsMkfsCmd   = "mkfs.xfs"

	// Flags
	selinuxRelabel  = "false"
	blackboxTesting = "false"
	// useAuthorizedKeysFile specifies whether to sync user SSH key
	// fragments in .ssh/authorized_keys.d to .ssh/authorized_keys. Set
	// this to "true" for distros that do not support reading key fragments
	// from .ssh/authorized_keys.d.
	useAuthorizedKeysFile = "false"
)

func DiskByIDDir() string       { return diskByIDDir }
func DiskByLabelDir() string    { return diskByLabelDir }
func DiskByPartUUIDDir() string { return diskByPartUUIDDir }

func KernelCmdlinePath() string { return kernelCmdlinePath }
func SystemConfigDir() string   { return fromEnv("SYSTEM_CONFIG_DIR", systemConfigDir) }

func ChrootCmd() string     { return chrootCmd }
func GroupaddCmd() string   { return groupaddCmd }
func IdCmd() string         { return idCmd }
func MdadmCmd() string      { return mdadmCmd }
func MountCmd() string      { return mountCmd }
func SgdiskCmd() string     { return sgdiskCmd }
func UdevadmCmd() string    { return udevadmCmd }
func UsermodCmd() string    { return usermodCmd }
func UseraddCmd() string    { return useraddCmd }
func RestoreconCmd() string { return restoreconCmd }

func BtrfsMkfsCmd() string { return btrfsMkfsCmd }
func Ext4MkfsCmd() string  { return ext4MkfsCmd }
func SwapMkfsCmd() string  { return swapMkfsCmd }
func VfatMkfsCmd() string  { return vfatMkfsCmd }
func XfsMkfsCmd() string   { return xfsMkfsCmd }

func SelinuxRelabel() bool  { return bakedStringToBool(selinuxRelabel) }
func BlackboxTesting() bool { return bakedStringToBool(blackboxTesting) }
func UseAuthorizedKeysFile() bool {
	return bakedStringToBool(fromEnv("USE_AUTHORIZED_KEYS_FILE", useAuthorizedKeysFile))
}

func fromEnv(nameSuffix, defaultValue string) string {
	value := os.Getenv("IGNITION_" + nameSuffix)
	if value != "" {
		return value
	}
	return defaultValue
}

func bakedStringToBool(s string) bool {
	// the linker only supports string args, so do some basic bool sensing
	if s == "true" || s == "1" {
		return true
	} else if s == "false" || s == "0" {
		return false
	} else {
		// if we got a bad compile flag, just crash and burn rather than assume
		panic(fmt.Sprintf("value '%s' cannot be interpreted as a boolean", s))
	}
}
