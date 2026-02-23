// ABOUTME: Test harness for example image generation tests.
// ABOUTME: Provides golden image comparison with pixel-level tolerance for cross-platform compatibility.
package examples

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"testing"
)

// maxChannelDelta is the maximum allowed per-channel difference (0–255) between
// corresponding pixels in the generated and golden images. Cairo's geometric
// rendering can differ by 1–2 levels between platforms due to sub-pixel
// antialiasing; 3 provides a comfortable margin without masking real regressions.
const maxChannelDelta = 3

// maxDiffPixelFraction is the maximum fraction of pixels that may exceed
// maxChannelDelta before the comparison is considered a failure. Antialiasing
// noise affects only edge pixels, so 1% is generous for geometric shapes.
const maxDiffPixelFraction = 0.01

// updateGolden is a flag that controls whether to update golden images
var updateGolden = flag.Bool("update-golden", false, "update golden reference images")

// ImageGeneratorFunc is a function that generates an image at the given output path.
// It should create a complete image file (e.g., PNG) at the specified location.
type ImageGeneratorFunc func(outputPath string) error

// CompareImageToGolden tests an image generator function against a golden reference image.
//
// This function:
//  1. Creates a temporary directory (automatically cleaned up by testing framework)
//  2. Runs the generator function to create an image
//  3. Compares the generated image to the golden reference using pixel-level tolerance
//  4. Returns true if images match within tolerance, false otherwise
//
// Comparison tolerates minor per-pixel differences to handle sub-pixel antialiasing
// variation across platforms: up to maxDiffPixelFraction of pixels may differ, and
// differing pixels may vary by at most maxChannelDelta per channel (0–255). Image
// size mismatches are always hard failures regardless of tolerance.
//
// If the -update-golden flag is set, this function will copy the generated image
// to the golden path instead of comparing, making it easy to update reference images.
//
// Usage:
//
//	func TestMyImage(t *testing.T) {
//	    match := CompareImageToGolden(t, GenerateMyImage, "testdata/golden/my_image.png")
//	    if !match {
//	        t.Error("Generated image does not match golden reference")
//	    }
//	}
func CompareImageToGolden(t *testing.T, generator ImageGeneratorFunc, goldenPath string) bool {
	t.Helper()

	// Create temporary directory for test output (automatically cleaned up)
	tempDir := t.TempDir()
	tempPath := filepath.Join(tempDir, "output.png")

	// Generate the image
	if err := generator(tempPath); err != nil {
		t.Errorf("Failed to generate image: %v", err)
		return false
	}

	// Verify the generated file exists
	if _, err := os.Stat(tempPath); os.IsNotExist(err) {
		t.Errorf("Generated image file does not exist at %s", tempPath)
		return false
	}

	// If -update-golden flag is set, update the golden image
	if *updateGolden {
		if err := updateGoldenImage(tempPath, goldenPath); err != nil {
			t.Errorf("Failed to update golden image: %v", err)
			return false
		}
		t.Logf("Updated golden image at %s", goldenPath)
		return true
	}

	// Compare generated image to golden reference
	match, diagnostic, err := compareImageFiles(tempPath, goldenPath)
	if err != nil {
		t.Errorf("Failed to compare images: %v", err)
		return false
	}

	if !match {
		t.Errorf("Generated image does not match golden reference: %s", diagnostic)
		t.Logf("  Generated: %s", tempPath)
		t.Logf("  Golden:    %s", goldenPath)
		t.Logf("  Thresholds: max channel delta=%d, max differing pixels=%.1f%%",
			maxChannelDelta, maxDiffPixelFraction*100)
		t.Logf("  To update: go test ./examples -update-golden")
	}

	return match
}

// compareImageFiles compares two PNG files using pixel-level tolerance.
// Returns (true, "", nil) when images match within tolerance.
// Returns (false, diagnostic, nil) when images differ beyond tolerance; diagnostic
// describes the extent of the difference to help tune thresholds if needed.
// Returns (false, "", err) when the comparison cannot be performed.
func compareImageFiles(generatedPath, goldenPath string) (bool, string, error) {
	if _, err := os.Stat(goldenPath); os.IsNotExist(err) {
		return false, "", fmt.Errorf(
			"golden reference image does not exist at %s (run with -update-golden to create it)",
			goldenPath,
		)
	}

	generated, err := decodePNG(generatedPath)
	if err != nil {
		return false, "", fmt.Errorf("failed to decode generated image: %w", err)
	}

	golden, err := decodePNG(goldenPath)
	if err != nil {
		return false, "", fmt.Errorf("failed to decode golden image: %w", err)
	}

	// Size mismatch is always a hard failure
	genBounds := generated.Bounds()
	goldenBounds := golden.Bounds()
	if genBounds != goldenBounds {
		return false, "", fmt.Errorf(
			"image size mismatch: generated %v, golden %v",
			genBounds, goldenBounds,
		)
	}

	width := genBounds.Max.X - genBounds.Min.X
	height := genBounds.Max.Y - genBounds.Min.Y
	totalPixels := width * height

	diffPixels, maxDelta := scanPixelDiffs(generated, golden, genBounds)

	if diffPixels == 0 {
		return true, "", nil
	}

	diffFraction := float64(diffPixels) / float64(totalPixels)
	diagnostic := fmt.Sprintf(
		"%d/%d pixels differ (%.2f%%, max channel delta: %d)",
		diffPixels, totalPixels, diffFraction*100, maxDelta,
	)

	if diffFraction > maxDiffPixelFraction {
		return false, diagnostic, nil
	}

	return true, "", nil
}

// scanPixelDiffs iterates every pixel in bounds and counts those where any RGBA
// channel exceeds maxChannelDelta. Returns the count and the largest delta seen.
func scanPixelDiffs(a, b image.Image, bounds image.Rectangle) (diffCount int, maxDelta uint8) {
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			p1 := color.NRGBAModel.Convert(a.At(x, y)).(color.NRGBA)
			p2 := color.NRGBAModel.Convert(b.At(x, y)).(color.NRGBA)

			delta := maxUint8(
				absDiff(p1.R, p2.R),
				absDiff(p1.G, p2.G),
				absDiff(p1.B, p2.B),
				absDiff(p1.A, p2.A),
			)
			if delta > maxDelta {
				maxDelta = delta
			}
			if delta > maxChannelDelta {
				diffCount++
			}
		}
	}
	return diffCount, maxDelta
}

// decodePNG opens and decodes a PNG file into an image.Image.
func decodePNG(path string) (image.Image, error) {
	filepath.Clean(path)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()

	img, err := png.Decode(f)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// absDiff returns the absolute difference between two uint8 values.
func absDiff(a, b uint8) uint8 {
	if a > b {
		return a - b
	}
	return b - a
}

// maxUint8 returns the largest of four uint8 values.
func maxUint8(a, b, c, d uint8) uint8 {
	m := a
	if b > m {
		m = b
	}
	if c > m {
		m = c
	}
	if d > m {
		m = d
	}
	return m
}

// updateGoldenImage copies the generated image to the golden reference location.
func updateGoldenImage(generatedPath, goldenPath string) error {
	goldenDir := filepath.Dir(goldenPath)
	if err := os.MkdirAll(goldenDir, 0o750); err != nil {
		return fmt.Errorf("failed to create golden directory: %w", err)
	}

	filepath.Clean(generatedPath)
	src, err := os.Open(generatedPath)
	if err != nil {
		return fmt.Errorf("failed to open generated image: %w", err)
	}
	defer func() {
		_ = src.Close()
	}()

	dst, err := os.Create(goldenPath)
	if err != nil {
		return fmt.Errorf("failed to create golden image: %w", err)
	}
	defer func() {
		_ = dst.Close()
	}()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("failed to copy image: %w", err)
	}

	return nil
}

// CheckFileSize is a helper that verifies a file exists and has a size within the expected range.
// This is useful for sanity-checking generated images before doing byte-level comparison.
func CheckFileSize(t *testing.T, path string, minBytes, maxBytes int64) bool {
	t.Helper()

	info, err := os.Stat(path)
	if err != nil {
		t.Errorf("Failed to stat file %s: %v", path, err)
		return false
	}

	size := info.Size()
	if size < minBytes {
		t.Errorf("File %s is too small: %d bytes (minimum: %d bytes)", path, size, minBytes)
		return false
	}

	if size > maxBytes {
		t.Errorf("File %s is too large: %d bytes (maximum: %d bytes)", path, size, maxBytes)
		return false
	}

	t.Logf("File %s size is reasonable: %d bytes", path, size)
	return true
}

// RegionHasNonBackgroundPixels returns true if any pixel in the rectangle
// [x0,x1] × [y0,y1] differs significantly from white (the background color).
// Use this for structural tests on images with a white background.
func RegionHasNonBackgroundPixels(img image.Image, x0, y0, x1, y1 int) bool {
	bounds := img.Bounds()
	for y := y0; y <= y1 && y < bounds.Max.Y; y++ {
		for x := x0; x <= x1 && x < bounds.Max.X; x++ {
			px := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)
			if px.R < 250 || px.G < 250 || px.B < 250 {
				return true
			}
		}
	}
	return false
}
