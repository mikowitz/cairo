# Migration Guide: C Cairo to go-cairo

This guide is for developers already familiar with the C Cairo API who want to
use go-cairo. It covers the key differences between the C API and the Go wrapper.

## The Mental Model

In C, you work directly with opaque pointer types and free functions:

```c
cairo_surface_t *surface = cairo_image_surface_create(CAIRO_FORMAT_ARGB32, 400, 400);
cairo_t *cr = cairo_create(surface);
cairo_set_source_rgb(cr, 1.0, 0.0, 0.0);
cairo_rectangle(cr, 10, 10, 100, 100);
cairo_fill(cr);
cairo_destroy(cr);
cairo_surface_destroy(surface);
```

In go-cairo, the same operations are methods on typed structs:

```go
surface, err := cairo.NewImageSurface(cairo.FormatARGB32, 400, 400)
defer surface.Close()

ctx, err := cairo.NewContext(surface)
defer ctx.Close()

ctx.SetSourceRGB(1.0, 0.0, 0.0)
ctx.Rectangle(10, 10, 100, 100)
ctx.Fill()
```

## Key Differences

### Memory Management

| C Cairo | go-cairo |
|---------|----------|
| `cairo_reference()` / `cairo_destroy()` | `Close()` — no manual ref counting |
| `cairo_surface_destroy()` | `surface.Close()` |
| Must track reference counts manually | Use `defer resource.Close()` |

Finalizers are registered as a safety net, but explicit `Close()` is strongly
preferred for deterministic resource cleanup:

```go
surf, err := cairo.NewImageSurface(cairo.FormatARGB32, 400, 400)
if err != nil { return err }
defer surf.Close()   // Release the C memory when this function returns
```

### Error Handling

C Cairo uses a deferred-error model: functions rarely fail; you check status
after the fact with `cairo_status()`.

go-cairo uses a hybrid approach:
- **Constructors and I/O** return `(result, error)` — check immediately.
- **Drawing operations** (`Fill`, `Stroke`, etc.) set internal status, checked
  with `ctx.Status()`.

```go
// Constructor — check error immediately
surf, err := cairo.NewImageSurface(cairo.FormatARGB32, 400, 400)
if err != nil { ... }

// Drawing — check status after a batch of operations (optional, like C)
ctx.Arc(200, 200, 100, 0, 2*math.Pi)
ctx.Fill()
if st := ctx.Status(); st != status.Success {
    return fmt.Errorf("drawing failed: %v", st)
}

// I/O — check error immediately
if err := surf.WriteToPNG("out.png"); err != nil { ... }
```

### Function Name Mapping

C functions become methods on the relevant receiver. The `cairo_` prefix and
the type-name segment are dropped:

| C Function | go-cairo |
|------------|----------|
| `cairo_create(surface)` | `cairo.NewContext(surface)` |
| `cairo_destroy(cr)` | `ctx.Close()` |
| `cairo_set_source_rgb(cr, r, g, b)` | `ctx.SetSourceRGB(r, g, b)` |
| `cairo_set_source_rgba(cr, r, g, b, a)` | `ctx.SetSourceRGBA(r, g, b, a)` |
| `cairo_set_line_width(cr, w)` | `ctx.SetLineWidth(w)` |
| `cairo_move_to(cr, x, y)` | `ctx.MoveTo(x, y)` |
| `cairo_line_to(cr, x, y)` | `ctx.LineTo(x, y)` |
| `cairo_arc(cr, x, y, r, a1, a2)` | `ctx.Arc(x, y, r, a1, a2)` |
| `cairo_curve_to(cr, ...)` | `ctx.CurveTo(...)` |
| `cairo_rectangle(cr, x, y, w, h)` | `ctx.Rectangle(x, y, w, h)` |
| `cairo_close_path(cr)` | `ctx.ClosePath()` |
| `cairo_fill(cr)` | `ctx.Fill()` |
| `cairo_fill_preserve(cr)` | `ctx.FillPreserve()` |
| `cairo_stroke(cr)` | `ctx.Stroke()` |
| `cairo_stroke_preserve(cr)` | `ctx.StrokePreserve()` |
| `cairo_paint(cr)` | `ctx.Paint()` |
| `cairo_save(cr)` | `ctx.Save()` |
| `cairo_restore(cr)` | `ctx.Restore()` |
| `cairo_translate(cr, tx, ty)` | `ctx.Translate(tx, ty)` |
| `cairo_scale(cr, sx, sy)` | `ctx.Scale(sx, sy)` |
| `cairo_rotate(cr, angle)` | `ctx.Rotate(angle)` |

### Enum Names

C enums become Go constants. The long-form C prefix is shortened to a package-level
constant:

| C Constant | go-cairo |
|------------|----------|
| `CAIRO_FORMAT_ARGB32` | `cairo.FormatARGB32` |
| `CAIRO_FORMAT_RGB24` | `cairo.FormatRGB24` |
| `CAIRO_LINE_CAP_ROUND` | `cairo.LineCapRound` |
| `CAIRO_LINE_JOIN_BEVEL` | `cairo.LineJoinBevel` |
| `CAIRO_FILL_RULE_EVEN_ODD` | `cairo.FillRuleEvenOdd` |
| `CAIRO_OPERATOR_OVER` | `cairo.OperatorOver` |
| `CAIRO_FONT_SLANT_ITALIC` | `cairo.SlantItalic` |
| `CAIRO_FONT_WEIGHT_BOLD` | `cairo.WeightBold` |
| `CAIRO_EXTEND_REPEAT` | `cairo.ExtendRepeat` |

### Pattern / Gradient Creation

```c
// C
cairo_pattern_t *pat = cairo_pattern_create_linear(0, 0, 400, 0);
cairo_pattern_add_color_stop_rgb(pat, 0.0, 1, 0, 0);
cairo_pattern_add_color_stop_rgb(pat, 1.0, 0, 0, 1);
cairo_set_source(cr, pat);
cairo_rectangle(cr, 0, 0, 400, 400);
cairo_fill(cr);
cairo_pattern_destroy(pat);
```

```go
// go-cairo
grad, err := cairo.NewLinearGradient(0, 0, 400, 0)
if err != nil { return err }
defer grad.Close()
grad.AddColorStopRGB(0.0, 1, 0, 0)
grad.AddColorStopRGB(1.0, 0, 0, 1)
ctx.SetSource(grad)
ctx.Rectangle(0, 0, 400, 400)
ctx.Fill()
```

### PDF and SVG Surfaces

PDF and SVG backends require no special CGO flags; they are gated by build tags:

- PDF: available unless built with `-tags nopdf`
- SVG: available unless built with `-tags nosvg`

```go
// PDF (dimensions in points, 1pt = 1/72 inch)
pdfSurf, err := cairo.NewPDFSurface("output.pdf", 612, 792)
defer pdfSurf.Close()

ctx, _ := cairo.NewContext(pdfSurf)
// ... draw page 1 ...
pdfSurf.ShowPage()
// ... draw page 2 ...

// SVG (dimensions in points)
svgSurf, err := cairo.NewSVGSurface("output.svg", 600, 400)
defer svgSurf.Close()
```

### Text

The toy font API maps directly:

```c
cairo_select_font_face(cr, "sans-serif", CAIRO_FONT_SLANT_NORMAL, CAIRO_FONT_WEIGHT_BOLD);
cairo_set_font_size(cr, 24.0);
cairo_move_to(cr, 50, 200);
cairo_show_text(cr, "Hello");
```

```go
ctx.SelectFontFace("sans-serif", cairo.SlantNormal, cairo.WeightBold)
ctx.SetFontSize(24.0)
ctx.MoveTo(50, 200)
ctx.ShowText("Hello")
```

Text extents:

```go
extents := ctx.TextExtents("Hello")
fmt.Printf("width=%.1f height=%.1f\n", extents.Width, extents.Height)
```

## What Is Not Yet Wrapped

The following C Cairo features are not yet available in go-cairo:

- Scaled fonts (`cairo_scaled_font_t`) and glyph-level rendering
- User fonts (`cairo_user_font_face_t`)
- Recording surfaces
- Script surfaces
- Device API (`cairo_device_t`)
- Region API (`cairo_region_t`)

Contributions are welcome. See `CONTRIBUTING.md` for guidelines.
