package loader

import (
	"fmt"
	gotypes "go/types"

	"github.com/sidkurella/gunion/internal/config"
	"github.com/sidkurella/gunion/internal/types"
	"golang.org/x/tools/go/packages"
)

// Loads in the union details based on the input configuration.
type Loader struct {
	config config.InputConfig
}

func NewLoader(config config.InputConfig) *Loader {
	return &Loader{
		config: config,
	}
}

func (l *Loader) Load() (types.Named, error) {
	pkgs, err := packages.Load(&packages.Config{
		// TODO: Identify which of these are unnecessary.
		Mode: packages.NeedTypes | packages.NeedImports | packages.NeedSyntax | packages.NeedTypesInfo |
			packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedModule,
	}, "file="+l.config.Source)
	if err != nil {
		return types.Named{}, fmt.Errorf("failed to load packages for source file %s: %w", l.config.Source, err)
	}

	if len(pkgs) != 1 {
		return types.Named{}, fmt.Errorf("expected to load 1 package but got %d", len(pkgs))
	}

	pkg := pkgs[0]

	obj := pkg.Types.Scope().Lookup(l.config.Type)
	if obj == nil {
		return types.Named{}, fmt.Errorf("could not find type %s in package", l.config.Type)
	}

	namedType, ok := obj.Type().(*gotypes.Named)
	if !ok {
		return types.Named{}, fmt.Errorf("type %s must be a named type, but it was not", l.config.Type)
	}

	// TODO: Change to debug log.
	fmt.Printf("%s : %s : %s\n", pkg.Module.Path, pkg.Name, obj.Name())

	return parseNamed(namedType)
}

func parseType(t gotypes.Type) (types.Type, error) {
	switch typ := t.(type) {
	case *gotypes.Array:
		return parseArray(typ)
	case *gotypes.Basic:
		return parseBasic(typ)
	case *gotypes.Chan:
		return parseChan(typ)
	case *gotypes.Interface:
		return nil, fmt.Errorf("unimplemented")
	case *gotypes.Map:
		return parseMap(typ)
	case *gotypes.Named:
		return parseNamed(typ)
	case *gotypes.Pointer:
		return parsePointer(typ)
	case *gotypes.Slice:
		return parseSlice(typ)
	case *gotypes.Struct:
		return parseStruct(typ)
	case *gotypes.Signature:
		return nil, fmt.Errorf("unimplemented")
	case *gotypes.Union:
		return nil, fmt.Errorf("unimplemented")
	default:
		return nil, fmt.Errorf("don't know how to handle type %T", t)
	}
}

func parseStruct(t *gotypes.Struct) (types.Struct, error) {
	var ret []types.Field
	for i := 0; i < t.NumFields(); i++ {
		field := t.Field(i)
		tag := t.Tag(i)
		v, err := parseVar(field)
		if err != nil {
			return types.Struct{}, fmt.Errorf("failed to parse struct field %d: %w", i, err)
		}

		ret = append(ret, types.Field{
			Var: v,
			Tag: tag,
		})
	}

	return types.Struct{
		Fields: ret,
	}, nil
}

func parseVar(t *gotypes.Var) (types.Var, error) {
	typ, err := parseType(t.Type())
	if err != nil {
		return types.Var{}, fmt.Errorf("failed to parse var type for var %s: %w", t.Name(), err)
	}

	return types.Var{
		Name: t.Name(),
		Type: typ,
	}, nil
}

func parseChan(t *gotypes.Chan) (types.Chan, error) {
	var chanDir types.ChanDir
	switch t.Dir() {
	case gotypes.SendRecv:
		chanDir = types.SendRecv
	case gotypes.SendOnly:
		chanDir = types.SendOnly
	case gotypes.RecvOnly:
		chanDir = types.RecvOnly
	default:
		return types.Chan{}, fmt.Errorf("invalid chan direction %v", t.Dir())
	}

	elem, err := parseType(t.Elem())
	if err != nil {
		return types.Chan{}, fmt.Errorf("failed to parse chan element type: %w", err)
	}

	return types.Chan{
		Direction: chanDir,
		Elem:      elem,
	}, nil
}

func parseMap(t *gotypes.Map) (types.Map, error) {
	key, err := parseType(t.Key())
	if err != nil {
		return types.Map{}, fmt.Errorf("failed to parse map key type: %w", err)
	}
	value, err := parseType(t.Elem())
	if err != nil {
		return types.Map{}, fmt.Errorf("failed to parse map value type: %w", err)
	}

	return types.Map{
		Key:   key,
		Value: value,
	}, nil
}

func parsePointer(t *gotypes.Pointer) (types.Pointer, error) {
	elem, err := parseType(t.Elem())
	if err != nil {
		return types.Pointer{}, fmt.Errorf("failed to parse pointer element type: %w", err)
	}

	return types.Pointer{
		Elem: elem,
	}, nil
}

func parseArray(t *gotypes.Array) (types.Array, error) {
	elem, err := parseType(t.Elem())
	if err != nil {
		return types.Array{}, fmt.Errorf("failed to parse array element type: %w", err)
	}

	return types.Array{
		Len:  t.Len(),
		Elem: elem,
	}, nil
}

func parseSlice(t *gotypes.Slice) (types.Slice, error) {
	elem, err := parseType(t.Elem())
	if err != nil {
		return types.Slice{}, fmt.Errorf("failed to parse slice element type: %w", err)
	}

	return types.Slice{
		Elem: elem,
	}, nil
}

func parseBasic(t *gotypes.Basic) (types.Basic, error) {
	return types.Basic{
		Name: t.Name(),
	}, nil
}

func parseNamed(t *gotypes.Named) (types.Named, error) {
	name := t.Obj().Name()
	underlyingType, err := parseType(t.Underlying())
	if err != nil {
		return types.Named{}, fmt.Errorf("failed to parse underlying type for named type %s: %w", name, err)
	}
	typeParams, err := parseTypeParamList(t.TypeParams())
	if err != nil {
		return types.Named{}, fmt.Errorf("failed to parse type params for named type %s: %w", name, err)
	}
	typeArgs, err := parseTypeList(t.TypeArgs())
	if err != nil {
		return types.Named{}, fmt.Errorf("failed to parse type args for named type %s: %w", name, err)
	}

	// TODO: We likely don't need to resolve the entire type tree here.
	// TODO: We likely can treat all Named nodes as leaves and not resolve its underlying type.
	// TODO: We would still need to resolve type parameters and type arguments, however.
	return types.Named{
		Name:       name,
		Package:    t.Obj().Pkg().Path(),
		Type:       underlyingType,
		TypeParams: typeParams,
		TypeArgs:   typeArgs,
	}, nil
}

func parseTypeParamList(t *gotypes.TypeParamList) ([]types.TypeParam, error) {
	var ret []types.TypeParam
	for i := 0; i < t.Len(); i++ {
		tp := t.At(i)
		param, err := parseTypeParam(tp)
		if err != nil {
			return nil, err
		}
		ret = append(ret, param)
	}
	return ret, nil
}

func parseTypeParam(t *gotypes.TypeParam) (types.TypeParam, error) {
	name := t.Obj().Name()
	constraint, err := parseType(t.Constraint())
	if err != nil {
		return types.TypeParam{}, fmt.Errorf("failed to parse constraint for type parameter %s: %w", name, err)
	}

	return types.TypeParam{
		Name:       t.Obj().Name(),
		Constraint: constraint,
	}, nil
}

func parseTypeList(t *gotypes.TypeList) ([]types.Type, error) {
	var ret []types.Type
	for i := 0; i < t.Len(); i++ {
		typ := t.At(i)
		parsedType, err := parseType(typ)
		if err != nil {
			return nil, err
		}
		ret = append(ret, parsedType)
	}
	return ret, nil
}
