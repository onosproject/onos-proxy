// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-proxy/pkg/manager"
)

var log = logging.GetLogger()

func printVersion() {
	log.Info(fmt.Sprintf("Go Version: %s", runtime.Version()))
	log.Info(fmt.Sprintf("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH))
}

func main() {
	caPath := flag.String("caPath", "", "path to CA certificate")
	keyPath := flag.String("keyPath", "", "path to client private key")
	certPath := flag.String("certPath", "", "path to client certificate")

	//logf.SetLogger(zap.New())
	printVersion()

	log.Info("Starting onos-proxy")
	cfg := manager.Config{
		CAPath:   *caPath,
		KeyPath:  *keyPath,
		CertPath: *certPath,
		GRPCPort: 5151,
	}
	mgr := manager.NewManager(cfg)
	mgr.Run()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	mgr.Close()
}
