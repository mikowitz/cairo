package surface

import (
	"path/filepath"
	"testing"
)

// BenchmarkImageSurfaceCreation benchmarks the creation and destruction of image surfaces.
// This measures the overhead of allocating Cairo image surfaces via CGO.
func BenchmarkImageSurfaceCreation(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		surf, err := NewImageSurface(FormatARGB32, 400, 400)
		if err != nil {
			b.Fatalf("Failed to create surface: %v", err)
		}
		_ = surf.Close()
	}
}

// BenchmarkWriteToPNG benchmarks writing surfaces to PNG files.
// This measures the performance of PNG encoding and file I/O.
// The surface is created once and written repeatedly to test the PNG
// writing overhead separately from surface creation.
func BenchmarkWriteToPNG(b *testing.B) {
	// Create a surface once
	surf, err := NewImageSurface(FormatARGB32, 400, 400)
	if err != nil {
		b.Fatalf("Failed to create surface: %v", err)
	}
	defer func() {
		_ = surf.Close()
	}()

	// Flush the surface once
	surf.Flush()

	// Create temp directory for PNG files
	tempDir := b.TempDir()

	b.ReportAllocs()
	b.ResetTimer()

	// Benchmark writing to PNG
	for i := 0; i < b.N; i++ {
		outputPath := filepath.Join(tempDir, "benchmark.png")
		err := surf.WriteToPNG(outputPath)
		if err != nil {
			b.Fatalf("Failed to write PNG: %v", err)
		}
	}
}
