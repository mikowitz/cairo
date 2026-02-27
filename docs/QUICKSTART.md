# Quick Start Guide

Get up and running with go-cairo in five minutes.

## Prerequisites

- Go 1.23 or later
- Cairo 1.18 or later
- pkg-config
- A C compiler (gcc or clang)

On macOS:

```bash
brew install cairo pkg-config
```

On Ubuntu/Debian:

```bash
sudo apt-get install libcairo2-dev pkg-config build-essential
```

## Installation

```bash
go get github.com/mikowitz/cairo
```

## Your First Image

The following program draws a red circle on a white background and saves it to a PNG file.

```go
package main

import (
    "log"
    "math"

    "github.com/mikowitz/cairo"
)

func main() {
    // 1. Create a drawing surface (400x400 pixels)
    surf, err := cairo.NewImageSurface(cairo.FormatARGB32, 400, 400)
    if err != nil {
        log.Fatal(err)
    }
    defer surf.Close()

    // 2. Create a drawing context
    ctx, err := cairo.NewContext(surf)
    if err != nil {
        log.Fatal(err)
    }
    defer ctx.Close()

    // 3. Fill the background white
    ctx.SetSourceRGB(1, 1, 1)
    ctx.Paint()

    // 4. Draw a red circle
    ctx.SetSourceRGB(1, 0, 0)
    ctx.Arc(200, 200, 100, 0, 2*math.Pi)
    ctx.Fill()

    // 5. Save to PNG
    surf.Flush()
    if err := surf.WriteToPNG("circle.png"); err != nil {
        log.Fatal(err)
    }
}
```

Run it:

```bash
go run main.go
```

## Core Concepts

### Surfaces

A **surface** is the drawing target. The most common type is `ImageSurface`, which
holds pixels in memory and can be exported to PNG.

```go
surf, err := cairo.NewImageSurface(cairo.FormatARGB32, width, height)
```

Other surface types:
- `cairo.NewPDFSurface(filename, widthPt, heightPt)` — vector PDF output
- `cairo.NewSVGSurface(filename, widthPt, heightPt)` — vector SVG output

### Contexts

A **context** (`cairo.Context`) is the drawing state machine. You draw by calling
methods on the context. Always create a context from an existing surface:

```go
ctx, err := cairo.NewContext(surf)
```

### Colors

Set the drawing color before each operation:

```go
ctx.SetSourceRGB(r, g, b)        // Opaque color (values 0.0–1.0)
ctx.SetSourceRGBA(r, g, b, a)    // Color with alpha
```

### The Path–Render Loop

Cairo uses a two-phase drawing model:

1. **Build a path** using `MoveTo`, `LineTo`, `Arc`, `Rectangle`, `CurveTo`
2. **Render the path** with `Fill()`, `Stroke()`, or both

```go
ctx.Rectangle(50, 50, 200, 100)   // Define shape
ctx.SetSourceRGB(0, 0.5, 1)
ctx.FillPreserve()                 // Fill and keep path
ctx.SetSourceRGB(0, 0, 0)
ctx.SetLineWidth(3)
ctx.Stroke()                       // Outline the same shape
```

`Fill()` and `Stroke()` consume the path. Use `FillPreserve()` / `StrokePreserve()`
to reuse the path for a second operation.

### State Stack

Save and restore the graphics state with `Save()` and `Restore()`:

```go
ctx.SetSourceRGB(1, 0, 0)
ctx.Save()
    ctx.SetSourceRGB(0, 1, 0)  // Temporary change
    ctx.Rectangle(10, 10, 50, 50)
    ctx.Fill()
ctx.Restore()
// Color is red again here
```

### Gradients

```go
grad, err := cairo.NewLinearGradient(x0, y0, x1, y1)
grad.AddColorStopRGB(0.0, 1, 0, 0)   // Red at start
grad.AddColorStopRGB(1.0, 0, 0, 1)   // Blue at end
defer grad.Close()

ctx.SetSource(grad)
ctx.Rectangle(0, 0, 400, 400)
ctx.Fill()
```

### Text

```go
ctx.SelectFontFace("sans-serif", cairo.SlantNormal, cairo.WeightBold)
ctx.SetFontSize(24)
ctx.MoveTo(50, 200)
ctx.ShowText("Hello, Cairo!")
```

## Resource Management

Always close resources when done. `defer` ensures cleanup even on error paths:

```go
surf, err := cairo.NewImageSurface(cairo.FormatARGB32, 400, 400)
if err != nil { ... }
defer surf.Close()

ctx, err := cairo.NewContext(surf)
if err != nil { ... }
defer ctx.Close()
```

Finalizers are registered as a safety net, but explicit `Close()` gives deterministic
release of the underlying C memory.

## Next Steps

- **Examples**: See the `examples/` directory for complete, runnable examples
  covering gradients, patterns, compositing, transformations, text, animation, and more.
- **API Reference**: Run `go doc github.com/mikowitz/cairo` or visit pkg.go.dev.
- **Architecture**: See `ARCHITECTURE.md` for package organization and design decisions.
- **Migration from C**: See `docs/MIGRATION.md` if you are familiar with the C Cairo API.
- **Performance**: See `docs/PERFORMANCE.md` for tips on getting the most out of go-cairo.
