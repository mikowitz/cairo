// Package cairo provides a Go wrapper around the Cairo 2D graphics library.
//
// Cairo is a 2D graphics library with support for multiple output devices.
// This package provides an idiomatic Go interface to Cairo's C API, with
// proper memory management, thread safety, and error handling.
//
// # Architecture Overview
//
// This package is organized into several subpackages, each handling a specific
// aspect of Cairo's functionality:
//
//   - status: Error handling and status codes from Cairo operations
//   - matrix: 2D affine transformations for coordinate space conversions
//   - surface: Drawing targets (image buffers, PDF files, SVG files, etc.)
//   - context: The main drawing interface with graphics state and operations
//   - pattern: Sources for drawing operations (colors, gradients, images)
//   - font: Text rendering support (future)
//
// The typical usage flow is:
//
//  1. Create a Surface (your drawing target)
//  2. Create a Context associated with that Surface
//  3. Use Context methods to draw (set colors, create paths, fill/stroke)
//  4. Flush the Surface and save/export the result
//  5. Close resources when done (or use defer)
//
// All packages use CGO to interface with the underlying Cairo C library while
// presenting an idiomatic Go API. Memory management is handled automatically
// through finalizers, though explicit cleanup via Close() is recommended for
// deterministic resource release.
//
// Thread safety is built into all types via sync.RWMutex, making it safe to
// use Cairo objects from multiple goroutines. However, sharing a single context
// or surface across goroutines may result in lock contention affecting performance.
//
// # Surface Types
//
// Cairo draws to surfaces, which represent a destination for graphics operations.
// The most commonly used surface type is ImageSurface, which represents an image
// buffer in memory. Other surface types include PDF and SVG surfaces for vector
// output.
//
// # Resource Management
//
// All Cairo resources (surfaces, contexts, patterns) should be explicitly closed
// when no longer needed by calling their Close() method. While finalizers are
// registered as a safety net to prevent memory leaks, relying on them can delay
// resource cleanup. For best practices, use defer to ensure resources are cleaned up:
//
//	surf, err := cairo.NewImageSurface(cairo.FormatARGB32, 640, 480)
//	if err != nil {
//		return err
//	}
//	defer surf.Close()
//
// # Thread Safety
//
// All types in this package are safe for concurrent use. However, sharing a
// single surface or context across multiple goroutines may result in contention.
// For best performance, consider creating separate surfaces/contexts per goroutine
// when performing concurrent rendering.
//
// # PNG Export
//
// ImageSurface supports writing to PNG files via the WriteToPNG method.
// The surface should be flushed before writing to ensure all drawing operations
// are complete:
//
//	surf, err := cairo.NewImageSurface(cairo.FormatARGB32, 640, 480)
//	if err != nil {
//		return err
//	}
//	defer surf.Close()
//
//	// ... perform drawing operations ...
//
//	surf.Flush() // Ensure all drawing is complete
//	err = surf.WriteToPNG("output.png")
//	if err != nil {
//		return err
//	}
//
// PNG export is available for all image surface formats (ARGB32, RGB24, A8, A1,
// RGB16_565, RGB30). The resulting PNG file will preserve the pixel data according
// to the surface's format.
//
// # Basic Usage Example
//
// The following example demonstrates creating a surface and saving it to a PNG file:
//
//	// Create a 400x400 pixel image surface
//	surf, err := cairo.NewImageSurface(cairo.FormatARGB32, 400, 400)
//	if err != nil {
//		return err
//	}
//	defer surf.Close()
//
//	// Drawing operations will be added with Context package (Prompt 9)
//	// For now, we can create and save an empty surface
//
//	// Ensure all drawing is flushed before writing
//	surf.Flush()
//
//	// Write the surface to a PNG file
//	err = surf.WriteToPNG("output.png")
//	if err != nil {
//		return err
//	}
//
// # Error Handling
//
// Functions that can fail return an error as their last return value. These errors
// typically wrap a Status value from Cairo. Check errors using standard Go patterns:
//
//	surf, err := cairo.NewImageSurface(cairo.FormatARGB32, -100, 100)
//	if err != nil {
//		// Handle invalid dimensions
//		return err
//	}
package cairo
