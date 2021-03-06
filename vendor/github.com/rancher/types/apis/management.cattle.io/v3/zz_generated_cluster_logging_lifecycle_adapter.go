package v3

import (
	"github.com/rancher/norman/lifecycle"
	"k8s.io/apimachinery/pkg/runtime"
)

type ClusterLoggingLifecycle interface {
	Create(obj *ClusterLogging) (runtime.Object, error)
	Remove(obj *ClusterLogging) (runtime.Object, error)
	Updated(obj *ClusterLogging) (runtime.Object, error)
}

type clusterLoggingLifecycleAdapter struct {
	lifecycle ClusterLoggingLifecycle
}

func (w *clusterLoggingLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*ClusterLogging))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *clusterLoggingLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*ClusterLogging))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *clusterLoggingLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*ClusterLogging))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewClusterLoggingLifecycleAdapter(name string, clusterScoped bool, client ClusterLoggingInterface, l ClusterLoggingLifecycle) ClusterLoggingHandlerFunc {
	adapter := &clusterLoggingLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *ClusterLogging) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
