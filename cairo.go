package cairo

import "github.com/mikowitz/cairo/surface"

// Format is used to identify the memory format of image data.
type Format = surface.Format

// Format constants specify the pixel format for image surfaces.
const (
	// FormatInvalid represents an invalid format.
	FormatInvalid = surface.FormatInvalid

	// FormatARGB32 represents 32-bit ARGB format with alpha in the most significant byte,
	// then red, green, and blue. The 32-bit quantities are stored native-endian.
	// Pre-multiplied alpha is used (i.e. 50% transparent red is 0x80800000, not 0x80ff0000).
	FormatARGB32 = surface.FormatARGB32

	// FormatRGB24 represents 24-bit RGB format stored in 32-bit quantities. The unused
	// bits should be zero. The 32-bit quantities are stored native-endian.
	FormatRGB24 = surface.FormatRGB24

	// FormatA8 represents 8-bit alpha-only format.
	FormatA8 = surface.FormatA8

	// FormatA1 represents 1-bit alpha-only format.
	FormatA1 = surface.FormatA1

	// FormatRGB16_565 represents 16-bit RGB format with 5 bits for red, 6 for green,
	// and 5 for blue.
	FormatRGB16_565 = surface.FormatRGB16_565

	// FormatRGB30 represents 30-bit RGB format with 10 bits per color component,
	// stored in 32-bit quantities.
	FormatRGB30 = surface.FormatRGB30
)

// Surface represents a destination for drawing operations.
//
// A Surface is the abstract type representing all different drawing targets
// that Cairo can render to. The actual drawing is done using a Context
// (available in later prompts).
//
// All Surface implementations must be explicitly closed when finished to free
// Cairo resources, or the finalizer will clean them up during garbage collection.
type Surface = surface.Surface

// NewImageSurface creates an image surface of the specified format and dimensions.
// The initial contents of the surface are set to transparent black (all pixels are
// fully transparent with RGBA values of 0,0,0,0).
//
// The surface should be closed with Close() when finished to release Cairo resources.
// A finalizer is registered as a safety net, but explicit cleanup is recommended.
func NewImageSurface(format Format, width, height int) (*surface.ImageSurface, error) {
	return surface.NewImageSurface(format, width, height)
}
