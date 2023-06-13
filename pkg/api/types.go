// Package api provides interfaces for Kubernetes controllers and event handlers
package api

// Import required packages
import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/workqueue"
)

// Controller is an interface for Kubernetes controllers
type Controller interface {
	ControllerName() string
	ControlObject() runtime.Object
	ControlResourceName() string
	ControlNamespace() string
	ControlLabelSelector() string
	HandleObject(client *kubernetes.Clientset, object interface{}) error
}

// EventHandler is an interface for Kubernetes event handlers
type EventHandler interface {
	AddEventHandlerFunc(queue workqueue.RateLimitingInterface) func(obj interface{})
	UpdateEventHandlerFunc(queue workqueue.RateLimitingInterface) func(oldObj interface{}, newObj interface{})
	DeleteEventHandlerFunc(queue workqueue.RateLimitingInterface) func(obj interface{})
}
