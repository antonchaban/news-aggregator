package predicates

import (
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

func TestSourcePredicates(t *testing.T) {
	pred := Source()

	t.Run("CreateFunc", func(t *testing.T) {
		obj := &unstructured.Unstructured{}
		obj.SetGroupVersionKind(schema.GroupVersionKind{
			Group:   "com.teamdev",
			Version: "v1",
			Kind:    "Source",
		})

		createEvent := event.CreateEvent{
			Object: obj,
		}

		if !pred.Create(createEvent) {
			t.Errorf("CreateFunc should return true for any create event")
		}
	})

	t.Run("DeleteFunc", func(t *testing.T) {
		obj := &unstructured.Unstructured{}
		obj.SetGroupVersionKind(schema.GroupVersionKind{
			Group:   "com.teamdev",
			Version: "v1",
			Kind:    "Source",
		})

		deleteEventUnknown := event.DeleteEvent{
			Object:             obj,
			DeleteStateUnknown: true,
		}

		if pred.Delete(deleteEventUnknown) {
			t.Errorf("DeleteFunc should return false when DeleteStateUnknown is true")
		}

		deleteEventKnown := event.DeleteEvent{
			Object:             obj,
			DeleteStateUnknown: false,
		}

		if !pred.Delete(deleteEventKnown) {
			t.Errorf("DeleteFunc should return true when DeleteStateUnknown is false")
		}
	})

	t.Run("UpdateFunc", func(t *testing.T) {
		objOld := &unstructured.Unstructured{}
		objOld.SetGroupVersionKind(schema.GroupVersionKind{
			Group:   "com.teamdev",
			Version: "v1",
			Kind:    "Source",
		})
		objOld.SetGeneration(1)

		objNew := objOld.DeepCopy()
		objNew.SetGeneration(1) // Same generation as objOld

		updateEventSameGen := event.UpdateEvent{
			ObjectOld: objOld,
			ObjectNew: objNew,
		}

		if pred.Update(updateEventSameGen) {
			t.Errorf("UpdateFunc should return false when generations are the same")
		}

		objNew.SetGeneration(2)

		updateEventNewGen := event.UpdateEvent{
			ObjectOld: objOld,
			ObjectNew: objNew,
		}

		if !pred.Update(updateEventNewGen) {
			t.Errorf("UpdateFunc should return true when generations are different")
		}
	})
}
