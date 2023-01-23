package loader

import (
	"fmt"
	"go/types"
	"path"
	"strings"

	"github.com/sidkurella/gunion/internal/config"
	"golang.org/x/tools/go/packages"
)

// Loads in the union details based on the input configuration.
type Loader struct {
	config config.InputConfig
}

type TypeParam struct {
	// Name of this type parameter.
	Name string
	// Interface type of this type parameter.
	Type Type
}

// Represents the info we need to know about a type.
type Type struct {
	// Name of this type.
	Name string
	// Number of indirections (i.e., how many stars.)
	IndirectCount int
	// Originating package of this type.
	Source string
	// Type parameter list.
	TypeParams []TypeParam
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
	// Type of this union.
	Type Type
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
		// TODO: Identify which of these are unnecessary.
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

	obj := pkg.Types.Scope().Lookup(l.config.Type)
	if obj == nil {
		return Union{}, fmt.Errorf("could not find type %s in package", l.config.Type)
	}

	structType, ok := obj.Type().Underlying().(*types.Struct)
	if !ok {
		return Union{}, fmt.Errorf("type %s must be a struct", l.config.Type)
	}

	// TODO: Change to debug log.
	fmt.Printf("%s : %s : %s\n", pkg.Module.Path, pkg.Name, obj.Name())
	fmt.Printf("%s\n", structType.String())

	var variants []Variant
	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		fieldType := field.Type().String()
		typeInfo := parseType(fieldType)

		variant := Variant{
			Name: field.Name(),
			Type: typeInfo,
		}

		// TODO: Change to debug log.
		fmt.Printf("%s : %v : %#v\n", field.Name(), fieldType, typeInfo)

		variants = append(variants, variant)
	}

	namedType, ok := obj.Type().(*types.Named)
	if !ok {
		return Union{}, fmt.Errorf("type %s must be a named type, but it was not", l.config.Type)
	}

	return Union{
		Type:     parseUnionType(namedType),
		Variants: variants,
	}, nil
}

func parseUnionType(t *types.Named) Type {
	var typeParams []TypeParam
	typeParamList := t.TypeParams()
	for i := 0; i < typeParamList.Len(); i++ {
		tp := typeParamList.At(i)
		typeParams = append(typeParams, parseTypeParam(tp))
	}

	return Type{
		Name:          t.Obj().Name(),
		IndirectCount: 0,
		Source:        t.Obj().Pkg().Path(),
		TypeParams:    typeParams,
	}
}

func parseTypeParam(tp *types.TypeParam) TypeParam {
	return TypeParam{
		Name: tp.Obj().Name(),
		Type: parseType(tp.Constraint().String()),
	}
}

// This system of parsing types will only work for named types or basic types (primitives).
// TODO: Rework to support any other kind of struct member.
func parseType(t string) Type {
	rawType := strings.TrimLeft(t, "*")
	indirectCt := len(t) - len(rawType)

	ext := path.Ext(t)
	var typeName string
	var typeSource string
	if ext != "" {
		typeName = strings.TrimLeft(ext, ".")
		typeSource = strings.TrimSuffix(rawType, ext)
	} else {
		typeName = rawType
	}

	return Type{
		Name:          typeName,
		IndirectCount: indirectCt,
		Source:        typeSource,
	}
}
