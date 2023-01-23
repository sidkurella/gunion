package basic

import "io"

type myUnion[T any, U comparable, V io.Writer] struct {
	a T
	b U
	c V
}
