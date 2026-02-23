// ABOUTME: Tests for fill rule support in Cairo drawing contexts.
// ABOUTME: Covers SetFillRule, GetFillRule, and the default FillRuleWinding behavior.

package context

import (
	"testing"

	"github.com/mikowitz/cairo/surface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestContextFillRuleDefault verifies that the default fill rule is FillRuleWinding.
func TestContextFillRuleDefault(t *testing.T) {
	ctx, _ := newTestContext(t, 100, 100)

	assert.Equal(t, FillRuleWinding, ctx.GetFillRule())
}

// TestContextSetFillRule tests setting both fill rule variants.
func TestContextSetFillRule(t *testing.T) {
	tests := []struct {
		name     string
		fillRule FillRule
	}{
		{name: "FillRuleWinding", fillRule: FillRuleWinding},
		{name: "FillRuleEvenOdd", fillRule: FillRuleEvenOdd},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := newTestContext(t, 100, 100)

			ctx.SetFillRule(tt.fillRule)
			assert.Equal(t, tt.fillRule, ctx.GetFillRule())
		})
	}
}

// TestContextGetFillRuleAfterClose verifies GetFillRule returns FillRuleWinding on a closed context.
func TestContextGetFillRuleAfterClose(t *testing.T) {
	surf, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
	require.NoError(t, err)
	defer surf.Close()

	ctx, err := NewContext(surf)
	require.NoError(t, err)

	err = ctx.Close()
	require.NoError(t, err)

	assert.Equal(t, FillRuleWinding, ctx.GetFillRule())
}

// TestContextFillRuleString verifies the string representation of FillRule values.
func TestContextFillRuleString(t *testing.T) {
	assert.Equal(t, "FillRuleWinding", FillRuleWinding.String())
	assert.Equal(t, "FillRuleEvenOdd", FillRuleEvenOdd.String())
}
