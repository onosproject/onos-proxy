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
	e2api.RegisterSubscriptionAdminServiceServer(r, server)
}

// SubscriptionServer implements the gRPC service for E2 Subscription related functions.
type SubscriptionServer struct {
}

func (s *SubscriptionServer) GetChannel(ctx context.Context, request *e2api.GetChannelRequest) (*e2api.GetChannelResponse, error) {
	log.Debugf("Received GetChannelRequest %+v", request)
	log.Errorf("TODO: Not implemented yet")
	return nil, nil
}

func (s *SubscriptionServer) ListChannels(ctx context.Context, request *e2api.ListChannelsRequest) (*e2api.ListChannelsResponse, error) {
	log.Debugf("Received ListChannelsRequest %+v", request)
	log.Errorf("TODO: Not implemented yet")
	return nil, nil
}

func (s *SubscriptionServer) WatchChannels(request *e2api.WatchChannelsRequest, server e2api.SubscriptionAdminService_WatchChannelsServer) error {
	log.Debugf("Received WatchChannelsRequest %+v", request)
	log.Errorf("TODO: Not implemented yet")
	return nil
}

func (s *SubscriptionServer) GetSubscription(ctx context.Context, request *e2api.GetSubscriptionRequest) (*e2api.GetSubscriptionResponse, error) {
	log.Debugf("Received GetSubscriptionRequest %+v", request)
	log.Errorf("TODO: Not implemented yet")
	return nil, nil
}

func (s *SubscriptionServer) ListSubscriptions(ctx context.Context, request *e2api.ListSubscriptionsRequest) (*e2api.ListSubscriptionsResponse, error) {
	log.Debugf("Received ListSubscriptionsRequest %+v", request)
	log.Errorf("TODO: Not implemented yet")
	return nil, nil
}

func (s *SubscriptionServer) WatchSubscriptions(request *e2api.WatchSubscriptionsRequest, server e2api.SubscriptionAdminService_WatchSubscriptionsServer) error {
	log.Debugf("Received WatchSubscriptionsRequest %+v", request)
	log.Errorf("TODO: Not implemented yet")
	return nil
}

func (s *SubscriptionServer) Subscribe(request *e2api.SubscribeRequest, server e2api.SubscriptionService_SubscribeServer) error {
	log.Debugf("Received SubscribeRequest %+v", request)
	log.Errorf("TODO: Not implemented yet")
	return nil
}

func (s *SubscriptionServer) Unsubscribe(ctx context.Context, request *e2api.UnsubscribeRequest) (*e2api.UnsubscribeResponse, error) {
	log.Debugf("Received UnsubscribeRequest %+v", request)
	log.Errorf("TODO: Not implemented yet")
	return nil, nil
}
