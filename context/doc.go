/*
Package context defines [Context], the main object used when drawing with
cairo. To draw with cairo, you create a [Context], set the target surface,
and drawing options for the [Context], create shapes with functions like
[MoveTo] and [LineTo], and then draw shapes with [Stroke] or [Fill].

[Context] can be pushed to a stack via [Context.Save]. They may then safely be changed,
without losing the current state. Use [Context.Restore] to restore to the saved state.
*/
package context
