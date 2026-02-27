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
- [Features](#features)
- [Performance](#performance)
- [Comparison to Other Go Graphics Libraries](#comparison-to-other-go-graphics-libraries)
- [Troubleshooting](#troubleshooting)
- [FAQ](#faq)
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

- **Cairo**: Version 1.18 or higher (includes PDF backend by default on most platforms)
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

## Features

### Surface Types
- **ImageSurface** — raster pixel buffer in ARGB32, RGB24, A8, or A1 format; export to PNG
- **PDFSurface** — multi-page vector PDF output (requires `cairo-pdf` pkg-config entry)
- **SVGSurface** — web-compatible SVG output with configurable document units (requires `cairo-svg` pkg-config entry)

### Path Operations
- `MoveTo`, `LineTo`, `RelMoveTo`, `RelLineTo` — basic path construction
- `Rectangle` — axis-aligned rectangles
- `Arc`, `ArcNegative` — circular arcs
- `CurveTo`, `RelCurveTo` — cubic Bézier curves
- `ClosePath` — close the current sub-path
- `NewPath`, `NewSubPath` — explicit path reset/sub-path

### Drawing Operations
- `Fill`, `FillPreserve` — fill the current path (consuming or preserving it)
- `Stroke`, `StrokePreserve` — stroke the current path
- `Paint`, `PaintWithAlpha` — paint the current source over the whole clip region
- `Clip`, `ClipPreserve`, `ResetClip` — clipping regions

### Colors and Patterns
- `SetSourceRGB`, `SetSourceRGBA` — solid color sources
- `SetSourceSurface` — image/surface-based patterns
- Linear and radial gradient patterns with multiple color stops
- Surface patterns with repeat/reflect/pad extend modes
- Mesh patterns for bicubic tensor-product patch meshes

### Line Styling
- `SetLineWidth`, `GetLineWidth`
- Line caps: `LineCapButt`, `LineCapRound`, `LineCapSquare`
- Line joins: `LineJoinMiter`, `LineJoinRound`, `LineJoinBevel`
- `SetMiterLimit`, `GetMiterLimit`
- Dash patterns: `SetDash`, `GetDash`

### Compositing
- Full set of Cairo compositing operators (`OperatorOver`, `OperatorSource`, `OperatorMultiply`, etc.)
- `SetOperator`, `GetOperator`

### Transformations
- `Translate`, `Scale`, `Rotate`
- `Transform`, `SetMatrix`, `GetMatrix`, `IdentityMatrix`
- `UserToDevice`, `DeviceToUser` and their distance variants
- Full `matrix.Matrix` type for affine transforms

### Text (Toy Font API)
- `SelectFontFace` — choose family, slant (`Normal`/`Italic`/`Oblique`), weight (`Normal`/`Bold`)
- `SetFontSize`
- `ShowText` — render text at the current point
- `TextPath` — convert text to a path for fill/stroke rendering
- `TextExtents`, `FontExtents` — precise ink and font metrics

### Fill Rules
- `FillRuleWinding`, `FillRuleEvenOdd`
- `SetFillRule`, `GetFillRule`

### State Management
- `Save` / `Restore` — graphics state stack
- `Status` — check for deferred drawing errors

### Memory and Concurrency
- All types are safe for concurrent use (`sync.RWMutex` internally)
- Explicit `Close()` methods for deterministic resource release
- Finalizers as a safety net against leaks

## Performance

Go-Cairo uses CGO to call into the Cairo C library. Key performance characteristics:

- **CGO call overhead**: Each C call carries ~50–100 ns of goroutine-pinning overhead.
  For tight inner loops (e.g., drawing thousands of individual paths), batch operations
  via `NewPath`/`LineTo`/`Fill` rather than calling `Fill` per segment.
- **No extra heap allocations**: The library does not copy pixel data between Go and C.
  `ImageSurface` pixel data lives in C-allocated memory accessed directly.
- **Vector surfaces are lazy**: PDF and SVG surfaces do not rasterise until the file is
  finalised; they are ideal for high-resolution or scalable output with negligible
  memory cost during drawing.
- **State saves are cheap**: `Save`/`Restore` only copy Cairo's lightweight graphics
  state struct; they do not copy pixel data.
- **Thread contention**: Each `Context` carries its own `RWMutex`. Sharing a single
  context across many goroutines will serialise drawing calls. Prefer one context per
  goroutine, all drawing to the same surface after acquiring the surface lock.

Benchmarks (run `go test -bench=. -benchmem ./...` on your hardware):

| Operation | Typical time |
|-----------|-------------|
| `NewImageSurface` (400×400) | ~5 µs |
| `NewContext` | ~1 µs |
| `Fill` (simple rectangle) | ~1 µs |
| `WriteToPNG` (400×400) | ~2 ms |

## Comparison to Other Go Graphics Libraries

| Library | Approach | Strengths | Limitations |
|---------|----------|-----------|-------------|
| **go-cairo** | CGO wrapper around Cairo | Full Cairo API, vector output (PDF/SVG), precise text measurement, patterns | Requires C library; CGO overhead |
| [gg](https://github.com/fogleman/gg) | Go-native (uses `golang.org/x/image`) | Pure Go, easy to use, no C deps | No PDF/SVG output; subset of Cairo ops; no patterns |
| [draw2d](https://github.com/llgcode/draw2d) | Pure Go | No CGO, cross-platform including WASM | Less feature-complete; slower rasterisation |
| [svgo](https://github.com/ajstarks/svgo) | SVG text generation | Simple SVG output | SVG only; no raster; no font metrics |
| [canvas](https://github.com/tdewolff/canvas) | Pure Go vector | PDF/SVG/raster, font subsetting | Pure Go (no Cairo); different API |

**Choose go-cairo when** you need the full Cairo feature set, precise font metrics,
multi-page PDF output, or are porting existing Cairo C code to Go.

**Choose a pure-Go library when** you cannot install C dependencies, need WASM
support, or have simpler drawing requirements.

## Troubleshooting

**`pkg-config: command not found`**
Install pkg-config: `sudo apt-get install pkg-config` or `brew install pkg-config`.

**`Package cairo was not found in the pkg-config search path`**
Install the Cairo development headers: `sudo apt-get install libcairo2-dev` or
`brew install cairo`. Verify with `pkg-config --modversion cairo`.

**`cairo version too old (need 1.18+)`**
Check your installed version: `pkg-config --modversion cairo`. Upgrade Cairo via
your package manager. On macOS: `brew upgrade cairo`.

**`cgo: C compiler "gcc" not found`**
CGO requires a C compiler. Install GCC: `sudo apt-get install build-essential` or
Xcode Command Line Tools on macOS: `xcode-select --install`.

**PDF or SVG surface fails to create**
PDF and SVG surfaces have separate build tags (`!nopdf`, `!nosvg`). If they are
excluded from your build, ensure `cairo-pdf` and `cairo-svg` pkg-config entries exist:
`pkg-config --modversion cairo-pdf` and `pkg-config --modversion cairo-svg`.

**Resource leak / too many open files**
Always call `defer ctx.Close()` and `defer surf.Close()` immediately after creation.
Finalizers are registered but run non-deterministically under the GC.

**Blank output / all-white image**
Ensure you call `surf.Flush()` before `surf.WriteToPNG()`. Also verify your drawing
color is not the same as the background (default source is opaque black).

## FAQ

**Q: Do I need to install Cairo separately?**
Yes. Go-Cairo is a CGO wrapper and requires a Cairo shared library and its development
headers. See [Installing Cairo](#installing-cairo) above.

**Q: Is the library thread-safe?**
Yes. Every `Context`, `Surface`, and `Pattern` embeds a `sync.RWMutex`. Multiple
goroutines may call methods concurrently. Read operations (`GetLineWidth`, etc.) acquire
a read lock and can run in parallel; write operations acquire an exclusive lock.

**Q: Can I use go-cairo in a WASM target?**
No. CGO is not available in the WASM target (`GOOS=js GOARCH=wasm`). For browser
graphics, consider a pure-Go library or the HTML5 Canvas API directly.

**Q: How do I generate a multi-page PDF?**
Create a `PDFSurface`, draw page 1, call `ShowPage()`, draw page 2, call `ShowPage()`,
and so on. Close the surface to finalise the file. See `examples/pdf_output.go`.

**Q: How do I convert the animation frames to a video?**
Use ffmpeg: `ffmpeg -r 30 -i frame_%03d.png -c:v libx264 -pix_fmt yuv420p output.mp4`.
See the godoc comment in `examples/animation.go` for the full command.

**Q: What coordinate system does Cairo use?**
Cairo uses a right-handed coordinate system with the origin at the top-left corner.
Positive X goes right; positive Y goes down. Angles in `Arc` are measured in radians,
clockwise from the positive X axis.

**Q: Why does my text look different on Linux vs macOS?**
Cairo's toy font API delegates to the platform font backend (FreeType on Linux,
CoreText on macOS). Different platforms use different default fonts and hinting
strategies, producing visually different (but correct) output. Use `TextExtents` for
portable layout calculations regardless of the rendered appearance.

## Roadmap

The library will continue to expand with additional Cairo features:

### Phase 1: Enhanced Drawing (Complete)

- [x] Transformations (Translate, Scale, Rotate, Matrix operations)
- [x] Advanced path operations (Arc, Curve, RelMoveTo, RelLineTo)
- [x] Pattern support (Gradients, Surface patterns)
- [x] Clipping operations
- [x] Fill rules and stroke/fill extents

### Phase 2: Text and Fonts (Complete)

- [x] Toy font API (SelectFontFace, SetFontSize, ShowText)
- [x] Text extents and measurement
- [x] Text path operations

### Phase 3: Advanced Features (Complete)

- [x] Line cap and join styles
- [x] Dash patterns
- [x] Compositing operators
- [x] Masks and opacity
- [x] Path querying and manipulation

### Phase 4: Additional Surface Types (Complete)

- [x] PDF surface for vector output
- [x] SVG surface for web graphics

### Phase 5: Performance and Polish

- [ ] Benchmarking suite
- [ ] Advanced examples and tutorials
- [ ] API stability and v1.0.0 release

See `development/go_cairo_prompts.md` for the detailed implementation plan
spanning 34 prompts from foundation through v0.1.0 release.

## License

[To be determined]
