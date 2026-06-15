package money

import "testing"

func TestRoundToTwo(t *testing.T) {
	got := roundToTwo(100.2222)
	want := 100.22

	if got != want {
		t.Errorf("got %.4f want %.4f", got, want)
	}
}
