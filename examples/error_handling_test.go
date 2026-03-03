// ABOUTME: Tests for the error handling example in the cairo library.
// ABOUTME: Verifies that error handling patterns work correctly and produce expected errors.
package examples

import (
	"errors"
	"testing"

	"github.com/mikowitz/cairo"
	"github.com/mikowitz/cairo/status"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDemonstrateErrorHandlingSucceeds(t *testing.T) {
	result, err := DemonstrateErrorHandling()
	require.NoError(t, err)

	assert.NotNil(t, result.InvalidFormatErr, "should have captured an invalid format error")
	assert.NotNil(t, result.NilSurfaceErr, "should have captured a nil surface error")
	assert.NotNil(t, result.NoCurrentPointErr, "should have captured a no current point error")
}

func TestInvalidFormatErrorIsSurfaceError(t *testing.T) {
	result, err := DemonstrateErrorHandling()
	require.NoError(t, err)

	var surfErr *cairo.SurfaceError
	require.True(t, errors.As(result.InvalidFormatErr, &surfErr), "invalid format error should be a *SurfaceError")
	assert.Equal(t, "image", surfErr.SurfaceType)
	assert.Equal(t, status.InvalidFormat, surfErr.Status)
}

func TestInvalidFormatErrorIsStatusCheck(t *testing.T) {
	result, err := DemonstrateErrorHandling()
	require.NoError(t, err)

	assert.True(t, errors.Is(result.InvalidFormatErr, status.InvalidFormat),
		"errors.Is should match InvalidFormat through the SurfaceError chain")
}

func TestNilSurfaceErrorIsContextError(t *testing.T) {
	result, err := DemonstrateErrorHandling()
	require.NoError(t, err)

	var ctxErr *cairo.ContextError
	require.True(t, errors.As(result.NilSurfaceErr, &ctxErr), "nil surface error should be a *ContextError")
	assert.Equal(t, "create", ctxErr.Operation)
	assert.Equal(t, status.NullPointer, ctxErr.Status)
}

func TestNilSurfaceErrorIsStatusCheck(t *testing.T) {
	result, err := DemonstrateErrorHandling()
	require.NoError(t, err)

	assert.True(t, errors.Is(result.NilSurfaceErr, status.NullPointer),
		"errors.Is should match NullPointer through the ContextError chain")
}

func TestNoCurrentPointErrorIsStatusCheck(t *testing.T) {
	result, err := DemonstrateErrorHandling()
	require.NoError(t, err)

	assert.True(t, errors.Is(result.NoCurrentPointErr, status.NoCurrentPoint),
		"errors.Is should match NoCurrentPoint status")
}

func TestErrorMessageIncludesContext(t *testing.T) {
	result, err := DemonstrateErrorHandling()
	require.NoError(t, err)

	assert.Contains(t, result.InvalidFormatErr.Error(), "image",
		"surface error message should include the surface type")
	assert.Contains(t, result.NilSurfaceErr.Error(), "create",
		"context error message should include the operation name")
}
