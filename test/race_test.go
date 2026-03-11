// ABOUTME: race_test.go - concurrency and thread safety tests for Cairo types.
// ABOUTME: Verifies Context, Surface, and Pattern operations are safe under concurrent access.
package race_test

import (
	"math"
	"sync"
	"testing"

	cairo "github.com/mikowitz/cairo"
	"github.com/stretchr/testify/require"
)

const numGoroutines = 10
const numOps = 20

// TestConcurrentContextOperations verifies that many goroutines can safely
// draw to the same context concurrently without data races or deadlocks.
func TestConcurrentContextOperations(t *testing.T) {
	surf, err := cairo.NewImageSurface(cairo.FormatARGB32, 200, 200)
	require.NoError(t, err)
	defer surf.Close()

	ctx, err := cairo.NewContext(surf)
	require.NoError(t, err)
	defer ctx.Close()

	var wg sync.WaitGroup
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			r := float64(id) / float64(numGoroutines)
			for j := 0; j < numOps; j++ {
				ctx.SetSourceRGBA(r, 0.5, 1.0-r, 0.8)
				ctx.Rectangle(10, 10, 100, 100)
				ctx.Fill()
				ctx.Arc(100, 100, 40, 0, 2*math.Pi)
				ctx.Stroke()
			}
		}(i)
	}
	wg.Wait()
}

// TestConcurrentSurfaceOperations verifies that surfaces can be created and
// destroyed concurrently from multiple goroutines without races or crashes.
func TestConcurrentSurfaceOperations(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			surf, err := cairo.NewImageSurface(cairo.FormatARGB32, 100, 100)
			require.NoError(t, err)
			surf.Flush()
			surf.MarkDirty()
			err = surf.Close()
			require.NoError(t, err)
		}()
	}
	wg.Wait()
}

// TestConcurrentPatternOperations verifies that a shared pattern can be
// read concurrently by multiple goroutines without data races.
func TestConcurrentPatternOperations(t *testing.T) {
	grad, err := cairo.NewLinearGradient(0, 0, 100, 0)
	require.NoError(t, err)
	grad.AddColorStopRGB(0.0, 1.0, 0.0, 0.0)
	grad.AddColorStopRGB(1.0, 0.0, 0.0, 1.0)
	defer grad.Close()

	var wg sync.WaitGroup
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numOps; j++ {
				_, _ = grad.GetColorStopCount()
				_ = grad.Status()
				_ = grad.GetType()
			}
		}()
	}
	wg.Wait()
}

// TestContextSharedAcrossSurfaces verifies that a single context can safely
// be used by many goroutines concurrently, each setting a different source surface.
func TestContextSharedAcrossSurfaces(t *testing.T) {
	mainSurf, err := cairo.NewImageSurface(cairo.FormatARGB32, 200, 200)
	require.NoError(t, err)
	defer mainSurf.Close()

	ctx, err := cairo.NewContext(mainSurf)
	require.NoError(t, err)
	defer ctx.Close()

	sources := make([]cairo.Surface, numGoroutines)
	for i := range sources {
		s, err := cairo.NewImageSurface(cairo.FormatARGB32, 50, 50)
		require.NoError(t, err)
		sources[i] = s
	}
	defer func() {
		for _, s := range sources {
			s.Close()
		}
	}()

	var wg sync.WaitGroup
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			srcPat, err := cairo.NewSurfacePattern(sources[idx])
			require.NoError(t, err)
			defer srcPat.Close()
			for j := 0; j < numOps; j++ {
				ctx.SetSource(srcPat)
				ctx.Rectangle(0, 0, 50, 50)
				ctx.Fill()
			}
		}(i)
	}
	wg.Wait()
}

// TestPatternSharedAcrossContexts verifies that a single pattern can be safely
// used as a source by multiple contexts on different goroutines concurrently.
func TestPatternSharedAcrossContexts(t *testing.T) {
	sharedPat, err := cairo.NewLinearGradient(0, 0, 100, 100)
	require.NoError(t, err)
	sharedPat.AddColorStopRGBA(0.0, 1.0, 0.0, 0.0, 1.0)
	sharedPat.AddColorStopRGBA(1.0, 0.0, 1.0, 0.0, 1.0)
	defer sharedPat.Close()

	var wg sync.WaitGroup
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			surf, err := cairo.NewImageSurface(cairo.FormatARGB32, 100, 100)
			require.NoError(t, err)
			defer surf.Close()

			ctx, err := cairo.NewContext(surf)
			require.NoError(t, err)
			defer ctx.Close()

			for j := 0; j < numOps; j++ {
				ctx.SetSource(sharedPat)
				ctx.Rectangle(0, 0, 100, 100)
				ctx.Fill()
			}
		}()
	}
	wg.Wait()
}
