package pattern

// PatternType identifies the type of a pattern.
//
// Cairo supports multiple pattern types, each with different characteristics
// for how they paint:
//   - Solid: Single uniform color
//   - Linear: Color gradient along a line
//   - Radial: Color gradient in a circular pattern
//   - Surface: Texturing with an image
//   - Mesh: Complex multi-point gradients
//   - RasterSource: Procedural pattern generation
//
// Use Pattern.GetType() to query a pattern's type at runtime.
//
// Example:
//
//	pattern, err := pattern.NewSolidPatternRGB(1.0, 0.0, 0.0)
//	if err != nil {
//	    return err
//	}
//	defer pattern.Close()
//
//	if pattern.GetType() == pattern.PatternTypeSolid {
//	    fmt.Println("This is a solid color pattern")
//	}
//
//go:generate stringer -type=PatternType -trimprefix=PatternType
type PatternType int

const (
	// PatternTypeSolid represents a pattern with a single, uniform color.
	// Created with NewSolidPatternRGB or NewSolidPatternRGBA.
	PatternTypeSolid PatternType = iota

	// PatternTypeSurface represents a pattern based on a Cairo surface (image).
	// Used for texturing with images or other rendered content.
	// (Planned - not yet implemented)
	PatternTypeSurface

	// PatternTypeLinear represents a linear gradient pattern.
	// Colors transition smoothly along a line between control points.
	// (Planned - not yet implemented)
	PatternTypeLinear

	// PatternTypeRadial represents a radial gradient pattern.
	// Colors transition smoothly in circles radiating from a center point.
	// (Planned - not yet implemented)
	PatternTypeRadial

	// PatternTypeMesh represents a mesh gradient pattern.
	// Complex gradients defined by a patch mesh with multiple control points.
	// (Planned - not yet implemented)
	PatternTypeMesh

	// PatternTypeRasterSource represents a procedural pattern.
	// Pattern content is generated programmatically on demand.
	// (Planned - not yet implemented)
	PatternTypeRasterSource
)
