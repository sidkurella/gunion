package aliasedimport

import (
	ctx "context"
	ioalias "io"
)

// myUnion tests how aliased imports are handled.
type myUnion struct {
	a ctx.Context    // context aliased as ctx
	b ioalias.Writer // io aliased as ioalias
	c ioalias.Reader // same alias, different type
}
