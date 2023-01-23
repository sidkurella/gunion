package loader

import (
	"fmt"
	"go/types"

	"github.com/sidkurella/gunion/internal/config"
	"golang.org/x/tools/go/packages"
)

// Loads in the union details based on the input configuration.
type Loader struct {
	config config.InputConfig
}

// Represents the info we need to know about a type.
type Type struct {
	// Name of this type.
	Name string
	// Number of indirections (i.e., how many stars.)
	IndirectCount int
	// Originating package of this type.
	Source string
}

// Represents one of the possible variants of the union.
type Variant struct {
	// Name of this variant.
	Name string
	// Type of this variant.
	Type Type
}

// Represents our union details.
type Union struct {
	// List of possible variants of this union.
	Variants []Variant
}

func NewLoader(config config.InputConfig) *Loader {
	return &Loader{
		config: config,
	}
}

func (l *Loader) Load() (Union, error) {
	pkgs, err := packages.Load(&packages.Config{
		// TODO: May need NeedDeps.
		Mode: packages.NeedTypes | packages.NeedImports | packages.NeedSyntax | packages.NeedTypesInfo |
			packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedModule,
	}, "file="+l.config.Source)
	if err != nil {
		return Union{}, fmt.Errorf("failed to load packages for source file %s: %w", l.config.Source, err)
	}

	if len(pkgs) != 1 {
		return Union{}, fmt.Errorf("expected to load 1 package but got %d", len(pkgs))
	}

	pkg := pkgs[0]
	// fileAst := pkg.Syntax[0]

	obj := pkg.Types.Scope().Lookup(l.config.Type)
	if obj == nil {
		return Union{}, fmt.Errorf("could not find type %s in package", l.config.Type)
	}

	structType, ok := obj.Type().Underlying().(*types.Struct)
	if !ok {
		return Union{}, fmt.Errorf("type %s must be a struct", l.config.Type)
	}

	fmt.Printf("%s : %s : %s\n", pkg.Module.Path, pkg.Name, obj.Name())
	fmt.Printf("%s\n", structType.String())

	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		fmt.Printf("%s : %v\n", field.Name(), field.Type())
	}

	return Union{}, nil
}
