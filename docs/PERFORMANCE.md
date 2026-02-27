# Performance Guide

This guide covers how to get the best performance out of go-cairo and explains
the underlying cost model.

## Cost Model

go-cairo calls into the Cairo C library via CGO. Every call that crosses the
CGO boundary pins the calling goroutine to its OS thread while the C code runs.
This pinning adds roughly **50–100 ns per call** regardless of how fast the C
function itself is.

Implications:
- A drawing loop that calls `ctx.LineTo` a million times pays ~50–100 ms in
  CGO overhead alone, independent of Cairo's work.
- Batch as much as possible into a single path before calling `Fill` or `Stroke`.
- Prefer operations that do more work per call (e.g., draw a full sub-path, then
  fill it once).

## Pixel Data

`ImageSurface` pixel data lives in C-allocated memory. The Go side holds only a
small header struct. Reading or writing pixels with `GetData()` does not copy the
buffer — it returns a Go slice backed by the C allocation. Avoid retaining this
slice beyond the surface's lifetime.

## Practical Tips

### 1. Batch paths before filling

```go
// Slow: one CGO fill call per segment
for _, seg := range segments {
    ctx.MoveTo(seg.X1, seg.Y1)
    ctx.LineTo(seg.X2, seg.Y2)
    ctx.Stroke()   // CGO + Cairo work per iteration
}

// Fast: accumulate entire path, stroke once
for _, seg := range segments {
    ctx.MoveTo(seg.X1, seg.Y1)
    ctx.LineTo(seg.X2, seg.Y2)
}
ctx.Stroke()   // one CGO + Cairo call for all segments
```

### 2. Use Save/Restore for style changes

`Save` and `Restore` copy only Cairo's lightweight state struct, not pixel data.
Use them freely to bracket temporary style changes:

```go
ctx.Save()
ctx.SetLineWidth(5)
ctx.SetSourceRGBA(1, 0, 0, 0.5)
// ... draw temporary shapes ...
ctx.Restore()   // line width and color are restored
```

### 3. Reuse patterns

`Pattern` objects (gradients, surface patterns) are cheap once created. Create
them once and reuse across frames or drawing calls rather than re-creating them
each time:

```go
grad, _ := cairo.NewLinearGradient(0, 0, 400, 0)
grad.AddColorStopRGB(0, 1, 0, 0)
grad.AddColorStopRGB(1, 0, 0, 1)
defer grad.Close()

for i := 0; i < numFrames; i++ {
    ctx.SetSource(grad)   // reuse — no CGO allocation
    ctx.Rectangle(0, 0, 400, 400)
    ctx.Fill()
}
```

### 4. One context per goroutine

Each `Context` carries a `sync.RWMutex`. Sharing one context across goroutines
serialises every call. For parallel rendering, give each goroutine its own surface
and context, then composite the results:

```go
var wg sync.WaitGroup
results := make([]*cairo.ImageSurface, numWorkers)

for i := range results {
    wg.Add(1)
    go func(idx int) {
        defer wg.Done()
        surf, _ := cairo.NewImageSurface(cairo.FormatARGB32, tileW, tileH)
        ctx, _ := cairo.NewContext(surf)
        // ... draw tile idx ...
        ctx.Close()
        results[idx] = surf
    }(i)
}
wg.Wait()
// composite results into final surface
```

### 5. Prefer vector surfaces for scalable output

PDF and SVG surfaces record drawing commands lazily. They use negligible memory
during drawing and produce resolution-independent output. Use them when generating
documents or graphics that will be printed or scaled.

```go
// ~0 MB of pixel memory during drawing
surf, _ := cairo.NewPDFSurface("report.pdf", 612, 792)
```

### 6. Flush before pixel access

`surf.Flush()` ensures all pending drawing commands are applied to the pixel
buffer before you read raw pixel data. Forgetting this can return stale pixels.

```go
surf.Flush()
pixels := surf.GetData()   // safe to read after Flush
```

### 7. Clip to reduce work

Restricting the clip region before drawing complex paths lets Cairo skip
rendering outside the clip area. This is especially effective for partial
screen redraws:

```go
ctx.Rectangle(dirtyX, dirtyY, dirtyW, dirtyH)
ctx.Clip()
// ... redraw only the dirty area ...
ctx.ResetClip()
```

## Running Benchmarks

The library includes benchmarks for core operations:

```bash
go test -bench=. -benchmem ./...
```

Representative timings on a 2021 M1 MacBook Pro (your numbers will vary):

| Operation | Time | Allocations |
|-----------|------|-------------|
| `NewImageSurface` (400×400) | ~5 µs | 1 alloc |
| `NewContext` | ~1 µs | 1 alloc |
| `Fill` (simple rectangle) | ~1 µs | 0 allocs |
| `WriteToPNG` (400×400) | ~2 ms | varies |

`Fill` on a simple path shows 0 Go-heap allocations; all work happens in C memory.

## Profiling CGO Code

Standard Go profiling tools (`pprof`, `go test -cpuprofile`) attribute CGO time
to the calling goroutine. Use the `-cgo` tag to enable CGO-aware profiling in
tools like `gperftools` if you need sub-microsecond resolution inside Cairo.

For Go-level profiling:

```bash
go test -bench=BenchmarkContextFill -cpuprofile=cpu.out ./context
go tool pprof cpu.out
```

Look for `runtime.cgocall` in the profile to quantify CGO boundary overhead
versus Cairo's internal work.
