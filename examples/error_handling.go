// ABOUTME: Example demonstrating proper error handling patterns for the cairo library.
// ABOUTME: Shows errors.Is/As usage, typed error inspection, and drawing status checking.
package examples

import (
	"errors"
	"fmt"

	"github.com/mikowitz/cairo"
	"github.com/mikowitz/cairo/status"
)

// ErrorHandlingResult captures results from the error handling demonstration.
type ErrorHandlingResult struct {
	// InvalidFormatErr is the error from creating a surface with an invalid format.
	InvalidFormatErr error
	// NilSurfaceErr is the error from creating a context with a nil surface.
	NilSurfaceErr error
	// NoCurrentPointErr is the error from querying current point when no path is active.
	NoCurrentPointErr error
}

// DemonstrateErrorHandling exercises the cairo error handling APIs and returns
// the errors encountered so callers can inspect them.
//
// This example shows three error handling patterns:
//
//  1. Constructor errors: NewImageSurface, NewContext return (value, error).
//     Use errors.Is to check the specific status and errors.As to inspect the
//     typed error fields (SurfaceType, Operation, PatternType).
//
//  2. Drawing status: drawing operations do not return errors. After a sequence
//     of drawing calls, inspect ctx.Status() to detect failures.
//
//  3. Getter errors: methods like GetCurrentPoint return (value, error) when
//     the operation can fail due to invalid state.
func DemonstrateErrorHandling() (ErrorHandlingResult, error) {
	var result ErrorHandlingResult

	// --- Pattern 1: constructor errors with errors.Is and errors.As ---
	// Attempt to create a surface with an invalid format.
	// NewImageSurface wraps the failure in a *cairo.SurfaceError.
	_, err := cairo.NewImageSurface(cairo.Format(-1), 100, 100)
	if err != nil {
		result.InvalidFormatErr = err

		// errors.Is traverses the error chain; SurfaceError.Unwrap returns the
		// underlying status.Status, so this check reaches the raw status code.
		if !errors.Is(err, status.InvalidFormat) {
			return result, fmt.Errorf("expected InvalidFormat status, got: %w", err)
		}

		// errors.As extracts the typed wrapper to read the SurfaceType field.
		var surfErr *cairo.SurfaceError
		if errors.As(err, &surfErr) {
			_ = surfErr.SurfaceType // "image" — identifies which surface type failed
		}
	}

	// Attempt to create a context with a nil surface.
	// NewContext wraps the failure in a *cairo.ContextError.
	_, err = cairo.NewContext(nil)
	if err != nil {
		result.NilSurfaceErr = err

		if !errors.Is(err, status.NullPointer) {
			return result, fmt.Errorf("expected NullPointer status, got: %w", err)
		}

		var ctxErr *cairo.ContextError
		if errors.As(err, &ctxErr) {
			_ = ctxErr.Operation // "create" — identifies the failing operation
		}
	}

	// Create a valid surface and context for the remaining patterns.
	surf, err := cairo.NewImageSurface(cairo.FormatARGB32, 200, 200)
	if err != nil {
		return result, fmt.Errorf("creating surface: %w", err)
	}
	defer func() { _ = surf.Close() }()

	ctx, err := cairo.NewContext(surf)
	if err != nil {
		return result, fmt.Errorf("creating context: %w", err)
	}
	defer func() { _ = ctx.Close() }()

	// --- Pattern 2: drawing operation status checking ---
	// Drawing calls like MoveTo, LineTo, Stroke do not return errors.
	// Check ctx.Status() after a sequence to detect any failure.
	ctx.MoveTo(10, 10)
	ctx.LineTo(190, 190)
	ctx.SetLineWidth(2)
	ctx.Stroke()

	if s := ctx.Status(); s != status.Success {
		return result, fmt.Errorf("drawing failed: %v", s)
	}

	// --- Pattern 3: getter errors for invalid state ---
	// Stroke clears the current path and current point, so GetCurrentPoint
	// now returns NoCurrentPoint — the same failure as before any MoveTo.
	_, _, err = ctx.GetCurrentPoint()
	if err != nil {
		result.NoCurrentPointErr = err

		if !errors.Is(err, status.NoCurrentPoint) {
			return result, fmt.Errorf("expected NoCurrentPoint status, got: %w", err)
		}
	}

	return result, nil
}
