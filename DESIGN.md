# Design Documentation

## Architecture Overview

This library provides Go bindings to the Cairo 2D graphics library using CGO.
The design follows the guidelines outlined in the [Cairo Language Bindings documentation](https://www.cairographics.org/manual/language-bindings.html).

## Core Design Principles

1. **API Parity**: Maintain compatibility with the C Cairo API
2. **Go-Idiomatic**: Follow Go conventions where possible
3. **Thread-Safe**: All operations protected by mutexes
4. **Memory-Safe**: Proper resource cleanup with finalizers and Close() methods

## Reference Documentation

### Cairo C API Documentation

- **Main Documentation**: [Cairo Graphics Manual](https://www.cairographics.org/manual/)
- **Language Bindings Guide**: [Language Bindings Documentation](https://www.cairographics.org/manual/language-bindings.html)
- **API Reference**: [Cairo API Reference](https://www.cairographics.org/manual/cairo-cairo-t.html)

### Key Cairo Concepts

- **Context** (`cairo_t`): The main drawing object
- **Surface** (`cairo_surface_t`): The target for drawing operations
- **Pattern** (`cairo_pattern_t`): Color sources (solid, gradient, image)
- **Matrix** (`cairo_matrix_t`): Transformation matrices
- **Path**: Vector paths constructed from lines, curves, and arcs

## Implementation Details

For complete implementation specifications, see `development/go-cairo-spec.md`.

### Package Structure

The library is organized into subpackages that mirror Cairo's conceptual organization:

- `status`: Error handling and status codes
- `surface`: Surface types (Image, PDF, SVG, etc.)
- `context`: Drawing context and operations
- `pattern`: Color sources and patterns
- `matrix`: Transformation matrices
- `font`: Text rendering (planned)

### Type Mapping

Cairo C types map to Go types following these conventions:

| C Type | Go Type |
|--------|---------|
| `cairo_t*` | `*context.Context` |
| `cairo_surface_t*` | `surface.Surface` (interface) |
| `cairo_pattern_t*` | `pattern.Pattern` (interface) |
| `cairo_matrix_t` | `matrix.Matrix` (struct) |
| `cairo_status_t` | `status.Status` (enum) |

### Error Handling

The library uses a hybrid error handling approach:

- **Critical operations** (creation, I/O) return `(result, error)`
- **Drawing operations** set internal status, check via `Status()` method
- **Getters that can fail** return `(value, error)`

This mirrors Cairo's design while providing Go-idiomatic error handling where appropriate.

## CGO Integration

CGO code is isolated in `*_cgo.go` files within each package. The library uses
`pkg-config` to locate Cairo headers and libraries at build time.

Build tags handle platform-specific differences (Linux, macOS, Windows).

## Development Phases

See `development/go-cairo-spec.md` for the complete phased implementation plan,
from MVP through advanced features.
