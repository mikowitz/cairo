// ABOUTME: Package font provides types for Cairo's toy font API.
// ABOUTME: Defines Slant and Weight types used with SelectFontFace.

// Package font provides font-related types for Cairo text rendering.
//
// # Toy Font API vs. Scaled Font API
//
// Cairo provides two text rendering APIs:
//
// The toy font API (this package) offers simple text rendering with minimal
// configuration. It is sufficient for basic needs such as labels, annotations,
// and simple overlays. The toy API selects fonts by family name, slant, and
// weight using the host platform's font system, so results are
// platform-dependent—the same code may render differently on macOS, Linux,
// and Windows due to differences in available fonts, hinting engines, and
// antialiasing strategies.
//
// The scaled font API (not yet implemented) provides fine-grained control over
// font metrics, glyph placement, and transformation. It is required for
// professional typography, internationalization, and precise text layout.
// Advanced text rendering that needs consistent cross-platform output typically
// uses Pango, which wraps Cairo's scaled font API (future integration).
//
// For most use cases involving simple text overlays, the toy font API is
// sufficient and much easier to use.
package font
