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

package manager

import (
	"context"
	"fmt"
	"github.com/onosproject/onos-lib-go/pkg/grpc/retry"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-lib-go/pkg/northbound"
	e2v1beta1service "github.com/onosproject/onos-proxy/pkg/e2/v1beta1"
	"github.com/onosproject/onos-proxy/pkg/e2/v1beta1/balancer"
	"github.com/onosproject/onos-proxy/pkg/utils/creds"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
)

var log = logging.GetLogger("onos", "proxy", "manager")

// Config is a manager configuration
type Config struct {
	CAPath   string
	KeyPath  string
	CertPath string
	GRPCPort int
}

// NewManager creates a new manager
func NewManager(config Config) *Manager {
	log.Info("Creating Manager")
	return &Manager{
		Config: config,
	}
}

// Manager is a manager for the E2T service
type Manager struct {
	Config Config
}

// Run starts the manager and the associated services
func (m *Manager) Run() {
	log.Info("Running Manager")
	if err := m.Start(); err != nil {
		log.Fatal("Unable to run Manager", err)
	}
}

// Start starts the manager
func (m *Manager) Start() error {
	err := m.startNorthboundServer()
	if err != nil {
		return err
	}
	return nil
}

// startSouthboundServer starts the northbound gRPC server
func (m *Manager) startNorthboundServer() error {
	s := northbound.NewServer(&northbound.ServerConfig{
		CaPath:      &m.Config.CAPath,
		KeyPath:     &m.Config.KeyPath,
		CertPath:    &m.Config.CertPath,
		Port:        int16(m.Config.GRPCPort),
		Insecure:    true,
		SecurityCfg: &northbound.SecurityConfig{},
	})

	conn, err := m.connect(context.Background())
	if err != nil {
		log.Errorf("Unable to connect to E2T service")
		return err
	}

	s.AddService(logging.Service{})
	s.AddService(e2v1beta1service.NewProxyService(conn))

	doneCh := make(chan error)
	go func() {
		err := s.Serve(func(started string) {
			log.Info("Started NBI on ", started)
			close(doneCh)
		})
		if err != nil {
			doneCh <- err
		}
	}()
	return <-doneCh
}

func (m *Manager) connect(ctx context.Context) (*grpc.ClientConn, error) {
	clientCreds, _ := creds.GetClientCredentials()
	conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:///%s", balancer.ResolverName, "onos-e2t:5150"),
		grpc.WithTransportCredentials(credentials.NewTLS(clientCreds)),
		grpc.WithUnaryInterceptor(retry.RetryingUnaryClientInterceptor(retry.WithRetryOn(codes.Unavailable))),
		grpc.WithStreamInterceptor(retry.RetryingStreamClientInterceptor(retry.WithRetryOn(codes.Unavailable))))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// Close kills the connections and manager related objects
func (m *Manager) Close() {
	log.Info("Closing Manager")
	if err := m.Stop(); err != nil {
		log.Fatal("Unable to Close Manager", err)
	}
}

// Stop stops the manager
func (m *Manager) Stop() error {
	return nil
}
