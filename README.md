# Go-Cairo

A Golang wrapper for the Cairo 2D graphics library, providing idiomatic
Go bindings to Cairo's powerful vector graphics capabilities.

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
    "math"
    "github.com/mikowitz/cairo"
)

func main() {
    // Create a 256x256 image surface
    surface, err := cairo.NewImageSurface(cairo.FormatARGB32, 256, 256)
    if err != nil {
        panic(err)
    }
    defer surface.Close()

    // Create a drawing context
    ctx := cairo.NewContext(surface)
    defer ctx.Close()

    // Draw a circle
    ctx.Arc(128, 128, 64, 0, 2*math.Pi)
    ctx.Fill()

    // Save to PNG
    if err := surface.WriteToPNG("circle.png"); err != nil {
        panic(err)
    }
}
```

## Development Status

This project is currently in early development. The initial focus is on
implementing core functionality:

- Basic surface creation (Image, PDF, SVG)
- Context operations (path construction, rendering)
- Solid colors and gradients
- Basic transformations

See `development/go-cairo-spec.md` for the complete implementation plan.

## License

[To be determined]
