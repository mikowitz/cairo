// Package cairo provides a Go wrapper around the Cairo 2D graphics library.
//
// Cairo is a 2D graphics library with support for multiple output devices.
// This package provides an idiomatic Go interface to Cairo's C API, with
// proper memory management, thread safety, and error handling.
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
// # Basic Usage Example
//
// The following example demonstrates creating a surface, drawing a simple shape,
// and saving the result to a PNG file:
//
//	// Create a 400x400 pixel image surface
//	surf, err := cairo.NewImageSurface(cairo.FormatARGB32, 400, 400)
//	if err != nil {
//		return err
//	}
//	defer surf.Close()
//
//	// TODO: Drawing operations will be added with Context package (Prompt 9)
//	// For now, we can work directly with the surface
//
//	// Ensure all drawing is flushed before writing
//	surf.Flush()
//
//	// Write the surface to a PNG file (available in Prompt 8)
//	// err = surf.WriteToPNG("output.png")
//	// if err != nil {
//	//     return err
//	// }
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
