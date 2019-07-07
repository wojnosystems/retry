package retry

import "testing"

func TestErrorList_Last(t *testing.T) {
	nothing := newErrorList()
	if nothing.Last() != nil {
		t.Error("expected an empty error list to return nil and not panic when asked for the last item")
	}
}
