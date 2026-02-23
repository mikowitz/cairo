// ABOUTME: Tests for compositing operator support in Cairo drawing contexts.
// ABOUTME: Covers SetOperator, GetOperator, and the default OperatorOver behavior.

package context

import (
	"testing"

	"github.com/mikowitz/cairo/surface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestContextOperatorDefault verifies that the default compositing operator is OperatorOver.
func TestContextOperatorDefault(t *testing.T) {
	ctx := newTestContext(t, 100, 100)

	assert.Equal(t, OperatorOver, ctx.GetOperator())
}

// TestContextSetOperator tests setting various compositing operators.
func TestContextSetOperator(t *testing.T) {
	tests := []struct {
		name string
		op   Operator
	}{
		{name: "OperatorClear", op: OperatorClear},
		{name: "OperatorSource", op: OperatorSource},
		{name: "OperatorOver", op: OperatorOver},
		{name: "OperatorIn", op: OperatorIn},
		{name: "OperatorOut", op: OperatorOut},
		{name: "OperatorAtop", op: OperatorAtop},
		{name: "OperatorDest", op: OperatorDest},
		{name: "OperatorDestOver", op: OperatorDestOver},
		{name: "OperatorDestIn", op: OperatorDestIn},
		{name: "OperatorDestOut", op: OperatorDestOut},
		{name: "OperatorDestAtop", op: OperatorDestAtop},
		{name: "OperatorXor", op: OperatorXor},
		{name: "OperatorAdd", op: OperatorAdd},
		{name: "OperatorSaturate", op: OperatorSaturate},
		{name: "OperatorMultiply", op: OperatorMultiply},
		{name: "OperatorScreen", op: OperatorScreen},
		{name: "OperatorOverlay", op: OperatorOverlay},
		{name: "OperatorDarken", op: OperatorDarken},
		{name: "OperatorLighten", op: OperatorLighten},
		{name: "OperatorColorDodge", op: OperatorColorDodge},
		{name: "OperatorColorBurn", op: OperatorColorBurn},
		{name: "OperatorHardLight", op: OperatorHardLight},
		{name: "OperatorSoftLight", op: OperatorSoftLight},
		{name: "OperatorDifference", op: OperatorDifference},
		{name: "OperatorExclusion", op: OperatorExclusion},
		{name: "OperatorHslHue", op: OperatorHslHue},
		{name: "OperatorHslSaturation", op: OperatorHslSaturation},
		{name: "OperatorHslColor", op: OperatorHslColor},
		{name: "OperatorHslLuminosity", op: OperatorHslLuminosity},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := newTestContext(t, 100, 100)

			ctx.SetOperator(tt.op)
			assert.Equal(t, tt.op, ctx.GetOperator())
		})
	}
}

// TestContextGetOperator tests getting the current operator after setting it.
func TestContextGetOperator(t *testing.T) {
	ctx := newTestContext(t, 100, 100)

	ctx.SetOperator(OperatorAdd)
	op := ctx.GetOperator()
	assert.Equal(t, OperatorAdd, op)

	ctx.SetOperator(OperatorMultiply)
	op = ctx.GetOperator()
	assert.Equal(t, OperatorMultiply, op)
}

// TestContextGetOperatorAfterClose verifies GetOperator returns OperatorOver on a closed context.
func TestContextGetOperatorAfterClose(t *testing.T) {
	surf, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
	require.NoError(t, err)
	defer surf.Close()

	ctx, err := NewContext(surf)
	require.NoError(t, err)

	err = ctx.Close()
	require.NoError(t, err)

	assert.Equal(t, OperatorOver, ctx.GetOperator())
}
