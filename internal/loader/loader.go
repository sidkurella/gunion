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

	return parseNamedWithDepth(namedType, true)
}

func parseType(t gotypes.Type) (types.Type, error) {
	switch typ := t.(type) {
	case *gotypes.Alias:
		return parseAlias(typ)
	case *gotypes.Array:
		return parseArray(typ)
	case *gotypes.Basic:
		return parseBasic(typ)
	case *gotypes.Chan:
		return parseChan(typ)
	case *gotypes.Interface:
		return parseInterface(typ)
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
		return parseSignature(typ)
	case *gotypes.TypeParam:
		return parseTypeParamAsType(typ)
	case *gotypes.Union:
		return parseUnion(typ)
	default:
		return nil, fmt.Errorf("don't know how to handle type %T", t)
	}
}

func parseSignature(t *gotypes.Signature) (types.Signature, error) {
	var receiver types.Var
	var receiverTypeParams []types.TypeParam

	// Parse receiver if this is a method.
	if recv := t.Recv(); recv != nil {
		var err error
		receiver, err = parseVar(recv)
		if err != nil {
			return types.Signature{}, fmt.Errorf("failed to parse receiver: %w", err)
		}

		// Parse receiver type params.
		receiverTypeParams, err = parseTypeParamList(t.RecvTypeParams())
		if err != nil {
			return types.Signature{}, fmt.Errorf("failed to parse receiver type params: %w", err)
		}
	}

	// Parse input parameters.
	params, err := parseTuple(t.Params())
	if err != nil {
		return types.Signature{}, fmt.Errorf("failed to parse params: %w", err)
	}

	// Parse return values.
	returns, err := parseTuple(t.Results())
	if err != nil {
		return types.Signature{}, fmt.Errorf("failed to parse returns: %w", err)
	}

	// Parse type parameters.
	typeParams, err := parseTypeParamList(t.TypeParams())
	if err != nil {
		return types.Signature{}, fmt.Errorf("failed to parse type params: %w", err)
	}

	return types.Signature{
		Receiver:           receiver,
		ReceiverTypeParams: receiverTypeParams,
		Params:             params,
		Returns:            returns,
		TypeParams:         typeParams,
		Variadic:           t.Variadic(),
	}, nil
}

func parseTuple(t *gotypes.Tuple) ([]types.Var, error) {
	if t == nil {
		return nil, nil
	}
	var ret []types.Var
	for i := 0; i < t.Len(); i++ {
		v, err := parseVar(t.At(i))
		if err != nil {
			return nil, fmt.Errorf("failed to parse tuple element %d: %w", i, err)
		}
		ret = append(ret, v)
	}
	return ret, nil
}

func parseFunc(t *gotypes.Func) (types.Func, error) {
	name := t.Name()
	sig, ok := t.Type().(*gotypes.Signature)
	if !ok {
		return types.Func{},
			fmt.Errorf("expected function %s to have underlying type Signature, but it did not", name)
	}
	signature, err := parseSignature(sig)
	if err != nil {
		return types.Func{}, fmt.Errorf("failed to parse signature for function %s: %w", name, err)
	}

	return types.Func{
		Name:      name,
		Signature: signature,
	}, nil
}

func parseInterface(t *gotypes.Interface) (types.Interface, error) {
	var embeds []types.Type
	// Parse embedded types.
	for i := 0; i < t.NumEmbeddeds(); i++ {
		typ, err := parseType(t.EmbeddedType(i))
		if err != nil {
			return types.Interface{}, fmt.Errorf("failed to parse embedded type %d of interface: %w", i, err)
		}
		embeds = append(embeds, typ)
	}

	// Parse explicit methods.
	var methods []types.Func
	for i := 0; i < t.NumExplicitMethods(); i++ {
		ifaceMethod := t.ExplicitMethod(i)
		method, err := parseFunc(ifaceMethod)
		if err != nil {
			return types.Interface{}, fmt.Errorf("failed to parse method %d of interface: %w", i, err)
		}
		methods = append(methods, method)
	}

	return types.Interface{
		Embeds:  embeds,
		Methods: methods,
	}, nil
}

func parseUnion(t *gotypes.Union) (types.Union, error) {
	var ret []types.UnionMember
	for i := 0; i < t.Len(); i++ {
		term := t.Term(i)
		member, err := parseTerm(term)
		if err != nil {
			return types.Union{}, fmt.Errorf("failed to parse union member %d: %w", i, err)
		}

		ret = append(ret, member)
	}

	return types.Union{
		Members: ret,
	}, nil
}

func parseTerm(t *gotypes.Term) (types.UnionMember, error) {
	typ, err := parseType(t.Type())
	if err != nil {
		return types.UnionMember{}, fmt.Errorf("failed to parse union member: %w", err)
	}

	return types.UnionMember{
		Approximate: t.Tilde(),
		Type:        typ,
	}, nil
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

func parseAlias(t *gotypes.Alias) (types.Named, error) {
	// Type aliases (e.g., `type any = interface{}`) are represented as Named types.
	// We don't expand the underlying type to avoid issues with built-in aliases.
	name := t.Obj().Name()
	var pkg string
	if t.Obj().Pkg() != nil {
		pkg = t.Obj().Pkg().Path()
	}
	return types.Named{
		Name:    name,
		Package: pkg,
	}, nil
}

func parseBasic(t *gotypes.Basic) (types.Basic, error) {
	return types.Basic{
		Name: t.Name(),
	}, nil
}

func parseNamed(t *gotypes.Named) (types.Named, error) {
	return parseNamedWithDepth(t, false)
}

func parseNamedWithDepth(t *gotypes.Named, expandUnderlying bool) (types.Named, error) {
	name := t.Obj().Name()

	var underlyingType types.Type
	if expandUnderlying {
		// Only expand the underlying type for the top-level type to avoid infinite recursion
		// on types with circular references.
		var err error
		underlyingType, err = parseType(t.Underlying())
		if err != nil {
			return types.Named{}, fmt.Errorf("failed to parse underlying type for named type %s: %w", name, err)
		}
	}

	typeParams, err := parseTypeParamList(t.TypeParams())
	if err != nil {
		return types.Named{}, fmt.Errorf("failed to parse type params for named type %s: %w", name, err)
	}
	typeArgs, err := parseTypeList(t.TypeArgs())
	if err != nil {
		return types.Named{}, fmt.Errorf("failed to parse type args for named type %s: %w", name, err)
	}

	var pkg string
	if t.Obj().Pkg() != nil {
		pkg = t.Obj().Pkg().Path()
	}

	return types.Named{
		Name:       name,
		Package:    pkg,
		Type:       underlyingType,
		TypeParams: typeParams,
		TypeArgs:   typeArgs,
	}, nil
}

func parseTypeParamList(t *gotypes.TypeParamList) ([]types.TypeParam, error) {
	var ret []types.TypeParam
	for tp := range t.TypeParams() {
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

// parseTypeParamAsType handles type parameters when they appear as types (e.g., field `a T`).
// We represent them as Named types with just the name, since the actual type is determined by instantiation.
func parseTypeParamAsType(t *gotypes.TypeParam) (types.Named, error) {
	var pkg string
	if t.Obj().Pkg() != nil {
		pkg = t.Obj().Pkg().Path()
	}
	return types.Named{
		Name:    t.Obj().Name(),
		Package: pkg,
	}, nil
}

func parseTypeList(t *gotypes.TypeList) ([]types.Type, error) {
	var ret []types.Type
	for typ := range t.Types() {
		parsedType, err := parseType(typ)
		if err != nil {
			return nil, err
		}
		ret = append(ret, parsedType)
	}
	return ret, nil
}
