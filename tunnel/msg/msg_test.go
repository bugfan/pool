package msg

import "testing"

func TestTypeMap(t *testing.T) {
	for name, typ := range TypeMap {
		if name != typ.Name() {
			t.Fatalf("Type error %s is not %s\n", typ.Name(), name)
		}
	}
}
