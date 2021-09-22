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

// NewControlService creates a new control service
func NewControlService() northbound.Service {
	return &ControlService{}
}

// ControlService is a Service implementation for control requests
type ControlService struct {
	northbound.Service
}

// Register registers the Service with the gRPC server.
func (s ControlService) Register(r *grpc.Server) {
	server := &ControlServer{}
	e2api.RegisterControlServiceServer(r, server)
}

// ControlServer implements the gRPC service for control
type ControlServer struct {
}

func (s *ControlServer) Control(ctx context.Context, request *e2api.ControlRequest) (*e2api.ControlResponse, error) {
	log.Infof("Received E2 Control Request %v", request)
	conn, err := grpc.Dial("onos-e2t:5150", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := e2api.NewControlServiceClient(conn)
	return client.Control(ctx, request)
}
