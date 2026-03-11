// ABOUTME: Custom error types for Cairo surface, context, and pattern operations.
// ABOUTME: Provides structured errors wrapping status codes with additional context.
package cairo

import (
	"fmt"

	"github.com/mikowitz/cairo/status"
)

// SurfaceError represents an error that occurred during a surface operation.
// It wraps a status.Status value and includes the surface type for additional context.
type SurfaceError struct {
	// Status is the underlying Cairo status code.
	Status status.Status
	// SurfaceType identifies the kind of surface (e.g., "image", "pdf", "svg").
	SurfaceType string
}

// Error implements the error interface.
func (e *SurfaceError) Error() string {
	if e.SurfaceType != "" {
		return fmt.Sprintf("cairo surface error (%s): %v", e.SurfaceType, e.Status)
	}
	return fmt.Sprintf("cairo surface error: %v", e.Status)
}

// Unwrap returns the underlying status error for use with errors.Is and errors.As.
func (e *SurfaceError) Unwrap() error {
	return e.Status
}

// Is reports whether target matches this error.
// It matches if target is a *SurfaceError with the same Status,
// and either target.SurfaceType is empty or equals e.SurfaceType.
func (e *SurfaceError) Is(target error) bool {
	t, ok := target.(*SurfaceError)
	if !ok {
		return false
	}
	if t.SurfaceType != "" && t.SurfaceType != e.SurfaceType {
		return false
	}
	return e.Status == t.Status
}

// ContextError represents an error that occurred during a context drawing operation.
// It wraps a status.Status value and includes the operation name for additional context.
type ContextError struct {
	// Status is the underlying Cairo status code.
	Status status.Status
	// Operation names the drawing operation that failed (e.g., "stroke", "fill").
	Operation string
}

// Error implements the error interface.
func (e *ContextError) Error() string {
	if e.Operation != "" {
		return fmt.Sprintf("cairo context error (%s): %v", e.Operation, e.Status)
	}
	return fmt.Sprintf("cairo context error: %v", e.Status)
}

// Unwrap returns the underlying status error for use with errors.Is and errors.As.
func (e *ContextError) Unwrap() error {
	return e.Status
}

// Is reports whether target matches this error.
// It matches if target is a *ContextError with the same Status,
// and either target.Operation is empty or equals e.Operation.
func (e *ContextError) Is(target error) bool {
	t, ok := target.(*ContextError)
	if !ok {
		return false
	}
	if t.Operation != "" && t.Operation != e.Operation {
		return false
	}
	return e.Status == t.Status
}

// PatternError represents an error that occurred during a pattern operation.
// It wraps a status.Status value and includes the pattern type for additional context.
type PatternError struct {
	// Status is the underlying Cairo status code.
	Status status.Status
	// PatternType identifies the kind of pattern (e.g., "solid", "linear", "radial").
	PatternType string
}

// Error implements the error interface.
func (e *PatternError) Error() string {
	if e.PatternType != "" {
		return fmt.Sprintf("cairo pattern error (%s): %v", e.PatternType, e.Status)
	}
	return fmt.Sprintf("cairo pattern error: %v", e.Status)
}

// Unwrap returns the underlying status error for use with errors.Is and errors.As.
func (e *PatternError) Unwrap() error {
	return e.Status
}

// Is reports whether target matches this error.
// It matches if target is a *PatternError with the same Status,
// and either target.PatternType is empty or equals e.PatternType.
func (e *PatternError) Is(target error) bool {
	t, ok := target.(*PatternError)
	if !ok {
		return false
	}
	if t.PatternType != "" && t.PatternType != e.PatternType {
		return false
	}
	return e.Status == t.Status
}

// wrapSurfaceErr converts a status.Status error into a *SurfaceError with the
// given surface type. Non-status errors pass through unchanged.
func wrapSurfaceErr(err error, surfaceType string) error {
	if st, ok := err.(status.Status); ok {
		return &SurfaceError{Status: st, SurfaceType: surfaceType}
	}
	return err
}

// wrapContextErr converts a status.Status error into a *ContextError with the
// given operation name. Non-status errors pass through unchanged.
func wrapContextErr(err error, operation string) error {
	if st, ok := err.(status.Status); ok {
		return &ContextError{Status: st, Operation: operation}
	}
	return err
}

// wrapPatternErr converts a status.Status error into a *PatternError with the
// given pattern type. Non-status errors pass through unchanged.
func wrapPatternErr(err error, patternType string) error {
	if st, ok := err.(status.Status); ok {
		return &PatternError{Status: st, PatternType: patternType}
	}
	return err
}
