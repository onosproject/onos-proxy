// Copyright 2021-present Open Networking Foundation.
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

package v1beta1

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

func init() {
	balancer.Register(base.NewBalancerBuilder(resolverName, &PickerBuilder{}, base.Config{}))
}

// PickerBuilder :
type PickerBuilder struct{}

// Build :
func (p *PickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	var master balancer.SubConn
	var backups []balancer.SubConn
	for sc, scInfo := range info.ReadySCs {
		isMaster := scInfo.Address.Attributes.Value("is_master").(bool)
		if isMaster {
			master = sc
			continue
		}
		backups = append(backups, sc)
	}
	log.Debugf("Built new picker. Master: %s, Backups: %s", master, backups)
	return &Picker{
		master: master,
	}
}

var _ base.PickerBuilder = (*PickerBuilder)(nil)

// Picker :
type Picker struct {
	master balancer.SubConn
}

// Pick :
func (p *Picker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	var result balancer.PickResult
	if p.master == nil {
		return result, balancer.ErrNoSubConnAvailable
	}
	result.SubConn = p.master
	return result, nil
}

var _ balancer.Picker = (*Picker)(nil)
