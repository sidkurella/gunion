package externalimport

import (
	"context"

	"golang.org/x/tools/go/packages"
)

type myUnion struct {
	a int
	b *packages.Package
	c context.Context
}
