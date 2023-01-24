package types

// Represents an arbitrary Go-type.
type Type interface{}

type TypeParam struct {
	// Name of this type parameter
	Name string
	// Constraint type. Should be either Interface, Union, or Named.
	Constraint Type
}

// Represents a named variable.
type Var struct {
	// Name of the variable.
	Name string
	// Type of the variable.
	Type Type
}

// Represents a basic type (primitive).
type Basic struct {
	// Name of this primitive type.
	Name string
}

// Represents a type with a given name.
type Named struct {
	// Package-local name of this type.
	Name string
	// Path of the enclosing package where this type name was declared.
	Package string
	// Actual type of this named type.
	Type Type
	// Type parameters for this named type.
	TypeParams []TypeParam
	// Concrete type arguments provided to instantiate the named type, if any.
	TypeArgs []Type
}

// Represents an array with fixed size.
type Array struct {
	// Length of the array.
	Len int64
	// Element type.
	Elem Type
}

// Represents a slice.
type Slice struct {
	// Element type.
	Elem Type
}

type ChanDir int

const (
	SendRecv ChanDir = iota
	SendOnly
	RecvOnly
)

// Represents a channel.
type Chan struct {
	// Direction of the channel.
	Direction ChanDir
	// Element type.
	Elem Type
}

// Represents a map.
type Map struct {
	// Key type of the map.
	Key Type
	// Value type of the map.
	Value Type
}

// Represents a pointer to another type.
type Pointer struct {
	// Type this pointer refers to.
	Elem Type
}

// Represents a struct field.
type Field struct {
	// Name and type information of this field.
	Var
	// Tag on this field, if any.
	Tag string
}

// Represents a struct.
type Struct struct {
	// The fields of the struct.
	Fields []Field
}

// Member of a union type-set.
type UnionMember struct {
	// Does this union member include aliased types? (i.e., does it start with a tilde?)
	Approximate bool
	// Base type.
	Type Type
}

// Represents a union of types. Used for generics type constraints.
type Union struct {
	// Members of this type-set union.
	Members []UnionMember
}

// Represents an arbitrary function signature.
type Signature struct {
	// Type of the receiver, if any. If this is not a method, Receiver.Type will be nil.
	Receiver Var
	// Type parameters for the receiver of this function.
	ReceiverTypeParams []TypeParam
	// Input parameters of the signature.
	Params []Var
	// Output return values of the signature.
	Returns []Var
	// Type parameters for this function signature.
	TypeParams []TypeParam
	// Is this function variadic?
	Variadic bool
}

// Represents a named function.
type Func struct {
	// Name of this function.
	Name string
	// Signature of this function.
	Signature Signature
}

// Represents an interface.
type Interface struct {
	// Embedded types of this interface.
	Embeds []Type
	// Methods belonging to this interface.
	Methods []Func
}
