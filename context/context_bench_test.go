package context

import (
	"testing"

	"github.com/mikowitz/cairo/surface"
)

// BenchmarkContextCreation benchmarks the overhead of creating and destroying contexts.
// This measures the CGO call overhead and resource allocation for context creation.
func BenchmarkContextCreation(b *testing.B) {
	// Create a surface once for all iterations
	surf, err := surface.NewImageSurface(surface.FormatARGB32, 400, 400)
	if err != nil {
		b.Fatalf("Failed to create surface: %v", err)
	}
	defer func() {
		_ = surf.Close()
	}()

	b.ReportAllocs()
	b.ResetTimer()

	// Benchmark creating and closing contexts
	for i := 0; i < b.N; i++ {
		ctx, err := NewContext(surf)
		if err != nil {
			b.Fatalf("Failed to create context: %v", err)
		}
		_ = ctx.Close()
	}
}

// BenchmarkContextPathOperations benchmarks path construction operations.
// This measures the performance of MoveTo and LineTo calls, which are
// fundamental building blocks for complex drawings.
func BenchmarkContextPathOperations(b *testing.B) {
	// Create surface and context once
	surf, err := surface.NewImageSurface(surface.FormatARGB32, 400, 400)
	if err != nil {
		b.Fatalf("Failed to create surface: %v", err)
	}
	defer func() {
		_ = surf.Close()
	}()

	ctx, err := NewContext(surf)
	if err != nil {
		b.Fatalf("Failed to create context: %v", err)
	}
	defer func() {
		_ = ctx.Close()
	}()

	b.ReportAllocs()
	b.ResetTimer()

	// Benchmark many path operations
	for i := 0; i < b.N; i++ {
		ctx.NewPath()

		// Create a complex path with 100 operations
		ctx.MoveTo(0, 0)
		for j := 0; j < 100; j++ {
			x := float64(j * 4)
			y := float64((j * 7) % 400)
			ctx.LineTo(x, y)
		}
	}
}

// BenchmarkContextFillOperations benchmarks fill operations.
// This measures the performance of creating paths and filling them,
// which is one of the most common drawing operations.
func BenchmarkContextFillOperations(b *testing.B) {
	// Create surface and context once
	surf, err := surface.NewImageSurface(surface.FormatARGB32, 400, 400)
	if err != nil {
		b.Fatalf("Failed to create surface: %v", err)
	}
	defer func() {
		_ = surf.Close()
	}()

	ctx, err := NewContext(surf)
	if err != nil {
		b.Fatalf("Failed to create context: %v", err)
	}
	defer func() {
		_ = ctx.Close()
	}()

	b.ReportAllocs()
	b.ResetTimer()

	// Benchmark fill operations
	for i := 0; i < b.N; i++ {
		// Create a rectangle
		ctx.Rectangle(50, 50, 100, 100)

		// Set source color
		ctx.SetSourceRGB(1.0, 0.0, 0.0)

		// Fill the rectangle
		ctx.Fill()
	}
}
