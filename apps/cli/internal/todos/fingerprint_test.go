package todos

import "testing"

func TestComputeID_Stability(t *testing.T) {
	id1 := computeID("main.go", "TODO", "implement this")
	id2 := computeID("main.go", "TODO", "implement this")
	if id1 != id2 {
		t.Errorf("expected stable ID, got %q and %q", id1, id2)
	}
}

func TestComputeID_Length(t *testing.T) {
	id := computeID("main.go", "TODO", "implement this")
	if len(id) != 12 {
		t.Errorf("expected 12 hex chars, got %d: %q", len(id), id)
	}
}

func TestComputeID_Uniqueness(t *testing.T) {
	id1 := computeID("main.go", "TODO", "first thing")
	id2 := computeID("main.go", "TODO", "second thing")
	id3 := computeID("other.go", "TODO", "first thing")
	id4 := computeID("main.go", "FIXME", "first thing")

	ids := []string{id1, id2, id3, id4}
	seen := make(map[string]bool)
	for _, id := range ids {
		if seen[id] {
			t.Errorf("duplicate ID: %s", id)
		}
		seen[id] = true
	}
}

func TestComputeID_DifferentMarkerSameText(t *testing.T) {
	id1 := computeID("f.go", "TODO", "do something")
	id2 := computeID("f.go", "FIXME", "do something")
	if id1 == id2 {
		t.Error("expected different IDs for different markers")
	}
}
