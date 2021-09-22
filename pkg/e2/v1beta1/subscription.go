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
	e2api "github.com/onosproject/onos-api/go/onos/e2t/e2/v1beta1"
	"github.com/onosproject/onos-lib-go/pkg/northbound"
	"google.golang.org/grpc"
)

// NewSubscriptionService creates a new E2T subscription service
func NewSubscriptionService() northbound.Service {
	return &SubscriptionService{}
}

// SubscriptionService is a Service implementation for E2 Subscription service.
type SubscriptionService struct {
	northbound.Service
}

// Register registers the SubscriptionService with the gRPC server.
func (s SubscriptionService) Register(r *grpc.Server) {
	server := &SubscriptionServer{}
	e2api.RegisterSubscriptionServiceServer(r, server)
}

// SubscriptionServer implements the gRPC service for E2 Subscription related functions.
type SubscriptionServer struct {
}

func (s *SubscriptionServer) Subscribe(request *e2api.SubscribeRequest, server e2api.SubscriptionService_SubscribeServer) error {
	log.Debugf("Received SubscribeRequest %+v", request)

	conn, err := grpc.Dial("onos-e2t:5150", grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := e2api.NewSubscriptionServiceClient(conn)
	clientStream, err := client.Subscribe(server.Context(), request)
	if err != nil {
		return err
	}

	for {
		response, err := clientStream.Recv()
		if err != nil {
			return err
		}
		err = server.Send(response)
		if err != nil {
			return err
		}
	}
}

func (s *SubscriptionServer) Unsubscribe(ctx context.Context, request *e2api.UnsubscribeRequest) (*e2api.UnsubscribeResponse, error) {
	log.Debugf("Received UnsubscribeRequest %+v", request)
	conn, err := grpc.Dial("onos-e2t:5150", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := e2api.NewSubscriptionServiceClient(conn)
	return client.Unsubscribe(ctx, request)
}


/*
	e2client := NewClient(WithAppID(AppID(request.Headers.AppID)),
						WithServiceModel(ServiceModelName(request.Headers.ServiceModel.Name),
										 ServiceModelVersion(request.Headers.ServiceModel.Version)))
	node := e2client.Node(NodeID(request.Headers.E2NodeID))

	node.Subscribe(request.Subscription, ch)
}
 */
