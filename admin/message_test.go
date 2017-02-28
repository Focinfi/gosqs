package admin

import "testing"

func TestMessageIndex(t *testing.T) {
	idx1 := messageIndex()
	idx2 := messageIndex()

	if idx1 == idx2 {
		t.Fatal("idx1 can not be idx2")
	}
	t.Log(idx1, idx2)
}
