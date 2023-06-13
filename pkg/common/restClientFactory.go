package common

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

// Constants for Kubernetes resource names
const (
	ConfigMaps             = "configmaps"
	Endpoints              = "endpoints"
	Events                 = "events"
	LimitRanges            = "limitranges"
	Namespaces             = "namespaces"
	Nodes                  = "nodes"
	PersistentVolumes      = "persistentvolumes"
	PersistentVolumeClaims = "persistentvolumeclaims"
	Pods                   = "pods"
	Secrets                = "secrets"
	Services               = "services"
	ServiceAccounts        = "serviceaccounts"
	CronJobs               = "cronjobs"
	DaemonSets             = "daemonsets"
	Deployments            = "deployments"
	StatefulSets           = "statefulsets"
	Ingresses              = "ingresses"
	NetworkPolicies        = "networkpolicies"
	RoleBindings           = "rolebindings"
	Roles                  = "roles"
	ClusterRoles           = "clusterroles"
	ClusterRolebindings    = "clusterrolebindings"
)

// CacheGetterMap maps resource names to functions that return a cache getter for that resource.
var CacheGetterMap = map[string]func(clientset *kubernetes.Clientset) cache.Getter{
	ConfigMaps:             getCoreV1RESTClient,
	Endpoints:              getCoreV1RESTClient,
	Events:                 getCoreV1RESTClient,
	LimitRanges:            getCoreV1RESTClient,
	Namespaces:             getCoreV1RESTClient,
	Nodes:                  getCoreV1RESTClient,
	PersistentVolumes:      getCoreV1RESTClient,
	PersistentVolumeClaims: getCoreV1RESTClient,
	Pods:                   getCoreV1RESTClient,
	Secrets:                getCoreV1RESTClient,
	Services:               getCoreV1RESTClient,
	ServiceAccounts:        getCoreV1RESTClient,
	CronJobs:               getBatchV1beta1RESTClient,
	DaemonSets:             getAppsV1RESTClient,
	Deployments:            getAppsV1RESTClient,
	StatefulSets:           getAppsV1RESTClient,
	Ingresses:              getExtensionsV1beta1RESTClient,
	NetworkPolicies:        getNetworkingV1RESTClient,
	RoleBindings:           getRbacV1RESTClient,
	Roles:                  getRbacV1RESTClient,
	ClusterRoles:           getRbacV1RESTClient,
	ClusterRolebindings:    getRbacV1RESTClient,
}

// getCoreV1RESTClient returns a cache getter for the CoreV1 API group.
func getCoreV1RESTClient(clientset *kubernetes.Clientset) cache.Getter {
	return clientset.CoreV1().RESTClient()
}

// getAppsV1RESTClient returns a cache getter for the AppsV1 API group.
func getAppsV1RESTClient(clientset *kubernetes.Clientset) cache.Getter {
	return clientset.AppsV1().RESTClient()
}

// getBatchV1beta1RESTClient returns a cache getter for the BatchV1beta1 API group.
func getBatchV1beta1RESTClient(clientset *kubernetes.Clientset) cache.Getter {
	return clientset.BatchV1beta1().RESTClient()
}

// getExtensionsV1beta1RESTClient returns a cache getter for the ExtensionsV1beta1 API group.
func getExtensionsV1beta1RESTClient(clientset *kubernetes.Clientset) cache.Getter {
	return clientset.ExtensionsV1beta1().RESTClient()
}

// getNetworkingV1RESTClient returns a cache getter for the NetworkingV1 API group.
func getNetworkingV1RESTClient(clientset *kubernetes.Clientset) cache.Getter {
	return clientset.NetworkingV1().RESTClient()
}

// getRbacV1RESTClient returns a cache getter for the RbacV1 API group.
func getRbacV1RESTClient(clientset *kubernetes.Clientset) cache.Getter {
	return clientset.RbacV1().RESTClient()
}
