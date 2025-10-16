# Go-Cairo: Golang Wrapper for Cairo Graphics Library

## Project Overview

A Golang wrapper for the Cairo 2D graphics library using CGO bindings to the C library directly. The wrapper maintains API parity with C Cairo while providing a Go-native API experience, following the guidelines at <https://www.cairographics.org/manual/language-bindings.html>.

### Target Audience

General-purpose graphics developers who are either already working in Go or prefer Go's modern tooling and ecosystem over C, while still needing Cairo's powerful 2D graphics capabilities.

### Core Principles

- API parity with Cairo C library
- Go-idiomatic where it doesn't conflict with Cairo's model
- Performance-conscious (minimize overhead)
- Thread-safe by default
- Start minimal, iterate toward completeness

## Architecture

### Package Structure

```
cairo/
  .git/
  go.mod
  go.sum
  Taskfile.yaml
  README.md
  DESIGN.md
  cairo.go           # Main package with re-exports
  cairo_test.go

  context/
    context.go       # Public API
    context_cgo.go   # CGO bindings
    context_test.go

  surface/
    surface.go
    surface_cgo.go
    surface_test.go
    image.go         # ImageSurface implementation
    pdf.go           # PDFSurface implementation
    svg.go           # SVGSurface implementation

  status/
    status.go
    status_cgo.go
    status_test.go

  pattern/
    pattern.go
    pattern_cgo.go
    pattern_test.go
    solid.go
    gradient.go
    surface_pattern.go

  matrix/
    matrix.go
    matrix_cgo.go
    matrix_test.go

  font/
    font.go
    font_cgo.go
    font_test.go

  examples/
    basic_shapes.go
    gradients.go
    text.go
    transformations.go
```

### Import Path

Initial: `github.com/[username]/cairo`
(Designed to potentially migrate to community ownership)

### Main Package Re-exports

The top-level `cairo` package should re-export:

- All Status constants and types
- All Format constants (FormatARGB32, etc.)
- Common enums (Operator, LineCap, LineJoin)
- Constructor functions:
  - `NewContext(surface Surface) *Context`
  - `NewImageSurface(format Format, width, height int) *ImageSurface`
  - `NewPDFSurface(filename string, widthPt, heightPt float64) *PDFSurface`
  - `NewSVGSurface(filename string, widthPt, heightPt float64) *SVGSurface`
  - `NewMatrix() *Matrix`
  - `NewMatrixIdentity() *Matrix`

## Type System

### Naming Conventions

Following Cairo language binding guidelines:

- `cairo_context_t` → `context.Context`
- `cairo_surface_t` → `surface.Surface` (interface)
- `cairo_matrix_t` → `matrix.Matrix`
- Method names use Go casing: `cairo_line_to()` → `LineTo()`

### Core Types

```go
// context/context.go
type Context struct {
    sync.RWMutex
    ptr *C.cairo_t
}

// surface/surface.go
type Surface interface {
    Close() error
    Status() Status
    GetReferenceCount() uint32
    Reference()
    WriteToPNG(filename string) error
    // ... common surface methods
}

type BaseSurface struct {
    sync.RWMutex
    ptr *C.cairo_surface_t
}

type ImageSurface struct {
    BaseSurface
    // image-specific fields if needed
}

// matrix/matrix.go
type Matrix struct {
    sync.RWMutex
    ptr *C.cairo_matrix_t
    XX, YX float64
    XY, YY float64
    X0, Y0 float64
}

// pattern/pattern.go
type Pattern interface {
    Close() error
    Status() Status
    GetReferenceCount() uint32
    Reference()
    // ... common pattern methods
}

type BasePattern struct {
    sync.RWMutex
    ptr *C.cairo_pattern_t
}

type SolidPattern struct {
    BasePattern
}

type LinearGradient struct {
    BasePattern
}
```

### Enums and Constants

```go
// status/status.go
type Status int

const (
    StatusSuccess Status = iota
    StatusNoMemory
    StatusInvalidRestore
    // ... all cairo_status_t values
)

//go:generate stringer -type=Status

// surface/format.go
type Format int

const (
    FormatInvalid Format = -1
    FormatARGB32 Format = 0
    FormatRGB24 Format = 1
    // ... all cairo_format_t values
)

//go:generate stringer -type=Format
```

## Error Handling

### Hybrid Approach

**Critical operations return `(result, error)`:**

- Surface creation
- Context creation
- File I/O operations (WriteToPNG, etc.)
- Operations that can fail due to external factors

```go
func NewImageSurface(format Format, width, height int) (*ImageSurface, error)
func (s *Surface) WriteToPNG(filename string) error
```

**Drawing operations set status internally:**

- Path construction (MoveTo, LineTo, etc.)
- Rendering operations (Fill, Stroke, etc.)
- Transformations

```go
func (c *Context) LineTo(x, y float64)
func (c *Context) Fill()
func (c *Context) Status() Status
```

**Getters that can fail return errors:**

```go
func (c *Context) GetCurrentPoint() (x, y float64, err error)
func (c *Context) GetDash() (dashes []float64, offset float64, err error)
```

## Memory Management

### Lifecycle Management

- All objects implement `Close()` method (implements `io.Closer`)
- Runtime finalizers set as safety net
- Reference counting exposed but discouraged in documentation

```go
func (c *Context) Close() error {
    c.Lock()
    defer c.Unlock()
    if c.ptr != nil {
        C.cairo_destroy(c.ptr)
        c.ptr = nil
    }
    return nil
}

func newContext(ptr *C.cairo_t) *Context {
    c := &Context{ptr: ptr}
    runtime.SetFinalizer(c, (*Context).Close)
    return c
}
```

### String Handling

- Use Go `string` type in public API
- Convert to C strings in CGO layer
- Free C strings after use

```go
// In surface_cgo.go
func writeToPNG(surface *C.cairo_surface_t, filename string) error {
    cFilename := C.CString(filename)
    defer C.free(unsafe.Pointer(cFilename))
    status := C.cairo_surface_write_to_png(surface, cFilename)
    return statusToError(status)
}
```

### Image Data Access

- Use zero-copy approach by default
- Return `[]byte` slice pointing to Cairo's memory
- Document lifetime constraints clearly

```go
func (s *ImageSurface) GetData() []byte {
    // Must call Flush() before accessing
    // Must call MarkDirty() after modifications
    // Slice invalid after surface.Close()
    ptr := C.cairo_image_surface_get_data(s.ptr)
    stride := C.cairo_image_surface_get_stride(s.ptr)
    height := C.cairo_image_surface_get_height(s.ptr)
    return unsafe.Slice((*byte)(ptr), stride*height)
}
```

## Thread Safety

### Embedded RWMutex

All stateful types embed `sync.RWMutex`:

- Read operations use `RLock()/RUnlock()`
- Write operations use `Lock()/Unlock()`

```go
func (c *Context) LineTo(x, y float64) {
    c.Lock()
    defer c.Unlock()
    C.cairo_line_to(c.ptr, C.double(x), C.double(y))
}

func (c *Context) Status() Status {
    c.RLock()
    defer c.RUnlock()
    return Status(C.cairo_status(c.ptr))
}
```

## CGO Implementation

### Organization

- CGO code isolated in `*_cgo.go` files within each package
- Direct C function calls for performance
- Type conversions handled in CGO layer

```go
// context_cgo.go
// #include <cairo.h>
// #cgo pkg-config: cairo
import "C"

func (c *Context) lineTo(x, y float64) {
    C.cairo_line_to(c.ptr, C.double(x), C.double(y))
}
```

### Build Configuration

Platform-specific build tags and pkg-config:

```go
// surface_cgo_linux.go
// +build linux

// #cgo pkg-config: cairo cairo-png cairo-pdf cairo-svg cairo-ps
// #include <cairo.h>
import "C"
```

```go
// surface_cgo_darwin.go
// +build darwin

// #cgo pkg-config: cairo cairo-png cairo-pdf cairo-svg cairo-ps
// #include <cairo.h>
import "C"
```

## Implementation Plan

### Phase 1: Foundation (MVP)

1. **Status package** - No dependencies
   - Status type and constants
   - Stringer generation
   - Error conversion utilities

2. **Matrix package** - Depends on Status
   - Matrix type with C struct layout
   - Basic operations (Init, InitIdentity, Multiply)

3. **Pattern package (minimal)** - Depends on Status, Matrix
   - Pattern interface
   - SolidPattern only initially

4. **Surface package (minimal)** - Depends on Status
   - Surface interface
   - BaseSurface implementation
   - ImageSurface with PNG write support

5. **Context package (minimal)** - Depends on all above
   - Basic path operations: MoveTo, LineTo, Rectangle
   - Basic rendering: Fill, Stroke
   - Solid color source setting

### MVP Validation Criteria

Ability to:

- Create an image surface
- Create a context
- Set a solid color source
- Draw a rectangle
- Fill/stroke the path
- Save to PNG

### Phase 2: Core Graphics

- Complete pattern types (LinearGradient, RadialGradient, SurfacePattern)
- Complete path operations (Arc, Curve, ClosePath, etc.)
- Transformations (Translate, Rotate, Scale)
- State management (Save, Restore)
- Line styles and caps

### Phase 3: Text and Fonts

- Toy text API
- Font face management
- Text extents
- Platform-specific font rendering considerations

### Phase 4: Additional Surfaces

- PDF surface
- SVG surface
- PostScript surface
- Platform-specific surfaces (XLib, Win32, Quartz)

### Phase 5: Advanced Features

- Pango integration for complex text
- Cairo script support
- OpenGL surface support
- Convenience methods and Go-idiomatic helpers

## Testing Strategy

### Test Types

**Unit Tests**

- Getters/setters verification
- Error status checking
- Reference counting
- Thread safety validation

**Integration Tests**

- Port Cairo's example programs
- End-to-end drawing operations
- Cross-package interactions

**Visual Regression Tests**

- Compare rendered output against reference images
- Per-platform baselines (font rendering differences)
- Use image comparison with tolerance

### Test Structure

```go
// context_test.go
func TestContextCreation(t *testing.T) { }
func TestContextStatusPropagation(t *testing.T) { }

// examples/example_test.go
func Example_DrawCircle() {
    surface, _ := cairo.NewImageSurface(cairo.FormatARGB32, 256, 256)
    defer surface.Close()

    ctx := cairo.NewContext(surface)
    defer ctx.Close()

    ctx.Arc(128, 128, 64, 0, 2*math.Pi)
    ctx.Fill()

    surface.WriteToPNG("circle.png")
    // Output: Creates circle.png
}
```

### Visual Test Framework

```go
func TestVisualRegression(t *testing.T) {
    testCases := []struct{
        name string
        draw func(*Context)
        reference string
    }{
        {"simple_rectangle", drawRectangle, "testdata/rectangle.png"},
        // ...
    }

    for _, tc := range testCases {
        surface := createTestSurface()
        ctx := cairo.NewContext(surface)
        tc.draw(ctx)

        compareImages(t, surface, tc.reference)
    }
}
```

## Documentation

### Approach

- Adapt Cairo's C documentation to Go
- Maintain consistency with Cairo's descriptions
- Update function signatures and examples for Go

### Example Documentation

```go
// LineTo adds a line to the path from the current point to position (x, y)
// in user-space coordinates. After this call the current point will be (x, y).
//
// If there is no current point before the call to LineTo this function will
// behave as MoveTo(x, y).
//
// This wraps cairo_line_to().
func (c *Context) LineTo(x, y float64)
```

## Build and Development

### Requirements

- Go 1.23+ (targeting 1.25, supporting 1.23-1.25)
- Cairo 1.18
- pkg-config

### Taskfile Commands

```yaml
version: '3'

tasks:
  generate:
    desc: Generate stringer code
    cmds:
      - go generate ./...

  test:
    desc: Run all tests
    cmds:
      - go test -race ./...

  test-visual:
    desc: Run visual regression tests
    cmds:
      - go test -tags=visual ./examples

  build-examples:
    desc: Build all examples
    cmds:
      - go build ./examples/...

  lint:
    desc: Run linters
    cmds:
      - golangci-lint run
```

### CI/CD Configuration

**Platforms**: Linux (Ubuntu latest), macOS (latest)
**Go Versions**: 1.23, 1.24, 1.25
**Cairo Version**: 1.18

```yaml
# .github/workflows/test.yml
name: Test
on: [push, pull_request]

jobs:
  test:
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest]
        go: ['1.23', '1.24', '1.25']
    runs-on: ${{ matrix.os }}
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: ${{ matrix.go }}
      - name: Install Cairo (Ubuntu)
        if: matrix.os == 'ubuntu-latest'
        run: |
          sudo apt-get update
          sudo apt-get install -y libcairo2-dev
      - name: Install Cairo (macOS)
        if: matrix.os == 'macos-latest'
        run: brew install cairo
      - run: go generate ./...
      - run: go test -race ./...
```

## Version Information

### Minimum Requirements

- Cairo: 1.18
- Features requiring newer Cairo versions documented per-method
- Version detection function provided

```go
func Version() (major, minor, micro int) {
    cmajor := C.int(0)
    cminor := C.int(0)
    cmicro := C.int(0)
    C.cairo_version(&cmajor, &cminor, &cmicro)
    return int(cmajor), int(cminor), int(cmicro)
}
```

## Future Enhancements (Post-MVP)

- [ ] Functional options for constructors
- [ ] `WithState(func())` convenience wrapper for Save/Restore
- [ ] `color.Color` interface compatibility
- [ ] Builder pattern for complex surface creation
- [ ] Structured types for Point, Rectangle, etc.
- [ ] Advanced debugging and logging features
- [ ] Fine-grained build tags for optional Cairo features
- [ ] Comprehensive benchmark suite

## Success Criteria

1. **Functional**: Can reproduce Cairo's basic examples
2. **Performant**: Minimal overhead over C Cairo
3. **Safe**: Thread-safe, no memory leaks
4. **Idiomatic**: Feels natural to Go developers
5. **Compatible**: Works on Linux and macOS with standard Cairo installations
6. **Documented**: Clear examples and API documentation
7. **Tested**: Visual regression tests pass

---

This specification provides a complete blueprint for implementing the Go-Cairo wrapper. The phased approach ensures quick validation of core concepts while building toward a comprehensive graphics library.

