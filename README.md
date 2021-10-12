# onos-proxy
ONOS side-car proxy for various subsystems, e.g. E2T

The main purpose of the sidecar proxy is to absorb the complexity of interfacing with various 
µONOS subsystems. This allows relatively easy implementations of the SDK in various languages, without
re-implementing the complex algorithms for each language.

Presently, the proxy implements only E2T service load-balancing and routing, but in future may be extended to accommodate sophisticated interactions with other
parts of teh µONOS platform.

## Deployment
The proxy is intended to be deployed as a sidecar container as part of an application pod. Such deployment
can be arranged explicitly by including the proxy container details in the application Helm chart, but an
easier way is to include the following metadata annotation as part of the `deployment.yaml` file.

```bigquery
annotations:
  proxy.onosproject.org/inject: "true"
```

This annotation will be detected by the `onos-app-operator` via its admission hook which will augment the 
deployment descriptor to include the proxy container as part of the application pod automatically.

## E2 Services
The proxy container exposes a locally accessible port on `localhost:5151` where it hosts the following services:

* E2 Control Service - allows issuing control requests to E2 nodes
* E2 Subscription Service - allows issuing subscribe and unsubscribe requests to E2 nodes

The E2 proxy tracks the E2T and E2 node mastership state via `onos-topo` information and appropriately forwards 
gRPC requests to the E2T instance which is presently the master for the given target E2 node. The target E2 node
ID is extracted from the E2AP request headers

The mastership information is derived from the `MastershipState` aspect of the E2 node topology entities and
from the `controls` topology relations setup between the E2T and E2 node topology entities. 

The proxy does not manipulate the messages passed between the application and the E2T instances in any manner.

## SDK Versions

The `onos-ric-sdk-go` version `0.7.30` or greater and `onos-ric-sdk-py` version `0.1.6` or greater expect
the sidecar proxy to be deployed to work correctly.