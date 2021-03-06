package v1alpha2

import (
	"github.com/lord/types/pkg/istio/apis/config/v1alpha2"
	"github.com/rancher/norman/lifecycle"
	"k8s.io/apimachinery/pkg/runtime"
)

type HandlerLifecycle interface {
	Create(obj *v1alpha2.Handler) (runtime.Object, error)
	Remove(obj *v1alpha2.Handler) (runtime.Object, error)
	Updated(obj *v1alpha2.Handler) (runtime.Object, error)
}

type handlerLifecycleAdapter struct {
	lifecycle HandlerLifecycle
}

func (w *handlerLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *handlerLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *handlerLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*v1alpha2.Handler))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *handlerLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*v1alpha2.Handler))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *handlerLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*v1alpha2.Handler))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewHandlerLifecycleAdapter(name string, clusterScoped bool, client HandlerInterface, l HandlerLifecycle) HandlerHandlerFunc {
	adapter := &handlerLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v1alpha2.Handler) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
