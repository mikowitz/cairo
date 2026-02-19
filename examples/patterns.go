// ABOUTME: Surface pattern examples demonstrating texture mapping and pattern fills.
// ABOUTME: Shows different extend modes (repeat, reflect, pad) and filter modes for surface patterns.
package examples

import (
	"fmt"
	"math"
	"unsafe"

	"github.com/mikowitz/cairo"
	"github.com/mikowitz/cairo/pattern"
	"github.com/mikowitz/cairo/status"
	"github.com/mikowitz/cairo/surface"
)

// surfaceAdapter adapts surface.ImageSurface to work with pattern.NewSurfacePattern.
// This is needed because pattern.Surface expects Ptr() interface{} but
// surface.Surface has Ptr() SurfacePtr.
type surfaceAdapter struct {
	*surface.ImageSurface
}

func (s surfaceAdapter) Ptr() interface{} {
	// Convert SurfacePtr to unsafe.Pointer for the pattern API
	return unsafe.Pointer(s.ImageSurface.Ptr())
}

func (s surfaceAdapter) Status() status.Status {
	return s.ImageSurface.Status()
}

// GeneratePatterns creates a 800x600 PNG image demonstrating Cairo surface patterns.
//
// The image contains multiple examples showing:
//   - Top left: Simple checker pattern with ExtendRepeat (tiling)
//   - Top right: Same pattern with ExtendReflect (mirroring)
//   - Middle left: Pattern with ExtendPad (edge colors extend)
//   - Middle right: Pattern with ExtendNone (transparent outside)
//   - Bottom left: Pattern with FilterNearest (pixelated scaling)
//   - Bottom right: Pattern with FilterBilinear (smooth scaling)
//
// This demonstrates:
//   - Creating a small ImageSurface to use as a texture
//   - Creating SurfacePattern from an ImageSurface
//   - Different Extend modes (None, Repeat, Reflect, Pad)
//   - Different Filter modes (Nearest, Bilinear)
//   - Using surface patterns as source for drawing
//   - Pattern transformations and tiling
//
// All resources are properly cleaned up using defer statements.
func GeneratePatterns(outputPath string) error {
	// Create the main 800x600 ARGB32 image surface
	surface, err := cairo.NewImageSurface(cairo.FormatARGB32, 800, 600)
	if err != nil {
		return fmt.Errorf("failed to create surface: %w", err)
	}
	defer func() {
		_ = surface.Close()
	}()

	// Create drawing context
	ctx, err := cairo.NewContext(surface)
	if err != nil {
		return fmt.Errorf("failed to create context: %w", err)
	}
	defer func() {
		_ = ctx.Close()
	}()

	// Fill the background with light gray
	ctx.SetSourceRGB(0.9, 0.9, 0.9)
	ctx.Paint()

	// Create a small 20x20 checker pattern surface to use as texture
	checkerSurface, err := createCheckerSurface(20, 20)
	if err != nil {
		return fmt.Errorf("failed to create checker surface: %w", err)
	}
	defer func() {
		_ = checkerSurface.Close()
	}()

	// 1. Top-left: ExtendRepeat (tiling)
	if err := drawPatternExample(ctx, checkerSurface, 20, 20, 240, 180,
		pattern.ExtendRepeat, pattern.FilterGood, "Repeat"); err != nil {
		return err
	}

	// 2. Top-right: ExtendReflect (mirroring)
	if err := drawPatternExample(ctx, checkerSurface, 280, 20, 240, 180,
		pattern.ExtendReflect, pattern.FilterGood, "Reflect"); err != nil {
		return err
	}

	// 3. Middle-left: ExtendPad (edge pixels extend)
	if err := drawPatternExample(ctx, checkerSurface, 540, 20, 240, 180,
		pattern.ExtendPad, pattern.FilterGood, "Pad"); err != nil {
		return err
	}

	// 4. Middle-right: ExtendNone (transparent outside)
	if err := drawPatternExample(ctx, checkerSurface, 20, 220, 240, 180,
		pattern.ExtendNone, pattern.FilterGood, "None"); err != nil {
		return err
	}

	// 5. Bottom-left: FilterNearest (pixelated when scaled)
	if err := drawPatternExample(ctx, checkerSurface, 280, 220, 240, 180,
		pattern.ExtendRepeat, pattern.FilterNearest, "Nearest"); err != nil {
		return err
	}

	// 6. Bottom-right: FilterBilinear (smooth when scaled)
	if err := drawPatternExample(ctx, checkerSurface, 540, 220, 240, 180,
		pattern.ExtendRepeat, pattern.FilterBilinear, "Bilinear"); err != nil {
		return err
	}

	// 7. Bottom section: Complex example with rotation and scaling
	if err := drawComplexPatternExample(ctx, checkerSurface, 20, 420); err != nil {
		return err
	}

	// Flush any pending operations
	surface.Flush()

	// Write the surface to a PNG file
	if err := surface.WriteToPNG(outputPath); err != nil {
		return fmt.Errorf("failed to write PNG: %w", err)
	}

	return nil
}

// createCheckerSurface creates a small checker pattern surface.
func createCheckerSurface(width, height int) (*surface.ImageSurface, error) {
	surf, err := cairo.NewImageSurface(cairo.FormatARGB32, width, height)
	if err != nil {
		return nil, err
	}

	ctx, err := cairo.NewContext(surf)
	if err != nil {
		_ = surf.Close()
		return nil, err
	}
	defer func() {
		_ = ctx.Close()
	}()

	// Draw a 2x2 checker pattern
	halfW := float64(width) / 2
	halfH := float64(height) / 2

	// Top-left: Red
	ctx.SetSourceRGB(1.0, 0.2, 0.2)
	ctx.Rectangle(0, 0, halfW, halfH)
	ctx.Fill()

	// Top-right: Blue
	ctx.SetSourceRGB(0.2, 0.2, 1.0)
	ctx.Rectangle(halfW, 0, halfW, halfH)
	ctx.Fill()

	// Bottom-left: Blue
	ctx.SetSourceRGB(0.2, 0.2, 1.0)
	ctx.Rectangle(0, halfH, halfW, halfH)
	ctx.Fill()

	// Bottom-right: Red
	ctx.SetSourceRGB(1.0, 0.2, 0.2)
	ctx.Rectangle(halfW, halfH, halfW, halfH)
	ctx.Fill()

	surf.Flush()
	return surf, nil
}

// drawPatternExample draws a single pattern example with label.
func drawPatternExample(ctx *cairo.Context, srcSurface *surface.ImageSurface,
	x, y, width, height float64, extend pattern.Extend, filter pattern.Filter, label string) error {
	// Create pattern from the checker surface (using adapter for interface compatibility)
	pat, err := pattern.NewSurfacePattern(surfaceAdapter{srcSurface})
	if err != nil {
		return fmt.Errorf("failed to create surface pattern: %w", err)
	}
	defer func() {
		_ = pat.Close()
	}()

	// Set extend and filter modes
	pat.SetExtend(extend)
	pat.SetFilter(filter)

	// Draw with the pattern
	ctx.SetSource(pat)
	ctx.Rectangle(x, y, width, height)
	ctx.Fill()

	// Draw a border around the rectangle
	ctx.SetSourceRGB(0.0, 0.0, 0.0)
	ctx.SetLineWidth(2.0)
	ctx.Rectangle(x, y, width, height)
	ctx.Stroke()

	// Draw label (simplified - in production you'd use text rendering)
	// For now we just draw a small indicator line
	_ = label // Label would require text rendering which isn't implemented yet

	return nil
}

// drawComplexPatternExample demonstrates pattern transformation.
func drawComplexPatternExample(ctx *cairo.Context, srcSurface *surface.ImageSurface, x, y float64) error {
	// Create pattern (using adapter for interface compatibility)
	pat, err := pattern.NewSurfacePattern(surfaceAdapter{srcSurface})
	if err != nil {
		return fmt.Errorf("failed to create surface pattern: %w", err)
	}
	defer func() {
		_ = pat.Close()
	}()

	pat.SetExtend(pattern.ExtendRepeat)
	pat.SetFilter(pattern.FilterGood)

	// Save the graphics state
	ctx.Save()

	// Translate to the center of where we want to draw
	centerX := x + 380
	centerY := y + 80

	ctx.Translate(centerX, centerY)
	ctx.Rotate(math.Pi / 6) // Rotate 30 degrees
	ctx.Scale(1.5, 1.5)     // Scale up

	// Set the pattern as source
	ctx.SetSource(pat)

	// Draw a circle with the transformed pattern
	ctx.Arc(0, 0, 80, 0, 2*math.Pi)
	ctx.Fill()

	// Draw border
	ctx.SetSourceRGB(0.0, 0.0, 0.0)
	ctx.SetLineWidth(2.0)
	ctx.Arc(0, 0, 80, 0, 2*math.Pi)
	ctx.Stroke()

	// Restore the graphics state
	ctx.Restore()

	return nil
}
