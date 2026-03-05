# gunion

Tagged union type generator for Go.

`gunion` reads a Go struct definition and generates a tagged union (discriminated union / sum type) with:

- A variant enum tracking which field is active
- Type-safe constructors for each variant
- `Is_`, `Unwrap_`, and `Get_` accessors per variant
- Exhaustive `Match` function with compile-time enforcement
- A `String()` method on the variant enum

## Installation

```sh
go install github.com/sidkurella/gunion@latest
```

Requires Go 1.24+.

## Quick start

Define a struct whose fields represent the variants of your union:

```go
// shape.go
package shape

type shape struct {
    circle    float64    // radius
    rectangle [2]float64 // width, height
    triangle  [3]float64 // sides
}
```

Run gunion:

```sh
gunion --type shape --src shape.go --out-pkg shape --out-type Shape
```

This generates `shape_gunion.go` containing a `Shape` union type with constructors, accessors, and a match function.

By default, the output type name is the capitalized input name suffixed with `Union` (e.g. `shape` → `ShapeUnion`). Use `--out-type` to choose a better name -- it appears in every generated symbol: the struct itself (`Shape`), constructors (`NewShape_circle`), match function (`Match_Shape`), etc.

### With `go:generate`

```go
//go:generate gunion --type shape --out-type Shape
type shape struct {
    circle    float64
    rectangle [2]float64
    triangle  [3]float64
}
```

When used with `go:generate`, `--src` and `--out-pkg` are automatically populated from the `GOFILE` and `GOPACKAGE` environment variables.

## Generated API

Given this input:

```go
type myUnion struct {
    a int
    b string
}
```

Running `gunion --type myUnion --src myunion.go --out-pkg mypkg --no-default` generates:

### Variant enum

```go
type _myUnionVariant int

const (
    _myUnionVariant_Invalid _myUnionVariant = 0
    _myUnionVariant_a       _myUnionVariant = 1
    _myUnionVariant_b       _myUnionVariant = 2
)
```

The `Invalid` variant represents the zero-value state (no variant has been set). It is included when `--no-default` is passed. Without it, the first field is the default and there is no `Invalid` variant.

### Union struct

```go
type MyUnionUnion struct {
    _variant _myUnionVariant
    _inner   myUnion
}
```

The original struct is embedded as `_inner`. The unexported fields prevent direct construction -- you must use the generated constructors.

### Constructors

```go
func NewMyUnionUnion_a(val int) MyUnionUnion
func NewMyUnionUnion_b(val string) MyUnionUnion
func NewMyUnionUnion_Invalid() MyUnionUnion
```

### `Is_` (variant check)

```go
func (u *MyUnionUnion) Is_a() bool
func (u *MyUnionUnion) Is_b() bool
func (u *MyUnionUnion) Is_Invalid() bool
```

### `Unwrap_` (panicking getter)

Returns the value if the variant matches, panics otherwise. Similar to Rust's `unwrap()`.

```go
func (u *MyUnionUnion) Unwrap_a() int
func (u *MyUnionUnion) Unwrap_b() string
```

### `Get_` (safe getter)

Returns `(value, true)` if the variant matches, `(zero, false)` otherwise.

```go
func (u *MyUnionUnion) Get_a() (int, bool)
func (u *MyUnionUnion) Get_b() (string, bool)
```

### `Match` (exhaustive pattern matching)

A generic function that requires a handler for every variant, enforced at compile time:

```go
func Match_MyUnionUnion[_R any](
    u *MyUnionUnion,
    on_a func(int) _R,
    on_b func(string) _R,
    on_Invalid func() _R,
) _R
```

Usage:

```go
result := Match_MyUnionUnion[string](u,
    func(a int) string { return fmt.Sprintf("got int: %d", a) },
    func(b string) string { return fmt.Sprintf("got string: %s", b) },
    func() string { return "invalid" },
)
```

### `String()` (variant name)

The variant enum type implements `fmt.Stringer`, returning the variant name (e.g. `"a"`, `"b"`, `"Invalid"`).

## Generics

gunion supports generic input structs. Given:

```go
type myUnion[T any, U comparable, V io.Writer] struct {
    a T
    b U
    c V
}
```

The generated union and all its methods are parameterized with the same type constraints:

```go
type MyUnionUnion[T any, U comparable, V io.Writer] struct { ... }

func NewMyUnionUnion_a[T any, U comparable, V io.Writer](val T) MyUnionUnion[T, U, V]

func Match_MyUnionUnion[T any, U comparable, V io.Writer, _R any](
    u *MyUnionUnion[T, U, V],
    on_a func(T) _R,
    on_b func(U) _R,
    on_c func(V) _R,
    on_Invalid func() _R,
) _R
```

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--type` | `-t` | (required) | Name of the source struct type |
| `--src` | | `$GOFILE` | Source file path. Falls back to `GOFILE` env var |
| `--out-type` | | `<Type>Union` | Name of the generated union type |
| `--out-file` | `-o` | `<src>_gunion.go` | Output file path |
| `--out-pkg` | | `$GOPACKAGE` | Output package name. Falls back to `GOPACKAGE` env var |
| `--no-getters` | | `false` | Omit `Unwrap_` and `Get_` methods |
| `--no-setters` | | `false` | Omit constructors (`New<OutType>_<Variant>`) |
| `--no-match` | | `false` | Omit the `Match` function |
| `--no-default` | | `false` | Insert an `Invalid` variant as the zero value. Without this flag, the first field is the default |

## License

MIT
