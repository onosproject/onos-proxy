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
	"context"
	"fmt"
	"github.com/onosproject/onos-api/go/onos/topo"
	"github.com/onosproject/onos-lib-go/pkg/grpc/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/serviceconfig"
)

const ResolverName = "e2"
const topoAddress = "onos-topo:5150"

// ResolverBuilder :
type ResolverBuilder struct {
}

// Scheme :
func (b *ResolverBuilder) Scheme() string {
	return ResolverName
}

// Build :
func (b *ResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	var dialOpts []grpc.DialOption
	if opts.DialCreds != nil {
		dialOpts = append(
			dialOpts,
			grpc.WithTransportCredentials(opts.DialCreds),
		)
	} else {
		dialOpts = append(dialOpts, grpc.WithInsecure())
	}
	dialOpts = append(dialOpts, grpc.WithUnaryInterceptor(retry.RetryingUnaryClientInterceptor(retry.WithRetryOn(codes.Unavailable, codes.Unknown))))
	dialOpts = append(dialOpts, grpc.WithStreamInterceptor(retry.RetryingStreamClientInterceptor(retry.WithRetryOn(codes.Unavailable, codes.Unknown))))
	dialOpts = append(dialOpts, grpc.WithContextDialer(opts.Dialer))

	topoConn, err := grpc.Dial(topoAddress, dialOpts...)
	if err != nil {
		return nil, err
	}

	serviceConfig := cc.ParseServiceConfig(
		fmt.Sprintf(`{"loadBalancingConfig":[{"%s":{}}]}`, ResolverName),
	)

	log.Infof("Built new resolver")

	resolver := &Resolver{
		clientConn:    cc,
		topoConn:      topoConn,
		serviceConfig: serviceConfig,
		nodes:         make(map[topo.ID]*topo.MastershipState),
		controls:      make(map[topo.ID]topo.ID),
		e2ts:          make(map[topo.ID]string),
	}
	err = resolver.start()
	if err != nil {
		return nil, err
	}
	return resolver, nil
}

var _ resolver.Builder = (*ResolverBuilder)(nil)

// Resolver :
type Resolver struct {
	clientConn    resolver.ClientConn
	topoConn      *grpc.ClientConn
	serviceConfig *serviceconfig.ParseResult
	nodes         map[topo.ID]*topo.MastershipState // E2 node to mastership (controls relation ID)
	controls      map[topo.ID]topo.ID               // controls relation to E2T ID
	e2ts          map[topo.ID]string                // E2T ID to address
}

func (r *Resolver) start() error {
	log.Infof("Starting resolver")

	client := topo.NewTopoClient(r.topoConn)
	request := &topo.WatchRequest{}
	stream, err := client.Watch(context.Background(), request)
	if err != nil {
		return err
	}
	go func() {
		for {
			response, err := stream.Recv()
			if err != nil {
				return
			}
			r.handleEvent(response.Event)
		}
	}()
	return nil
}

func (r *Resolver) handleEvent(event topo.Event) {
	object := event.Object
	if entity, ok := object.Obj.(*topo.Object_Entity); ok && entity.Entity.KindID == topo.E2NODE {
		// Track changes in E2 nodes
		switch event.Type {
		case topo.EventType_REMOVED:
			delete(r.nodes, object.ID)
		default:
			var m topo.MastershipState
			_ = object.GetAspect(&m)
			if node, ok := r.nodes[object.ID]; !ok || m.Term > node.Term {
				r.nodes[object.ID] = &m
			}
		}
		r.updateState()

	} else if entity, ok := object.Obj.(*topo.Object_Entity); ok && entity.Entity.KindID == topo.E2T {
		// Track changes in E2T instances
		switch event.Type {
		case topo.EventType_REMOVED:
			delete(r.e2ts, object.ID)
			r.updateState()
		default:
			var info topo.E2TInfo
			_ = object.GetAspect(&info)
			newAddress := r.e2ts[object.ID]
			for _, iface := range info.Interfaces {
				if iface.Type == topo.Interface_INTERFACE_E2T {
					newAddress = fmt.Sprintf("%s:%d", iface.IP, iface.Port)
				}
			}
			if r.e2ts[object.ID] != newAddress {
				r.e2ts[object.ID] = newAddress
				r.updateState()
			}
		}

	} else if relation, ok := object.Obj.(*topo.Object_Relation); ok && relation.Relation.KindID == topo.CONTROLS {
		// Track changes in E2T/E2Node controls relations
		switch event.Type {
		case topo.EventType_REMOVED:
			delete(r.controls, object.ID)
		default:
			r.controls[object.ID] = relation.Relation.SrcEntityID
		}
		r.updateState()
	}
}

func (r *Resolver) updateState() {
	// Produce list of addresses for available E2T instances
	// Annotate each address with a list of nodes for which this instances is presently the master
	e2tMastership := make(map[topo.ID][]string)

	// Scan over all nodes and insert their ID into the list of nodes of its master E2T instance
	for nodeID, mastership := range r.nodes {
		if e2tID, ok := r.controls[topo.ID(mastership.NodeId)]; ok {
			var nodes []string
			if nodes, ok = e2tMastership[e2tID]; !ok {
				nodes = make([]string, 0)
			}
			e2tMastership[e2tID] = append(nodes, string(nodeID))
		}
	}

	// Transpose the map of E2T node IDs into a list of addresses with nodes attribute
	addresses := make([]resolver.Address, 0, len(r.e2ts))
	for e2tID, addr := range r.e2ts {
		var nodes []string
		var ok bool
		if nodes, ok = e2tMastership[e2tID]; !ok {
			nodes = make([]string, 0)
		}
		addresses = append(addresses, resolver.Address{
			Addr: addr,
			Attributes: attributes.New(
				"nodes",
				nodes,
			),
		})
		log.Debugf("New resolver address: %s => %+v", addr, nodes)
	}

	log.Infof("New resolver addresses: %+v", addresses)

	// Update the resolver state with list of E2T addresses annotated by nodes for which they are masters
	r.clientConn.UpdateState(resolver.State{
		Addresses:     addresses,
		ServiceConfig: r.serviceConfig,
	})
}

// ResolveNow :
func (r *Resolver) ResolveNow(resolver.ResolveNowOptions) {}

// Close :
func (r *Resolver) Close() {
	if err := r.topoConn.Close(); err != nil {
		log.Error("failed to close conn", err)
	}
}

var _ resolver.Resolver = (*Resolver)(nil)
