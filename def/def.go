package def

import (
	"fmt"
	"strings"

	"github.com/nathanborror/startapp/graphql"
	"github.com/nathanborror/startapp/graphql/common"
)

type Definition struct {
	Queries    []FuncDef
	Mutations  []FuncDef
	Scalars    []TypeDef
	Objects    []TypeDef
	Interfaces []TypeDef
	Unions     []TypeDef
	Enums      []TypeDef
	Inputs     []TypeDef
}

type FuncDef struct {
	Name      string
	Arguments ArgDefs
	Return    TypeDef
}

type ArgDef struct {
	Name string
	Type TypeDef
}

type ArgDefs []ArgDef

type TypeDef struct {
	Name          string
	Fields        []FieldDef
	Interfaces    []string
	IsScalar      bool
	IsOptional    bool
	IsInterface   bool
	IsEnum        bool
	IsList        bool
	EnumValues    []string
	PossibleTypes []TypeDef
}

type FieldDef struct {
	Name string
	Type TypeDef
}

func New(s *graphql.Schema) Definition {
	def := Definition{}
	if s == nil {
		return def
	}
	for _, v := range s.Queries {
		fn := NewFunction(v)
		def.Queries = append(def.Queries, fn)
	}
	for _, v := range s.Mutations {
		fn := NewFunction(v)
		def.Mutations = append(def.Mutations, fn)
	}
	for _, v := range s.Scalars {
		tp := NewType(v)
		def.Scalars = append(def.Scalars, tp)
	}
	for _, v := range s.Objects {
		tp := NewType(v)
		def.Objects = append(def.Objects, tp)
	}
	for _, v := range s.Interfaces {
		tp := NewType(v)
		def.Interfaces = append(def.Interfaces, tp)
	}
	for _, v := range s.Unions {
		tp := NewType(v)
		def.Unions = append(def.Unions, tp)
	}
	for _, v := range s.Enums {
		tp := NewType(v)
		def.Enums = append(def.Enums, tp)
	}
	for _, v := range s.Inputs {
		tp := NewType(v)
		def.Inputs = append(def.Inputs, tp)
	}
	return def
}

func NewFunction(in *graphql.Field) FuncDef {
	out := FuncDef{}
	out.Name = in.Name
	out.Return = NewType(in.Type)
	for _, arg := range in.Args {
		out.Arguments = append(out.Arguments, ArgDef{
			Name: arg.Name.Name,
			Type: NewType(arg.Type),
		})
	}
	return out
}

func NewType(in common.Type) TypeDef {
	out := TypeDef{}
	out.IsOptional = true

	switch t := in.(type) {
	case *graphql.Scalar:
		out.Name = t.Name
		out.IsScalar = true
	case *graphql.Object:
		out.Name = t.Name
		out.Fields = NewFields(t.Fields)
		out.Interfaces = NewInterfaces(t.Interfaces)
	case *graphql.Interface:
		out.Name = t.Name
		out.Fields = NewFields(t.Fields)
		out.IsInterface = true
		for _, v := range t.PossibleTypes {
			obj := NewType(v)
			out.PossibleTypes = append(out.PossibleTypes, obj)
		}
	case *graphql.Union:
		out.Name = t.Name
		for _, v := range t.PossibleTypes {
			obj := NewType(v)
			out.PossibleTypes = append(out.PossibleTypes, obj)
		}
	case *graphql.Enum:
		out.Name = t.Name
		out.IsEnum = true
		for _, v := range t.Values {
			out.EnumValues = append(out.EnumValues, v.Name)
		}
	case *graphql.InputObject:
		out.Name = t.Name
		out.Fields = NewFieldInputs(t.Values)
	case *common.List:
		out.Name = t.OfType.String()
		out.IsList = true
	case *common.NonNull:
		out = NewType(t.OfType)
		out.IsOptional = false
	default:
		fmt.Printf("> %#v\n", t)
		break
	}
	return out
}

func NewFields(in graphql.FieldList) []FieldDef {
	out := []FieldDef{}
	for _, f := range in {
		field := NewField(*f)
		out = append(out, field)
	}
	return out
}

func NewField(in graphql.Field) FieldDef {
	out := FieldDef{}
	out.Name = in.Name
	out.Type = NewType(in.Type)
	return out
}

func NewFieldInputs(in common.InputValueList) []FieldDef {
	out := []FieldDef{}
	for _, f := range in {
		field := NewFieldInput(*f)
		out = append(out, field)
	}
	return out
}

func NewFieldInput(in common.InputValue) FieldDef {
	out := FieldDef{}
	out.Name = in.Name.Name
	out.Type = NewType(in.Type)
	return out
}

func NewUnionValues(in []*graphql.Object) []FieldDef {
	out := []FieldDef{}
	for _, v := range in {
		field := NewUnionValue(v)
		out = append(out, field)
	}
	return out
}

func NewUnionValue(in *graphql.Object) FieldDef {
	out := FieldDef{}
	out.Name = in.Name
	out.Type = NewType(in)
	return out
}

func NewInterfaces(in []*graphql.Interface) []string {
	out := []string{}
	for _, v := range in {
		out = append(out, v.Name)
	}
	return out
}

// Strings

func (v Definition) String() string {

	// Queries & Mutations
	out := "type Query {\n"
	for _, query := range v.Queries {
		out = out + fmt.Sprintf("\t%s\n", query.String())
	}
	out = out + "}\n\ntype Mutation {\n"
	for _, mutation := range v.Mutations {
		out = out + fmt.Sprintf("\t%s\n", mutation.String())
	}
	out = out + "}\n\n"

	// Scalars
	for _, v := range v.Scalars {
		out = out + fmt.Sprintf("scalar %s\n\n", v.Name)
	}

	// Objects
	for _, v := range v.Objects {
		out += printType("type", v) + "\n\n"
	}

	// Inputs
	for _, v := range v.Inputs {
		out += printType("input", v) + "\n\n"
	}

	// Interfaces
	for _, v := range v.Interfaces {
		out += printInterface(v) + "\n\n"
	}

	// Enums
	for _, v := range v.Enums {
		out += printEnum(v) + "\n\n"
	}

	// Unions
	for _, v := range v.Unions {
		out += printType("union", v) + "\n\n"
	}
	return out
}

func printType(prefix string, in TypeDef) string {
	out := fmt.Sprintf("%s %s {\n", prefix, in.Name)
	for _, f := range in.Fields {
		out = out + fmt.Sprintf("\t%s\n", f.String())
	}
	return out + "}"
}

func printInterface(in TypeDef) string {
	out := fmt.Sprintf("interface %s {\n", in.Name)
	for _, f := range in.Fields {
		out = out + fmt.Sprintf("\t%s\n", f.String())
	}
	return out + "}"
}

func printEnum(in TypeDef) string {
	out := fmt.Sprintf("enum %s {\n", in.Name)
	for _, v := range in.EnumValues {
		out = out + fmt.Sprintf("\t%s\n", v)
	}
	return out + "}"
}

func (v FuncDef) String() string {
	return fmt.Sprintf("%s(%s): %s", v.Name, v.Arguments.String(), v.Return.String())
}

func (v ArgDefs) String() string {
	var args []string
	for _, arg := range v {
		args = append(args, arg.String())
	}
	return strings.Join(args, ", ")
}

func (v TypeDef) String() string {
	if v.IsOptional {
		return fmt.Sprintf("%s", v.Name)
	}
	return fmt.Sprintf("%s!", v.Name)
}

func (v ArgDef) String() string {
	return fmt.Sprintf("%s: %s", v.Name, v.Type.String())
}

func (v FieldDef) String() string {
	return fmt.Sprintf("%s: %s", v.Name, v.Type.String())
}

// Swift

// JoinArgsForSwift takes a list of ArgDef and returns a concatenated string
// suitable for Swift function definitions: `name: String, email: String?`
func JoinArgsForSwift(in []ArgDef) string {
	var out []string
	for _, arg := range in {
		optional := ""
		if arg.Type.IsOptional {
			optional = "?"
		}
		out = append(out, fmt.Sprintf("%s: %s%s", arg.Name, arg.Type.Name, optional))
	}
	return strings.Join(out, ", ")
}

// JoinInterfacesForSwift takes a list of InterfaceDef and returns a concatenated
// string suitable for Swift class or struct definitions: `Node, Human`
func JoinInterfacesForSwift(in []string) string {
	return strings.Join(in, ", ")
}

// ExcludeSwiftScalars returns a TypeDef list without any Swift native scalars.
func ExcludeSwiftScalars(defs []TypeDef) []TypeDef {
	var out []TypeDef
	exclude := map[string]bool{
		"Float":   true,
		"Int":     true,
		"String":  true,
		"Boolean": true,
	}
	for _, t := range defs {
		if exclude[t.Name] {
			continue
		}
		out = append(out, t)
	}
	return out
}

func GraphQLScalarToSwiftScalar(in TypeDef) string {
	if !in.IsScalar {
		return in.Name
	}
	return map[string]string{
		"ID":        "String",
		"Timestamp": "String",
	}[in.Name]
}

func ToSwiftScalar(in string) string {
	convert := map[string]string{
		"Boolean": "Bool",
	}
	if scalar, ok := convert[in]; ok {
		return scalar
	}
	return in
}

// Go

func ToGoScalar(in string) string {
	convert := map[string]string{
		"Boolean": "bool",
		"String":  "string",
	}
	if scalar, ok := convert[in]; ok {
		return scalar
	}
	return in
}

// GraphQL

func JoinArgsForGraphQL(in ArgDefs) string {
	var out []string
	for _, arg := range in {
		optional := "!"
		if arg.Type.IsOptional {
			optional = ""
		}
		out = append(out, fmt.Sprintf("$%s: %s%s", arg.Name, arg.Type.Name, optional))
	}
	return strings.Join(out, ", ")
}

func JoinArgsForGraphQLVars(in ArgDefs) string {
	var out []string
	for _, arg := range in {
		out = append(out, fmt.Sprintf("%s: $%s", arg.Name, arg.Name))
	}
	return strings.Join(out, ", ")
}
