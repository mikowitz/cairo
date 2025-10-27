# Go-Cairo Implementation Checklist

A comprehensive checklist for implementing a Go wrapper around the Cairo
graphics library through 34 incremental prompts.

## Project Overview

- **Total Prompts**: 34
- **MVP Milestone**: Prompt 13
- **Target**: v0.1.0 release-ready library
- **Coverage Goal**: >80% test coverage
- **Platforms**: Linux and macOS

---

## Phase 1: Foundation Setup (Prompts 1-3)

### Prompt 1: Project Foundation

- [x] Initialize Go module at `github.com/username/cairo`
- [x] Create `go.mod` requiring Go 1.23
- [x] Create `README.md` explaining CGO wrapper for Cairo
- [x] Create `Taskfile.yaml` with tasks:
  - [x] `test`: run `go test -race ./...`
  - [x] `generate`: run `go generate ./...`
  - [x] `lint`: run `golangci-lint run`
- [x] Create `.gitignore` for Go projects
- [x] Create `DESIGN.md` linking to Cairo C API docs

### Prompt 2: Status Package - Types and Constants

- [x] Create `status/status.go`:
  - [x] `Status` type as int
  - [x] Constants for all `cairo_status_t` values (minimum 23 constants)
  - [x] `//go:generate stringer -type=Status`
  - [x] `func (s Status) Error() string` method
  - [x] `func (s Status) ToError() error` method
  - [x] Package documentation
- [x] Create `status/status_test.go`:
  - [x] `TestStatusError`: verify Error() returns non-empty strings
  - [x] `TestStatusToError`: verify nil for StatusSuccess, error for others
  - [x] `TestStatusSuccess`: verify StatusSuccess equals 0
- [x] Verify complete test coverage
- [x] Verify code compiles without warnings

### Prompt 3: Status Package - CGO Integration

- [x] Create `status/status_cgo.go`:
  - [x] CGO preamble with `#include <cairo.h>`
  - [x] `#cgo pkg-config: cairo`
  - [x] `import "C"`
  - [x] `func statusFromC(cStatus C.cairo_status_t) Status`
  - [x] `func (s Status) toC() C.cairo_status_t`
- [x] Update `status/status_test.go`:
  - [x] `TestCGOStatusConversion`: verify round-trip conversion
  - [x] Test StatusSuccess, StatusNoMemory, StatusInvalidRestore
  - [x] Use `C.cairo_status_to_string()` to verify mapping
- [x] Add build tags if needed for platform-specific CGO flags
- [x] Verify tests pass on Linux and macOS
- [x] Run `go generate ./...` to generate stringer code

---

## Phase 2: Core Types (Prompts 4-8)

### Prompt 4: Matrix Package - Structure and Basic Operations

- [x] Create `matrix/matrix.go`:
  - [x] `Matrix` struct with `sync.RWMutex` and float64 fields
    (XX, YX, XY, YY, X0, Y0)
  - [x] `func NewMatrix(xx, yx, xy, yy, x0, y0 float64) *Matrix` (revised API)
  - [x] `func NewIdentityMatrix() *Matrix` (revised API)
  - [x] `func (m *Matrix) String() string`
  - [x] Package documentation explaining affine transformations
- [x] Create `matrix/matrix_cgo.go`:
  - [x] CGO preamble
  - [x] `func (m *Matrix) toC() *C.cairo_matrix_t`
  - [x] `func matrixFromC(cm *C.cairo_matrix_t) *Matrix`
- [x] Create `matrix/matrix_test.go`:
  - [x] `TestNewMatrix`: verify matrix with given values
  - [x] `TestNewIdentityMatrix`: verify identity matrix
  - [x] `TestMatrixThreadSafety`: concurrent reads/writes (test exists but skipped pending transformation methods)

Note: Unlike C Cairo, we do not provide Init/InitIdentity methods. The idiomatic
Go approach is to create a new matrix with NewMatrix() or NewIdentityMatrix()
rather than mutating an existing matrix in place.

### Prompt 5: Matrix Package - Transformations ✅ COMPLETE

- [x] Update `matrix/matrix.go`:
  - [x] `func (m *Matrix) Multiply(other *Matrix)` - returns new matrix
  - [x] `func (m *Matrix) TransformPoint(x, y float64) (float64, float64)`
  - [x] `func (m *Matrix) TransformDistance(dx, dy float64) (float64, float64)`
  - [x] `func (m *Matrix) Translate(tx, ty float64)`
  - [x] `func (m *Matrix) Scale(sx, sy float64)`
  - [x] `func (m *Matrix) Rotate(radians float64)`
  - [x] `func (m *Matrix) Invert() error`
  - [x] All methods use proper locking
  - [x] BONUS: `func NewTranslationMatrix(tx, ty float64) *Matrix`
  - [x] BONUS: `func NewScalingMatrix(sx, sy float64) *Matrix`
  - [x] BONUS: `func NewRotationMatrix(radians float64) *Matrix`
- [x] Update `matrix/matrix_cgo.go`:
  - [x] CGO wrappers for cairo_matrix_* functions
  - [x] `func matrixInvert(m *Matrix) error` - uses status package
  - [x] `func matrixInitTranslate(tx, ty float64) *Matrix`
  - [x] `func matrixInitScale(sx, sy float64) *Matrix`
  - [x] `func matrixInitRotate(radians float64) *Matrix`
  - [x] `func (m *Matrix) updateFromC()` helper method
- [x] Update `matrix/matrix_test.go`:
  - [x] `TestMatrixMultiply` - identity and simple multiplication
  - [x] `TestMatrixTransformPoint` - identity, translation, scaling
  - [x] `TestMatrixTransformDistance` - verifies translation has no effect
  - [x] `TestMatrixTranslate` - simple, negative, zero translation
  - [x] `TestMatrixScale` - uniform, non-uniform, fractional scaling
  - [x] `TestMatrixRotate` - 0°, 90°, 180° rotations
  - [x] `TestMatrixInvert` - identity, scaling, singular matrix error
  - [x] `TestMatrixOperationsCombined` - complex transformation chains
  - [x] BONUS: `TestNewTranslationMatrix` - 4 test cases
  - [x] BONUS: `TestNewScalingMatrix` - 5 test cases
  - [x] BONUS: `TestNewRotationMatrix` - 4 test cases with 45° included
- [x] All operations match Cairo's semantics exactly
- [x] Uses status package for error handling in Invert()

### Prompt 6: Surface Package - Interface and Base Types ✅ COMPLETE

- [x] Create `surface/surface.go`:
  - [x] `Surface` interface with methods: Close(), Status(), Flush(), MarkDirty(),
    MarkDirtyRectangle()
  - [x] `BaseSurface` struct with `sync.RWMutex`, ptr, closed flag
  - [x] `newBaseSurface()` helper function
  - [x] Finalizer setup with `runtime.SetFinalizer` in `newBaseSurface()` (line 29)
  - [x] Comprehensive documentation on all public methods explaining Cairo semantics
  - [x] Package documentation in `surface/doc.go`
- [x] Create `surface/format.go`:
  - [x] `Format` type as int
  - [x] Constants: FormatInvalid, FormatARGB32, FormatRGB24, FormatA8, FormatA1,
    FormatRGB16_565, FormatRGB30
  - [x] `//go:generate stringer -type=Format`
  - [x] `func (f Format) StrideForWidth(width int) int`
- [x] Create `surface/surface_cgo.go`:
  - [x] CGO preamble with `#cgo pkg-config: cairo`
  - [x] BaseSurface CGO methods for all interface methods
  - [x] All methods check closed flag (check for nil ptr)
- [x] Create `surface/surface_test.go`:
  - [x] Tests for Format constants (PASSING)
  - [x] `TestFormatStrideForWidth` (PASSING)
  - [x] Comprehensive test structure with 19 tests (ALL PASSING)
  - [x] Thread safety tests with race detector
  - [x] Finalizer tests
- [x] Run `go generate ./...` to generate stringer code

✅ **Status: 100% COMPLETE - All requirements met, all tests passing**

### Prompt 7: Surface Package - ImageSurface Creation ✅ COMPLETE

- [x] Create `surface/image_surface.go`:
  - [x] `ImageSurface` struct embedding `*BaseSurface` (uses pointer embedding)
  - [x] `func NewImageSurface(format Format, width, height int) (*ImageSurface, error)` ✅
    - [x] Returns error type as required
    - [x] Checks status after creation and returns error if not Success (line 13-17)
    - [x] Calls `newBaseSurface()` which sets up finalizer automatically (line 19)
  - [x] `func (s *ImageSurface) GetFormat() Format`
  - [x] `func (s *ImageSurface) GetWidth() int`
  - [x] `func (s *ImageSurface) GetHeight() int`
  - [x] `func (s *ImageSurface) GetStride() int`
  - [x] All methods use proper locking (RLock/RUnlock)
- [x] Update `surface/surface_cgo.go`:
  - [x] `func imageSurfaceCreate(format Format, width, height int) SurfacePtr`
  - [x] Calls `cairo_image_surface_create`
  - [x] Finalizer setup via `newBaseSurface()` helper (inherited)
  - [x] CGO implementations for surface operations (inherited from BaseSurface)
- [x] Update `surface/surface_test.go`:
  - [x] `TestNewImageSurface` (PASSING)
  - [x] `TestNewImageSurfaceInvalidFormat` (part of TestImageSurfaceInvalidParameters - PASSING)
  - [x] `TestNewImageSurfaceInvalidSize` (part of TestImageSurfaceInvalidParameters - PASSING)
  - [x] `TestImageSurfaceGetters` (PASSING)
  - [x] Close behavior tested via `TestBaseSurfaceClose` (PASSING)
  - [x] Status behavior tested via `TestBaseSurfaceStatus` (PASSING)
- [x] Create `cairo.go` at package root:
  - [x] Basic file exists with NewImageSurface wrapper
  - [x] Re-export Format type and constants (cairo.go:72-102)
  - [x] Re-export Surface interface (cairo.go:104-112)
  - [x] Add comprehensive package documentation (cairo.go:1-67)
    - [x] Surface Types section
    - [x] Resource Management section with defer pattern example
    - [x] Thread Safety section
    - [x] Basic Usage Example section
    - [x] Error Handling section
  - [x] Add usage example showing surface lifecycle (cairo.go:34-56)
  - [x] Document Close() requirement and finalizer behavior (cairo.go:14-25, 110-111)
- [x] Additional test coverage in `cairo_test.go`:
  - [x] `TestSurfaceInterfaceReexport` - verifies Surface interface works
  - [x] `TestSurfaceLifecycle` - demonstrates proper resource management

✅ **Status: 100% COMPLETE - All requirements met, all tests passing, comprehensive documentation added**

### Prompt 8: Surface Package - PNG Support

- [x] Update `surface/surface.go`:
  - [x] `func (s *ImageSurface) WriteToPNG(filename string) error` (implemented on BaseSurface - line 97)
  - [x] Document that surface must be flushed before writing (comprehensive docs at line 97-129)
- [x] Update `surface/surface_cgo.go`:
  - [x] `func surfaceWriteToPNG(ptr SurfacePtr, filepath string) error` (line 50-60)
  - [x] Converts filename to C string
  - [x] Calls `cairo_surface_write_to_png`
  - [x] Frees C string
  - [x] Returns error if surface is closed (via BaseSurface.WriteToPNG checking ptr == nil)
- [x] Update `surface/surface_test.go`:
  - [x] `TestImageSurfaceWriteToPNG` (line 416)
  - [x] `TestImageSurfaceWriteToPNGInvalidPath` (line 470)
  - [x] `TestImageSurfaceWriteToPNGAfterClose` (line 580)
  - [x] Use `t.TempDir()` for test files
  - [x] Verify file exists and is non-empty
  - [x] Test with different surface formats (TestImageSurfaceWriteToPNGWithDifferentFormats - line 625)
  - [x] BONUS: `TestImageSurfaceWriteToPNGInvalidFilename` - tests empty string, null bytes, long filenames (line 487)
  - [x] BONUS: `TestImageSurfaceWriteToPNGNullByteHandling` - documents C.CString truncation behavior (line 534)
  - [x] BONUS: `TestImageSurfaceWriteToPNGMultipleTimes` - tests writing same surface multiple times (line 600)
- [x] Update `doc.go`:
  - [x] Document PNG support (PNG Export section at line 34-56)
  - [x] Add example usage in package comment (Basic Usage Example updated at line 58-79)

✅ **Status: 100% COMPLETE - All requirements met, all tests passing, comprehensive documentation added**

---

## Phase 3: Context & Drawing (Prompts 9-12)

### Prompt 9: Context Package - Creation and Lifecycle ✅ COMPLETE

- [x] Create `context/context.go`:
  - [x] `Context` struct with `sync.RWMutex`, ptr, closed flag (line 70-74)
  - [x] `func NewContext(surface Surface) (*Context, error)` (line 76-93)
  - [x] `func (c *Context) Close() error` (line 108-110)
  - [x] `func (c *Context) Status() Status` (line 96-106)
  - [x] `func (c *Context) Save()` (line 112-121)
  - [x] `func (c *Context) Restore()` (line 123-132)
  - [x] Package documentation explaining Context purpose (lines 1-58)
    - [x] Explains Context as central drawing object
    - [x] Documents drawing pipeline (6-step workflow)
    - [x] Covers lifecycle and resource management with example
    - [x] Explains state management (Save/Restore)
    - [x] Documents thread safety guarantees
- [x] Create `context/context_cgo.go`:
  - [x] CGO preamble
  - [x] `func contextCreate(sPtr unsafe.Pointer) ContextPtr`
  - [x] Calls `cairo_create`
  - [x] Sets up finalizer (in NewContext)
  - [x] CGO implementations for Close, Status, Save, Restore
  - [x] Helper to extract C surface pointer (uses surface.Ptr() method)
- [x] Create `context/context_test.go`:
  - [x] `TestNewContext` - creates context from ImageSurface ✅ PASSING
  - [x] `TestNewContextNilSurface` - tests error handling for nil surface ✅ PASSING
  - [x] `TestContextClose` - verifies close and double-close safety ✅ PASSING
  - [x] `TestContextStatus` - verifies status reporting ✅ PASSING
  - [x] `TestContextSaveRestore` - verifies save/restore stack ✅ PASSING
  - [x] `TestContextSaveRestoreImbalance` - tests restore without save ✅ PASSING
  - [x] BONUS: `TestContextCloseIndependentOfSurface` - verifies context/surface independence ✅ PASSING
  - [x] BONUS: `TestContextMultipleContextsOnSameSurface` - tests multiple contexts on one surface ✅ PASSING
  - [x] BONUS: `TestContextCreationWithDifferentSurfaceFormats` - tests all surface formats ✅ PASSING
- [x] Update `cairo.go`:
  - [x] Import context package (line 4)
  - [x] Re-export Context type with comprehensive documentation (lines 60-125)
    - [x] Basic usage pattern section
    - [x] Resource management section with example
    - [x] State stack section
    - [x] Thread safety documentation
  - [x] Re-export NewContext function with documentation and example (lines 127-152)

✅ **Status: 100% COMPLETE - All requirements met, all tests passing, comprehensive documentation added**

### Prompt 10: Context Package - Source Colors ✅ COMPLETE

- [x] Update `context/context.go`:
  - [x] `func (c *Context) SetSourceRGB(r, g, b float64)` (line 150-158) ✅
  - [x] `func (c *Context) SetSourceRGBA(r, g, b, a float64)` (line 172-180) ✅
  - [x] Document that r, g, b, a are in range [0.0, 1.0] - comprehensive docs ✅
  - [x] Both methods use proper locking (Lock/Unlock with nil pointer checks) ✅
- [x] Update `context/context_cgo.go`:
  - [x] `func contextSetSourceRGB(ptr ContextPtr, r, g, b float64)` (line 35-40) ✅
  - [x] `func contextSetSourceRGBA(ptr ContextPtr, r, g, b, a float64)` (line 42-47) ✅
  - [x] Both call appropriate cairo_set_source_* functions (cairo_set_source_rgb/rgba) ✅
- [x] Update `context/context_test.go`:
  - [x] `TestContextSetSourceRGB` - tests 6 color combinations ✅ PASSING
  - [x] `TestContextSetSourceRGBA` - tests 5 color/alpha combinations ✅ PASSING
  - [x] `TestContextSetSourceAfterClose` - verifies safe no-op behavior ✅ PASSING
  - [x] Integration test included in above tests (create context, set source, check status) ✅
- [x] Update `cairo/cairo.go`:
  - [x] Add complete example showing source color setting (lines 91-95) ✅
    - [x] Shows SetSourceRGB usage (opaque red)
    - [x] Shows SetSourceRGBA usage (semi-transparent blue)

✅ **Status: 100% COMPLETE - All requirements met, all 3 tests passing, documentation complete**

### Prompt 11: Context Package - Basic Path Operations ✅ COMPLETE

- [x] Update `context/context.go`:
  - [x] `func (c *Context) MoveTo(x, y float64)` (context.go:138-142) ✅
  - [x] `func (c *Context) LineTo(x, y float64)` (context.go:158-162) ✅
  - [x] `func (c *Context) Rectangle(x, y, width, height float64)` (context.go:182-186) ✅
  - [x] `func (c *Context) ClosePath()` (context.go:284-288) ✅
  - [x] `func (c *Context) NewPath()` (context.go:257-261) ✅
  - [x] `func (c *Context) NewSubPath()` (context.go:307-311) ✅
  - [x] `func (c *Context) GetCurrentPoint() (x, y float64, err error)` (context.go:204-217) ✅
  - [x] BONUS: `func (c *Context) HasCurrentPoint() bool` (context.go:233-241) ✅
  - [x] All methods use proper locking (withLock or RLock/RUnlock) ✅
  - [x] Document coordinate system (user-space coordinates) - comprehensive godoc comments added ✅
    - [x] MoveTo: documents user-space coordinates and current point behavior (lines 126-137)
    - [x] LineTo: documents user-space coordinates and current point updates (lines 144-157)
    - [x] Rectangle: documents user-space, equivalence to path operations (lines 164-181)
    - [x] GetCurrentPoint: documents user-space return values (lines 188-203)
    - [x] HasCurrentPoint: documents current point lifecycle (lines 219-232)
    - [x] NewPath: documents clearing path and current point (lines 243-256)
    - [x] ClosePath: documents sub-path closing and current point (lines 263-283)
    - [x] NewSubPath: documents sub-path creation (lines 290-306)
- [x] Update `context/context_cgo.go`:
  - [x] CGO implementations for all path methods (lines 49-90) ✅
  - [x] GetCurrentPoint checks status via HasCurrentPoint (context_cgo.go:64-74) ✅
- [x] Update `context/context_test.go`:
  - [x] `TestContextMoveTo` - ✅ PASSING
  - [x] `TestContextLineTo` - ✅ PASSING
  - [x] `TestContextRectangle` - ✅ PASSING
  - [x] `TestContextClosePath` - ✅ PASSING
  - [x] `TestContextNewPath` - ✅ PASSING
  - [x] `TestContextGetCurrentPoint` - ✅ PASSING
  - [x] `TestContextHasCurrentPointNoPoint` - ✅ PASSING (bonus)
  - [x] `TestContextGetCurrentPointNoPoint` - ✅ PASSING
  - [x] `TestContextPathOperationsAfterClose` - ✅ PASSING
- [x] Update `cairo/cairo.go`:
  - [x] Document path construction basics in Context type documentation (lines 104-123) ✅
    - [x] Explains path-based drawing model
    - [x] Lists basic path operations with signatures
    - [x] Provides triangle drawing example
  - [x] Document current point concept in Context type documentation (lines 125-146) ✅
    - [x] Explains current point lifecycle
    - [x] Documents user-space coordinates
    - [x] Shows GetCurrentPoint/HasCurrentPoint usage example
  - [x] Updated basic usage example to include actual path operations (lines 94-102) ✅

✅ **Status: 100% COMPLETE - All requirements met, all 9 tests passing, comprehensive documentation added**

### Prompt 12: Context Package - Fill and Stroke Operations

- [ ] Update `context/context.go`:
  - [ ] `func (c *Context) Fill()`
  - [ ] `func (c *Context) FillPreserve()`
  - [ ] `func (c *Context) Stroke()`
  - [ ] `func (c *Context) StrokePreserve()`
  - [ ] `func (c *Context) Paint()`
  - [ ] `func (c *Context) SetLineWidth(width float64)`
  - [ ] All methods use proper locking
  - [ ] Document that Fill/Stroke consume path, Preserve variants don't
- [ ] Update `context/context_cgo.go`:
  - [ ] CGO implementations for all rendering methods
  - [ ] Check if context is closed
- [ ] Update `context/context_test.go`:
  - [ ] `TestContextFill`
  - [ ] `TestContextFillPreserve`
  - [ ] `TestContextStroke`
  - [ ] `TestContextStrokePreserve`
  - [ ] `TestContextPaint`
  - [ ] `TestContextSetLineWidth`
  - [ ] `TestContextRenderAfterClose`
  - [ ] Integration test: NewPath, Rectangle, SetSourceRGB, Fill
- [ ] Update `cairo/cairo.go`:
  - [ ] Add example of complete drawing operation
  - [ ] Document Fill vs Stroke semantics

---

## Phase 4: MVP Validation (Prompt 13) 🎯

### Prompt 13: MVP Integration Test

- [ ] Create `examples/basic_shapes.go`:
  - [ ] Example function creating 400x400 PNG
  - [ ] Draws filled red rectangle (100, 100, 200, 200)
  - [ ] Draws blue stroked rectangle outline (120, 120, 160, 160)
  - [ ] Saves to output.png
  - [ ] Complete error handling
  - [ ] Proper cleanup with defer
- [ ] Create `examples/basic_shapes_test.go`:
  - [ ] Test that runs the example
  - [ ] Verifies PNG file is created
  - [ ] Checks file size is reasonable
  - [ ] Uses `t.TempDir()` for output
- [ ] Update `README.md`:
  - [ ] Add "Quick Start" section
  - [ ] Add code example showing rectangle drawing
  - [ ] Add build instructions
  - [ ] Note about Cairo system dependency
- [ ] Update `cairo/cairo.go`:
  - [ ] Complete package example showing MVP usage
  - [ ] Shows full workflow from surface to PNG
- [ ] **SUCCESS CRITERIA**:
  - [ ] Example compiles and runs without errors
  - [ ] PNG file is created and viewable
  - [ ] No memory leaks
  - [ ] All defers properly clean up resources
- [ ] **MILESTONE**: MVP Complete! 🎉

---

## Phase 5: Expansion (Prompts 14-34)

### Prompt 14: Enhanced Testing and Documentation

- [ ] Add benchmarks:
  - [ ] Create `context/context_bench_test.go`:
    - [ ] `BenchmarkContextCreation`
    - [ ] `BenchmarkContextPathOperations`
    - [ ] `BenchmarkContextFillOperations`
  - [ ] Create `surface/surface_bench_test.go`:
    - [ ] `BenchmarkImageSurfaceCreation`
    - [ ] `BenchmarkWriteToPNG`
- [ ] Add example tests:
  - [ ] Create `examples/example_test.go`:
    - [ ] `Example_drawRectangle`
    - [ ] `Example_fillAndStroke`
    - [ ] `Example_colorBlending`
- [ ] Improve package documentation:
  - [ ] `cairo/cairo.go`: expand with architecture overview
  - [ ] `context/context.go`: add Context lifecycle explanation
  - [ ] `surface/surface.go`: add Surface types overview
  - [ ] `matrix/matrix.go`: add transformation math explanation
  - [ ] `status/status.go`: add error handling guide
- [ ] Create `CONTRIBUTING.md`:
  - [ ] How to build
  - [ ] How to run tests
  - [ ] Code style guidelines
  - [ ] How to add new features
- [ ] Update `README.md`:
  - [ ] Add badges (build status placeholder)
  - [ ] Add table of contents
  - [ ] Add "Current Status" section
  - [ ] Add "Roadmap" section
- [ ] Verify coverage is >80% for all packages

### Prompt 15: Pattern Package - Foundation

- [ ] Create `pattern/pattern.go`:
  - [ ] `Pattern` interface with Close(), Status(), SetMatrix(), GetMatrix()
  - [ ] `BasePattern` struct with `sync.RWMutex`, ptr, closed flag
  - [ ] `SolidPattern` struct embedding BasePattern
  - [ ] `func NewSolidPatternRGB(r, g, b float64) (*SolidPattern, error)`
  - [ ] `func NewSolidPatternRGBA(r, g, b, a float64) (*SolidPattern, error)`
- [ ] Create `pattern/pattern_cgo.go`:
  - [ ] CGO preamble
  - [ ] BasePattern CGO methods
  - [ ] SolidPattern CGO constructor functions
- [ ] Create `pattern/pattern_test.go`:
  - [ ] `TestNewSolidPatternRGB`
  - [ ] `TestNewSolidPatternRGBA`
  - [ ] `TestPatternClose`
  - [ ] `TestPatternStatus`
  - [ ] `TestPatternMatrix`
- [ ] Update `context/context.go`:
  - [ ] `func (c *Context) SetSource(pattern Pattern)`
- [ ] Update `cairo/cairo.go`:
  - [ ] Re-export Pattern interface
  - [ ] Re-export NewSolidPattern* functions
- [ ] Integration test: Create solid pattern, set as Context source, draw

### Prompt 16: Context Package - Transformations

- [ ] Update `context/context.go`:
  - [ ] `func (c *Context) Translate(tx, ty float64)`
  - [ ] `func (c *Context) Scale(sx, sy float64)`
  - [ ] `func (c *Context) Rotate(angle float64)`
  - [ ] `func (c *Context) Transform(matrix *matrix.Matrix)`
  - [ ] `func (c *Context) SetMatrix(matrix *matrix.Matrix)`
  - [ ] `func (c *Context) GetMatrix() (*matrix.Matrix, error)`
  - [ ] `func (c *Context) IdentityMatrix()`
  - [ ] `func (c *Context) UserToDevice(x, y float64) (float64, float64)`
  - [ ] `func (c *Context) UserToDeviceDistance(dx, dy float64) (float64, float64)`
  - [ ] `func (c *Context) DeviceToUser(x, y float64) (float64, float64)`
  - [ ] `func (c *Context) DeviceToUserDistance(dx, dy float64) (float64, float64)`
- [ ] Update `context/context_cgo.go`:
  - [ ] CGO implementations for all transformation methods
  - [ ] Proper matrix conversion between Go and C
- [ ] Update `context/context_test.go`:
  - [ ] `TestContextTranslate`
  - [ ] `TestContextScale`
  - [ ] `TestContextRotate`
  - [ ] `TestContextTransform`
  - [ ] `TestContextGetSetMatrix`
  - [ ] `TestContextIdentityMatrix`
  - [ ] `TestContextCoordinateConversion`
  - [ ] `TestContextTransformationsCombined`
- [ ] Create `examples/transformations.go`:
  - [ ] Example showing translate, scale, rotate
  - [ ] Draw same shape at different transformations
  - [ ] Save to PNG
- [ ] Document transformation order and CTM concept
- [ ] Ensure transformations work with Save/Restore

### Prompt 17: Context Package - Advanced Path Operations

- [ ] Update `context/context.go`:
  - [ ] `func (c *Context) Arc(xc, yc, radius, angle1, angle2 float64)`
  - [ ] `func (c *Context) ArcNegative(xc, yc, radius, angle1, angle2 float64)`
  - [ ] `func (c *Context) CurveTo(x1, y1, x2, y2, x3, y3 float64)`
  - [ ] `func (c *Context) RelMoveTo(dx, dy float64)`
  - [ ] `func (c *Context) RelLineTo(dx, dy float64)`
  - [ ] `func (c *Context) RelCurveTo(dx1, dy1, dx2, dy2, dx3, dy3 float64)`
  - [ ] Document angles are in radians
- [ ] Update `context/context_cgo.go`:
  - [ ] CGO implementations for all arc and curve methods
- [ ] Update `context/context_test.go`:
  - [ ] `TestContextArc`
  - [ ] `TestContextArcNegative`
  - [ ] `TestContextCurveTo`
  - [ ] `TestContextRelativeOperations`
  - [ ] `TestContextCircle`
- [ ] Update `examples/basic_shapes.go`:
  - [ ] Add circle drawing example
  - [ ] Add curved path example
  - [ ] Save to circles.png
- [ ] Create `examples/gradients.go` placeholder:
  - [ ] Comment explaining gradients come next
  - [ ] Basic structure ready

### Prompt 18: Pattern Package - Gradient Patterns

- [ ] Update `pattern/pattern.go`:
  - [ ] `LinearGradient` struct embedding BasePattern
  - [ ] `func NewLinearGradient(x0, y0, x1, y1 float64) (*LinearGradient, error)`
  - [ ] `RadialGradient` struct embedding BasePattern
  - [ ] `func NewRadialGradient(cx0, cy0, radius0, cx1, cy1, radius1 float64)
    (*RadialGradient, error)`
  - [ ] `func (g *LinearGradient) AddColorStopRGB(offset, r, g, b float64)`
  - [ ] `func (g *LinearGradient) AddColorStopRGBA(offset, r, g, b, a float64)`
  - [ ] Same for RadialGradient
- [ ] Update `pattern/pattern_cgo.go`:
  - [ ] CGO implementations for gradient creation
  - [ ] CGO implementations for color stop methods
- [ ] Update `pattern/pattern_test.go`:
  - [ ] `TestNewLinearGradient`
  - [ ] `TestNewRadialGradient`
  - [ ] `TestLinearGradientColorStops`
  - [ ] `TestRadialGradientColorStops`
  - [ ] `TestGradientWithContext`
- [ ] Update `examples/gradients.go`:
  - [ ] Implement linear gradient example (red to blue)
  - [ ] Implement radial gradient example (white center to blue edge)
  - [ ] Draw shapes with gradients
  - [ ] Save to gradients.png
- [ ] Update `cairo/cairo.go`:
  - [ ] Re-export LinearGradient, RadialGradient types
  - [ ] Re-export gradient constructor functions
- [ ] Document color stops must be in [0.0, 1.0]
- [ ] Test gradients with transformations

### Prompt 19: Context Package - Line Styles

- [ ] Create `context/line_style.go`:
  - [ ] `LineCap` type as int with constants: LineCapButt, LineCapRound, LineCapSquare
  - [ ] `LineJoin` type as int with constants: LineJoinMiter, LineJoinRound, LineJoinBevel
  - [ ] `//go:generate stringer` for both types
- [ ] Update `context/context.go`:
  - [ ] `func (c *Context) SetLineCap(lineCap LineCap)`
  - [ ] `func (c *Context) GetLineCap() LineCap`
  - [ ] `func (c *Context) SetLineJoin(lineJoin LineJoin)`
  - [ ] `func (c *Context) GetLineJoin() LineJoin`
  - [ ] `func (c *Context) SetDash(dashes []float64, offset float64)`
  - [ ] `func (c *Context) GetDash() (dashes []float64, offset float64, err error)`
  - [ ] `func (c *Context) SetMiterLimit(limit float64)`
  - [ ] `func (c *Context) GetMiterLimit() float64`
  - [ ] `func (c *Context) GetLineWidth() float64`
- [ ] Update `context/context_cgo.go`:
  - [ ] CGO implementations for all line style methods
  - [ ] GetDash needs to query dash count first
- [ ] Update `context/context_test.go`:
  - [ ] `TestContextLineCap`
  - [ ] `TestContextLineJoin`
  - [ ] `TestContextDash`
  - [ ] `TestContextDashEmpty`
  - [ ] `TestContextMiterLimit`
  - [ ] `TestContextLineStyleCombinations`
- [ ] Create `examples/line_styles.go`:
  - [ ] Draw lines with different caps
  - [ ] Draw angles with different joins
  - [ ] Draw dashed lines with different patterns
  - [ ] Save to line_styles.png
- [ ] Update `cairo/cairo.go`:
  - [ ] Re-export LineCap and LineJoin types and constants
- [ ] Run `go generate ./...` to generate stringer code
- [ ] Test line styles persist across Save/Restore

### Prompt 20: Context Package - Clipping

- [ ] Update `context/context.go`:
  - [ ] `func (c *Context) Clip()`
  - [ ] `func (c *Context) ClipPreserve()`
  - [ ] `func (c *Context) ClipExtents() (x1, y1, x2, y2 float64)`
  - [ ] `func (c *Context) InClip(x, y float64) bool`
  - [ ] `func (c *Context) ResetClip()`
  - [ ] Document clipping is intersective
- [ ] Update `context/context_cgo.go`:
  - [ ] CGO implementations for all clipping methods
- [ ] Update `context/context_test.go`:
  - [ ] `TestContextClip`
  - [ ] `TestContextClipPreserve`
  - [ ] `TestContextClipExtents`
  - [ ] `TestContextInClip`
  - [ ] `TestContextResetClip`
  - [ ] `TestContextNestedClips`
  - [ ] `TestContextClipWithTransform`
- [ ] Create `examples/clipping.go`:
  - [ ] Example showing circular clip region
  - [ ] Draw shapes that get clipped
  - [ ] Example showing nested clips
  - [ ] Save to clipping.png
- [ ] Update `examples/basic_shapes_test.go`:
  - [ ] Add visual regression test helper function
- [ ] Document clip region interaction with Save/Restore

### Prompt 21: Surface Package - Surface Pattern

- [ ] Update `pattern/pattern.go`:
  - [ ] `SurfacePattern` struct embedding BasePattern
  - [ ] `func NewSurfacePattern(surface Surface) (*SurfacePattern, error)`
  - [ ] `Extend` type as int with constants: ExtendNone, ExtendRepeat,
    ExtendReflect, ExtendPad
  - [ ] `Filter` type as int with constants: FilterFast, FilterGood, FilterBest,
    FilterNearest, FilterBilinear
  - [ ] `func (p *SurfacePattern) SetExtend(extend Extend)`
  - [ ] `func (p *SurfacePattern) GetExtend() Extend`
  - [ ] `func (p *SurfacePattern) SetFilter(filter Filter)`
  - [ ] `func (p *SurfacePattern) GetFilter() Filter`
- [ ] Update `pattern/pattern_cgo.go`:
  - [ ] CGO implementation for surface pattern creation
  - [ ] CGO implementations for extend and filter methods
- [ ] Update `pattern/pattern_test.go`:
  - [ ] `TestNewSurfacePattern`
  - [ ] `TestSurfacePatternExtend`
  - [ ] `TestSurfacePatternFilter`
  - [ ] `TestSurfacePatternWithContext`
- [ ] Create `examples/patterns.go`:
  - [ ] Create small ImageSurface with checker pattern
  - [ ] Use as surface pattern with ExtendRepeat
  - [ ] Fill large rectangle with pattern
  - [ ] Test different extend modes
  - [ ] Save to patterns.png
- [ ] Update `cairo/cairo.go`:
  - [ ] Re-export SurfacePattern type
  - [ ] Re-export Extend and Filter types and constants
  - [ ] Re-export NewSurfacePattern function
- [ ] Document source surface must remain valid while pattern is in use

### Prompt 22: Context Package - Operators and Compositing

- [ ] Create `context/operator.go`:
  - [ ] `Operator` type as int
  - [ ] Constants for all cairo_operator_t values (28+ constants)
  - [ ] `//go:generate stringer -type=Operator`
- [ ] Update `context/context.go`:
  - [ ] `func (c *Context) SetOperator(op Operator)`
  - [ ] `func (c *Context) GetOperator() Operator`
  - [ ] Document operator effects and use cases
- [ ] Update `context/context_cgo.go`:
  - [ ] CGO implementations for operator methods
- [ ] Update `context/context_test.go`:
  - [ ] `TestContextSetOperator`
  - [ ] `TestContextGetOperator`
  - [ ] `TestContextOperatorDefault`
- [ ] Create `examples/compositing.go`:
  - [ ] Draw overlapping shapes with different operators
  - [ ] Show OperatorOver, OperatorAdd, OperatorMultiply, OperatorXor
  - [ ] Create visual comparison grid
  - [ ] Save to compositing.png
- [ ] Update `cairo/cairo.go`:
  - [ ] Re-export Operator type and all constants
- [ ] Run `go generate ./...` to generate stringer code
- [ ] Document porter-duff operators vs blend modes
- [ ] Test operators with alpha transparency

### Prompt 23: Context Package - Fill and Stroke Rules

- [ ] Create `context/fill_rule.go`:
  - [ ] `FillRule` type as int
  - [ ] Constants: FillRuleWinding, FillRuleEvenOdd
  - [ ] `//go:generate stringer -type=FillRule`
- [ ] Update `context/context.go`:
  - [ ] `func (c *Context) SetFillRule(fillRule FillRule)`
  - [ ] `func (c *Context) GetFillRule() FillRule`
  - [ ] `func (c *Context) FillExtents() (x1, y1, x2, y2 float64)`
  - [ ] `func (c *Context) StrokeExtents() (x1, y1, x2, y2 float64)`
  - [ ] `func (c *Context) PathExtents() (x1, y1, x2, y2 float64)`
  - [ ] `func (c *Context) InFill(x, y float64) bool`
  - [ ] `func (c *Context) InStroke(x, y float64) bool`
  - [ ] Document fill rule differences
- [ ] Update `context/context_cgo.go`:
  - [ ] CGO implementations for all methods
- [ ] Update `context/context_test.go`:
  - [ ] `TestContextFillRule`
  - [ ] `TestContextFillExtents`
  - [ ] `TestContextStrokeExtents`
  - [ ] `TestContextPathExtents`
  - [ ] `TestContextInFill`
  - [ ] `TestContextInStroke`
  - [ ] `TestContextFillRuleWindingVsEvenOdd`
- [ ] Create `examples/fill_rules.go`:
  - [ ] Draw self-intersecting star with winding rule
  - [ ] Draw same star with even-odd rule
  - [ ] Show visual difference
  - [ ] Save to fill_rules.png
- [ ] Update `cairo/cairo.go`:
  - [ ] Re-export FillRule type and constants
- [ ] Run `go generate ./...` to generate stringer code
- [ ] Test fill rules with complex self-intersecting paths

### Prompt 24: Font Package - Toy Font API

- [ ] Create `font/font.go`:
  - [ ] `Slant` type as int with constants: SlantNormal, SlantItalic, SlantOblique
  - [ ] `Weight` type as int with constants: WeightNormal, WeightBold
  - [ ] `//go:generate stringer` for both types
  - [ ] Package documentation explaining toy vs scaled font APIs
- [ ] Update `context/context.go`:
  - [ ] `func (c *Context) SelectFontFace(family string, slant font.Slant,
    weight font.Weight)`
  - [ ] `func (c *Context) SetFontSize(size float64)`
  - [ ] `func (c *Context) ShowText(text string)`
  - [ ] `func (c *Context) TextPath(text string)`
  - [ ] Document text is positioned at current point
- [ ] Update `context/context_cgo.go`:
  - [ ] CGO implementations for font methods
  - [ ] String conversion for family parameter
- [ ] Update `context/context_test.go`:
  - [ ] `TestContextSelectFontFace`
  - [ ] `TestContextSetFontSize`
  - [ ] `TestContextShowText`
  - [ ] `TestContextTextPath`
- [ ] Create `examples/text.go`:
  - [ ] Draw text with different fonts
  - [ ] Draw text with different sizes
  - [ ] Draw text with different slants and weights
  - [ ] Position text at different locations
  - [ ] Save to text.png
- [ ] Update `cairo/cairo.go`:
  - [ ] Re-export Font Slant and Weight types and constants
- [ ] Run `go generate ./...` to generate stringer code
- [ ] Document toy font API is platform-dependent
- [ ] Test text rendering with transformations

### Prompt 25: Font Package - Text Extents

- [ ] Create `font/extents.go`:
  - [ ] `TextExtents` struct with fields: XBearing, YBearing, Width, Height,
    XAdvance, YAdvance
  - [ ] `FontExtents` struct with fields: Ascent, Descent, Height, MaxXAdvance, MaxYAdvance
  - [ ] Document what each field represents
- [ ] Update `context/context.go`:
  - [ ] `func (c *Context) TextExtents(text string) (*font.TextExtents, error)`
  - [ ] `func (c *Context) FontExtents() (*font.FontExtents, error)`
  - [ ] Document use cases
- [ ] Update `context/context_cgo.go`:
  - [ ] CGO implementations for extents methods
  - [ ] Convert C extents structs to Go structs
- [ ] Update `context/context_test.go`:
  - [ ] `TestContextTextExtents`
  - [ ] `TestContextFontExtents`
  - [ ] `TestContextTextExtentsEmpty`
  - [ ] `TestContextExtentsWithDifferentFonts`
- [ ] Update `examples/text.go`:
  - [ ] Text alignment example (center, right-align using extents)
  - [ ] Multi-line text with proper spacing
  - [ ] Draw bounding boxes around text
  - [ ] Save to text_extents.png
- [ ] Test extents with transformations
- [ ] Document coordinate system for bearings and advances

### Prompt 26: Surface Package - PDF Surface

- [ ] Create `surface/pdf.go`:
  - [ ] `PDFSurface` struct embedding BaseSurface
  - [ ] `func NewPDFSurface(filename string, widthPt, heightPt float64)
    (*PDFSurface, error)`
  - [ ] `func (s *PDFSurface) SetSize(widthPt, heightPt float64)`
  - [ ] `func (s *PDFSurface) ShowPage()`
  - [ ] Document dimensions are in points (1/72 inch)
- [ ] Create `surface/pdf_cgo.go`:
  - [ ] Build tag: `// +build !nopdf`
  - [ ] CGO preamble: `#cgo pkg-config: cairo-pdf`
  - [ ] CGO implementations for PDF surface methods
- [ ] Update `surface/surface_test.go`:
  - [ ] `TestNewPDFSurface`
  - [ ] `TestPDFSurfaceSetSize`
  - [ ] `TestPDFSurfaceMultiPage`
  - [ ] Use `t.TempDir()` for test PDF files
- [ ] Create `examples/pdf_output.go`:
  - [ ] Create multi-page PDF document
  - [ ] Draw different content on each page
  - [ ] Include text, shapes, and gradients
  - [ ] Save to output.pdf
- [ ] Update `cairo/cairo.go`:
  - [ ] Re-export PDFSurface type
  - [ ] Re-export NewPDFSurface function
- [ ] Update `README.md`:
  - [ ] Add PDF surface to features list
  - [ ] Note Cairo PDF backend requirement
- [ ] Document PDF coordinate system
- [ ] Verify file is valid PDF

### Prompt 27: Surface Package - SVG Surface

- [ ] Create `surface/svg.go`:
  - [ ] `SVGSurface` struct embedding BaseSurface
  - [ ] `func NewSVGSurface(filename string, widthPt, heightPt float64)
    (*SVGSurface, error)`
  - [ ] `func (s *SVGSurface) SetDocumentUnit(unit SVGUnit)`
  - [ ] `SVGUnit` type as int with 10 constants
  - [ ] `//go:generate stringer -type=SVGUnit`
- [ ] Create `surface/svg_cgo.go`:
  - [ ] Build tag: `// +build !nosvg`
  - [ ] CGO preamble: `#cgo pkg-config: cairo-svg`
  - [ ] CGO implementations for SVG surface methods
- [ ] Update `surface/surface_test.go`:
  - [ ] `TestNewSVGSurface`
  - [ ] `TestSVGSurfaceDocumentUnit`
  - [ ] Use `t.TempDir()` for test SVG files
  - [ ] Verify SVG file starts with proper XML header
- [ ] Create `examples/svg_output.go`:
  - [ ] Create SVG with various shapes
  - [ ] Include gradients and patterns
  - [ ] Test text rendering in SVG
  - [ ] Save to output.svg
- [ ] Update `cairo/cairo.go`:
  - [ ] Re-export SVGSurface type
  - [ ] Re-export SVGUnit type and constants
  - [ ] Re-export NewSVGSurface function
- [ ] Update `README.md`:
  - [ ] Add SVG surface to features list
  - [ ] Note Cairo SVG backend requirement
- [ ] Run `go generate ./...` to generate stringer code
- [ ] Test SVG file is valid XML
- [ ] **MILESTONE**: Vector output formats complete! 🎉

### Prompt 28: Advanced Examples and Documentation

- [ ] Create `examples/dashboard.go`:
  - [ ] Complex data dashboard example
  - [ ] Multiple charts: bar, line, pie
  - [ ] Use gradients, patterns, text, transformations
  - [ ] Generate to both PNG and PDF
  - [ ] Well-commented
- [ ] Create `examples/animation.go`:
  - [ ] Generate sequence of PNG frames
  - [ ] Simple animation (rotating shape, moving object)
  - [ ] Document frame rate considerations
  - [ ] Show batch-processing frames
- [ ] Create/update `examples/README.md`:
  - [ ] Describe each example
  - [ ] Show output images
  - [ ] Explain what each demonstrates
  - [ ] Build and run instructions
- [ ] Update main `README.md`:
  - [ ] Add comprehensive feature list
  - [ ] Add performance notes
  - [ ] Add comparison to other Go graphics libraries
  - [ ] Add troubleshooting section
  - [ ] Add FAQ section
- [ ] Create `ARCHITECTURE.md`:
  - [ ] Explain package organization
  - [ ] Explain CGO boundary design
  - [ ] Explain memory management strategy
  - [ ] Explain thread safety approach
  - [ ] Include diagrams if helpful
- [ ] Create `docs/` directory with:
  - [ ] `QUICKSTART.md`: fast intro for new users
  - [ ] `MIGRATION.md`: guide for C Cairo users
  - [ ] `PERFORMANCE.md`: performance tips and benchmarks
- [ ] Ensure all examples compile and have error handling
- [ ] Add godoc examples for key functions

### Prompt 29: Error Handling Improvements

- [ ] Create `cairo/errors.go`:
  - [ ] `SurfaceError`, `ContextError`, `PatternError` types
  - [ ] Each wraps Status with additional context
  - [ ] `func (e *SurfaceError) Unwrap() error`
  - [ ] Implement `Is()` and `As()` for errors.Is/As support
- [ ] Update `status/status.go`:
  - [ ] More descriptive error messages in Error() method
  - [ ] Include suggestions for common errors
- [ ] Update all packages to use wrapped errors:
  - [ ] surface: wrap with SurfaceError
  - [ ] context: wrap with ContextError
  - [ ] pattern: wrap with PatternError
- [ ] Update tests across packages:
  - [ ] `TestErrorUnwrapping`
  - [ ] `TestErrorContext`
  - [ ] `TestErrorTypes`
- [ ] Update examples:
  - [ ] Show proper error handling patterns
  - [ ] Using errors.Is for error checking
- [ ] Update documentation:
  - [ ] Add error handling guide to README
  - [ ] Document common errors and solutions
  - [ ] Show idiomatic error checking patterns
- [ ] Ensure backward compatibility

### Prompt 30: Thread Safety Validation and Race Testing

- [ ] Create `test/race_test.go`:
  - [ ] `TestConcurrentContextOperations`
  - [ ] `TestConcurrentSurfaceOperations`
  - [ ] `TestConcurrentPatternOperations`
  - [ ] `TestContextSharedAcrossSurfaces`
  - [ ] `TestPatternSharedAcrossContexts`
  - [ ] All tests run with -race flag
- [ ] Create `benchmarks/concurrent_bench_test.go`:
  - [ ] `BenchmarkConcurrentDrawing`
  - [ ] `BenchmarkConcurrentSurfaceCreation`
  - [ ] `BenchmarkLockContention`
  - [ ] Compare single-threaded vs multi-threaded
- [ ] Update documentation:
  - [ ] Add concurrency section to README
  - [ ] Document safe concurrent usage patterns
  - [ ] Document unsafe patterns to avoid
  - [ ] Note performance implications
- [ ] Add stress test:
  - [ ] `test/stress_test.go` with build tag `// +build stress`
  - [ ] Long-running tests with many goroutines
  - [ ] Memory leak detection under concurrent load
- [ ] Update CI configuration:
  - [ ] Run tests with -race flag
  - [ ] Run benchmarks to detect regressions
  - [ ] Add stress test job (optional)
- [ ] Verify no deadlocks or race conditions

### Prompt 31: Memory Management Validation

- [ ] Create `test/memory_test.go`:
  - [ ] `TestNoLeaksSimpleDrawing`
  - [ ] `TestNoLeaksWithoutClose`
  - [ ] `TestNoLeaksPatternReuse`
  - [ ] `TestNoLeaksSurfaceReuse`
  - [ ] `TestMemoryWithDeferredClose`
  - [ ] Use runtime.GC() and runtime.ReadMemStats
- [ ] Create `benchmarks/memory_bench_test.go`:
  - [ ] `BenchmarkMemoryAllocation`
  - [ ] `BenchmarkFinalizerOverhead`
  - [ ] `BenchmarkExplicitClose`
- [ ] Add memory profiling example:
  - [ ] `examples/profile_memory.go`
  - [ ] Show how to profile memory usage
  - [ ] Include pprof integration example
- [ ] Update documentation:
  - [ ] Add memory management guide
  - [ ] Document Close() vs finalizer tradeoffs
  - [ ] Show best practices for long-running programs
  - [ ] Add troubleshooting section for memory issues
- [ ] Create `tools/leak_detector.go`:
  - [ ] Helper tool to detect resource leaks
  - [ ] Tracks creation/destruction
  - [ ] Debugging mode logging
- [ ] Test with `GODEBUG=cgocheck=2`
- [ ] Document memory usage patterns

### Prompt 32: Build System and Distribution

- [ ] Update `go.mod`:
  - [ ] Ensure all dependencies properly versioned
  - [ ] Add go directive for minimum Go version
  - [ ] Clean up unused dependencies
- [ ] Create install scripts:
  - [ ] `scripts/install-deps-linux.sh`
  - [ ] `scripts/install-deps-macos.sh`
  - [ ] Make scripts idempotent and documented
- [ ] Create build configuration:
  - [ ] `.golangci.yml`: comprehensive linter config
  - [ ] `Makefile`: alternative to Taskfile
  - [ ] Support different Cairo configurations
- [ ] Update `Taskfile.yaml`:
  - [ ] Add coverage task
  - [ ] Add bench task
  - [ ] Add docs task
  - [ ] Add install task
- [ ] Create GitHub templates:
  - [ ] `.github/ISSUE_TEMPLATE/bug_report.md`
  - [ ] `.github/ISSUE_TEMPLATE/feature_request.md`
  - [ ] `.github/PULL_REQUEST_TEMPLATE.md`
- [ ] Update CI/CD:
  - [ ] Add code coverage reporting
  - [ ] Add automated release workflow
  - [ ] Add dependency update automation
  - [ ] Add automatic documentation deployment
- [ ] Create `RELEASE.md`:
  - [ ] Document release process
  - [ ] Version numbering scheme
  - [ ] Changelog generation process
- [ ] Test installation on clean systems
- [ ] Verify pkg-config detection works

### Prompt 33: Performance Optimization Pass

- [ ] Profile and optimize:
  - [ ] Run CPU profiling on examples
  - [ ] Identify hot paths in CGO transitions
  - [ ] Optimize frequent operations
  - [ ] Reduce allocations in tight loops
- [ ] Add fast paths:
  - [ ] `context/fast.go` with optimized operations
  - [ ] Batch operations where possible
  - [ ] Consider unsafe optimizations
  - [ ] Document performance characteristics
- [ ] Update benchmarks:
  - [ ] Add micro-benchmarks for individual operations
  - [ ] Add macro-benchmarks for realistic workloads
  - [ ] Compare against C Cairo directly
  - [ ] Set performance budgets (max overhead %)
- [ ] Optimize memory usage:
  - [ ] Reduce allocations in hot paths
  - [ ] Pool frequently allocated objects if beneficial
  - [ ] Review slice and string conversions
- [ ] Add performance documentation:
  - [ ] `PERFORMANCE.md` with optimization guide
  - [ ] Document when to use batch operations
  - [ ] Show performance comparison tables
  - [ ] List known performance gotchas
- [ ] Create profiling examples:
  - [ ] `examples/profile_cpu.go`
  - [ ] `examples/benchmark_drawing.go`
  - [ ] Document how to profile user code
- [ ] Measure before and after each optimization
- [ ] Ensure optimizations don't break thread safety

### Prompt 34: Final Polish and Release Preparation

- [ ] Documentation audit:
  - [ ] Review all package documentation
  - [ ] Ensure all public functions have examples
  - [ ] Fix typos or unclear explanations
  - [ ] Verify all links work
  - [ ] Ensure consistent terminology
- [ ] API review:
  - [ ] Review naming consistency
  - [ ] Verify Go idioms are followed
  - [ ] Check for awkward APIs
  - [ ] Ensure error handling is consistent
  - [ ] Verify all exported functions are necessary
- [ ] Test coverage review:
  - [ ] Ensure >80% coverage across all packages
  - [ ] Add tests for edge cases
  - [ ] Add integration tests for common workflows
  - [ ] Verify all examples work and are tested
- [ ] Create migration aids:
  - [ ] `API_REFERENCE.md`: complete API listing
  - [ ] `EXAMPLES_INDEX.md`: all examples categorized
  - [ ] `COMPARISON.md`: compare to other Go graphics libraries
- [ ] Prepare release:
  - [ ] Create `CHANGELOG.md` for v0.1.0
  - [ ] Tag appropriate features as v0.1.0 scope
  - [ ] Create GitHub release with notes
- [ ] Create announcement materials:
  - [ ] README badge for version
  - [ ] Social media announcement text
  - [ ] Blog post outline
  - [ ] Hacker News submission template
- [ ] Final testing:
  - [ ] Clean checkout test
  - [ ] Fresh system test (Linux and macOS)
  - [ ] Documentation test: follow README from scratch
  - [ ] Example test: run all examples and verify output
- [ ] **MILESTONE**: v0.1.0 Release Ready! 🎉🎉🎉

---

## Success Criteria

### Functional

- [ ] Can reproduce Cairo's basic examples
- [ ] All surface types work (Image, PDF, SVG)
- [ ] All drawing operations work correctly

### Performance

- [ ] Minimal overhead over C Cairo (<10%)
- [ ] Benchmarks show acceptable performance
- [ ] No obvious performance regressions

### Safety

- [ ] Thread-safe operations
- [ ] No memory leaks detected
- [ ] No data races found with `-race` flag
- [ ] Proper resource cleanup

### Quality

- [ ] >80% test coverage across all packages
- [ ] All tests pass on Linux and macOS
- [ ] Visual regression tests pass
- [ ] Examples compile and run

### Idiomatic Go

- [ ] Feels natural to Go developers
- [ ] Follows Go best practices
- [ ] Clear error handling
- [ ] Good documentation

### Documentation

- [ ] Clear package documentation
- [ ] Working examples for all major features
- [ ] Architecture documentation
- [ ] Installation and usage guides

---

## Testing Strategy

### Continuous Testing

- [ ] Run `go test -race ./...` after each prompt
- [ ] Run benchmarks periodically to catch regressions
- [ ] Use `go test -cover ./...` to monitor coverage
- [ ] Test on both Linux and macOS throughout

### Key Test Phases

- [ ] After Prompt 13: Run MVP example and verify PNG output
- [ ] After Prompt 18: Test all gradient types
- [ ] After Prompt 27: Test all surface types comprehensively
- [ ] After Prompt 34: Full test suite passes, ready for release

---

## Build Requirements

- [ ] Go 1.23 or later installed
- [ ] Cairo 1.18 or later (system library) installed
- [ ] pkg-config installed
- [ ] Standard Go tooling:
  - [ ] gofmt
  - [ ] golangci-lint
  - [ ] stringer (`go install golang.org/x/tools/cmd/stringer@latest`)

---

## Notes

- Each prompt builds incrementally on previous work
- No orphaned code - everything integrates immediately
- Test-driven development throughout
- Early validation at Prompt 13 (MVP)
- All code must compile and pass tests before proceeding to next prompt

---

**Current Status**: Ready to begin
**Next Step**: Prompt 1 - Project Foundation
