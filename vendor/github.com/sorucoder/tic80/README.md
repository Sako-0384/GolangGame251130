# tic80

This is the Go package for the WASM binding for TIC-80.

## General Usage

This package follows the native [TIC-80 API](https://github.com/nesbox/TIC-80/wiki/API) as closely as possible, including optional arguments.
For functions that have optional arguments, you can either use the defaults by passing `nil`, like so:

```go
tic80.Print("HELLO WORLD FROM GO!", 65, 84, nil)
```

Or, you can pass an instance of the corresponding `tic80.<APIName>Options`, chaining its methods to configure it, like so:

```go
tic80.Spr(1+t%60/30*2, x, y, tic80.NewSpriteOptions().AddTransparentColor(14).SetScale(3).SetSize(2, 2))
```