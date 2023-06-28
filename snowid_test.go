package snowid

import (
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"
)

var (
	testEpoch = time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
)

func TestConstants(t *testing.T) {
	require.Equal(t, Uint41Mask, Uint40Mask|Uint40Bit)
}

type testClock struct {
	since func(t time.Time) time.Duration
	sleep func(d time.Duration)
}

func (tc *testClock) Since(t time.Time) time.Duration {
	return tc.since(t)
}

func (tc *testClock) Sleep(d time.Duration) {
	tc.sleep(d)
}

func TestNew(t *testing.T) {
	t.Run("error-on-invalid-epoch", func(t *testing.T) {
		_, err := New(Options{})
		require.Equal(t, ErrInvalidEpoch, err)
	})

	t.Run("error-on-invalid-id", func(t *testing.T) {
		_, err := New(Options{
			Epoch: time.Now(),
			ID:    Uint41Mask,
		})
		require.Equal(t, ErrInvalidID, err)
	})

	t.Run("default", func(t *testing.T) {
		g, err := New(Options{
			Epoch: testEpoch,
		})
		require.NoError(t, err)
		defer g.Stop()

		require.NotEqualValues(t, 0, g.NewID())
	})
}

func TestGenerator(t *testing.T) {
	t.Run("error-on-stopped", func(t *testing.T) {
		g, err := New(Options{
			Epoch: testEpoch,
		})
		require.NoError(t, err)
		require.PanicsWithError(t, ErrStopped.Error(), func() {
			g.Stop()
			g.NewID()
		})
	})

	t.Run("standard-and-with-overflow", func(t *testing.T) {
		var sleepInvoked int64

		customGrain := time.Duration(rand.Intn(20)+10) * time.Millisecond

		tc := &testClock{
			since: func(t time.Time) time.Duration {
				return 11*customGrain + time.Duration(sleepInvoked)*customGrain
			},
			sleep: func(d time.Duration) {
				require.Equal(t, customGrain, d)
				sleepInvoked++
			},
		}

		g, err := New(Options{
			Epoch: testEpoch,
			ID:    0b111111,
			Grain: customGrain,
			Clock: tc,
		})
		require.NoError(t, err)
		defer g.Stop()

		for i := uint64(0); i < 4096; i++ {
			require.Equal(t, 11<<22|0b111111<<12|i, g.NewID())
		}
		for i := uint64(0); i < 4096; i++ {
			require.Equal(t, 12<<22|0b111111<<12|i, g.NewID())
		}
		require.Equal(t, int64(1), sleepInvoked)
		require.Equal(t, uint64(4096*2), g.Count())
	})

	t.Run("standard-and-with-overflow-leading-bit", func(t *testing.T) {
		var sleepInvoked int64

		customGrain := time.Duration(rand.Intn(20)+10) * time.Millisecond

		tc := &testClock{
			since: func(t time.Time) time.Duration {
				return 11*customGrain + time.Duration(sleepInvoked)*customGrain
			},
			sleep: func(d time.Duration) {
				require.Equal(t, customGrain, d)
				sleepInvoked++
			},
		}

		g, err := New(Options{
			Epoch:      testEpoch,
			ID:         0b111111,
			Grain:      customGrain,
			Clock:      tc,
			LeadingBit: true,
		})
		require.NoError(t, err)
		defer g.Stop()

		for i := uint64(0); i < 4096; i++ {
			require.Equal(t, 1<<62|11<<22|0b111111<<12|i, g.NewID())
		}
		for i := uint64(0); i < 4096; i++ {
			require.Equal(t, 1<<62|12<<22|0b111111<<12|i, g.NewID())
		}
		require.Equal(t, int64(1), sleepInvoked)
		require.Equal(t, uint64(4096*2), g.Count())
	})

	t.Run("samples", func(t *testing.T) {
		g, err := New(Options{
			Epoch: testEpoch,
			ID:    0b1010101,
		})
		require.NoError(t, err)
		defer g.Stop()

		for i := 0; i < 5; i++ {
			t.Log("ID:", g.NewID())
		}
	})

	t.Run("samples-leading-bit", func(t *testing.T) {
		g, err := New(Options{
			Epoch:      testEpoch,
			ID:         0b1010101,
			LeadingBit: true,
		})
		require.NoError(t, err)
		defer g.Stop()

		for i := 0; i < 5; i++ {
			t.Log("ID:", g.NewID())
		}
	})

}

func BenchmarkGenerator_NewID(b *testing.B) {
	g, err := New(Options{
		Epoch:      testEpoch,
		ID:         0b1010101,
		LeadingBit: true,
	})
	require.NoError(b, err)
	defer g.Stop()

	for n := 0; n < b.N; n++ {
		g.NewID()
	}
}
