package main

import (
	"encoding/json"
	"io/ioutil"
	"strings"
)

type GraphQL struct {
	Data struct {
		Schema Schema `json:"__schema"`
	} `json:"data"`
}

type Schema struct {
	Types     []Type
	Mutations []Field
	Queries   []Field
	QueryType struct {
		Name string
	}
	MutationType struct {
		Name string
	}
}

type Type struct {
	Kind          TypeKind
	Name          string
	Description   *string
	Fields        []Field
	InputFields   []Field
	Interfaces    []Field
	PossibleTypes []Field
}

func (t *Type) IsMeta() bool {
	return strings.HasPrefix(t.Name, "__")
}

func (t *Type) IsEdge() bool {
	return strings.HasSuffix(t.Name, "Edge")
}

func (t *Type) IsConnection() bool {
	return strings.HasSuffix(t.Name, "Connection")
}

type Field struct {
	Name         string
	Description  *string
	Type         FieldType
	Args         []Field
	DefaultValue *string
}

func (f *Field) IsFunction() bool {
	return len(f.Args) > 0
}

func (f *Field) IsScalar() bool {
	if f.Type.OfType != nil {
		return f.Type.OfType.Kind == ScalarKind
	}
	return f.Type.Kind == ScalarKind
}

type FieldType struct {
	Kind   TypeKind
	Name   *string
	OfType *FieldType
}

type TypeKind string

const (
	ScalarKind        TypeKind = "SCALAR"
	ObjectKind                 = "OBJECT"
	UnionKind                  = "UNION"
	InterfaceKind              = "INTERFACE"
	EnumKind                   = "ENUM"
	InputObjectKind            = "INPUT_OBJECT"
	ListKind                   = "LIST"
	NonNullKind                = "NON_NULL"
	PayloadObjectKind          = "PAYLOAD_OBJECT"
)

func StringToTypeKind(in string) TypeKind {
	switch in {
	case "SCALAR":
		return ScalarKind
	case "OBJECT":
		return ObjectKind
	case "UNION":
		return UnionKind
	case "INTERFACE":
		return InterfaceKind
	case "ENUM":
		return EnumKind
	case "INPUT_OBJECT":
		return InputObjectKind
	case "LIST":
		return ListKind
	case "NON_NULL":
		return NonNullKind
	case "PAYLOAD_OBJECT":
		return PayloadObjectKind
	}
	return ""
}

func parseSchema(filename string) (*Schema, error) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var parsed GraphQL
	if err := json.Unmarshal(contents, &parsed); err != nil {
		return nil, err
	}
	var schema Schema
	for _, t := range parsed.Data.Schema.Types {
		newType := t

		// Assign Payload suffixed types a new kind
		if strings.HasSuffix(t.Name, "Payload") {
			newType.Kind = PayloadObjectKind
		}

		// Hoist Query types into their own property
		if t.Name == parsed.Data.Schema.QueryType.Name {
			schema.Queries = cleanFields(t.Fields)
			continue
		}

		// Hoist Mutation types into their own property
		if t.Name == parsed.Data.Schema.MutationType.Name {
			schema.Mutations = cleanFields(t.Fields)
			continue
		}

		newType.Fields = cleanFields(t.Fields)
		newType.InputFields = cleanFields(t.InputFields)
		schema.Types = append(schema.Types, newType)
	}
	return &schema, nil
}

func cleanFields(in []Field) (out []Field) {
	for _, field := range in {
		fieldType := field.Type
		fieldType.Name = fieldTypeName(fieldType)
		f := field
		f.Type = fieldType
		if len(f.Args) > 0 {
			f.Args = cleanFields(f.Args)
		}
		out = append(out, f)
	}
	return
}

func fieldTypeName(field FieldType) *string {
	if field.Name != nil {
		return field.Name
	}
	if field.OfType != nil {
		return fieldTypeName(*field.OfType)
	}
	return nil
}
