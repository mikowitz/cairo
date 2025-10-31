package examples

import (
	"fmt"
	"math"

	"github.com/mikowitz/cairo"
	"github.com/mikowitz/cairo/pattern"
)

// GenerateGradients creates a 600x600 PNG image demonstrating Cairo gradient patterns.
//
// The image contains:
//   - Top left: Simple linear gradient (red to blue, horizontal)
//   - Top right: Linear gradient with multiple color stops (rainbow)
//   - Middle left: Simple radial gradient (white center to blue edge)
//   - Middle right: Radial gradient with transparency (fading effect)
//   - Bottom left: Vertical linear gradient with semi-transparency
//   - Bottom right: Radial gradient with offset centers
//
// This demonstrates:
//   - Creating LinearGradient patterns
//   - Creating RadialGradient patterns
//   - Adding color stops with AddColorStopRGB
//   - Adding color stops with AddColorStopRGBA (transparency)
//   - Using gradients as source patterns
//   - Multiple color stops for complex gradients
//
// All resources are properly cleaned up using defer statements.
func GenerateGradients(outputPath string) error {
	// Create a 600x600 ARGB32 image surface
	surface, err := cairo.NewImageSurface(cairo.FormatARGB32, 600, 600)
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

	// Fill the background with white
	ctx.SetSourceRGB(1.0, 1.0, 1.0)
	ctx.Paint()

	// 1. Top-left: Simple horizontal linear gradient (red to blue)
	linearSimple, err := pattern.NewLinearGradient(20, 20, 280, 20)
	if err != nil {
		return fmt.Errorf("failed to create linear gradient: %w", err)
	}
	defer func() {
		_ = linearSimple.Close()
	}()
	linearSimple.AddColorStopRGB(0.0, 1.0, 0.0, 0.0) // Red at start
	linearSimple.AddColorStopRGB(1.0, 0.0, 0.0, 1.0) // Blue at end
	ctx.SetSource(linearSimple)
	ctx.Rectangle(20, 20, 260, 160)
	ctx.Fill()

	// 2. Top-right: Linear gradient with multiple color stops (rainbow)
	linearRainbow, err := pattern.NewLinearGradient(320, 20, 580, 20)
	if err != nil {
		return fmt.Errorf("failed to create rainbow gradient: %w", err)
	}
	defer func() {
		_ = linearRainbow.Close()
	}()
	linearRainbow.AddColorStopRGB(0.0, 1.0, 0.0, 0.0)   // Red
	linearRainbow.AddColorStopRGB(0.25, 1.0, 1.0, 0.0)  // Yellow
	linearRainbow.AddColorStopRGB(0.5, 0.0, 1.0, 0.0)   // Green
	linearRainbow.AddColorStopRGB(0.75, 0.0, 0.0, 1.0)  // Blue
	linearRainbow.AddColorStopRGB(1.0, 0.5, 0.0, 0.5)   // Purple
	ctx.SetSource(linearRainbow)
	ctx.Rectangle(320, 20, 260, 160)
	ctx.Fill()

	// 3. Middle-left: Simple radial gradient (white center to blue edge)
	radialSimple, err := pattern.NewRadialGradient(150, 290, 10, 150, 290, 120)
	if err != nil {
		return fmt.Errorf("failed to create radial gradient: %w", err)
	}
	defer func() {
		_ = radialSimple.Close()
	}()
	radialSimple.AddColorStopRGB(0.0, 1.0, 1.0, 1.0) // White center
	radialSimple.AddColorStopRGB(1.0, 0.0, 0.0, 1.0) // Blue edge
	ctx.SetSource(radialSimple)
	ctx.Arc(150, 290, 120, 0, 2*math.Pi)
	ctx.Fill()

	// 4. Middle-right: Radial gradient with transparency (fading effect)
	radialFade, err := pattern.NewRadialGradient(450, 290, 0, 450, 290, 120)
	if err != nil {
		return fmt.Errorf("failed to create fading gradient: %w", err)
	}
	defer func() {
		_ = radialFade.Close()
	}()
	radialFade.AddColorStopRGBA(0.0, 1.0, 0.5, 0.0, 1.0) // Opaque orange center
	radialFade.AddColorStopRGBA(0.7, 1.0, 0.0, 0.5, 0.5) // Semi-transparent pink
	radialFade.AddColorStopRGBA(1.0, 1.0, 0.0, 0.0, 0.0) // Transparent red edge
	ctx.SetSource(radialFade)
	ctx.Arc(450, 290, 120, 0, 2*math.Pi)
	ctx.Fill()

	// 5. Bottom-left: Vertical linear gradient with semi-transparency
	linearVertical, err := pattern.NewLinearGradient(20, 420, 20, 580)
	if err != nil {
		return fmt.Errorf("failed to create vertical gradient: %w", err)
	}
	defer func() {
		_ = linearVertical.Close()
	}()
	linearVertical.AddColorStopRGBA(0.0, 0.0, 0.8, 0.0, 1.0) // Opaque green at top
	linearVertical.AddColorStopRGBA(0.5, 1.0, 1.0, 0.0, 0.6) // Semi-transparent yellow middle
	linearVertical.AddColorStopRGBA(1.0, 0.8, 0.0, 0.8, 0.3) // More transparent purple at bottom
	ctx.SetSource(linearVertical)
	ctx.Rectangle(20, 420, 260, 160)
	ctx.Fill()

	// 6. Bottom-right: Radial gradient with offset centers (creates directional lighting effect)
	radialOffset, err := pattern.NewRadialGradient(420, 460, 20, 480, 520, 100)
	if err != nil {
		return fmt.Errorf("failed to create offset radial gradient: %w", err)
	}
	defer func() {
		_ = radialOffset.Close()
	}()
	radialOffset.AddColorStopRGB(0.0, 1.0, 1.0, 0.8)  // Pale yellow (highlight)
	radialOffset.AddColorStopRGB(0.4, 1.0, 0.6, 0.0)  // Orange
	radialOffset.AddColorStopRGB(0.7, 0.8, 0.2, 0.0)  // Dark orange
	radialOffset.AddColorStopRGB(1.0, 0.4, 0.1, 0.0)  // Deep red-orange
	ctx.SetSource(radialOffset)
	ctx.Rectangle(320, 420, 260, 160)
	ctx.Fill()

	// Flush any pending operations
	surface.Flush()

	// Write the surface to a PNG file
	if err := surface.WriteToPNG(outputPath); err != nil {
		return fmt.Errorf("failed to write PNG: %w", err)
	}

	return nil
}
