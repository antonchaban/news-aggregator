package predicates

import (
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

func TestSourcePredicatesWithNamespace(t *testing.T) {
	namespace := "test-namespace"
	pred := Source(namespace)

	t.Run("CreateFunc Namespace Match", func(t *testing.T) {
		obj := &unstructured.Unstructured{}
		obj.SetNamespace(namespace)
		createEvent := event.CreateEvent{
			Object: obj,
		}

		if !pred.Create(createEvent) {
			t.Errorf("CreateFunc should return true when namespace matches")
		}
	})

	t.Run("CreateFunc Namespace Mismatch", func(t *testing.T) {
		obj := &unstructured.Unstructured{}
		obj.SetNamespace("other-namespace")
		createEvent := event.CreateEvent{
			Object: obj,
		}

		if pred.Create(createEvent) {
			t.Errorf("CreateFunc should return false when namespace does not match")
		}
	})

	t.Run("DeleteFunc Namespace Match", func(t *testing.T) {
		obj := &unstructured.Unstructured{}
		obj.SetNamespace(namespace)
		deleteEvent := event.DeleteEvent{
			Object:             obj,
			DeleteStateUnknown: false,
		}

		if !pred.Delete(deleteEvent) {
			t.Errorf("DeleteFunc should return true when namespace matches and DeleteStateUnknown is false")
		}
	})

	t.Run("DeleteFunc Namespace Mismatch", func(t *testing.T) {
		obj := &unstructured.Unstructured{}
		obj.SetNamespace("other-namespace")
		deleteEvent := event.DeleteEvent{
			Object:             obj,
			DeleteStateUnknown: false,
		}

		if pred.Delete(deleteEvent) {
			t.Errorf("DeleteFunc should return false when namespace does not match")
		}
	})

	t.Run("DeleteFunc DeleteStateUnknown", func(t *testing.T) {
		obj := &unstructured.Unstructured{}
		obj.SetNamespace(namespace)
		deleteEvent := event.DeleteEvent{
			Object:             obj,
			DeleteStateUnknown: true,
		}

		if pred.Delete(deleteEvent) {
			t.Errorf("DeleteFunc should return false when DeleteStateUnknown is true")
		}
	})

	t.Run("UpdateFunc Namespace Match Generation Changed", func(t *testing.T) {
		objOld := &unstructured.Unstructured{}
		objOld.SetNamespace(namespace)
		objOld.SetGeneration(1)

		objNew := objOld.DeepCopy()
		objNew.SetGeneration(2)

		updateEvent := event.UpdateEvent{
			ObjectOld: objOld,
			ObjectNew: objNew,
		}

		if !pred.Update(updateEvent) {
			t.Errorf("UpdateFunc should return true when namespace matches and generation has changed")
		}
	})

	t.Run("UpdateFunc Namespace Mismatch", func(t *testing.T) {
		objOld := &unstructured.Unstructured{}
		objOld.SetNamespace("other-namespace")
		objOld.SetGeneration(1)

		objNew := objOld.DeepCopy()
		objNew.SetGeneration(2)

		updateEvent := event.UpdateEvent{
			ObjectOld: objOld,
			ObjectNew: objNew,
		}

		if pred.Update(updateEvent) {
			t.Errorf("UpdateFunc should return false when namespace does not match")
		}
	})

	t.Run("UpdateFunc Generation Unchanged", func(t *testing.T) {
		objOld := &unstructured.Unstructured{}
		objOld.SetNamespace(namespace)
		objOld.SetGeneration(1)

		objNew := objOld.DeepCopy()
		objNew.SetGeneration(1)

		updateEvent := event.UpdateEvent{
			ObjectOld: objOld,
			ObjectNew: objNew,
		}

		if pred.Update(updateEvent) {
			t.Errorf("UpdateFunc should return false when generation has not changed")
		}
	})

	t.Run("GenericFunc Always True", func(t *testing.T) {
		obj := &unstructured.Unstructured{}
		obj.SetNamespace(namespace)
		genericEvent := event.GenericEvent{
			Object: obj,
		}

		if !pred.Generic(genericEvent) {
			t.Errorf("GenericFunc should always return true")
		}
	})
}
