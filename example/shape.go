// Package example demonstrates using gunion with go:generate.
//
// Run `go generate ./...` from the repository root to regenerate the union type.
package example

//go:generate go run .. --type shape --no-default

type shape struct {
	circle    float64
	rectangle [2]float64
	triangle  [3]float64
}
