# Go-Cairo

[![Go Version](https://img.shields.io/badge/Go-1.23%2B-blue.svg)](https://golang.org/doc/devel/release.html)
[![CI](https://github.com/mikowitz/go-cairo/actions/workflows/ci.yml/badge.svg)](https://github.com/mikowitz/go-cairo/actions/workflows/ci.yml)
<!-- [![Go Report Card](https://goreportcard.com/badge/github.com/mikowitz/cairo)](https://goreportcard.com/report/github.com/mikowitz/cairo) -->
[![License](https://img.shields.io/badge/license-TBD-lightgrey.svg)](#license)
<!-- [![GoDoc](https://pkg.go.dev/badge/github.com/mikowitz/cairo.svg)](https://pkg.go.dev/github.com/mikowitz/cairo) -->

A Golang wrapper for the Cairo 2D graphics library, providing idiomatic
Go bindings to Cairo's powerful vector graphics capabilities.

## Table of Contents

- [Purpose](#purpose)
- [Goals](#goals)
- [Target Audience](#target-audience)
- [Requirements](#requirements)
- [Installation](#installation)
- [Quick Example](#quick-example)
- [Current Status](#current-status)
- [Roadmap](#roadmap)
- [License](#license)

## Purpose

Go-Cairo enables Go developers to leverage Cairo's 2D graphics capabilities
without leaving the Go ecosystem. The library provides API parity with the
C Cairo library while offering a Go-native experience with proper error
handling, thread safety, and memory management.

## Goals

- **API Parity**: Maintain compatibility with the Cairo C library API
- **Go-Idiomatic**: Follow Go conventions where they don't conflict with
  Cairo's design
- **Performance-Conscious**: Minimize overhead over native C Cairo
- **Thread-Safe**: All operations are safe for concurrent use
- **Progressive Development**: Start minimal, iterate toward completeness

## Target Audience

General-purpose graphics developers who prefer Go's modern tooling and ecosystem
while needing Cairo's 2D graphics capabilities.

## Requirements

### System Dependencies

- **Cairo**: Version 1.18 or higher
- **pkg-config**: Required for build configuration

### Go Version

- **Go**: 1.23 or higher (targeting 1.23-1.25)

### Development Tools

- **golangci-lint**: Required for code linting during development

### Installing Cairo

**Ubuntu/Debian:**

```bash
sudo apt-get update
sudo apt-get install -y libcairo2-dev
```

**macOS:**

```bash
brew install cairo
```

## Installation

```bash
go get github.com/mikowitz/cairo
```

## Quick Example

```go
package main

import (
    "github.com/mikowitz/cairo/context"
    "github.com/mikowitz/cairo/surface"
)

func main() {
    // Create a 400x400 image surface
    surf, err := surface.NewImageSurface(surface.FormatARGB32, 400, 400)
    if err != nil {
        panic(err)
    }
    defer surf.Close()

    // Create a drawing context
    ctx, err := context.NewContext(surf)
    if err != nil {
        panic(err)
    }
    defer ctx.Close()

    // Set white background
    ctx.SetSourceRGB(1.0, 1.0, 1.0)
    ctx.Paint()

    // Draw a red filled rectangle
    ctx.SetSourceRGB(1.0, 0.0, 0.0)
    ctx.Rectangle(100, 100, 200, 150)
    ctx.Fill()

    // Draw a blue stroked rectangle
    ctx.SetSourceRGB(0.0, 0.0, 1.0)
    ctx.SetLineWidth(5.0)
    ctx.Rectangle(150, 150, 100, 100)
    ctx.Stroke()

    // Flush and save to PNG
    surf.Flush()
    if err := surf.WriteToPNG("rectangles.png"); err != nil {
        panic(err)
    }
}
```

## Current Status

**MVP Complete! ✓**

The core functionality is now fully implemented and tested. You can:

- ✅ **Create ImageSurfaces** with various pixel formats (ARGB32, RGB24, A8, etc.)
- ✅ **Draw with Context** using complete path operations
- ✅ **Render paths** with Fill, Stroke, and Paint operations
- ✅ **Set colors** using RGB/RGBA solid colors
- ✅ **Export to PNG** from ImageSurface
- ✅ **Basic path operations** including MoveTo, LineTo, Rectangle, ClosePath
- ✅ **Line styling** with SetLineWidth and GetLineWidth
- ✅ **State management** via Save/Restore
- ✅ **Thread-safe operations** with proper locking throughout
- ✅ **Memory management** with finalizers and explicit Close() methods
- ✅ **Comprehensive tests** with >80% coverage
- ✅ **Full documentation** with examples and guides

All core drawing operations work correctly and the library is ready for basic
2D graphics tasks. See the `examples/` directory for working code samples.

## Roadmap

The library will continue to expand with additional Cairo features:

### Phase 1: Enhanced Drawing (Next)

- [ ] Transformations (Translate, Scale, Rotate, Matrix operations)
- [ ] Advanced path operations (Arc, Curve, RelMoveTo, RelLineTo)
- [ ] Pattern support (Gradients, Surface patterns)
- [ ] Clipping operations
- [ ] Fill rules and stroke/fill extents

### Phase 2: Text and Fonts

- [ ] Toy font API (SelectFontFace, SetFontSize, ShowText)
- [ ] Text extents and measurement
- [ ] Text path operations
- [ ] Font options and rendering quality

### Phase 3: Advanced Features

- [ ] Line cap and join styles
- [ ] Dash patterns
- [ ] Compositing operators
- [ ] Masks and opacity
- [ ] Path querying and manipulation

### Phase 4: Additional Surface Types

- [ ] PDF surface for vector output
- [ ] SVG surface for web graphics
- [ ] Recording surface for operation capture
- [ ] Platform-specific surfaces (if needed)

### Phase 5: Performance and Polish

- [ ] Performance optimization pass
- [ ] Benchmarking suite
- [ ] Advanced examples and tutorials
- [ ] API stability and v1.0.0 release

See `development/go_cairo_prompts.md` for the detailed implementation plan
spanning 34 prompts from foundation through v0.1.0 release.

## License

[To be determined]
