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

package options

type PortLayerOptionsType struct {
	SDK      string `long:"sdk" description:"SDK URL or proxy" env:"VC_URL" required:"true"`
	Cert     string `long:"cert" description:"Client certificate" env:"VC_CERTIFICATE"`
	Key      string `long:"key" description:"Private key file" env:"VC_PRIVATE_KEY"`
	Insecure bool   `long:"insecure" description:"Skip verification of server certificate" env:"VC_INSECURE"`

	DatacenterPath string `long:"datacenter" description:"Datacenter path" env:"DC_PATH" required:"true"`
	ClusterPath    string `long:"cluster" description:"Cluster path" env:"CS_PATH" required:"true"`
	DatastorePath  string `long:"datastore" description:"Datastore path" env:"DS_PATH" required:"true"`
	NetworkPath    string `long:"network" description:"Network path" env:"NET_PATH"`
}

var (
	PortLayerOptions = new(PortLayerOptionsType)
)