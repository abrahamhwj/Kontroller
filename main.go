// Package main provides the entry point to the Kontroller application.
package main

import (
	"Kontroller/controllers/cfgReloader"
	"Kontroller/logging"
	"Kontroller/pkg/api"
	"Kontroller/pkg/manager"
	"flag"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

var log *logging.Logging

func init() {
	log = logging.NewLogging("main")
}
func main() {
	// Define variables.
	var kubeconfig *string
	var err error
	config := &rest.Config{}
	// Set the kubeconfig file path.
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) abs path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "abs path to the kubeconfig file")
	}
	flag.Parse()
	// Get the kubernetes config.
	if config, err = rest.InClusterConfig(); err != nil {
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			log.Fatalf("get kube config failed")
			panic(err)
		}
	}
	// Create a new manager.
	mgr := manager.NewManager(config)
	// Register the cfgReloader controller.
	var _ api.Controller = (*cfgReloader.Reloader)(nil)
	r := cfgReloader.NewReloader("reloader")
	mgr.RegisController(r)
	// Run the controllers.
	stopper := make(chan struct{})
	defer close(stopper)
	mgr.RunControllers(stopper)
	// Wait for the stopper to close.
	select {
	case <-stopper:
		log.Infof("manager closed")
	}
}
