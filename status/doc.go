/*
Package status is used by Cairo to represent all kinds of errors.
A status value of `status.Success" represents no error and has an
integer value of zero. All other status values represent an error.

Cairo's error handling is designed to be easy to use and safe.
All major cairo objects retain an error status internally which
can be queried anytime by the users using "*.Status()" calls.
In the mean time, it is safe to call all cairo functions normally
even if the underlying object is in an error status.
This means that no error handling code is required before or after
each individual cairo function call.

# Error Handling Guide

Error handling in this Go wrapper for Cairo follows both Cairo's conventions
and Go's idiomatic error handling patterns. Understanding both is key to
writing robust code.

# Cairo's Error Model

Cairo objects maintain internal error state:

  1. Once an error occurs on an object, that object enters an error state
  2. Subsequent operations on that object become no-ops
  3. The error persists until the object is closed
  4. Calling methods on an errored object is safe (won't crash)

This "sticky error" model means you can chain operations without checking
errors after each call:

  ctx.SetSourceRGB(1, 0, 0)
  ctx.Rectangle(10, 10, 50, 50)
  ctx.Fill()
  if ctx.Status() != status.Success {
      // Handle error - one of the above operations failed
  }

However, the Go wrapper provides a more idiomatic approach by returning
errors directly from constructors and certain operations.

# Go Error Handling Patterns

Constructor functions return (object, error):

  surf, err := surface.NewImageSurface(surface.FormatARGB32, -100, 100)
  if err != nil {
      // Handle invalid dimensions - surf is nil
      return err
  }
  defer surf.Close()

This is the preferred pattern because it catches errors immediately at
construction time, preventing work with invalid objects.

# Two-Level Error Checking

For maximum safety, you can check errors at two levels:

  1. Construction: Check the returned error
  2. Operation: Periodically check Status()

Example:

  surf, err := surface.NewImageSurface(format, width, height)
  if err != nil {
      return fmt.Errorf("creating surface: %w", err)
  }
  defer surf.Close()

  ctx, err := context.NewContext(surf)
  if err != nil {
      return fmt.Errorf("creating context: %w", err)
  }
  defer ctx.Close()

  // Perform many drawing operations
  ctx.SetSourceRGB(1, 0, 0)
  ctx.Rectangle(10, 10, 100, 100)
  ctx.Fill()
  // ... many more operations ...

  // Check status before saving
  if status := ctx.Status(); status != status.Success {
      return fmt.Errorf("drawing error: %v", status)
  }

  surf.Flush()
  if err := surf.WriteToPNG("output.png"); err != nil {
      return fmt.Errorf("writing PNG: %w", err)
  }

# Common Error Scenarios

NullPointer:
  - Calling methods on a closed object
  - Passing nil to a constructor
  - Solution: Check for nil, don't use closed objects

NoMemory:
  - System out of memory
  - Surface/context allocation failed
  - Solution: Reduce dimensions, close unused resources

InvalidRestore:
  - Calling Restore() without matching Save()
  - Solution: Balance Save/Restore calls

NoCurrentPoint:
  - Calling GetCurrentPoint() when no point is defined
  - Calling LineTo() without prior MoveTo()
  - Solution: Use HasCurrentPoint() or call MoveTo() first

InvalidMatrix:
  - Matrix operation resulted in invalid matrix
  - Often from Invert() on singular matrix
  - Solution: Check matrix determinant before inverting

FileNotFound / WriteError:
  - Invalid file path for WriteToPNG()
  - Insufficient permissions
  - Solution: Validate paths, check permissions

# Error Recovery

Once an object is in an error state, it cannot be recovered:

  ctx.Status() == status.InvalidRestore  // Error occurred
  ctx.Save()                              // This is now a no-op
  ctx.Status() == status.InvalidRestore  // Error persists

The only recovery is to close the object and create a new one.

# Best Practices

  1. Always check constructor errors
  2. Use defer for Close() immediately after successful construction
  3. Check Status() before expensive operations or before saving results
  4. Wrap errors with context using fmt.Errorf("context: %w", err)
  5. Use Status.ToError() to convert status to standard error if needed
*/
package status
