package pattern

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestPatternTypeConstants verifies all PatternType constants are defined.
func TestPatternTypeConstants(t *testing.T) {
	tests := []struct {
		name        string
		patternType PatternType
		expected    int
	}{
		{"Solid", PatternTypeSolid, 0},
		{"Surface", PatternTypeSurface, 1},
		{"Linear", PatternTypeLinear, 2},
		{"Radial", PatternTypeRadial, 3},
		{"Mesh", PatternTypeMesh, 4},
		{"RasterSource", PatternTypeRasterSource, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, int(tt.patternType),
				"PatternType%s should have value %d", tt.name, tt.expected)
		})
	}
}

// TestPatternTypeString verifies the String() method returns correct values.
func TestPatternTypeString(t *testing.T) {
	tests := []struct {
		name        string
		patternType PatternType
	}{
		{"Solid", PatternTypeSolid},
		{"Surface", PatternTypeSurface},
		{"Linear", PatternTypeLinear},
		{"Radial", PatternTypeRadial},
		{"Mesh", PatternTypeMesh},
		{"RasterSource", PatternTypeRasterSource},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str := tt.patternType.String()
			assert.Equal(t, tt.name, str,
				"PatternType%s.String() should return %q", tt.name, tt.name)
		})
	}
}

// TestPatternTypeStringInvalid verifies String() handles invalid types gracefully.
func TestPatternTypeStringInvalid(t *testing.T) {
	tests := []struct {
		name        string
		patternType PatternType
	}{
		{"negative", PatternType(-1)},
		{"large_value", PatternType(999)},
		{"out_of_range", PatternType(100)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str := tt.patternType.String()
			// stringer typically returns "PatternType(N)" for unknown values
			assert.NotEmpty(t, str, "String() should return non-empty string for invalid type")
			assert.Contains(t, str, "PatternType", "Invalid type string should contain 'PatternType'")
		})
	}
}

// TestPatternTypeUniqueness verifies all PatternType constants are unique.
func TestPatternTypeUniqueness(t *testing.T) {
	types := []PatternType{
		PatternTypeSolid,
		PatternTypeSurface,
		PatternTypeLinear,
		PatternTypeRadial,
		PatternTypeMesh,
		PatternTypeRasterSource,
	}

	// Check that all values are unique
	seen := make(map[PatternType]bool)
	for _, pt := range types {
		assert.False(t, seen[pt], "PatternType value %d should be unique", pt)
		seen[pt] = true
	}

	// Should have 6 unique types
	assert.Equal(t, 6, len(seen), "Should have 6 unique PatternType values")
}

// TestPatternTypeComparison verifies PatternType values can be compared.
func TestPatternTypeComparison(t *testing.T) {
	t.Run("equality", func(t *testing.T) {
		assert.Equal(t, PatternTypeSolid, PatternTypeSolid, "Same types should be equal")
		assert.NotEqual(t, PatternTypeSolid, PatternTypeSurface, "Different types should not be equal")
	})

	t.Run("ordering", func(t *testing.T) {
		assert.True(t, PatternTypeSolid < PatternTypeSurface, "PatternTypeSolid should be less than PatternTypeSurface")
		assert.True(t, PatternTypeLinear < PatternTypeRadial, "PatternTypeLinear should be less than PatternTypeRadial")
		assert.True(t, PatternTypeRasterSource > PatternTypeSolid, "PatternTypeRasterSource should be greater than PatternTypeSolid")
	})

	t.Run("switch_statement", func(t *testing.T) {
		// Verify PatternType can be used in switch statements
		testType := PatternTypeLinear
		var result string

		switch testType {
		case PatternTypeSolid:
			result = "solid"
		case PatternTypeSurface:
			result = "surface"
		case PatternTypeLinear:
			result = "linear"
		case PatternTypeRadial:
			result = "radial"
		case PatternTypeMesh:
			result = "mesh"
		case PatternTypeRasterSource:
			result = "raster"
		default:
			result = "unknown"
		}

		assert.Equal(t, "linear", result, "Switch should match PatternTypeLinear")
	})
}

// TestPatternTypeZeroValue verifies the zero value behavior.
func TestPatternTypeZeroValue(t *testing.T) {
	var pt PatternType

	// Zero value should be PatternTypeSolid (0)
	assert.Equal(t, PatternTypeSolid, pt, "Zero value should be PatternTypeSolid")
	assert.Equal(t, 0, int(pt), "Zero value should be 0")
	assert.Equal(t, "Solid", pt.String(), "Zero value String() should return 'Solid'")
}

// TestPatternTypeInMap verifies PatternType can be used as map key.
func TestPatternTypeInMap(t *testing.T) {
	typeMap := make(map[PatternType]string)

	// Add entries
	typeMap[PatternTypeSolid] = "solid"
	typeMap[PatternTypeSurface] = "surface"
	typeMap[PatternTypeLinear] = "linear"
	typeMap[PatternTypeRadial] = "radial"
	typeMap[PatternTypeMesh] = "mesh"
	typeMap[PatternTypeRasterSource] = "raster"

	// Verify retrieval
	assert.Equal(t, "solid", typeMap[PatternTypeSolid])
	assert.Equal(t, "surface", typeMap[PatternTypeSurface])
	assert.Equal(t, "linear", typeMap[PatternTypeLinear])
	assert.Equal(t, "radial", typeMap[PatternTypeRadial])
	assert.Equal(t, "mesh", typeMap[PatternTypeMesh])
	assert.Equal(t, "raster", typeMap[PatternTypeRasterSource])

	// Map should have 6 entries
	assert.Equal(t, 6, len(typeMap))
}

// TestPatternTypeInSlice verifies PatternType can be used in slices.
func TestPatternTypeInSlice(t *testing.T) {
	types := []PatternType{
		PatternTypeSolid,
		PatternTypeSurface,
		PatternTypeLinear,
		PatternTypeRadial,
		PatternTypeMesh,
		PatternTypeRasterSource,
	}

	assert.Equal(t, 6, len(types), "Should have 6 types in slice")
	assert.Equal(t, PatternTypeSolid, types[0])
	assert.Equal(t, PatternTypeRasterSource, types[5])

	// Verify iteration
	count := 0
	for _, pt := range types {
		assert.True(t, pt >= PatternTypeSolid && pt <= PatternTypeRasterSource,
			"PatternType %v should be in valid range", pt)
		count++
	}
	assert.Equal(t, 6, count)
}

// TestPatternTypeDefaultInStruct verifies PatternType as struct field.
func TestPatternTypeDefaultInStruct(t *testing.T) {
	type TestStruct struct {
		Type PatternType
		Name string
	}

	// Zero value struct
	var s TestStruct
	assert.Equal(t, PatternTypeSolid, s.Type, "Default PatternType in struct should be PatternTypeSolid")

	// Initialized struct
	s2 := TestStruct{Type: PatternTypeLinear, Name: "test"}
	assert.Equal(t, PatternTypeLinear, s2.Type)
	assert.Equal(t, "test", s2.Name)
}

// TestPatternTypeSequence verifies PatternType constants are sequential.
func TestPatternTypeSequence(t *testing.T) {
	assert.Equal(t, 0, int(PatternTypeSolid), "PatternTypeSolid should be 0")
	assert.Equal(t, 1, int(PatternTypeSurface), "PatternTypeSurface should be 1")
	assert.Equal(t, 2, int(PatternTypeLinear), "PatternTypeLinear should be 2")
	assert.Equal(t, 3, int(PatternTypeRadial), "PatternTypeRadial should be 3")
	assert.Equal(t, 4, int(PatternTypeMesh), "PatternTypeMesh should be 4")
	assert.Equal(t, 5, int(PatternTypeRasterSource), "PatternTypeRasterSource should be 5")

	// Verify they are consecutive
	assert.Equal(t, int(PatternTypeSurface), int(PatternTypeSolid)+1)
	assert.Equal(t, int(PatternTypeLinear), int(PatternTypeSurface)+1)
	assert.Equal(t, int(PatternTypeRadial), int(PatternTypeLinear)+1)
	assert.Equal(t, int(PatternTypeMesh), int(PatternTypeRadial)+1)
	assert.Equal(t, int(PatternTypeRasterSource), int(PatternTypeMesh)+1)
}

// TestPatternTypeIsValid verifies a helper to check valid types (if implemented).
func TestPatternTypeIsValid(t *testing.T) {
	t.Run("valid_types", func(t *testing.T) {
		validTypes := []PatternType{
			PatternTypeSolid,
			PatternTypeSurface,
			PatternTypeLinear,
			PatternTypeRadial,
			PatternTypeMesh,
			PatternTypeRasterSource,
		}

		for _, pt := range validTypes {
			// Check if type is in valid range
			isValid := pt >= PatternTypeSolid && pt <= PatternTypeRasterSource
			assert.True(t, isValid, "PatternType %v should be valid", pt)
		}
	})

	t.Run("invalid_types", func(t *testing.T) {
		invalidTypes := []PatternType{
			PatternType(-1),
			PatternType(6),
			PatternType(100),
			PatternType(999),
		}

		for _, pt := range invalidTypes {
			// Check if type is out of valid range
			isValid := pt >= PatternTypeSolid && pt <= PatternTypeRasterSource
			assert.False(t, isValid, "PatternType %v should be invalid", pt)
		}
	})
}

// TestPatternTypeStringFormatting verifies String() output format.
func TestPatternTypeStringFormatting(t *testing.T) {
	tests := []struct {
		name        string
		patternType PatternType
	}{
		{"Solid", PatternTypeSolid},
		{"Surface", PatternTypeSurface},
		{"Linear", PatternTypeLinear},
		{"Radial", PatternTypeRadial},
		{"Mesh", PatternTypeMesh},
		{"RasterSource", PatternTypeRasterSource},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str := tt.patternType.String()
			assert.NotEmpty(t, str, "String() should not be empty")
			assert.Contains(t, str, tt.name, "String() should contain prefix")

			// Should not have spaces or special characters
			assert.NotContains(t, str, " ", "String() should not contain spaces")
			assert.NotContains(t, str, "\n", "String() should not contain newlines")
		})
	}
}
