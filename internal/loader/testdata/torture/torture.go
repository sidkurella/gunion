package torture

import (
	"context"
	"io"
	"sync"
)

// myUnion is a torture test with many nested and edge-case types.
type myUnion struct {
	// Basic types
	a int
	b string

	// Pointer to basic
	c *float64

	// Slice of pointers
	d []*int

	// Array of fixed size
	e [5]byte

	// Map with complex key and value
	f map[string][]int

	// Nested map
	g map[int]map[string]bool

	// Channel types
	h chan int
	i <-chan string
	j chan<- bool

	// Function with no params or returns
	k func()

	// Function with multiple params and returns
	l func(a int, b string) (int, error)

	// Function with variadic params
	m func(format string, args ...any) string

	// Function with named return values (appears same as unnamed in go/types)
	n func(x, y int) (sum int, diff int)

	// Function returning a function
	o func(multiplier int) func(int) int

	// Function taking a function as param
	p func(callback func(int) bool) error

	// Function with channel params
	q func(in <-chan int, out chan<- int)

	// Function with context (common pattern)
	r func(ctx context.Context) error

	// Function with interface params
	s func(w io.Writer, r io.Reader) (int64, error)

	// Pointer to function
	t *func(int) int

	// Slice of functions
	u []func() error

	// Map with function values
	v map[string]func(int) int

	// Deeply nested pointer
	w ***int

	// Pointer to slice of pointers to arrays
	x *[]*[3]int

	// Function with pointer receiver style params
	y func(*int, **string) *bool

	// Empty interface (any)
	z any

	// Slice of any
	aa []any

	// Map of any to any
	bb map[any]any

	// Function returning multiple errors
	cc func() (int, int, error, error)

	// Function with only variadic param
	dd func(...string)

	// Complex nested function type
	ee func(func(func(int) int) func(int) int) func(func(int) int) func(int) int

	// Alias preservation tests - verify we use the alias name, not underlying type
	ff rune  // alias for int32
	gg int32 // NOT an alias, should stay int32
	hh byte  // alias for uint8 (also tested in field e)
	ii uint8 // NOT an alias, should stay uint8

	// Generic instantiation tests
	jj Generic[int]          // local generic with one type arg
	kk Generic[string]       // same generic, different type arg
	ll Generic[*float64]     // generic with pointer type arg
	mm TwoParam[int, string] // generic with two type args
	nn Generic[Generic[int]] // nested generic instantiation
	oo sync.Map              // external non-generic (for comparison)
	pp sync.Pool             // another external non-generic
}

// Generic is a generic type for testing instantiation.
type Generic[T any] struct {
	value T
}

// TwoParam is a generic type with two type parameters.
type TwoParam[K comparable, V any] struct {
	key   K
	value V
}
