// ABOUTME: stress_test.go - long-running stress tests for concurrent Cairo operations.
// ABOUTME: Detects memory leaks and validates stability under heavy concurrent load.

//go:build stress

package race_test

import (
	"fmt"
	"math"
	"runtime"
	"sync"
	"testing"

	cairo "github.com/mikowitz/cairo"
	"github.com/stretchr/testify/require"
)

const (
	stressGoroutines = 50
	stressIterations = 1000
	memLeakCycles    = 500
)

// TestStressConcurrentDrawing runs many goroutines doing many draw operations
// to detect races, deadlocks, and instability under sustained concurrent load.
func TestStressConcurrentDrawing(t *testing.T) {
	surf, err := cairo.NewImageSurface(cairo.FormatARGB32, 400, 400)
	require.NoError(t, err)
	defer surf.Close()

	ctx, err := cairo.NewContext(surf)
	require.NoError(t, err)
	defer ctx.Close()

	var wg sync.WaitGroup
	for i := 0; i < stressGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			r := float64(id) / float64(stressGoroutines)
			for j := 0; j < stressIterations; j++ {
				ctx.SetSourceRGBA(r, 0.5, 1.0-r, 0.8)
				ctx.Rectangle(10, 10, 100, 100)
				ctx.Fill()
				ctx.Arc(200, 200, 50, 0, 2*math.Pi)
				ctx.Stroke()
			}
		}(i)
	}
	wg.Wait()
}

// TestStressSurfaceCreation creates and destroys many surfaces concurrently
// to verify no resource leaks or crashes under heavy lifecycle churn.
func TestStressSurfaceCreation(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < stressGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < stressIterations; j++ {
				surf, err := cairo.NewImageSurface(cairo.FormatARGB32, 50, 50)
				require.NoError(t, err)
				surf.Flush()
				require.NoError(t, surf.Close())
			}
		}()
	}
	wg.Wait()
}

// TestStressMemoryLeak verifies that Go heap usage does not grow unboundedly
// after running many create/draw/close cycles with forced garbage collection.
func TestStressMemoryLeak(t *testing.T) {
	// Warm up to stabilize allocator state before measuring.
	for i := 0; i < 10; i++ {
		surf, err := cairo.NewImageSurface(cairo.FormatARGB32, 100, 100)
		require.NoError(t, err)
		ctx, err := cairo.NewContext(surf)
		require.NoError(t, err)
		ctx.Close()
		surf.Close()
	}

	runtime.GC()
	runtime.GC()
	var before runtime.MemStats
	runtime.ReadMemStats(&before)

	for i := 0; i < memLeakCycles; i++ {
		surf, err := cairo.NewImageSurface(cairo.FormatARGB32, 100, 100)
		require.NoError(t, err)
		ctx, err := cairo.NewContext(surf)
		require.NoError(t, err)
		ctx.SetSourceRGBA(0.5, 0.5, 0.5, 1.0)
		ctx.Rectangle(10, 10, 80, 80)
		ctx.Fill()
		ctx.Close()
		surf.Close()
	}

	runtime.GC()
	runtime.GC()
	var after runtime.MemStats
	runtime.ReadMemStats(&after)

	// Go heap in use should not grow by more than 10 MB after all resources are closed.
	const maxGrowthBytes uint64 = 10 * 1024 * 1024
	if after.HeapInuse > before.HeapInuse+maxGrowthBytes {
		t.Errorf("possible memory leak: heap grew from %d to %d bytes (%d MB growth)",
			before.HeapInuse, after.HeapInuse,
			(after.HeapInuse-before.HeapInuse)/(1024*1024))
	}
}

// TestStressGOMAXPROCS runs concurrent draw workloads with varying GOMAXPROCS
// settings to verify stability with different OS thread counts.
func TestStressGOMAXPROCS(t *testing.T) {
	for _, procs := range []int{1, 2, 4, 8} {
		procs := procs
		t.Run(fmt.Sprintf("GOMAXPROCS%d", procs), func(t *testing.T) {
			old := runtime.GOMAXPROCS(procs)
			defer runtime.GOMAXPROCS(old)

			surf, err := cairo.NewImageSurface(cairo.FormatARGB32, 200, 200)
			require.NoError(t, err)
			defer surf.Close()

			ctx, err := cairo.NewContext(surf)
			require.NoError(t, err)
			defer ctx.Close()

			var wg sync.WaitGroup
			const goroutines = 20
			const iterations = 100
			for i := 0; i < goroutines; i++ {
				wg.Add(1)
				go func(id int) {
					defer wg.Done()
					for j := 0; j < iterations; j++ {
						ctx.SetSourceRGB(float64(id)/goroutines, 0.5, 0.5)
						ctx.Rectangle(0, 0, 200, 200)
						ctx.Fill()
					}
				}(i)
			}
			wg.Wait()
		})
	}
}
