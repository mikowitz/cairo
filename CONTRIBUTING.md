# Contributing to Go-Cairo

Thank you for your interest in contributing to Go-Cairo! This document provides
guidelines and instructions for contributing to the project.

## Table of Contents

- [Getting Started](#getting-started)
- [Development Environment Setup](#development-environment-setup)
- [How to Build](#how-to-build)
- [How to Run Tests](#how-to-run-tests)
- [Code Style Guidelines](#code-style-guidelines)
- [How to Add New Features](#how-to-add-new-features)
- [Pull Request Process](#pull-request-process)
- [Reporting Issues](#reporting-issues)

## Getting Started

Go-Cairo is a Go wrapper around the Cairo 2D graphics library. The project
follows a structured, incremental development approach with each feature
fully tested before moving forward.

### Prerequisites

Before contributing, ensure you have:

- **Go 1.23 or higher** - The project targets Go 1.23-1.25
- **Cairo 1.18 or higher** - System library for graphics operations
- **pkg-config** - Required for CGO build configuration
- **golangci-lint** - For code quality checks
- **git** - For version control

## Development Environment Setup

### Installing Cairo

**Ubuntu/Debian:**
```bash
sudo apt-get update
sudo apt-get install -y libcairo2-dev pkg-config
```

**macOS:**
```bash
brew install cairo pkg-config
```

### Installing Go Tools

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Install stringer (for code generation)
go install golang.org/x/tools/cmd/stringer@latest

# Install Task (optional, for using Taskfile)
go install github.com/go-task/task/v3/cmd/task@latest
```

### Cloning the Repository

```bash
git clone https://github.com/mikowitz/cairo.git
cd cairo
```

### Verify Setup

```bash
# Check Cairo installation
pkg-config --modversion cairo

# Verify Go version
go version

# Run initial build
go build ./...
```

## How to Build

### Building the Entire Project

```bash
# Using go command
go build ./...

# Or using Task
task build
```

### Building Specific Packages

```bash
# Build just the surface package
go build ./surface

# Build just the context package
go build ./context
```

### Running Code Generation

Some packages use code generation (e.g., `stringer` for enum types):

```bash
# Using go generate
go generate ./...

# Or using Task
task generate
```

### Building Examples

```bash
# Build all examples
go build ./examples/...
go build ./cmd/...

# Build specific example
go build ./cmd/basic_shapes
```

## How to Run Tests

### Running All Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with race detector (recommended)
go test -race ./...

# Or using Task
task test
```

### Running Tests for Specific Packages

```bash
# Test just the surface package
go test ./surface

# Test with verbose output
go test -v ./surface

# Test with race detector
go test -race ./surface
```

### Running Coverage Analysis

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View coverage report in browser
go tool cover -html=coverage.out

# Or using Task
task coverage
```

### Running Benchmarks

```bash
# Run all benchmarks
go test -bench=. ./...

# Run benchmarks for specific package
go test -bench=. ./context

# Run benchmarks with memory allocation stats
go test -bench=. -benchmem ./...
```

### Test Organization

Tests are organized as follows:

- `*_test.go` - Unit tests for each package
- `*_bench_test.go` - Benchmark tests (when available)
- `examples/*_test.go` - Example tests demonstrating usage

## Code Style Guidelines

### General Principles

1. **Follow Go idioms** - Write idiomatic Go code
2. **Be explicit** - Clarity over cleverness
3. **Document everything** - All exported functions, types, and constants
4. **Test thoroughly** - Aim for >80% coverage
5. **Thread safety** - All types must be safe for concurrent use

### Formatting

All code must be formatted with `gofmt` and pass `golangci-lint`:

```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run

# Or using Task
task lint
```

### Naming Conventions

**Packages:**
- Use lowercase, single-word names when possible
- Match Cairo's conceptual organization (surface, context, pattern, etc.)

**Types:**
- Use PascalCase for exported types: `ImageSurface`, `Context`, `Matrix`
- Use camelCase for private types: `baseSurface`, `contextPtr`

**Functions:**
- Use PascalCase for exported functions: `NewImageSurface`, `SetLineWidth`
- Use camelCase for private functions: `surfaceCreate`, `contextSetSource`
- Constructors should be named `New*`: `NewContext`, `NewMatrix`

**Constants:**
- Use PascalCase with descriptive names: `FormatARGB32`, `StatusSuccess`
- Group related constants in const blocks

**Variables:**
- Use camelCase: `currentPoint`, `lineWidth`
- Keep names concise but descriptive
- Avoid single-letter names except for:
  - Loop indices: `i`, `j`
  - Coordinates: `x`, `y`
  - Common abbreviations: `ctx`, `err`

### Documentation

**Package Documentation:**
- Every package must have a `doc.go` file
- Include overview, usage patterns, and examples
- Explain key concepts and gotchas

**Function Documentation:**
- All exported functions must have doc comments
- Start with the function name: "NewContext creates..."
- Explain parameters, return values, and side effects
- Include examples for complex functions
- Document thread safety considerations

**Example:**
```go
// SetLineWidth sets the current line width for the Context. The line width
// value specifies the diameter of the pen used for stroking paths, in user-space
// units.
//
// The line width affects all subsequent Stroke and StrokePreserve operations.
// The default line width is 2.0.
//
// Note: The line width is transformed by the current transformation matrix (CTM),
// so scaling transformations will affect the actual rendered line width.
func (c *Context) SetLineWidth(width float64) {
    // implementation
}
```

### CGO Patterns

When adding CGO bindings, follow these patterns:

**File Organization:**
- Pure Go code in `package.go`
- CGO code in `package_cgo.go`
- Tests in `package_test.go`
- Benchmarks in `package_bench_test.go` (if needed)

**CGO Function Naming:**
- Use lowercase camelCase for CGO wrapper functions
- Prefix with package concept: `surfaceCreate`, `contextSetSource`
- Keep CGO surface area minimal

**Example:**
```go
// In surface_cgo.go
/*
#cgo pkg-config: cairo
#include <cairo.h>
*/
import "C"

func surfaceCreate(format Format, width, height int) (SurfacePtr, error) {
    ptr := C.cairo_image_surface_create(
        C.cairo_format_t(format),
        C.int(width),
        C.int(height),
    )

    st := C.cairo_surface_status(ptr)
    if st != C.CAIRO_STATUS_SUCCESS {
        C.cairo_surface_destroy(ptr)
        return nil, statusFromC(st)
    }

    return SurfacePtr(ptr), nil
}
```

### Thread Safety

All types must be thread-safe:

1. **Embed sync.RWMutex** in all types that wrap Cairo objects
2. **Use RLock/RUnlock** for read operations
3. **Use Lock/Unlock** for write operations
4. **Check closed flag** before accessing C pointers
5. **Document concurrency behavior** in function comments

**Example:**
```go
type Context struct {
    sync.RWMutex
    ptr    ContextPtr
    closed bool
}

func (c *Context) GetLineWidth() float64 {
    c.RLock()
    defer c.RUnlock()

    if c.ptr == nil {
        return 0.0
    }

    return contextGetLineWidth(c.ptr)
}
```

### Error Handling

1. **Return errors explicitly** - Use `(result, error)` pattern for constructors
2. **Check status** - Always check Cairo status after C calls
3. **Provide context** - Wrap errors with meaningful messages
4. **No panics** - Never panic in library code (tests are OK)

### Testing Requirements

Every feature must include:

1. **Unit tests** - Test all code paths
2. **Error tests** - Test error conditions
3. **Thread safety tests** - Test concurrent access
4. **Integration tests** - Test with other packages
5. **Examples** - At least one working example

**Test naming:**
- `TestFunctionName` - Basic functionality test
- `TestFunctionNameError` - Error condition test
- `TestFunctionNameThreadSafety` - Concurrency test

## How to Add New Features

### Following the Development Plan

This project follows a structured development plan outlined in
`development/go_cairo_prompts.md`. When adding features:

1. **Check the prompts** - See if your feature is already planned
2. **Follow the sequence** - Implement prompts in order when possible
3. **Maintain integration** - Ensure all code integrates immediately
4. **No orphaned code** - Every function should be used or tested

### Test-Driven Development

Follow TDD principles:

1. **Write tests first** (or alongside implementation)
2. **Make tests fail** - Verify tests actually test something
3. **Implement the feature** - Make tests pass
4. **Refactor** - Clean up while keeping tests green
5. **Document** - Add comprehensive documentation

### Adding a New Function

**Step 1: Define the Go API**

```go
// In context/context.go

// Arc adds a circular arc to the current path. The arc is centered at (xc, yc)
// with the given radius, and extends from angle1 to angle2 in radians.
// Angles are measured from the positive X axis, with positive angles extending
// in the direction from the positive X axis toward the positive Y axis.
func (c *Context) Arc(xc, yc, radius, angle1, angle2 float64) {
    c.Lock()
    defer c.Unlock()

    if c.ptr == nil {
        return
    }
    contextArc(c.ptr, xc, yc, radius, angle1, angle2)
}
```

**Step 2: Add CGO wrapper**

```go
// In context/context_cgo.go

func contextArc(ptr ContextPtr, xc, yc, radius, angle1, angle2 float64) {
    C.cairo_arc(
        (*C.cairo_t)(ptr),
        C.double(xc),
        C.double(yc),
        C.double(radius),
        C.double(angle1),
        C.double(angle2),
    )
}
```

**Step 3: Write tests**

```go
// In context/context_test.go

func TestContextArc(t *testing.T) {
    surf, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
    if err != nil {
        t.Fatalf("Failed to create surface: %v", err)
    }
    defer surf.Close()

    ctx, err := NewContext(surf)
    if err != nil {
        t.Fatalf("Failed to create context: %v", err)
    }
    defer ctx.Close()

    // Test arc operation
    ctx.Arc(50, 50, 30, 0, 2*math.Pi)

    // Verify no error occurred
    if status := ctx.Status(); status != status.Success {
        t.Errorf("Arc operation failed with status: %v", status)
    }

    // Verify current point is set
    if !ctx.HasCurrentPoint() {
        t.Error("Arc should set current point")
    }
}
```

**Step 4: Add example**

```go
// In examples/example_test.go or examples/circles.go

func ExampleContext_Arc() {
    surf, _ := surface.NewImageSurface(surface.FormatARGB32, 200, 200)
    defer surf.Close()

    ctx, _ := context.NewContext(surf)
    defer ctx.Close()

    // Draw a circle
    ctx.SetSourceRGB(1.0, 0.0, 0.0)
    ctx.Arc(100, 100, 50, 0, 2*math.Pi)
    ctx.Fill()

    surf.Flush()
    surf.WriteToPNG("circle.png")
}
```

**Step 5: Update documentation**

Add the function to package documentation if it represents a significant feature.

### Adding a New Package

When adding a completely new package:

1. **Create directory structure**
2. **Add doc.go** with package documentation
3. **Implement core types** with thread safety
4. **Add CGO bindings** in separate file
5. **Write comprehensive tests**
6. **Add benchmarks** if performance-critical
7. **Create examples**
8. **Update main package** to re-export if needed

## Pull Request Process

### Before Submitting

1. **Run all tests** - `go test -race ./...`
2. **Run linter** - `golangci-lint run`
3. **Check coverage** - Aim for >80%
4. **Update documentation** - Add/update relevant docs
5. **Add examples** - Include working examples
6. **Run examples** - Verify they actually work

### PR Guidelines

1. **One feature per PR** - Keep PRs focused and reviewable
2. **Write clear descriptions** - Explain what and why
3. **Reference issues** - Link to related issues if applicable
4. **Keep commits clean** - Squash unnecessary commits
5. **Be responsive** - Address review feedback promptly

### PR Template

```markdown
## Description
Brief description of the changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Checklist
- [ ] Tests added/updated
- [ ] Documentation added/updated
- [ ] All tests pass
- [ ] Linter passes
- [ ] Examples work
- [ ] Follows code style guidelines

## Related Issues
Closes #XXX
```

## Reporting Issues

### Bug Reports

When reporting bugs, include:

1. **Go version** - `go version` output
2. **Cairo version** - `pkg-config --modversion cairo`
3. **Operating system** - OS and version
4. **Minimal reproduction** - Smallest code that reproduces the issue
5. **Expected vs actual** - What you expected and what happened
6. **Error messages** - Complete error output

### Feature Requests

When requesting features:

1. **Check existing issues** - Avoid duplicates
2. **Describe the use case** - Why do you need this?
3. **Propose an API** - What should the function look like?
4. **Check Cairo docs** - Ensure Cairo supports it
5. **Consider alternatives** - Any workarounds?

## Code of Conduct

This project follows the Go Community Code of Conduct. Please be respectful,
inclusive, and professional in all interactions.

## Questions?

If you have questions about contributing:

1. Check existing documentation
2. Search closed issues
3. Ask in a new issue with the "question" label

Thank you for contributing to Go-Cairo!
