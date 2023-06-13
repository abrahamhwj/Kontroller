package manager

import (
	"Kontroller/config"
	"Kontroller/pkg/api"
	"Kontroller/pkg/common"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type ConcreteController struct {
	Controller api.Controller
	Queue      workqueue.RateLimitingInterface
	Client     *kubernetes.Clientset
	ListWatch  *cache.ListWatch
	Indexer    cache.Indexer
	Informer   cache.Controller
}

// AddEventHandlerFunc returns a function that adds an object to the queue
func (c *ConcreteController) AddEventHandlerFunc(queue workqueue.RateLimitingInterface) func(obj interface{}) {
	return func(obj interface{}) {
		key, err := cache.MetaNamespaceKeyFunc(obj)
		if err == nil {
			queue.Add(key)
		}
	}
}

// UpdateEventHandlerFunc returns a function that updates an object in the queue
func (c *ConcreteController) UpdateEventHandlerFunc(queue workqueue.RateLimitingInterface) func(oldObj interface{}, newObj interface{}) {
	return func(oldObj interface{}, newObj interface{}) {
		key, err := cache.MetaNamespaceKeyFunc(newObj)
		if err == nil {
			queue.Add(key)
		}
	}
}

// DeleteEventHandlerFunc returns a function that deletes an object from the queue
func (c *ConcreteController) DeleteEventHandlerFunc(queue workqueue.RateLimitingInterface) func(obj interface{}) {
	return func(obj interface{}) {
		key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
		if err == nil {
			queue.Add(key)
		}
	}
}

// Run starts the controller with the specified number of threads and stopper channel
func (c *ConcreteController) Run(stopper <-chan struct{}, threads int) {
	defer runtime.HandleCrash()
	defer c.Queue.ShuttingDown()
	name := c.Controller.ControllerName()
	log.Infof("start controller: %s\n", name)
	go c.Informer.Run(stopper)
	if !cache.WaitForCacheSync(stopper, c.Informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("time out wait for cache to sync of controller:%s\n, please check if the controller run property\n", name))
		return
	}
	for i := 0; i < threads; i++ {
		go wait.Until(func() {
			for c.ProcessNextItem() {
			}
		}, config.Cfg.Manager.ThreadTimeout, stopper)
	}
	<-stopper
	log.Infof("stop controller: %s\n", name)
}

// processNextItem processes the next item in the queue
func (c *ConcreteController) ProcessNextItem() bool {
	key, shutdown := c.Queue.Get()
	if shutdown {
		return false
	}
	defer c.Queue.Done(key)
	name := c.Controller.ControllerName()
	obj, exists, err := c.Indexer.GetByKey(key.(string))
	if err != nil {
		log.Errorf("controller %s fetching object %s from local cache failed with err: %v\n, internal err, you may need restart\n", name, key, err)
	}
	if !exists {
		log.Infof("controller %s object %s has been deleted\n", name, key)
	}
	handleErr := c.Controller.HandleObject(c.Client, obj)
	if handleErr != nil {
		if c.Queue.NumRequeues(key) < int(config.Cfg.Manager.ControllerMaxRetryTimes) {
			c.Queue.AddRateLimited(key)
			log.Errorf("controller %s handle obj %s failed %d times with err:%v\n", name, key, err)
			return false
		}
		log.Errorf("controller %s handle obj %s failed finally: %v\n", name, key, err)
		return false
	}
	return true
}

// build ConcreteController with fluentApi style
type (
	ConcreteControllerBuilder struct {
		ConcreteController *ConcreteController
	}
	ControllerBuilder interface {
		Controller(controller api.Controller) QueueBuilder
	}
	QueueBuilder interface {
		Queue() ClientBuilder
	}
	ClientBuilder interface {
		Client(config *rest.Config) ListWatchBuilder
	}
	ListWatchBuilder interface {
		ListWatch() IndexerInformerBuilder
	}
	IndexerInformerBuilder interface {
		IndexerInformer() EndBuilder
	}
	EndBuilder interface {
		Build() *ConcreteController
	}
)

func (c *ConcreteControllerBuilder) Controller(controller api.Controller) QueueBuilder {
	if c.ConcreteController == nil {
		c.ConcreteController = &ConcreteController{}
	}
	c.ConcreteController.Controller = controller
	return c
}
func (c *ConcreteControllerBuilder) Queue() ClientBuilder {
	c.ConcreteController.Queue = workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	return c
}
func (c *ConcreteControllerBuilder) Client(config *rest.Config) ListWatchBuilder {
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("NewForConfig err with:%s\n, please check if the config file is right\n", err)
		panic(err)
	}
	c.ConcreteController.Client = client
	return c
}
func (c *ConcreteControllerBuilder) ListWatch() IndexerInformerBuilder {
	resourceName := c.ConcreteController.Controller.ControlResourceName()
	restClient, ok := common.CacheGetterMap[resourceName]
	if !ok {
		log.Fatalf("rest client get err, please check if the resource name incorrect. or check CacheGetterMap factory, if the resource supported.\n")
		panic(fmt.Errorf("resource not supported: %s\n", resourceName))
	}
	listOptions := func(options *metav1.ListOptions) {
		options.LabelSelector = c.ConcreteController.Controller.ControlLabelSelector()
	}
	getter := restClient(c.ConcreteController.Client)
	namespace := c.ConcreteController.Controller.ControlNamespace()
	c.ConcreteController.ListWatch = cache.NewFilteredListWatchFromClient(getter, resourceName, namespace, listOptions)
	return c
}
func (c *ConcreteControllerBuilder) IndexerInformer() EndBuilder {
	queue := c.ConcreteController.Queue
	controller := c.ConcreteController.Controller
	addFunc := c.ConcreteController.AddEventHandlerFunc(queue)
	updateFunc := c.ConcreteController.UpdateEventHandlerFunc(queue)
	deleteFunc := c.ConcreteController.DeleteEventHandlerFunc(queue)
	if _, ok := interface{}(controller).(api.EventHandler); ok {
		addFunc = interface{}(controller).(api.EventHandler).AddEventHandlerFunc(queue)
		updateFunc = interface{}(controller).(api.EventHandler).UpdateEventHandlerFunc(queue)
		deleteFunc = interface{}(controller).(api.EventHandler).DeleteEventHandlerFunc(queue)
	}
	lw := c.ConcreteController.ListWatch
	obj := c.ConcreteController.Controller.ControlObject()
	Period := config.Cfg.Manager.ReSyncPeriod
	indexer, informer := cache.NewIndexerInformer(lw, obj, Period, cache.ResourceEventHandlerFuncs{
		AddFunc:    addFunc,
		UpdateFunc: updateFunc,
		DeleteFunc: deleteFunc,
	}, cache.Indexers{})
	c.ConcreteController.Indexer = indexer
	c.ConcreteController.Informer = informer
	return c
}
func (c *ConcreteControllerBuilder) Build() *ConcreteController {
	return c.ConcreteController
}
func NewConcreteControllerBuilder() ControllerBuilder {
	return &ConcreteControllerBuilder{}
}
