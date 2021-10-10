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

package balancer

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/metadata"
)

const e2NodeIDHeader = "e2-node-id"

func init() {
	balancer.Register(base.NewBalancerBuilder(ResolverName, &PickerBuilder{}, base.Config{}))
}

// PickerBuilder :
type PickerBuilder struct{}

// Build :
func (p *PickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	masters := make(map[string]balancer.SubConn)

	for sc, scInfo := range info.ReadySCs {
		nodes := scInfo.Address.Attributes.Value("nodes").([]string)
		for _, node := range nodes {
			masters[node] = sc
		}
	}
	log.Infof("Built new picker for E2T instances: %+v", masters)
	return &Picker{
		masters: masters,
	}
}

var _ base.PickerBuilder = (*PickerBuilder)(nil)

// Picker :
type Picker struct {
	masters map[string]balancer.SubConn // NodeID string to connection mapping
}

// Pick :
func (p *Picker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	var result balancer.PickResult
	if md, ok := metadata.FromIncomingContext(info.Ctx); ok {
		ids := md.Get(e2NodeIDHeader)
		if len(ids) > 0 {
			if subConn, ok := p.masters[ids[0]]; ok {
				log.Debugf("Picked subconn for %s: %+v", ids[0], subConn)
				result.SubConn = subConn
				return result, nil
			}
		}
	}
	log.Warn("No subconn available")
	return result, balancer.ErrNoSubConnAvailable
}

var _ balancer.Picker = (*Picker)(nil)
