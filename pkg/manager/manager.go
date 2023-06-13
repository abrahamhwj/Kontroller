// Package manager provides functionality to manage controllers.
package manager

import (
	"Kontroller/config"
	"Kontroller/pkg/api"
	"fmt"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/rest"
)

// Manager represents a controller manager.
type Manager struct {
	Items  map[string]*ConcreteController
	Config *rest.Config
}

// NewManager creates a new instance of Manager.
func NewManager(config *rest.Config) *Manager {
	return &Manager{Config: config}
}

// RegisController registers a controller with the manager.
func (m *Manager) RegisController(controller api.Controller) {
	// Initialize the Items map if it is nil.
	if m.Items == nil {
		m.Items = make(map[string]*ConcreteController)
	}
	name := controller.ControllerName()
	// Check if the controller is already registered.
	if _, ok := m.Items[name]; ok {
		log.Infof("controller %s already registered\n", name)
		return
	}
	// Create a new ConcreteController and add it to the Items map.
	concreteController := NewConcreteControllerBuilder().Controller(controller).Queue().Client(m.Config).ListWatch().IndexerInformer().Build()
	m.Items[name] = concreteController
	log.Infof("controller %s registered successfully\n", name)
	return
}

// DeregisController deregisters a controller from the manager.
func (m *Manager) DeregisController(name string) {
	// Check if the Items map is nil.
	if m.Items == nil {
		log.Infof("no controller has registered")
		return
	}
	// Check if the controller is registered.
	if _, ok := m.Items[name]; !ok {
		log.Infof("controller %s not registered yet\n", name)
		return
	}
	// Remove the controller from the Items map.
	delete(m.Items, name)
	log.Infof("controller %s deregistered successfully.\n", name)
	return
}

// RunControllers runs all registered controllers.
func (m *Manager) RunControllers(stopper <-chan struct{}) {
	// Check if there are any controllers to run.
	if m.Items == nil {
		log.Infof("no controllers in manager to run")
		runtime.HandleError(fmt.Errorf("no controllers in manager to run"))
		return
	}
	// Run each controller in a separate goroutine.
	for _, concreteController := range m.Items {
		go concreteController.Run(stopper, int(config.Cfg.Manager.ThreadNumber))
	}
	return
}
