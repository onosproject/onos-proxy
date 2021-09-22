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
	"github.com/onosproject/onos-lib-go/pkg/logging"
)

var log = logging.GetLogger("e2", "v1beta1")

// Client is an E2 client
type Client interface {
	// Node returns a Node with the given NodeID
	Node(nodeID NodeID) Node
}

// NewClient creates a new E2 client
func NewClient(opts ...Option) Client {
	return &e2Client{
		opts: opts,
	}
}

// e2Client is the default E2 client implementation
type e2Client struct {
	opts []Option
}

func (c *e2Client) Node(nodeID NodeID) Node {
	return NewNode(nodeID, c.opts...)
}

var _ Client = &e2Client{}
