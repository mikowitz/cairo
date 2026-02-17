# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Go-Cairo is a Go wrapper for the Cairo 2D graphics library using CGO. The project provides idiomatic Go bindings to Cairo's C API with proper memory management, thread safety, and error handling. The library maintains API parity with the C Cairo library while following Go conventions.

**Target:** Go 1.23-1.25 with Cairo 1.18+

## Essential Build Commands

### Building
```bash
# Build all packages
go build ./...

# Build specific package
go build ./surface
go build ./context

# Build examples
go build ./examples/...
go build ./cmd/...

# Using Task
task build
```

### Testing
```bash
# Run all tests with race detector (recommended)
go test -race ./...

# Run tests for specific package
go test -race ./surface
go test -race ./context

# Run with verbose output
go test -v ./...

# Coverage analysis
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Benchmarks
go test -bench=. -benchmem ./...

# Using Task
task test          # Runs tests with race detector
task watch         # Watch mode for tests
```

### Code Generation
```bash
# Generate code (stringer for enums)
go generate ./...

# Using Task
task generate
```

### Linting
```bash
# Run golangci-lint
golangci-lint run

# Using Task
task lint
```

## Architecture

### Package Structure

The library mirrors Cairo's conceptual organization into focused subpackages:

- **`status/`** - Error handling and status codes from Cairo operations
- **`surface/`** - Drawing targets (ImageSurface, PDF, SVG surfaces)
- **`context/`** - Main drawing interface with graphics state and operations
- **`pattern/`** - Sources for drawing (solid colors, gradients, textures)
- **`matrix/`** - 2D affine transformations for coordinate space conversions
- **Root package** - Re-exports common types for convenience

### CGO Integration Pattern

CGO code is strictly isolated in `*_cgo.go` files within each package. The pattern is:

- **`package.go`** - Pure Go API with thread-safe wrappers
- **`package_cgo.go`** - CGO bindings to Cairo C functions
- **`package_test.go`** - Unit tests
- **`doc.go`** - Package documentation

CGO functions use lowercase camelCase naming (e.g., `surfaceCreate`, `contextSetSource`) to distinguish them from public Go APIs.

### Type Mapping

| C Type | Go Type | Package |
|--------|---------|---------|
| `cairo_t*` | `*context.Context` | context |
| `cairo_surface_t*` | `surface.Surface` (interface) | surface |
| `cairo_pattern_t*` | `pattern.Pattern` (interface) | pattern |
| `cairo_matrix_t` | `matrix.Matrix` (struct) | matrix |
| `cairo_status_t` | `status.Status` (enum) | status |

### Thread Safety Pattern

All types wrapping Cairo objects embed `sync.RWMutex` and follow this pattern:

```go
type Context struct {
    sync.RWMutex
    ptr ContextPtr
}

// Read operations use RLock
func (c *Context) GetLineWidth() float64 {
    c.RLock()
    defer c.RUnlock()
    if c.ptr == nil {
        return 0.0
    }
    return contextGetLineWidth(c.ptr)
}

// Write operations use Lock
func (c *Context) SetLineWidth(width float64) {
    c.Lock()
    defer c.Unlock()
    if c.ptr == nil {
        return
    }
    contextSetLineWidth(c.ptr, width)
}
```

### Memory Management

- All Cairo resource types (Surface, Context, Pattern) have:
  - `Close()` method for explicit cleanup
  - Finalizers registered as safety net
  - Thread-safe access to underlying C pointers
- Always use `defer resource.Close()` after creation
- Check for `nil` pointers before accessing C resources

### Error Handling Approach

The library uses a hybrid error handling model matching Cairo's design:

- **Constructors and I/O**: Return `(result, error)` pattern
- **Drawing operations**: Set internal status, check via `Status()` method
- **Getters that can fail**: Return `(value, error)` pattern

Status checking example:
```go
ctx.Arc(50, 50, 30, 0, 2*math.Pi)
if status := ctx.Status(); status != status.Success {
    return fmt.Errorf("drawing failed: %v", status)
}
```

## Development Workflow

### Typical Usage Flow

1. Create a Surface (drawing target)
2. Create a Context for that Surface
3. Use Context methods to draw (paths, colors, transforms)
4. Flush the Surface
5. Export/save result (e.g., WriteToPNG)
6. Close resources (Context, then Surface)

### Code Generation

Some packages use `go generate` with `stringer` for enum types:
- `context/linecap_string.go`
- `context/linejoin_string.go`
- `pattern/patterntype_string.go`
- `surface/format_string.go`
- `status/status_string.go`

Always run `go generate ./...` after modifying enum types.

### Testing Organization

- `*_test.go` - Unit tests with >80% coverage target
- `*_bench_test.go` - Benchmark tests
- `examples/*_test.go` - Example tests demonstrating usage patterns
- `examples/test_harness.go` - Shared testing utilities for visual output validation

### Running Single Tests

```bash
# Run specific test
go test -v ./context -run TestContextArc

# Run specific test with race detector
go test -race -v ./surface -run TestImageSurface

# Run specific benchmark
go test -bench=BenchmarkContextCreate ./context
```

## Important Conventions

### Naming

- **Constructors**: Always named `New*` (e.g., `NewContext`, `NewImageSurface`)
- **CGO wrappers**: Lowercase camelCase with package prefix (e.g., `surfaceCreate`, `contextArc`)
- **Enums**: PascalCase with package context (e.g., `FormatARGB32`, `LineCapRound`)

### Documentation Requirements

- All exported types, functions, and constants must have doc comments
- Package-level documentation in `doc.go` files
- Start function docs with the function name: "NewContext creates..."
- Include usage examples for complex APIs
- Document thread safety considerations

### No Panics

Library code must never panic. Use error returns instead. Panics are acceptable in test code.

### Testing

Always use the stretchr/testify libraries when writing tests

### Integration Requirements

Per the development plan (`development/go_cairo_prompts.md`):
- All code must integrate immediately - no orphaned functions
- Every new function should be used in tests or examples
- Follow TDD principles: write tests first or alongside implementation
- Maintain backwards compatibility within the same major version

## Cairo-Specific Concepts

### Path-Based Drawing Model

Cairo uses paths constructed via operations like `MoveTo`, `LineTo`, `Arc`, then rendered with `Fill`/`Stroke`:

```go
ctx.MoveTo(50, 10)
ctx.LineTo(90, 90)
ctx.LineTo(10, 90)
ctx.ClosePath()
ctx.Fill()  // Path is consumed
```

### Preserve Variants

- `Fill()` and `Stroke()` consume the path
- `FillPreserve()` and `StrokePreserve()` keep the path for subsequent operations
- Use Preserve variants when you need both fill and stroke on the same shape

### Current Point

Cairo maintains a "current point" used as the starting point for path operations. It's set by `MoveTo` and updated by operations like `LineTo`. Check with `HasCurrentPoint()` or retrieve with `GetCurrentPoint()`.

### State Stack

Context maintains a graphics state stack:
- `Save()` pushes current state
- `Restore()` pops saved state
- Affects transformations, colors, line styles, etc.

## Reference Documentation

- **Cairo C Manual**: https://www.cairographics.org/manual/
- **Language Bindings Guide**: https://www.cairographics.org/manual/language-bindings.html
- **API Reference**: https://www.cairographics.org/manual/cairo-cairo-t.html
- **Development Plan**: `development/go_cairo_prompts.md`
- **Design Decisions**: `DESIGN.md`
- **Contributing Guide**: `CONTRIBUTING.md`

## Common Gotchas

1. **Always close resources**: Use `defer` to ensure `Close()` is called on all Cairo objects
2. **Check status after drawing**: Drawing operations don't return errors; check `ctx.Status()`
3. **Flush before export**: Call `surf.Flush()` before writing to PNG or other output
4. **Thread safety overhead**: While all types are thread-safe, sharing across goroutines may cause lock contention
5. **Premultiplied alpha**: Cairo uses premultiplied alpha internally, but constructors handle this for you
6. **Path consumption**: Remember that `Fill()` and `Stroke()` clear the current path
7. **Coordinate spaces**: Transformations affect both drawing and line widths
