# Architecture

This document explains the internal design of go-cairo: how packages are organized,
where the CGO boundary sits, how memory is managed, and how thread safety is enforced.

## Package Organization

The library mirrors Cairo's conceptual model in focused subpackages:

```
cairo/                  ← root package (re-exports for user convenience)
├── status/             ← error codes (cairo_status_t)
├── matrix/             ← 2D affine transforms (cairo_matrix_t)
├── font/               ← font type enums (Slant, Weight)
├── surface/            ← drawing targets (ImageSurface, PDFSurface, SVGSurface)
├── context/            ← drawing operations (cairo_t)
├── pattern/            ← paint sources (solid, gradients, surface)
└── examples/           ← runnable demonstrations
```

### Root Package Re-exports

The root `cairo` package provides type aliases and constructor wrappers so callers can
use a single import instead of importing every subpackage:

```go
import "github.com/mikowitz/cairo"

surf, _ := cairo.NewImageSurface(cairo.FormatARGB32, 640, 480)
ctx, _  := cairo.NewContext(surf)
```

Internally, `cairo.Format` is `= surface.Format`, `cairo.NewContext` delegates to
`context.NewContext`, etc. There is no logic in the root package beyond delegation.

### Build-Tag-Gated Packages

PDF and SVG surface support is optional. Files guarded by `//go:build !nopdf` and
`//go:build !nosvg` are compiled by default but can be excluded:

```bash
go build -tags nopdf,nosvg ./...   # ImageSurface only, no external backends
```

This keeps the core library buildable on systems without the Cairo PDF or SVG backends
installed.

---

## CGO Boundary Design

CGO code is strictly isolated in `*_cgo.go` files. Each package contains exactly two
kinds of Go files:

| File pattern | Contents |
|---|---|
| `package.go`, `type.go`, etc. | Pure Go: exported types, methods, constructors |
| `package_cgo.go` | CGO: `import "C"`, type aliases, C function wrappers |

### Naming Convention

CGO wrapper functions use **lowercase camelCase with a package prefix** to distinguish
them from the public API:

```
surfaceCreate(...)          ← CGO wrapper (unexported)
NewImageSurface(...)        ← public Go constructor
```

### Type Aliases at the Boundary

Each CGO file defines an opaque pointer type alias:

```go
// surface/surface_cgo.go
type SurfacePtr *C.cairo_surface_t

// context/context_cgo.go
type ContextPtr *C.cairo_t

// pattern/pattern_cgo.go
type PatternPtr *C.cairo_pattern_t
```

These aliases let pure-Go files refer to the C pointers without importing `"C"` themselves.
Only `*_cgo.go` files contain `import "C"`.

### pkg-config Integration

Cairo headers and link flags are resolved at build time via `pkg-config`:

```c
// #cgo pkg-config: cairo
// #include <cairo.h>
// #include <stdlib.h>
import "C"
```

Backend-specific files use their own pkg-config targets:
- `surface/pdf_cgo.go`: `#cgo pkg-config: cairo-pdf`
- `surface/svg_cgo.go`: `#cgo pkg-config: cairo-svg`

### String Conversion

C functions that accept strings require `C.CString` (heap-allocated) and explicit
`C.free`. This is always handled inside the `*_cgo.go` file:

```go
func surfaceWriteToPNG(ptr SurfacePtr, filepath string) error {
    cPath := C.CString(filepath)
    defer C.free(unsafe.Pointer(cPath))
    st := status.Status(C.cairo_surface_write_to_png(ptr, cPath))
    return st.ToError()
}
```

---

## Memory Management

Cairo objects are reference-counted in C. The Go wrapper takes ownership of each
newly-created pointer and is responsible for calling the matching destroy function.

### Two-Layer Cleanup

Every resource type uses both an explicit `Close()` method and a runtime finalizer:

```
User calls Close()           ← preferred, deterministic
    │
    └─► close() internal
            ├── calls C destroy (e.g., cairo_surface_destroy)
            ├── sets ptr = nil  (guards against double-free)
            └── removes finalizer (runtime.SetFinalizer(b, nil))

GC runs finalizer            ← safety net if Close() was forgotten
    └─► close() internal     (same path as above)
```

The finalizer is set in the constructor and cleared on first close:

```go
func newBaseSurface(ptr SurfacePtr) *BaseSurface {
    b := &BaseSurface{ptr: ptr}
    runtime.SetFinalizer(b, (*BaseSurface).close)
    return b
}

func (b *BaseSurface) close() error {
    b.Lock()
    defer b.Unlock()
    if b.ptr != nil {
        surfaceClose(b.ptr)
        runtime.SetFinalizer(b, nil)  // prevent double-free
        b.ptr = nil
    }
    return nil
}
```

### Nil Pointer After Close

Every method checks `ptr == nil` before calling into C. This means calling any method on
a closed object is a safe no-op rather than a crash.

### Embedding for Code Reuse

Concrete types embed a base type to inherit memory management without duplication:

```
BaseSurface          (manages ptr, implements Surface interface)
    ↑ embedded by
ImageSurface         (adds format, width, height fields)
PDFSurface           (adds SetSize, ShowPage)
SVGSurface           (adds SetDocumentUnit)

BasePattern          (manages ptr, implements Pattern interface)
    ↑ embedded by
SolidPattern
LinearGradient
RadialGradient
SurfacePattern
```

---

## Thread Safety

All types that wrap a Cairo C pointer embed `sync.RWMutex` directly:

```go
type BaseSurface struct {
    sync.RWMutex
    ptr SurfacePtr
}

type Context struct {
    sync.RWMutex
    ptr ContextPtr
}

type BasePattern struct {
    sync.RWMutex
    ptr         PatternPtr
    patternType PatternType
}
```

### Lock Discipline

- **Read-only queries** (Status, getters) use `RLock/RUnlock`
- **Mutations** (all drawing ops, setters, Close) use `Lock/Unlock`
- **Nil check always happens inside the lock** to avoid TOCTOU races

```go
// Read operation
func (b *BaseSurface) Status() status.Status {
    b.RLock()
    defer b.RUnlock()
    if b.ptr == nil {
        return status.NullPointer
    }
    return surfaceStatus(b.ptr)
}

// Write operation
func (b *BaseSurface) Flush() {
    b.Lock()
    defer b.Unlock()
    if b.ptr == nil {
        return
    }
    surfaceFlush(b.ptr)
}
```

### withLock Helper

The `context` package provides a `withLock` helper for write operations with a
nil-pointer guard:

```go
func (c *Context) withLock(fn func()) {
    c.Lock()
    defer c.Unlock()
    if c.ptr == nil {
        return
    }
    fn()
}
```

### Contention Note

Cairo itself is not thread-safe at the C level. The Go mutexes protect the Go-side
pointer from concurrent access, but sharing a single `Context` across goroutines will
serialize all drawing calls through the mutex. For parallel rendering, use separate
surfaces and contexts per goroutine.

---

## Error Handling Model

The library uses a hybrid model that matches Cairo's own design:

| Situation | Return style |
|---|---|
| Constructors and I/O | `(result, error)` |
| Drawing operations | no return; check `ctx.Status()` |
| Getters that can fail | `(value, error)` |

`status.Status` implements the `error` interface. `status.Success` converts to `nil`;
all other values produce a descriptive error string via the generated `stringer`.

```go
// Constructor: error returned immediately
surf, err := surface.NewImageSurface(surface.FormatARGB32, 800, 600)

// Drawing: status checked after the fact
ctx.Arc(50, 50, 30, 0, 2*math.Pi)
if st := ctx.Status(); st != status.Success {
    return st
}

// Getter: error returned alongside value
m, err := pattern.GetMatrix()
```

---

## Code Generation

Enum types use `go generate` with the `stringer` tool. The generated `*_string.go`
files are committed alongside the source:

```
font/slant_string.go
font/weight_string.go
surface/format_string.go
surface/svgunit_string.go
status/status_string.go
pattern/patterntype_string.go
context/linecap_string.go
context/linejoin_string.go
```

Run `go generate ./...` after modifying any enum type.
