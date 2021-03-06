// Copyright 2016-2017 VMware, Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fs

import (
	"os/exec"
	"strings"

	"github.com/docker/docker/pkg/mount"

	"github.com/vmware/vic/pkg/trace"

	log "github.com/Sirupsen/logrus"
)

// Ext4 satisfies the Filesystem interface
type Ext4 struct{}

func NewExt4() *Ext4 {
	return &Ext4{}
}

// Mkfs creates an ext4 fs on the given device and applices the given label
func (e *Ext4) Mkfs(devPath, label string) error {
	defer trace.End(trace.Begin(devPath))

	log.Infof("Creating ext4 filesystem on device %s", devPath)

	// -v is verbose - this is only useful when things go wrong,
	// -F is needed to use the entire disk without prompting
	// we can't use -V as well for fs specific stuff as that prevents it actually being done.
	// #nosec: Subprocess launching with variable
	cmd := exec.Command("/sbin/mkfs.ext4", "-L", label, "-vF", devPath)

	if output, err := cmd.CombinedOutput(); err != nil {
		log.Errorf("vmdk storage driver failed to format disk %s: %s", devPath, err)
		log.Errorf("mkfs output: %s", string(output))
		return err
	}
	log.Debugf("Filesystem created on device %s", devPath)

	return nil
}

// Mount mounts an ext4 formatted device at the given path.  From the Docker
// mount pkg, args must in the form arg=val.
func (e *Ext4) Mount(devPath, targetPath string, options []string) error {
	defer trace.End(trace.Begin(devPath))
	log.Infof("Mounting %s to %s", devPath, targetPath)
	return mount.Mount(devPath, targetPath, "ext4", strings.Join(options, ","))
}

// Unmount unmounts the disk.
// path can be a device path or a mount point
func (e *Ext4) Unmount(path string) error {
	defer trace.End(trace.Begin(path))
	log.Infof("Unmounting %s", path)
	return mount.Unmount(path)
}

// SetLabel sets the label of an ext4 formated device
func (e *Ext4) SetLabel(devPath, labelName string) error {
	defer trace.End(trace.Begin(devPath))

	// #nosec: Subprocess launching with variable
	cmd := exec.Command("/sbin/e2label", devPath, labelName)
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Errorf("failed to set label %s: %s", devPath, err)
		log.Error(string(output))
		return err
	}

	return nil
}
