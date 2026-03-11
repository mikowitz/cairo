// ABOUTME: concurrent_bench_test.go - benchmarks for concurrent Cairo drawing operations.
// ABOUTME: Measures parallel drawing performance, surface creation, and lock contention overhead.
package benchmarks_test

import (
	"fmt"
	"math"
	"testing"

	cairo "github.com/mikowitz/cairo"
)

// BenchmarkConcurrentDrawing measures drawing performance under three concurrency
// models: single-threaded, shared context with lock contention, and per-goroutine
// contexts with no contention. Compare results to quantify locking overhead.
func BenchmarkConcurrentDrawing(b *testing.B) {
	b.Run("SingleThreaded", func(b *testing.B) {
		surf, err := cairo.NewImageSurface(cairo.FormatARGB32, 400, 400)
		if err != nil {
			b.Fatalf("Failed to create surface: %v", err)
		}
		defer surf.Close()

		ctx, err := cairo.NewContext(surf)
		if err != nil {
			b.Fatalf("Failed to create context: %v", err)
		}
		defer ctx.Close()

		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			ctx.SetSourceRGBA(0.5, 0.5, 0.5, 1.0)
			ctx.Rectangle(10, 10, 100, 100)
			ctx.Fill()
			ctx.Arc(200, 200, 50, 0, 2*math.Pi)
			ctx.Stroke()
		}
	})

	// SharedContextParallel: all goroutines compete for the same context mutex.
	b.Run("SharedContextParallel", func(b *testing.B) {
		surf, err := cairo.NewImageSurface(cairo.FormatARGB32, 400, 400)
		if err != nil {
			b.Fatalf("Failed to create surface: %v", err)
		}
		defer surf.Close()

		ctx, err := cairo.NewContext(surf)
		if err != nil {
			b.Fatalf("Failed to create context: %v", err)
		}
		defer ctx.Close()

		b.ReportAllocs()
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				ctx.SetSourceRGBA(0.5, 0.5, 0.5, 1.0)
				ctx.Rectangle(10, 10, 100, 100)
				ctx.Fill()
				ctx.Arc(200, 200, 50, 0, 2*math.Pi)
				ctx.Stroke()
			}
		})
	})

	// PerGoroutineContextParallel: each goroutine has its own context; no lock contention.
	// This represents the recommended pattern for maximum parallel throughput.
	b.Run("PerGoroutineContextParallel", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			surf, err := cairo.NewImageSurface(cairo.FormatARGB32, 400, 400)
			if err != nil {
				b.Fatalf("Failed to create surface: %v", err)
			}
			defer surf.Close()

			ctx, err := cairo.NewContext(surf)
			if err != nil {
				b.Fatalf("Failed to create context: %v", err)
			}
			defer ctx.Close()

			for pb.Next() {
				ctx.SetSourceRGBA(0.5, 0.5, 0.5, 1.0)
				ctx.Rectangle(10, 10, 100, 100)
				ctx.Fill()
				ctx.Arc(200, 200, 50, 0, 2*math.Pi)
				ctx.Stroke()
			}
		})
	})
}

// BenchmarkConcurrentSurfaceCreation measures the cost of ImageSurface allocation
// in single-threaded and parallel scenarios. Surface creation involves CGO calls
// and memory allocation on both the Go heap and the C heap.
func BenchmarkConcurrentSurfaceCreation(b *testing.B) {
	b.Run("SingleThreaded", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			surf, err := cairo.NewImageSurface(cairo.FormatARGB32, 100, 100)
			if err != nil {
				b.Fatalf("Failed to create surface: %v", err)
			}
			_ = surf.Close()
		}
	})

	b.Run("Parallel", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				surf, err := cairo.NewImageSurface(cairo.FormatARGB32, 100, 100)
				if err != nil {
					b.Fatalf("Failed to create surface: %v", err)
				}
				_ = surf.Close()
			}
		})
	})
}

// BenchmarkLockContention measures how mutex contention scales with the number
// of concurrent goroutines drawing to a single shared context. Higher parallelism
// multipliers reveal the per-operation cost of lock waiting.
func BenchmarkLockContention(b *testing.B) {
	parallelismLevels := []int{1, 2, 4, 8}

	for _, p := range parallelismLevels {
		p := p
		b.Run(fmt.Sprintf("Parallelism%d", p), func(b *testing.B) {
			surf, err := cairo.NewImageSurface(cairo.FormatARGB32, 200, 200)
			if err != nil {
				b.Fatalf("Failed to create surface: %v", err)
			}
			defer surf.Close()

			ctx, err := cairo.NewContext(surf)
			if err != nil {
				b.Fatalf("Failed to create context: %v", err)
			}
			defer ctx.Close()

			b.SetParallelism(p)
			b.ReportAllocs()
			b.ResetTimer()

			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					ctx.SetSourceRGB(0.5, 0.5, 0.5)
					ctx.Rectangle(10, 10, 50, 50)
					ctx.Fill()
				}
			})
		})
	}
}
