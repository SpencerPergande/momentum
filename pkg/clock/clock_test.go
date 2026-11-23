package clock

import (
	"testing"
	"time"
)

func TestFakeSetAndAdvance(t *testing.T) {
	start := time.Date(2026, 11, 22, 8, 0, 0, 0, time.UTC)
	c := NewFake(start)

	if got := c.Now(); !got.Equal(start) {
		t.Fatalf("Now() = %v, want %v", got, start)
	}

	got := c.Advance(25 * time.Hour)
	want := start.Add(25 * time.Hour)
	if !got.Equal(want) || !c.Now().Equal(want) {
		t.Fatalf("after Advance, Now() = %v, want %v", c.Now(), want)
	}

	later := time.Date(2027, 1, 2, 0, 0, 0, 0, time.UTC)
	c.Set(later)
	if got := c.Now(); !got.Equal(later) {
		t.Fatalf("after Set, Now() = %v, want %v", got, later)
	}
}

func TestRealIsCurrent(t *testing.T) {
	before := time.Now()
	got := Real{}.Now()
	if got.Before(before) || got.Sub(before) > time.Second {
		t.Fatalf("Real.Now() = %v, not close to %v", got, before)
	}
}

func TestClockInterface(t *testing.T) {
	var _ Clock = Real{}
	var _ Clock = (*Fake)(nil)
}
