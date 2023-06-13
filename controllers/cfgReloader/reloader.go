package cfgReloader

import (
	"Kontroller/logging"
	"Kontroller/pkg/common"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"strings"
)

var log *logging.Logging

const (
	ReloaderName         = "reloader"
	DefaultLabelSelector = "kontroller/reloader=true"
)

func init() {
	log = logging.NewLogging(ReloaderName)
}

type Reloader struct {
	Name          string
	Resource      string
	Object        runtime.Object
	Namespace     string
	LabelSelector string
}

func (r *Reloader) ControllerName() string {
	return r.Name
}

func (r *Reloader) ControlObject() runtime.Object {
	return r.Object
}

func (r *Reloader) ControlResourceName() string {
	return r.Resource
}

func (r *Reloader) ControlNamespace() string {
	return r.Namespace
}

func (r *Reloader) ControlLabelSelector() string {
	return r.LabelSelector
}

func (r *Reloader) HandleObject(client *kubernetes.Clientset, object interface{}) error {
	configmaps := object.(*corev1.ConfigMap)
	log.Infof("configmap %s/%s changed", configmaps.Namespace, configmaps.Name)
	//TODO: your logic here
	return nil
}

type Option func(reloader *Reloader)

func Namespace(n string) Option {
	return func(reloader *Reloader) {
		reloader.Namespace = n
	}
}

func LabelSelector(l string) Option {
	return func(reloader *Reloader) {
		reloader.LabelSelector = strings.ReplaceAll(strings.ToLower(l), " ", "")
	}
}

func NewReloader(name string, options ...Option) *Reloader {
	r := &Reloader{
		Name:          name,
		Resource:      common.ConfigMaps,
		Object:        &corev1.ConfigMap{},
		Namespace:     corev1.NamespaceAll,
		LabelSelector: DefaultLabelSelector,
	}
	for _, option := range options {
		option(r)
	}
	return r
}
