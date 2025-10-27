package examples

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
)

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
//  3. Compares the generated image to the golden reference using SHA256 hashes
//  4. Returns true if images match, false otherwise
//
// If the -update-golden flag is set, this function will copy the generated image
// to the golden path instead of comparing, making it easy to update reference images.
//
// Usage:
//
//	func TestMyImage(t *testing.T) {
//	    generator := func(path string) error {
//	        return GenerateMyImage(path)
//	    }
//	    match := CompareImageToGolden(t, generator, "testdata/golden/my_image.png")
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
	match, err := compareImageFiles(tempPath, goldenPath)
	if err != nil {
		t.Errorf("Failed to compare images: %v", err)
		return false
	}

	if !match {
		t.Errorf("Generated image does not match golden reference")
		t.Logf("  Generated: %s", tempPath)
		t.Logf("  Golden:    %s", goldenPath)
		t.Logf("To update the golden image, run: go test -update-golden")
	}

	return match
}

// compareImageFiles compares two image files using SHA256 hashes.
// Returns true if the files have identical content, false otherwise.
func compareImageFiles(generatedPath, goldenPath string) (bool, error) {
	// Check if golden file exists
	if _, err := os.Stat(goldenPath); os.IsNotExist(err) {
		return false, fmt.Errorf("golden reference image does not exist at %s (run with -update-golden to create it)", goldenPath)
	}

	// Compute hash of generated image
	generatedHash, err := computeFileHash(generatedPath)
	if err != nil {
		return false, fmt.Errorf("failed to hash generated image: %w", err)
	}

	// Compute hash of golden image
	goldenHash, err := computeFileHash(goldenPath)
	if err != nil {
		return false, fmt.Errorf("failed to hash golden image: %w", err)
	}

	// Compare hashes
	return generatedHash == goldenHash, nil
}

// computeFileHash computes the SHA256 hash of a file.
func computeFileHash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// updateGoldenImage copies the generated image to the golden reference location.
func updateGoldenImage(generatedPath, goldenPath string) error {
	// Ensure the directory for the golden image exists
	goldenDir := filepath.Dir(goldenPath)
	if err := os.MkdirAll(goldenDir, 0755); err != nil {
		return fmt.Errorf("failed to create golden directory: %w", err)
	}

	// Open source file
	src, err := os.Open(generatedPath)
	if err != nil {
		return fmt.Errorf("failed to open generated image: %w", err)
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(goldenPath)
	if err != nil {
		return fmt.Errorf("failed to create golden image: %w", err)
	}
	defer dst.Close()

	// Copy content
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
