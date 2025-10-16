# Go-Cairo Implementation Prompts

## Overview

This document contains a complete series of 34 prompts for implementing a Go
wrapper around the Cairo graphics library. Each prompt builds incrementally
on previous work, ensuring safe, test-driven development with no orphaned code.

## Implementation Strategy

### Key Principles

- **Incremental Progress**: Each step compiles and passes tests
- **Test-Driven Development**: Tests written alongside implementation
- **No Orphaned Code**: Everything integrates immediately
- **Early Validation**: MVP complete early to validate approach
- **Comprehensive Coverage**: >80% test coverage goal

### Phase Breakdown

1. **Foundation Setup** (Prompts 1-3): Project structure, build tooling, status package
2. **Core Types** (Prompts 4-8): Matrix, surfaces, basic lifecycle
3. **Context & Drawing** (Prompts 9-12): Drawing context with basic operations
4. **MVP Validation** (Prompt 13): Complete rectangle-to-PNG example
5. **Expansion** (Prompts 14-34): Additional features, optimization, polish

### Key Milestones

- **Prompt 13**: MVP Complete (draw rectangle to PNG)
- **Prompt 21**: Core drawing features complete
- **Prompt 27**: Vector output formats complete
- **Prompt 34**: v0.1.0 release ready

---

## Prompt 1: Project Foundation

```markdown
Create the initial project structure for a Go wrapper around the Cairo graphics library.

Requirements:
1. Initialize a Go module at github.com/username/cairo
2. Create the following directory structure:
   - cairo/ (root)
   - status/
   - matrix/
   - surface/
   - context/
   - pattern/
   - font/
   - examples/
3. Create a go.mod file requiring Go 1.23
4. Create a basic README.md explaining this is a CGO wrapper for Cairo
5. Create a Taskfile.yaml with tasks for: test, generate, lint
6. Create .gitignore for Go projects (binaries, coverage, IDE files)
7. Add a simple DESIGN.md that links to the Cairo C API docs

Do not write any Go code yet. Just set up the project structure and build tooling.

The Taskfile should have these tasks:
- test: run `go test -race ./...`
- generate: run `go generate ./...`
- lint: run `golangci-lint run` (document that it needs to be installed)

Ensure all markdown files are properly formatted and professional.
```

---

## Prompt 2: Status Package - Types and Constants

```markdown
Implement the status package which handles Cairo error codes.

Context: This is step 1 of building a Go wrapper for Cairo. The status package
has no dependencies and provides error handling primitives for all other packages.

Requirements:
1. Create status/status.go with:
   - Status type as an int
   - Constants for all cairo_status_t values (at minimum: StatusSuccess, StatusNoMemory, StatusInvalidRestore, StatusInvalidPopGroup, StatusNoCurrentPoint, StatusInvalidMatrix, StatusInvalidStatus, StatusNullPointer, StatusInvalidString, StatusInvalidPathData, StatusReadError, StatusWriteError, StatusSurfaceFinished, StatusSurfaceTypeMismatch, StatusPatternTypeMismatch, StatusInvalidContent, StatusInvalidFormat, StatusInvalidVisual, StatusFileNotFound, StatusInvalidDash, StatusInvalidDscComment, StatusInvalidIndex, StatusClipNotRepresentable)
   - Add go:generate comment for stringer: `//go:generate stringer -type=Status`
   - func (s Status) Error() string method that returns s.String()
   - func (s Status) ToError() error that returns nil for StatusSuccess, otherwise returns s as error

2. Create status/status_test.go with:
   - TestStatusError: verify Error() returns non-empty string for all non-success statuses
   - TestStatusToError: verify ToError() returns nil for StatusSuccess and error for others
   - TestStatusSuccess: verify StatusSuccess equals 0

All code must:
- Include package documentation
- Follow Go best practices
- Have complete test coverage
- Compile without warnings

Do not implement CGO yet - these are just Go type definitions.
```

---

## Prompt 3: Status Package - CGO Integration

```markdown
Add CGO bindings to the status package to enable conversion from Cairo's C status codes.

Context: Building on the previous status package, we now add the CGO layer to interface with Cairo's C library.

Requirements:
1. Create status/status_cgo.go with:
   - CGO preamble with #include <cairo.h> and #cgo pkg-config: cairo
   - import "C"
   - func statusFromC(cStatus C.cairo_status_t) Status that converts C status to Go Status
   - func (s Status) toC() C.cairo_status_t that converts Go Status to C status

2. Update status/status_test.go to add:
   - TestCGOStatusConversion: verify round-trip conversion between C and Go status codes
   - Test at least StatusSuccess, StatusNoMemory, StatusInvalidRestore
   - Use C.cairo_status_to_string() to verify we're mapping correctly

3. Add build tags if needed for platform-specific CGO flags

Requirements:
- All tests must pass
- CGO code must compile on Linux and macOS
- Add comments explaining the CGO integration
- Keep the CGO surface area minimal - just conversion functions

Testing note: The test should verify that our Go Status constants match Cairo's
C values.
```

---

## Prompt 4: Matrix Package - Structure and Basic Operations

```markdown
Implement the matrix package for 2D affine transformations.

Context: Matrices are used throughout Cairo for transformations. This package
wraps cairo_matrix_t.

Requirements:
1. Create matrix/matrix.go with:
   - Matrix struct with:
     - sync.RWMutex embedded (for thread safety)
     - XX, YX, XY, YY, X0, Y0 float64 fields matching cairo_matrix_t layout
   - func NewMatrix() *Matrix - returns zero matrix
   - func NewMatrixIdentity() *Matrix - returns identity matrix
   - func (m *Matrix) Init(xx, yx, xy, yy, x0, y0 float64) - initializes matrix
   - func (m *Matrix) InitIdentity() - sets to identity
   - func (m *Matrix) String() string - returns formatted matrix values
   - Package documentation explaining affine transformations

2. Create matrix/matrix_cgo.go with:
   - CGO preamble: #include <cairo.h>, #cgo pkg-config: cairo
   - func (m *Matrix) toC() *C.cairo_matrix_t - converts Go matrix to C
   - func matrixFromC(cm *C.cairo_matrix_t) *Matrix - converts C matrix to Go
   - Internal helper to sync Go fields with C struct

3. Create matrix/matrix_test.go with:
   - TestNewMatrix: verify zero matrix
   - TestNewMatrixIdentity: verify identity matrix (diagonal 1s, rest 0s)
   - TestMatrixInit: verify Init sets fields correctly
   - TestMatrixInitIdentity: verify InitIdentity sets correct values
   - TestMatrixThreadSafety: concurrent reads/writes don't race

All matrix methods must use appropriate locking (Lock/RLock).
Matrix must not allocate C memory yet - just provide conversion functions.
```

---

## Prompt 5: Matrix Package - Transformations

```markdown
Add transformation operations to the matrix package.

Context: Building on the basic Matrix structure, add operations for combining and applying transformations.

Requirements:
1. Update matrix/matrix.go to add:
   - func (m *Matrix) Multiply(other *Matrix) - post-multiply this matrix by other (m = m * other)
   - func (m *Matrix) TransformPoint(x, y float64) (float64, float64) - transform point by matrix
   - func (m *Matrix) TransformDistance(dx, dy float64) (float64, float64) - transform distance vector
   - func (m *Matrix) Translate(tx, ty float64) - apply translation
   - func (m *Matrix) Scale(sx, sy float64) - apply scaling
   - func (m *Matrix) Rotate(radians float64) - apply rotation
   - func (m *Matrix) Invert() error - invert matrix (returns error if singular)
   - All methods must use proper locking

2. Update matrix/matrix_cgo.go to add:
   - CGO wrappers calling cairo_matrix_multiply, cairo_matrix_transform_point, etc.
   - func (m *Matrix) invert() Status - wraps cairo_matrix_invert

3. Update matrix/matrix_test.go to add:
   - TestMatrixMultiply: verify matrix multiplication
   - TestMatrixTransformPoint: verify point transformation
   - TestMatrixTransformDistance: verify distance transformation
   - TestMatrixTranslate: verify translation matrix
   - TestMatrixScale: verify scaling matrix
   - TestMatrixRotate: verify rotation (test with 90 degrees for easy validation)
   - TestMatrixInvert: verify inversion and singular matrix error
   - TestMatrixOperationsCombined: test combining translate, scale, rotate

Ensure all operations match Cairo's semantics exactly.
Use the status package for error handling in Invert().
```

---

## Prompt 6: Surface Package - Interface and Base Types

```markdown
Create the surface package foundation with interfaces and base types.

Context: Surfaces are drawing targets in Cairo. This establishes the core surface abstraction before implementing specific surface types.

Requirements:
1. Create surface/surface.go with:
   - Surface interface with methods:
     - Close() error
     - Status() Status
     - Flush()
     - MarkDirty()
     - MarkDirtyRectangle(x, y, width, height int)
   - BaseSurface struct with:
     - sync.RWMutex embedded
     - ptr *C.cairo_surface_t (will be in CGO file)
     - closed bool flag
   - Package documentation explaining surface lifetime

2. Create surface/format.go with:
   - Format type as int
   - Constants: FormatInvalid, FormatARGB32, FormatRGB24, FormatA8, FormatA1, FormatRGB16_565, FormatRGB30
   - go:generate stringer -type=Format
   - func (f Format) StrideForWidth(width int) int - calculates stride

3. Create surface/surface_cgo.go with:
   - CGO preamble
   - BaseSurface CGO methods:
     - func (s *BaseSurface) close() error
     - func (s *BaseSurface) status() Status
     - func (s *BaseSurface) flush()
     - func (s *BaseSurface) markDirty()
     - func (s *BaseSurface) markDirtyRectangle(x, y, width, height int)

4. Create surface/surface_test.go with:
   - Tests for Format constants
   - TestFormatStrideForWidth: verify stride calculations
   - Test structure ready for surface implementations

The BaseSurface should implement Surface interface.
Include runtime finalizer setup pattern (document but don't implement yet).
All methods must check the closed flag and return appropriate errors.
```

---

## Prompt 7: Surface Package - ImageSurface Creation

```markdown
Implement ImageSurface creation and basic lifecycle management.

Context: ImageSurface is the simplest surface type and foundation for MVP. It stores pixels in memory.

Requirements:
1. Update surface/surface.go to add:
   - ImageSurface struct embedding BaseSurface
   - func NewImageSurface(format Format, width, height int) (*ImageSurface, error)
   - func (s *ImageSurface) GetFormat() Format
   - func (s *ImageSurface) GetWidth() int
   - func (s *ImageSurface) GetHeight() int
   - func (s *ImageSurface) GetStride() int
   - All methods use proper locking

2. Update surface/surface_cgo.go to add:
   - func newImageSurface(format Format, width, height int) (*ImageSurface, error)
     - Calls cairo_image_surface_create
     - Checks status
     - Sets up finalizer with runtime.SetFinalizer
     - Returns error on failure
   - CGO implementations for GetFormat, GetWidth, GetHeight, GetStride

3. Update surface/surface_test.go to add:
   - TestNewImageSurface: test successful creation
   - TestNewImageSurfaceInvalidFormat: test error handling
   - TestNewImageSurfaceInvalidSize: test error handling (zero/negative dimensions)
   - TestImageSurfaceGetters: verify format, width, height, stride
   - TestImageSurfaceClose: verify close works and double-close is safe
   - TestImageSurfaceStatusAfterClose: verify status after close

4. Wire up to main package:
   - Update cairo/cairo.go to re-export:
     - Format type and constants
     - NewImageSurface function
     - Surface interface

Integration: Ensure the surface can be created and closed without leaks.
Use the status package for error conversion.
Document that Close() must be called or rely on finalizer.
```

---

## Prompt 8: Surface Package - PNG Support

```markdown
Add PNG writing capability to ImageSurface.

Context: PNG export is essential for MVP validation and testing. This completes the basic ImageSurface implementation.

Requirements:
1. Update surface/surface.go to add:
   - func (s *ImageSurface) WriteToPNG(filename string) error
   - Document that surface must be flushed before writing

2. Update surface/surface_cgo.go to add:
   - func (s *ImageSurface) writeToPNG(filename string) error
     - Converts filename to C string
     - Calls cairo_surface_write_to_png
     - Frees C string
     - Converts status to error
     - Returns error if surface is closed

3. Update surface/surface_test.go to add:
   - TestImageSurfaceWriteToPNG: create surface, write PNG, verify file exists
   - TestImageSurfaceWriteToPNGInvalidPath: test error handling for invalid path
   - TestImageSurfaceWriteToPNGAfterClose: verify error after close
   - Add helper to create temp directory for test files
   - Add cleanup to remove test files

4. Update cairo/cairo.go to document:
   - PNG support in package documentation
   - Example usage in package comment

Testing requirements:
- Use t.TempDir() for test file output
- Verify file is created and non-empty
- Test with different surface formats
- Ensure no memory leaks with repeated write operations
```

---

## Prompt 9: Context Package - Creation and Lifecycle

```markdown
Implement the Context type with creation and lifecycle management.

Context (meta): The Context is Cairo's main drawing object. This step creates it but doesn't add drawing operations yet.

Requirements:
1. Create context/context.go with:
   - Context struct with:
     - sync.RWMutex embedded
     - ptr *C.cairo_t (declared in CGO file)
     - closed bool
   - func NewContext(surface Surface) (*Context, error)
   - func (c *Context) Close() error
   - func (c *Context) Status() Status
   - func (c *Context) Save()
   - func (c *Context) Restore()
   - Package documentation explaining Context purpose

2. Create context/context_cgo.go with:
   - CGO preamble
   - func newContext(surface Surface) (*Context, error)
     - Extract C surface pointer from Surface interface
     - Call cairo_create
     - Check status
     - Set up finalizer
     - Return Context or error
   - CGO implementations for Close, Status, Save, Restore
   - Helper to extract C surface pointer from BaseSurface

3. Create context/context_test.go with:
   - TestNewContext: create context from ImageSurface
   - TestNewContextNilSurface: test error handling
   - TestContextClose: verify close and double-close safety
   - TestContextStatus: verify status reporting
   - TestContextSaveRestore: verify save/restore stack
   - TestContextSaveRestoreImbalance: test restore without save

4. Update cairo/cairo.go to re-export:
   - NewContext function
   - Document Context usage pattern

Integration:
- Context must work with ImageSurface from previous step
- Use proper type assertions to get C pointer from Surface interface
- Ensure Context.Close() is independent of Surface.Close()
```

---

## Prompt 10: Context Package - Source Colors

```markdown
Add solid color source setting to Context.

Context: Before drawing, we need to set the source color. This adds the simplest source type - solid colors.

Requirements:
1. Update context/context.go to add:
   - func (c *Context) SetSourceRGB(r, g, b float64)
   - func (c *Context) SetSourceRGBA(r, g, b, a float64)
   - Document that r, g, b, a are in range [0.0, 1.0]
   - Both methods use proper locking

2. Update context/context_cgo.go to add:
   - func (c *Context) setSourceRGB(r, g, b float64)
   - func (c *Context) setSourceRGBA(r, g, b, a float64)
   - Both call appropriate cairo_set_source_* functions

3. Update context/context_test.go to add:
   - TestContextSetSourceRGB: set color and verify no error status
   - TestContextSetSourceRGBA: set color with alpha and verify
   - TestContextSetSourceAfterClose: verify error handling
   - Integration test: create context, set source, check status is success

4. Update cairo/cairo.go documentation to include:
   - Example showing source color setting

Do NOT implement pattern sources yet - just solid colors.
Colors should be validated to be in [0.0, 1.0] range in documentation but not enforced in code (Cairo handles it).
```

---

## Prompt 11: Context Package - Basic Path Operations

```markdown
Add basic path construction operations to Context.

Context: Path operations define what to draw. This adds the minimum operations needed for MVP.

Requirements:
1. Update context/context.go to add:
   - func (c *Context) MoveTo(x, y float64)
   - func (c *Context) LineTo(x, y float64)
   - func (c *Context) Rectangle(x, y, width, height float64)
   - func (c *Context) ClosePath()
   - func (c *Context) NewPath()
   - func (c *Context) NewSubPath()
   - func (c *Context) GetCurrentPoint() (x, y float64, err error)
   - All methods use proper locking
   - Document coordinate system (user-space coordinates)

2. Update context/context_cgo.go to add:
   - CGO implementations for all above methods
   - GetCurrentPoint should check status and return error if no current point

3. Update context/context_test.go to add:
   - TestContextMoveTo: verify MoveTo sets current point
   - TestContextLineTo: verify LineTo adds line and updates current point
   - TestContextRectangle: verify rectangle path is created
   - TestContextClosePath: verify path closing
   - TestContextNewPath: verify path clearing
   - TestContextGetCurrentPoint: verify getting current point
   - TestContextGetCurrentPointNoPoint: verify error when no current point
   - TestContextPathOperationsAfterClose: verify error handling

4. Document in cairo/cairo.go:
   - Path construction basics
   - Current point concept

Integration test idea: MoveTo, LineTo several points, GetCurrentPoint to verify.
Do NOT add arc, curve, or text paths yet - just straight lines and rectangles.
```

---

## Prompt 12: Context Package - Fill and Stroke Operations

```markdown
Add fill and stroke rendering operations to Context.

Context: These operations actually render the path that was constructed. This completes the basic drawing pipeline.

Requirements:
1. Update context/context.go to add:
   - func (c *Context) Fill()
   - func (c *Context) FillPreserve()
   - func (c *Context) Stroke()
   - func (c *Context) StrokePreserve()
   - func (c *Context) Paint()
   - func (c *Context) SetLineWidth(width float64)
   - All methods use proper locking
   - Document that Fill/Stroke consume the path, Preserve variants don't

2. Update context/context_cgo.go to add:
   - CGO implementations for all above methods
   - Each should check if context is closed

3. Update context/context_test.go to add:
   - TestContextFill: verify Fill operation
   - TestContextFillPreserve: verify path is preserved
   - TestContextStroke: verify Stroke operation
   - TestContextStrokePreserve: verify path is preserved
   - TestContextPaint: verify Paint operation
   - TestContextSetLineWidth: verify line width setting
   - TestContextRenderAfterClose: verify error handling

4. Update cairo/cairo.go documentation:
   - Add example of complete drawing operation
   - Document Fill vs Stroke semantics

Integration is key: These operations must work with paths from previous step.
Add a test that does: NewPath, Rectangle, SetSourceRGB, Fill - check status is success.
```

---

## Prompt 13: MVP Integration Test

```markdown
Create a complete end-to-end test demonstrating MVP functionality.

Context: This validates that all pieces work together to accomplish the core goal: draw a rectangle and save to PNG.

Requirements:
1. Create examples/basic_shapes.go with:
   - Example function that creates a 400x400 PNG
   - Draws a filled red rectangle (100, 100, 200, 200)
   - Draws a blue stroked rectangle outline (120, 120, 160, 160)
   - Saves to output.png
   - Includes complete error handling
   - Has proper cleanup with defer

2. Create examples/basic_shapes_test.go with:
   - Test that runs the example
   - Verifies PNG file is created
   - Checks file size is reasonable
   - Uses t.TempDir() for output

3. Update README.md to add:
   - "Quick Start" section
   - Code example showing rectangle drawing
   - Build instructions: go build ./examples/basic_shapes.go
   - Note about Cairo system dependency

4. Update cairo/cairo.go to add:
   - Complete package example showing MVP usage:
     - Create ImageSurface
     - Create Context
     - Set source color
     - Draw rectangle
     - Fill
     - Save PNG
     - Proper cleanup

Success criteria:
- Example compiles and runs without errors
- PNG file is created and viewable
- No memory leaks (can verify with repeated runs)
- All defers properly clean up resources

This is a major milestone - the MVP is complete!
```

---

## Prompt 14: Enhanced Testing and Documentation

```markdown
Add comprehensive testing and improve documentation across all packages.

Context: With MVP working, solidify the foundation with better tests and docs before expanding features.

Requirements:
1. Add benchmarks:
   - context/context_bench_test.go:
     - BenchmarkContextCreation
     - BenchmarkContextPathOperations (many MoveTo/LineTo)
     - BenchmarkContextFillOperations
   - surface/surface_bench_test.go:
     - BenchmarkImageSurfaceCreation
     - BenchmarkWriteToPNG

2. Add example tests (use Example_ pattern):
   - examples/example_test.go:
     - Example_drawRectangle (testable example)
     - Example_fillAndStroke (both operations)
     - Example_colorBlending (multiple shapes)

3. Improve package documentation:
   - cairo/cairo.go: expand package doc with architecture overview
   - context/context.go: add detailed Context lifecycle explanation
   - surface/surface.go: add Surface types overview
   - matrix/matrix.go: add transformation math explanation
   - status/status.go: add error handling guide

4. Add CONTRIBUTING.md:
   - How to build
   - How to run tests
   - Code style guidelines
   - How to add new features

5. Update README.md:
   - Add badges (build status placeholder)
   - Add table of contents
   - Add "Current Status" section showing MVP complete
   - Add "Roadmap" section for future features

No new functionality - just testing and documentation polish.
Run all tests and ensure coverage is >80% for all packages.
```

---

## Prompt 15: Pattern Package - Foundation

```markdown
Implement the pattern package foundation with solid patterns.

Context: Patterns are sources for drawing operations. Starting with solid patterns (which duplicate SetSourceRGB functionality) to establish the pattern infrastructure.

Requirements:
1. Create pattern/pattern.go with:
   - Pattern interface with:
     - Close() error
     - Status() Status
     - SetMatrix(m *matrix.Matrix)
     - GetMatrix() (*matrix.Matrix, error)
   - BasePattern struct with:
     - sync.RWMutex embedded
     - ptr *C.cairo_pattern_t
     - closed bool
   - SolidPattern struct embedding BasePattern
   - func NewSolidPatternRGB(r, g, b float64) (*SolidPattern, error)
   - func NewSolidPatternRGBA(r, g, b, a float64) (*SolidPattern, error)

2. Create pattern/pattern_cgo.go with:
   - CGO preamble
   - BasePattern CGO methods for Close, Status, SetMatrix, GetMatrix
   - SolidPattern CGO constructor functions

3. Create pattern/pattern_test.go with:
   - TestNewSolidPatternRGB: test creation
   - TestNewSolidPatternRGBA: test creation with alpha
   - TestPatternClose: verify close and double-close
   - TestPatternStatus: verify status reporting
   - TestPatternMatrix: verify matrix get/set

4. Update context/context.go to add:
   - func (c *Context) SetSource(pattern Pattern)
   - Updates source to use pattern instead of direct color

5. Update cairo/cairo.go to re-export:
   - Pattern interface
   - NewSolidPattern* functions

Integration test: Create solid pattern, set as Context source, draw and verify.
Document that SetSourceRGB is convenience wrapper around solid patterns.
```

---

## Prompt 16: Context Package - Transformations

```markdown
Add transformation operations to Context.

Context: Transformations modify the coordinate system for subsequent drawing operations.

Requirements:
1. Update context/context.go to add:
   - func (c *Context) Translate(tx, ty float64)
   - func (c *Context) Scale(sx, sy float64)
   - func (c *Context) Rotate(angle float64)
   - func (c *Context) Transform(matrix *matrix.Matrix)
   - func (c *Context) SetMatrix(matrix *matrix.Matrix)
   - func (c *Context) GetMatrix() (*matrix.Matrix, error)
   - func (c *Context) IdentityMatrix()
   - func (c *Context) UserToDevice(x, y float64) (float64, float64)
   - func (c *Context) UserToDeviceDistance(dx, dy float64) (float64, float64)
   - func (c *Context) DeviceToUser(x, y float64) (float64, float64)
   - func (c *Context) DeviceToUserDistance(dx, dy float64) (float64, float64)

2. Update context/context_cgo.go to add:
   - CGO implementations for all transformation methods
   - Proper matrix conversion between Go and C

3. Update context/context_test.go to add:
   - TestContextTranslate: verify translation
   - TestContextScale: verify scaling
   - TestContextRotate: verify rotation
   - TestContextTransform: verify matrix transformation
   - TestContextGetSetMatrix: verify matrix get/set round-trip
   - TestContextIdentityMatrix: verify identity reset
   - TestContextCoordinateConversion: verify user/device conversions
   - TestContextTransformationsCombined: test multiple transformations

4. Create examples/transformations.go:
   - Example showing translate, scale, rotate
   - Draw same shape at different transformations
   - Save to PNG

Document transformation order and current transformation matrix (CTM) concept.
Ensure transformations work correctly with Save/Restore.
```

---

## Prompt 17: Context Package - Advanced Path Operations

```markdown
Add arc and curve operations to Context for complex paths.

Context: Expands path construction beyond straight lines to include curves and arcs.

Requirements:
1. Update context/context.go to add:
   - func (c *Context) Arc(xc, yc, radius, angle1, angle2 float64)
   - func (c *Context) ArcNegative(xc, yc, radius, angle1, angle2 float64)
   - func (c *Context) CurveTo(x1, y1, x2, y2, x3, y3 float64)
   - func (c *Context) RelMoveTo(dx, dy float64)
   - func (c *Context) RelLineTo(dx, dy float64)
   - func (c *Context) RelCurveTo(dx1, dy1, dx2, dy2, dx3, dy3 float64)
   - Document angles are in radians

2. Update context/context_cgo.go to add:
   - CGO implementations for all arc and curve methods

3. Update context/context_test.go to add:
   - TestContextArc: verify arc path creation
   - TestContextArcNegative: verify negative arc
   - TestContextCurveTo: verify Bezier curve
   - TestContextRelativeOperations: verify relative path operations
   - TestContextCircle: use Arc to create complete circle

4. Update examples/basic_shapes.go to add:
   - Circle drawing example
   - Curved path example
   - Save to circles.png

5. Create examples/gradients.go placeholder:
   - Comment explaining gradients come next phase
   - Basic structure ready

Integration: Test arc operations with fill and stroke.
Verify relative operations work correctly with transformations.
```

---

## Prompt 18: Pattern Package - Gradient Patterns

```markdown
Implement gradient patterns (linear and radial).

Context: Gradients provide smooth color transitions and are essential for complex graphics.

Requirements:
1. Update pattern/pattern.go to add:
   - LinearGradient struct embedding BasePattern
   - func NewLinearGradient(x0, y0, x1, y1 float64) (*LinearGradient, error)
   - RadialGradient struct embedding BasePattern
   - func NewRadialGradient(cx0, cy0, radius0, cx1, cy1, radius1 float64) (*RadialGradient, error)
   - func (g *LinearGradient) AddColorStopRGB(offset, r, g, b float64)
   - func (g *LinearGradient) AddColorStopRGBA(offset, r, g, b, a float64)
   - Same for RadialGradient

2. Update pattern/pattern_cgo.go to add:
   - CGO implementations for gradient creation
   - CGO implementations for color stop methods

3. Update pattern/pattern_test.go to add:
   - TestNewLinearGradient: test creation
   - TestNewRadialGradient: test creation
   - TestLinearGradientColorStops: verify adding stops
   - TestRadialGradientColorStops: verify adding stops
   - TestGradientWithContext: integration test using gradient as source

4. Update examples/gradients.go:
   - Implement linear gradient example (left red to right blue)
   - Implement radial gradient example (center white to edge blue)
   - Draw shapes with gradients
   - Save to gradients.png

5. Update cairo/cairo.go to re-export:
   - LinearGradient, RadialGradient types
   - Gradient constructor functions

Document that color stops must be in [0.0, 1.0] for offset.
Test gradients with transformations.
Verify gradient patterns can be reused across multiple drawing operations.
```

---

## Prompt 19: Context Package - Line Styles

```markdown
Add line styling capabilities to Context.

Context: Control stroke appearance with caps, joins, dashes, and miter limits.

Requirements:
1. Create context/line_style.go with:
   - LineCap type as int with constants:
     - LineCapButt, LineCapRound, LineCapSquare
   - LineJoin type as int with constants:
     - LineJoinMiter, LineJoinRound, LineJoinBevel
   - go:generate stringer for both types

2. Update context/context.go to add:
   - func (c *Context) SetLineCap(lineCap LineCap)
   - func (c *Context) GetLineCap() LineCap
   - func (c *Context) SetLineJoin(lineJoin LineJoin)
   - func (c *Context) GetLineJoin() LineJoin
   - func (c *Context) SetDash(dashes []float64, offset float64)
   - func (c *Context) GetDash() (dashes []float64, offset float64, err error)
   - func (c *Context) SetMiterLimit(limit float64)
   - func (c *Context) GetMiterLimit() float64
   - func (c *Context) GetLineWidth() float64

3. Update context/context_cgo.go to add:
   - CGO implementations for all line style methods
   - GetDash needs to query dash count first, allocate slice, then get dashes

4. Update context/context_test.go to add:
   - TestContextLineCap: test setting/getting line caps
   - TestContextLineJoin: test setting/getting line joins
   - TestContextDash: test setting/getting dash patterns
   - TestContextDashEmpty: test empty dash pattern (solid line)
   - TestContextMiterLimit: test miter limit
   - TestContextLineStyleCombinations: test various combinations

5. Create examples/line_styles.go:
   - Draw lines with different caps (butt, round, square)
   - Draw angles with different joins (miter, round, bevel)
   - Draw dashed lines with different patterns
   - Save to line_styles.png

6. Update cairo/cairo.go to re-export:
   - LineCap and LineJoin types and constants

Test that line styles persist across Save/Restore.
Verify dash patterns work with transformations (scale affects dash lengths).
```

---

## Prompt 20: Context Package - Clipping

```markdown
Add clipping operations to Context.

Context: Clipping restricts drawing to specific regions, essential for complex layouts.

Requirements:
1. Update context/context.go to add:
   - func (c *Context) Clip()
   - func (c *Context) ClipPreserve()
   - func (c *Context) ClipExtents() (x1, y1, x2, y2 float64)
   - func (c *Context) InClip(x, y float64) bool
   - func (c *Context) ResetClip()
   - Document that clipping is intersective (clips intersect with previous clips)

2. Update context/context_cgo.go to add:
   - CGO implementations for all clipping methods

3. Update context/context_test.go to add:
   - TestContextClip: verify clipping restricts drawing
   - TestContextClipPreserve: verify path preservation
   - TestContextClipExtents: verify extents calculation
   - TestContextInClip: verify point-in-clip testing
   - TestContextResetClip: verify clip clearing
   - TestContextNestedClips: verify clip intersection behavior
   - TestContextClipWithTransform: verify clipping with transformations

4. Create examples/clipping.go:
   - Example showing circular clip region
   - Draw shapes that get clipped
   - Example showing nested clips
   - Save to clipping.png

5. Update examples/basic_shapes_test.go:
   - Add visual regression test helper function
   - Structure for comparing generated PNGs against references

Document clip region interaction with Save/Restore (clips are saved/restored).
Test that clip regions work correctly with transformations.
```

---

## Prompt 21: Surface Package - Surface Pattern

```markdown
Implement surface patterns for texture mapping and pattern fills.

Context: Surface patterns allow using images as sources for drawing operations.

Requirements:
1. Update pattern/pattern.go to add:
   - SurfacePattern struct embedding BasePattern
   - func NewSurfacePattern(surface Surface) (*SurfacePattern, error)
   - Extend type as int with constants:
     - ExtendNone, ExtendRepeat, ExtendReflect, ExtendPad
   - Filter type as int with constants:
     - FilterFast, FilterGood, FilterBest, FilterNearest, FilterBilinear
   - func (p *SurfacePattern) SetExtend(extend Extend)
   - func (p *SurfacePattern) GetExtend() Extend
   - func (p *SurfacePattern) SetFilter(filter Filter)
   - func (p *SurfacePattern) GetFilter() Filter

2. Update pattern/pattern_cgo.go to add:
   - CGO implementation for surface pattern creation
   - CGO implementations for extend and filter methods

3. Update pattern/pattern_test.go to add:
   - TestNewSurfacePattern: test creation from ImageSurface
   - TestSurfacePatternExtend: test extend modes
   - TestSurfacePatternFilter: test filter modes
   - TestSurfacePatternWithContext: integration test drawing with surface pattern

4. Create examples/patterns.go:
   - Create small ImageSurface with checker pattern
   - Use as surface pattern with ExtendRepeat
   - Fill large rectangle with pattern
   - Test different extend modes
   - Save to patterns.png

5. Update cairo/cairo.go to re-export:
   - SurfacePattern type
   - Extend and Filter types and constants
   - NewSurfacePattern function

Test surface patterns with transformations (pattern matrix).
Verify pattern works with different surface formats.
Document that source surface must remain valid while pattern is in use.
```

---

## Prompt 22: Context Package - Operators and Compositing

```markdown
Add compositing operators to Context for advanced blending.

Context: Operators control how drawing operations combine with existing content.

Requirements:
1. Create context/operator.go with:
   - Operator type as int
   - Constants for all cairo_operator_t values:
     - OperatorClear, OperatorSource, OperatorOver (default)
     - OperatorIn, OperatorOut, OperatorAtop
     - OperatorDest, OperatorDestOver, OperatorDestIn, OperatorDestOut, OperatorDestAtop
     - OperatorXor, OperatorAdd, OperatorSaturate
     - OperatorMultiply, OperatorScreen, OperatorOverlay
     - OperatorDarken, OperatorLighten
     - OperatorColorDodge, OperatorColorBurn
     - OperatorHardLight, OperatorSoftLight
     - OperatorDifference, OperatorExclusion
     - OperatorHslHue, OperatorHslSaturation, OperatorHslColor, OperatorHslLuminosity
   - go:generate stringer -type=Operator

2. Update context/context.go to add:
   - func (c *Context) SetOperator(op Operator)
   - func (c *Context) GetOperator() Operator
   - Document operator effects and common use cases

3. Update context/context_cgo.go to add:
   - CGO implementations for operator methods

4. Update context/context_test.go to add:
   - TestContextSetOperator: test setting operators
   - TestContextGetOperator: test getting current operator
   - TestContextOperatorDefault: verify default is OperatorOver

5. Create examples/compositing.go:
   - Draw overlapping shapes with different operators
   - Show OperatorOver, OperatorAdd, OperatorMultiply, OperatorXor
   - Create visual comparison grid
   - Save to compositing.png

6. Update cairo/cairo.go to re-export:
   - Operator type and all constants

Document that OperatorOver is the default and most common.
Explain porter-duff operators vs blend modes.
Test operators with alpha transparency.
```

---

## Prompt 23: Context Package - Fill and Stroke Rules

```markdown
Add fill rule control and stroke/fill extents calculation.

Context: Fill rules determine what's "inside" complex paths. Extents provide bounding boxes.

Requirements:
1. Create context/fill_rule.go with:
   - FillRule type as int
   - Constants:
     - FillRuleWinding (default)
     - FillRuleEvenOdd
   - go:generate stringer -type=FillRule

2. Update context/context.go to add:
   - func (c *Context) SetFillRule(fillRule FillRule)
   - func (c *Context) GetFillRule() FillRule
   - func (c *Context) FillExtents() (x1, y1, x2, y2 float64)
   - func (c *Context) StrokeExtents() (x1, y1, x2, y2 float64)
   - func (c *Context) PathExtents() (x1, y1, x2, y2 float64)
   - func (c *Context) InFill(x, y float64) bool
   - func (c *Context) InStroke(x, y float64) bool
   - Document fill rule differences (winding vs even-odd)

3. Update context/context_cgo.go to add:
   - CGO implementations for all methods

4. Update context/context_test.go to add:
   - TestContextFillRule: test setting/getting fill rule
   - TestContextFillExtents: verify extents calculation
   - TestContextStrokeExtents: verify stroke extents
   - TestContextPathExtents: verify path extents
   - TestContextInFill: test point-in-fill detection
   - TestContextInStroke: test point-in-stroke detection
   - TestContextFillRuleWindingVsEvenOdd: compare results

5. Create examples/fill_rules.go:
   - Draw self-intersecting star with winding rule
   - Draw same star with even-odd rule
   - Show visual difference
   - Save to fill_rules.png

6. Update cairo/cairo.go to re-export:
   - FillRule type and constants

Test fill rules with complex self-intersecting paths.
Verify extents calculations are accurate with transformations.
```

---

## Prompt 24: Font Package - Toy Font API

```markdown
Implement the toy font API for basic text rendering.

Context: The toy font API provides simple text rendering. It's limited but sufficient for basic needs.

Requirements:
1. Create font/font.go with:
   - Slant type as int with constants:
     - SlantNormal, SlantItalic, SlantOblique
   - Weight type as int with constants:
     - WeightNormal, WeightBold
   - go:generate stringer for both types
   - Package documentation explaining toy vs scaled font APIs

2. Update context/context.go to add:
   - func (c *Context) SelectFontFace(family string, slant font.Slant, weight font.Weight)
   - func (c *Context) SetFontSize(size float64)
   - func (c *Context) ShowText(text string)
   - func (c *Context) TextPath(text string)
   - Document that text is positioned at current point

3. Update context/context_cgo.go to add:
   - CGO implementations for font methods
   - String conversion for family parameter

4. Update context/context_test.go to add:
   - TestContextSelectFontFace: test font selection
   - TestContextSetFontSize: test font size setting
   - TestContextShowText: test text rendering
   - TestContextTextPath: test text path creation

5. Create examples/text.go:
   - Draw text with different fonts
   - Draw text with different sizes
   - Draw text with different slants and weights
   - Position text at different locations
   - Save to text.png

6. Update cairo/cairo.go to re-export:
   - Font Slant and Weight types and constants

Document that toy font API is platform-dependent.
Note that advanced text rendering requires Pango integration (future).
Test text rendering with different transformations.
```

---

## Prompt 25: Font Package - Text Extents

```markdown
Add text measurement capabilities for layout.

Context: Text extents allow measuring text for proper positioning and alignment.

Requirements:
1. Create font/extents.go with:
   - TextExtents struct with fields:
     - XBearing, YBearing float64
     - Width, Height float64
     - XAdvance, YAdvance float64
   - FontExtents struct with fields:
     - Ascent, Descent float64
     - Height float64
     - MaxXAdvance, MaxYAdvance float64
   - Document what each field represents

2. Update context/context.go to add:
   - func (c *Context) TextExtents(text string) (*font.TextExtents, error)
   - func (c *Context) FontExtents() (*font.FontExtents, error)
   - Document use cases for each type of extents

3. Update context/context_cgo.go to add:
   - CGO implementations for extents methods
   - Convert C extents structs to Go structs

4. Update context/context_test.go to add:
   - TestContextTextExtents: test text measurement
   - TestContextFontExtents: test font metrics
   - TestContextTextExtentsEmpty: test empty string
   - TestContextExtentsWithDifferentFonts: verify different fonts give different metrics

5. Update examples/text.go to add:
   - Text alignment example (center, right-align using extents)
   - Multi-line text with proper spacing using font extents
   - Draw bounding boxes around text using text extents
   - Save to text_extents.png

Test extents with transformations (scaling affects measurements).
Verify extents are consistent with actual rendered text.
Document coordinate system for bearings and advances.
```

---

## Prompt 26: Surface Package - PDF Surface

```markdown
Implement PDF surface for vector output.

Context: PDF surfaces write to PDF files instead of raster images, enabling scalable output.

Requirements:
1. Create surface/pdf.go with:
   - PDFSurface struct embedding BaseSurface
   - func NewPDFSurface(filename string, widthPt, heightPt float64) (*PDFSurface, error)
   - func (s *PDFSurface) SetSize(widthPt, heightPt float64)
   - func (s *PDFSurface) ShowPage()
   - Document that dimensions are in points (1/72 inch)

2. Create surface/pdf_cgo.go with:
   - Build tag: // +build !nopdf
   - CGO preamble: #cgo pkg-config: cairo-pdf
   - CGO implementations for PDF surface methods

3. Update surface/surface_test.go to add:
   - TestNewPDFSurface: test creation
   - TestPDFSurfaceSetSize: test size changes
   - TestPDFSurfaceMultiPage: test multi-page document
   - Use t.TempDir() for test PDF files

4. Create examples/pdf_output.go:
   - Create multi-page PDF document
   - Draw different content on each page
   - Include text, shapes, and gradients
   - Save to output.pdf

5. Update cairo/cairo.go to re-export:
   - PDFSurface type
   - NewPDFSurface function

6. Update README.md:
   - Add PDF surface to features list
   - Note Cairo PDF backend requirement

Document PDF coordinate system (origin at top-left).
Test that PDF surface works with all drawing operations.
Verify file is valid PDF (can open in reader).
```

---

## Prompt 27: Surface Package - SVG Surface

```markdown
Implement SVG surface for web-compatible vector output.

Context: SVG surfaces generate scalable vector graphics files, ideal for web use.

Requirements:
1. Create surface/svg.go with:
   - SVGSurface struct embedding BaseSurface
   - func NewSVGSurface(filename string, widthPt, heightPt float64) (*SVGSurface, error)
   - func (s *SVGSurface) SetDocumentUnit(unit SVGUnit)
   - SVGUnit type as int with constants:
     - SVGUnitUser, SVGUnitEm, SVGUnitEx
     - SVGUnitPx, SVGUnitIn, SVGUnitCm, SVGUnitMm
     - SVGUnitPt, SVGUnitPc, SVGUnitPercent
   - go:generate stringer -type=SVGUnit

2. Create surface/svg_cgo.go with:
   - Build tag: // +build !nosvg
   - CGO preamble: #cgo pkg-config: cairo-svg
   - CGO implementations for SVG surface methods

3. Update surface/surface_test.go to add:
   - TestNewSVGSurface: test creation
   - TestSVGSurfaceDocumentUnit: test unit setting
   - Use t.TempDir() for test SVG files
   - Verify SVG file starts with proper XML header

4. Create examples/svg_output.go:
   - Create SVG with various shapes
   - Include gradients and patterns
   - Test text rendering in SVG
   - Save to output.svg

5. Update cairo/cairo.go to re-export:
   - SVGSurface type
   - SVGUnit type and constants
   - NewSVGSurface function

6. Update README.md:
   - Add SVG surface to features list
   - Note Cairo SVG backend requirement

Document SVG coordinate system and units.
Test that SVG file is valid XML (basic parsing check).
Verify SVG renders correctly in browsers.
```

---

## Prompt 28: Advanced Examples and Documentation

```markdown
Create comprehensive examples and finalize documentation.

Context: With most features implemented, create polished examples and documentation showing real-world usage.

Requirements:
1. Create examples/dashboard.go:
   - Complex example creating a data dashboard
   - Multiple charts: bar chart, line chart, pie chart
   - Use gradients, patterns, text, transformations
   - Generate to both PNG and PDF
   - Well-commented explaining each section

2. Create examples/animation.go:
   - Generate sequence of PNG frames
   - Simple animation (rotating shape, moving object)
   - Document frame rate considerations
   - Show how to batch-process frames

3. Update examples/README.md:
   - Describe each example
   - Show output images (reference outputs)
   - Explain what each demonstrates
   - Build and run instructions

4. Update main README.md:
   - Add comprehensive feature list
   - Add performance notes
   - Add comparison to other Go graphics libraries
   - Add troubleshooting section
   - Add FAQ section

5. Create ARCHITECTURE.md:
   - Explain package organization
   - Explain CGO boundary design
   - Explain memory management strategy
   - Explain thread safety approach
   - Include diagrams if helpful

6. Create docs/ directory with:
   - QUICKSTART.md: fast intro for new users
   - MIGRATION.md: guide for users coming from C Cairo
   - PERFORMANCE.md: performance tips and benchmarks

All examples should be production-quality with error handling.
Add godoc examples for key functions.
Ensure all examples in documentation actually compile.
```

---

## Prompt 29: Error Handling Improvements

```markdown
Enhance error handling with wrapped errors and detailed context.

Context: Improve error messages to make debugging easier while maintaining API compatibility.

Requirements:
1. Create cairo/errors.go with:
   - Custom error types:
     - SurfaceError, ContextError, PatternError
   - Each wraps Status with additional context
   - func (e *SurfaceError) Unwrap() error
   - Implement Is() and As() for errors.Is/As support

2. Update status/status.go to add:
   - More descriptive error messages in Error() method
   - Include suggestions for common errors

3. Update all packages to use wrapped errors:
   - surface: wrap with SurfaceError including surface type
   - context: wrap with ContextError including operation
   - pattern: wrap with PatternError including pattern type

4. Update tests across packages:
   - TestErrorUnwrapping: verify errors.Is works
   - TestErrorContext: verify error messages include context
   - TestErrorTypes: verify type assertions work

5. Update examples to show:
   - Proper error handling patterns
   - Using errors.Is for error checking
   - Error message inspection

6. Update documentation:
   - Add error handling guide to README
   - Document common errors and solutions
   - Show idiomatic error checking patterns

All error improvements must maintain backward compatibility.
Existing error checking code should continue to work.
```

---

## Prompt 30: Thread Safety Validation and Race Testing

```markdown
Comprehensive thread safety testing and validation.

Context: Verify thread safety implementation is correct and performs well under concurrent access.

Requirements:
1. Create test/race_test.go with:
   - TestConcurrentContextOperations: many goroutines drawing to same context
   - TestConcurrentSurfaceOperations: concurrent surface creation/destruction
   - TestConcurrentPatternOperations: concurrent pattern usage
   - TestContextSharedAcrossSurfaces: one context used with multiple surfaces
   - TestPatternSharedAcrossContexts: one pattern used by multiple contexts
   - All tests run with -race flag

2. Create benchmarks/concurrent_bench_test.go:
   - BenchmarkConcurrentDrawing: measure parallel drawing performance
   - BenchmarkConcurrentSurfaceCreation: measure parallel surface creation
   - BenchmarkLockContention: measure lock contention overhead
   - Compare single-threaded vs multi-threaded performance

3. Update documentation:
   - Add concurrency section to README
   - Document safe concurrent usage patterns
   - Document unsafe patterns to avoid
   - Note performance implications of locking

4. Add stress test:
   - test/stress_test.go with build tag: // +build stress
   - Long-running tests with many goroutines
   - Memory leak detection under concurrent load
   - Run with different GOMAXPROCS values

5. Update CI configuration:
   - Run tests with -race flag
   - Run benchmarks to detect performance regressions
   - Add stress test job (separate, optional)

Document recommended concurrency patterns.
Verify no deadlocks or race conditions.
Measure performance impact of thread safety.
```

---

## Prompt 31: Memory Management Validation

```markdown
Validate memory management and add leak detection testing.

Context: Ensure proper cleanup and no memory leaks in CGO boundary.

Requirements:
1. Create test/memory_test.go with:
   - TestNoLeaksSimpleDrawing: repeated draw cycles, check memory doesn't grow
   - TestNoLeaksWithoutClose: verify finalizers prevent leaks
   - TestNoLeaksPatternReuse: patterns created and released repeatedly
   - TestNoLeaksSurfaceReuse: surfaces created and released repeatedly
   - TestMemoryWithDeferredClose: verify defer Close() pattern works
   - Use runtime.GC() and runtime.ReadMemStats to detect leaks

2. Create benchmarks/memory_bench_test.go:
   - BenchmarkMemoryAllocation: measure allocation overhead
   - BenchmarkFinalizerOverhead: measure finalizer cost
   - BenchmarkExplicitClose: compare explicit close vs finalizer

3. Add memory profiling example:
   - examples/profile_memory.go
   - Show how to profile memory usage
   - Include pprof integration example

4. Update documentation:
   - Add memory management guide
   - Document Close() vs finalizer tradeoffs
   - Show best practices for long-running programs
   - Add troubleshooting section for memory issues

5. Create tools/leak_detector.go:
   - Helper tool to detect resource leaks
   - Tracks creation/destruction of Cairo objects
   - Debugging mode that logs all allocations

Run tests with increased iterations to detect slow leaks.
Test with GODEBUG=cgocheck=2 for extra validation.
Document memory usage patterns and expectations.
```

---

## Prompt 32: Build System and Distribution

```markdown
Finalize build system, distribution, and installation.

Context: Make the library easy to install and use in different environments.

Requirements:
1. Update go.mod:
   - Ensure all dependencies are properly versioned
   - Add go directive for minimum Go version
   - Clean up any unused dependencies

2. Create install scripts:
   - scripts/install-deps-linux.sh: install Cairo on various Linux distros
   - scripts/install-deps-macos.sh: install Cairo on macOS
   - Make scripts idempotent and well-documented

3. Create build configuration:
   - .golangci.yml: comprehensive linter configuration
   - Makefile: alternative to Taskfile with common targets
   - Support for different Cairo configurations (minimal vs full)

4. Update Taskfile.yaml:
   - Add coverage task: generate and view coverage report
   - Add bench task: run all benchmarks
   - Add docs task: generate and serve godoc locally
   - Add install task: install necessary tools

5. Create GitHub templates:
   - .github/ISSUE_TEMPLATE/bug_report.md
   - .github/ISSUE_TEMPLATE/feature_request.md
   - .github/PULL_REQUEST_TEMPLATE.md

6. Update CI/CD:
   - Add code coverage reporting
   - Add automated release workflow
   - Add dependency update automation (Dependabot)
   - Add automatic documentation deployment

7. Create RELEASE.md:
   - Document release process
   - Version numbering scheme (semantic versioning)
   - Changelog generation process

Document installation on different platforms.
Test installation on clean systems.
Verify pkg-config detection works correctly.
```

---

## Prompt 33: Performance Optimization Pass

```markdown
Optimize performance and reduce overhead.

Context: Review implementation for optimization opportunities while maintaining correctness.

Requirements:
1. Profile and optimize:
   - Run CPU profiling on examples
   - Identify hot paths in CGO transitions
   - Optimize frequent operations (LineTo, MoveTo, etc.)
   - Reduce allocations in tight loops

2. Add fast paths:
   - context/fast.go with optimized common operations
   - Batch operations where possible (multiple LineTo calls)
   - Consider unsafe optimizations for high-frequency operations
   - Document performance characteristics

3. Update benchmarks:
   - Add micro-benchmarks for individual operations
   - Add macro-benchmarks for realistic workloads
   - Compare against C Cairo directly
   - Set performance budgets (max overhead %)

4. Optimize memory usage:
   - Reduce allocations in hot paths
   - Pool frequently allocated objects if beneficial
   - Review slice and string conversions

5. Add performance documentation:
   - PERFORMANCE.md with optimization guide
   - Document when to use batch operations
   - Show performance comparison tables
   - List known performance gotchas

6. Create profiling examples:
   - examples/profile_cpu.go: CPU profiling example
   - examples/benchmark_drawing.go: realistic benchmark
   - Document how to profile user code

Measure before and after each optimization.
Ensure optimizations don't break thread safety.
Document any unsafe optimizations clearly.
```

---

## Prompt 34: Final Polish and Release Preparation

```markdown
Final polish, documentation review, and prepare for v0.1.0 release.

Context: Complete the MVP and prepare for initial public release.

Requirements:
1. Documentation audit:
   - Review all package documentation
   - Ensure all public functions have examples
   - Fix any typos or unclear explanations
   - Verify all links work
   - Ensure consistent terminology

2. API review:
   - Review naming consistency across packages
   - Verify Go idioms are followed
   - Check for any awkward APIs
   - Ensure error handling is consistent
   - Verify all exported functions are necessary

3. Test coverage review:
   - Ensure >80% coverage across all packages
   - Add tests for edge cases
   - Add integration tests for common workflows
   - Verify all examples work and are tested

4. Create migration aids:
   - API_REFERENCE.md: complete API listing
   - EXAMPLES_INDEX.md: all examples categorized
   - COMPARISON.md: compare to other Go graphics libraries

5. Prepare release:
   - Create CHANGELOG.md for v0.1.0
   - Tag appropriate features as v0.1.0 scope
   - Create GitHub release with:
     - Release notes
     - Breaking changes (none for v0.1.0)
     - Known limitations
     - Installation instructions

6. Create announcement materials:
   - README badge for version
   - Social media announcement text
   - Blog post outline (if applicable)
   - Hacker News submission template

7. Final testing:
   - Clean checkout test: clone and build from scratch
   - Fresh system test: test on clean Linux and macOS
   - Documentation test: follow README from scratch
   - Example test: run all examples and verify output

All documentation must be polished and professional.
All tests must pass on supported platforms.
Version v0.1.0 is MVP-complete and production-ready for basic use cases.
```

---

## Usage Instructions

### For Code-Generation LLM

1. Start with Prompt 1 and work sequentially through all 34 prompts
2. Each prompt assumes context from all previous prompts
3. Implement test-driven: write tests first or alongside implementation
4. Ensure each step compiles and all tests pass before proceeding
5. Integrate all code immediately - no orphaned files or functions

### Key Checkpoints

- **After Prompt 13**: Run the MVP example and verify PNG output
- **After Prompt 18**: Test all gradient types
- **After Prompt 27**: Test all surface types (Image, PDF, SVG)
- **After Prompt 34**: Full test suite passes, ready for release

### Testing Strategy

- Run `go test -race ./...` after each prompt
- Run benchmarks periodically to catch performance regressions
- Use `go test -cover ./...` to monitor coverage
- Test on both Linux and macOS throughout

### Build Requirements

- Go 1.23 or later
- Cairo 1.18 or later (system library)
- pkg-config
- Standard Go tooling (gofmt, golangci-lint, stringer)

---

## Project Success Criteria

1. **Functional**: Can reproduce Cairo's basic examples
2. **Performant**: Minimal overhead over C Cairo (<10%)
3. **Safe**: Thread-safe, no memory leaks, no data races
4. **Idiomatic**: Feels natural to Go developers
5. **Compatible**: Works on Linux and macOS with standard Cairo
6. **Documented**: Clear examples and API documentation
7. **Tested**: >80% coverage, visual regression tests pass

---

## Appendix: Design Decisions

### Why CGO Direct Calls

- Minimum overhead for performance-critical operations
- Simplest maintenance burden
- Direct API parity with C Cairo

### Why Embedded RWMutex

- Thread safety by default
- Familiar Go pattern
- Low overhead for read-heavy workloads

### Why Finalizers + Close()

- Safety net for forgetting Close()
- Explicit Close() for deterministic cleanup
- Best of both worlds approach

### Why Phased Implementation

- Early validation of architecture (MVP at Prompt 13)
- Reduces risk of major refactoring
- Each phase delivers usable functionality

---

*This prompt series was designed to build a production-ready Go wrapper for Cairo through incremental, test-driven development. Each step is sized for safe implementation while maintaining forward progress toward a complete graphics library.*

