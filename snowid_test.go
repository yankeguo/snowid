package snowid

import (
	"testing"
	"time"
)

var (
	testInstanceID = uint64(1) | uint64(1)<<9
	testStartTime  = time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)
)

type testClockForSleep struct {
	defaultClock
}

func (tcfs testClockForSleep) Since(t time.Time) time.Duration {
	return tcfs.defaultClock.Since(t) / 100
}

func (tcfs testClockForSleep) Sleep(d time.Duration) {
	tcfs.defaultClock.Sleep(d)
}

func TestNew_Sleep(t *testing.T) {
	s, _ := New(Options{
		ID:    0b11111111,
		Epoch: testStartTime,
		Clock: testClockForSleep{},
	})
	defer s.Stop()
	for i := 0; i < 20000; i++ {
		_ = s.NewID()
	}
}

func TestNew_BadID_BadEpoch(t *testing.T) {
	_, err := New(Options{})
	if err == nil {
		t.Error("should failed")
	}
	_, err = New(Options{
		Epoch: testStartTime,
		ID:    0b11111111111,
	})
	if err == nil {
		t.Error("should failed")
	}
}

func BenchmarkGenerator_NewID(b *testing.B) {
	s, _ := New(Options{
		Epoch: testStartTime,
		ID:    testInstanceID,
	})
	defer s.Stop()
	for n := 0; n < b.N; n++ {
		s.NewID()
	}
}

func TestGenerator_NewID(t *testing.T) {
	s, _ := New(Options{
		Epoch: testStartTime,
		ID:    testInstanceID,
	})
	defer s.Stop()
	var id uint64
	for i := 0; i < 10; i++ {
		id = s.NewID()
	}
	if s.Count() != 10 {
		t.Fatal("bad number of count")
	}
	t.Logf("ins: %b, seq: %b, mask: %b, id: %b", testInstanceID, id&Uint12Mask, Uint12Mask, id)
	t.Logf("ins: %x, seq: %x, mask: %x, id: %x", testInstanceID, id&Uint12Mask, Uint12Mask, id)
	t.Logf("ins: %d, seq: %d, mask: %d, id: %d", testInstanceID, id&Uint12Mask, Uint12Mask, id)
	if id&Uint12Mask != 9 {
		t.Fatal("bad sequence id")
	}
	if (id>>12)&Uint10Mask != testInstanceID {
		t.Fatal("bad instance id")
	}
	if time.Since(testStartTime)/time.Second != time.Duration(id>>22)*time.Millisecond/time.Second {
		t.Fatal("bad timestamp")
	}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		} else {
			t.Logf("%v", r)
		}
	}()
	s2, _ := New(Options{
		Epoch: testStartTime,
		ID:    testInstanceID,
	})
	s2.Stop()
	s2.NewID()
}

func TestGenerator_CheckDup(t *testing.T) {
	s, _ := New(Options{
		Epoch: testStartTime,
		ID:    testInstanceID,
	})
	defer s.Stop()

	out := map[uint64]bool{}

	for i := 0; i < 100000; i++ {
		id := s.NewID()
		if out[id] {
			t.Fatal("duplicated")
		}
		out[id] = true
	}
}
