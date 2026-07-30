package typingstate

import "testing"

func TestInitialAndSentenceCapitalization(t *testing.T) {
	c := New(true, true, true)
	if got := c.BeginWord(); !got.UpperFirst || got.AllUpper {
		t.Fatalf("initial word case = %+v", got)
	}
	if got := c.BeginWord(); got.UpperFirst || got.AllUpper {
		t.Fatalf("second word should be normal = %+v", got)
	}
	c.ObserveDelimiter(".")
	if got := c.BeginWord(); !got.UpperFirst || got.AllUpper {
		t.Fatalf("after period = %+v", got)
	}
}

func TestSingleAndDoubleShift(t *testing.T) {
	c := New(false, true, true)
	c.ShiftDown()
	if !c.ShiftUp(1000) {
		t.Fatal("single shift must change mode")
	}
	if got := c.BeginWord(); !got.UpperFirst || got.AllUpper {
		t.Fatalf("single shift = %+v", got)
	}
	c.ShiftDown()
	c.ShiftUp(2000)
	c.ShiftDown()
	c.ShiftUp(2250)
	if got := c.BeginWord(); !got.UpperFirst || !got.AllUpper {
		t.Fatalf("double shift = %+v", got)
	}
	c.ShiftDown()
	c.ShiftUp(3000)
	if c.StickyUpper() {
		t.Fatal("single clean shift must turn sticky uppercase off")
	}
}

func TestShiftHeldWithLetterIsNotTap(t *testing.T) {
	c := New(false, true, true)
	c.ShiftDown()
	c.MarkShiftUsed()
	if c.ShiftUp(1000) {
		t.Fatal("held Shift used with another key must not trigger tap gesture")
	}
	if got := c.BeginWord(); got.UpperFirst || got.AllUpper {
		t.Fatalf("unexpected state after held Shift = %+v", got)
	}
}

func TestSnapshotRestore(t *testing.T) {
	c := New(true, true, true)
	_ = c.BeginWord()
	before := c.Snapshot()
	c.ObserveDelimiter(".")
	c.Restore(before)
	if got := c.BeginWord(); got.UpperFirst {
		t.Fatalf("restored pre-period state should not auto-cap: %+v", got)
	}
}

func BenchmarkCapitalizerWordCycle(b *testing.B) {
	c := New(true, true, true)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = c.BeginWord()
		c.ObserveDelimiter(".")
	}
