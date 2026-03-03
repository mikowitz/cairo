// ABOUTME: Tests for the animation frame generation example.
// ABOUTME: Verifies frame count, PNG validity, dimensions, and that frames differ over time.

package examples

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnimationGeneratesAllFrames(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, GenerateAnimation(dir))

	for i := range AnimationFrameCount {
		path := filepath.Join(dir, fmt.Sprintf("frame_%03d.png", i))
		_, err := os.Stat(path)
		require.NoErrorf(t, err, "frame %d should exist at %s", i, path)
	}
}

func TestAnimationFramesAreValidPNGs(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, GenerateAnimation(dir))

	for _, i := range []int{0, AnimationFrameCount / 2, AnimationFrameCount - 1} {
		path := filepath.Join(dir, fmt.Sprintf("frame_%03d.png", i))
		img, err := decodePNG(path)
		require.NoErrorf(t, err, "frame %d should be a valid PNG", i)
		assert.Equal(t, 400, img.Bounds().Dx(), "frame %d should be 400px wide", i)
		assert.Equal(t, 300, img.Bounds().Dy(), "frame %d should be 300px tall", i)
	}
}

func TestAnimationFramesDiffer(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, GenerateAnimation(dir))

	first, err := decodePNG(filepath.Join(dir, "frame_000.png"))
	require.NoError(t, err)

	mid, err := decodePNG(filepath.Join(dir, fmt.Sprintf("frame_%03d.png", AnimationFrameCount/2)))
	require.NoError(t, err)

	diffPixels := countDifferingPixels(first, mid)
	assert.Greater(t, diffPixels, 500, "frame 0 and the midpoint frame should differ meaningfully")
}

func TestAnimationRendersVisibleContent(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, GenerateAnimation(dir))

	img, err := decodePNG(filepath.Join(dir, "frame_000.png"))
	require.NoError(t, err)

	// The rotating hexagon is centered at (200, 150), so the center region should
	// contain non-background pixels.
	assert.True(t,
		RegionHasNonBackgroundPixels(img, 100, 50, 300, 250),
		"center region should contain the rotating shape",
	)
}

// countDifferingPixels counts pixels that differ between two same-sized images.
func countDifferingPixels(a, b image.Image) int {
	count := 0
	bounds := a.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if a.At(x, y) != b.At(x, y) {
				count++
			}
		}
	}
	return count
}
