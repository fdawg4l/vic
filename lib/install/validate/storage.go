// Copyright 2016 VMware, Inc. All Rights Reserved.
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

package validate

import (
	"fmt"
	"net/url"

	"context"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/lib/install/data"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/trace"
	"github.com/vmware/vic/pkg/vsphere/datastore"
)

func (v *Validator) storage(ctx context.Context, input *data.Data, conf *config.VirtualContainerHostConfigSpec) {
	defer trace.End(trace.Begin(""))

	// Image Store
	imageDSpath, ds, err := v.DatastoreHelper(ctx, input.ImageDatastorePath, "", "--image-store")

	if err != nil {
		v.NoteIssue(err)
		return
	}

	// provide a default path if only a DS name is provided
	if imageDSpath.Path == "" {
		imageDSpath.Path = input.DisplayName
	}

	if ds != nil {
		v.SetDatastore(ds, imageDSpath)
		conf.AddImageStore(imageDSpath)
	}

	if conf.VolumeLocations == nil {
		conf.VolumeLocations = make(map[string]*url.URL)
	}

	// TODO: add volume locations
	for label, volDSpath := range input.VolumeLocations {
		dsURL, _, err := v.DatastoreHelper(ctx, volDSpath, label, "--volume-store")
		v.NoteIssue(err)
		if dsURL != nil {
			conf.VolumeLocations[label] = dsURL
		}
	}
}

func (v *Validator) DatastoreHelper(ctx context.Context, path string, label string, flag string) (*object.DatastorePath, *object.Datastore, error) {
	defer trace.End(trace.Begin(path))

	var (
		dsURL  *object.DatastorePath
		dstore *object.Datastore
		err    error
	)

	if len(path) {
		dsURL, err = datastore.DatastorePathFromURLString(path)
		if err != nil {
			return nil, nil, errors.Errorf("error parsing datastore path: %s", err)
		}

		stores, err := v.Session.Finder.DatastoreList(ctx, dsURL.Datastore)
		if err != nil {
			log.Debugf("no such datastore %#v", dsURL)
			v.suggestDatastore(path, label, flag)
			// TODO: error message about no such match and how to get a datastore list
			// we return err directly here so we can check the type
			return nil, nil, err
		}

		if len(stores) > 1 {
			// TODO: error about required disabmiguation and list entries in stores
			v.suggestDatastore(path, label, flag)
			return nil, nil, errors.New("ambiguous datastore " + dsURL.Host)
		}

		dstore = stores[0]

	} else {

		// see if we can find a default datastore
		dstore, err = v.Session.Finder.DatastoreOrDefault(ctx, "*")
		if err != nil {
			v.suggestDatastore("*", label, flag)
			return nil, nil, errors.New("datastore empty")
		}

		dsURL = &object.DatastorePath{dstore.Name(), ""}
		log.Infof("Using default datastore: " + dsURL.String())
	}

	// temporary until session is extracted
	// FIXME: commented out until components can consume moid
	// dsURL.Host = stores[0].Reference().Value

	return dsURL, dstore, nil
}

func (v *Validator) SetDatastore(ds *object.Datastore, path *object.DatastorePath) {
	v.Session.Datastore = ds
	v.Session.DatastorePath = path
}

// suggestDatastore suggests all datastores present on target in datastore:label format if applicable
func (v *Validator) suggestDatastore(path string, label string, flag string) {
	defer trace.End(trace.Begin(""))

	var val string
	if label != "" {
		val = fmt.Sprintf("%s:%s", path, label)
	} else {
		val = path
	}
	log.Infof("Suggesting valid values for %s based on %q", flag, val)

	dss, err := v.Session.Finder.DatastoreList(v.Context, "*")
	if err != nil {
		log.Errorf("Unable to list datastores: %s", err)
		return
	}

	if len(dss) == 0 {
		log.Info("No datastores found")
		return
	}

	matches := make([]string, len(dss))
	for i, d := range dss {
		if label != "" {
			matches[i] = fmt.Sprintf("%s:%s", d.Name(), label)
		} else {
			matches[i] = d.Name()
		}
	}

	if matches != nil {
		log.Infof("Suggested values for %s:", flag)
		for _, d := range matches {
			log.Infof("  %q", d)
		}
	}
}
