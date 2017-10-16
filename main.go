package main

//go:generate esc -o static.go -prefix "/" static

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
)

var (
	defaultName = "Main"
	name        = flag.String("name", "App", "Name of the project")
	hasKit      = flag.Bool("kit", false, "Adds a Kit framework to project")
	hasTests    = flag.Bool("tests", false, "Adds Tests to project")
	hasUITests  = flag.Bool("uitests", false, "Adds UI Tests to project")
	teamID      = flag.String("teamid", "", "Apple Developer account Team ID")
	schemaFile  = flag.String("schema", "", "GraphQL schema used to write types")
	liveFlag    = flag.Bool("liveassets", false, "Serve Static Assets from Disk relative to CWD")
)

type fileContext struct {
	Name       string
	HasKit     bool
	HasTests   bool
	HasUITests bool
	TeamID     string
	Schema     *Schema
}

func main() {
	flag.Parse()

	ctx := fileContext{
		Name:       *name,
		HasKit:     *hasKit,
		HasTests:   *hasTests,
		HasUITests: *hasUITests,
		TeamID:     *teamID,
	}

	if schemaFile != nil && *schemaFile != "" {
		schema, err := parseSchema(*schemaFile)
		if err != nil {
			log.Fatal(err)
		}
		ctx.Schema = schema
	}

	dirs := []string{}
	manifest := make(map[string]string)

	for template, item := range _escData {
		filepath := strings.TrimPrefix(item.local, "static/")

		if filepath == "" {
			continue
		}
		if filepath == ".DS_Store" {
			continue
		}
		if filepath == "static" {
			continue
		}

		filepath = fmt.Sprintf("%s/%s", ctx.Name, filepath)
		filepath = strings.Replace(filepath, defaultName, ctx.Name, -1)

		if item.isDir {
			dirs = append(dirs, filepath)
		} else {
			manifest[template] = filepath
		}
	}

	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			os.MkdirAll(dir, os.ModePerm)
		}
	}
	for tmpl, filepath := range manifest {
		if err := writeFile(tmpl, filepath, ctx); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Created file: ", filepath)
	}
}

func writeFile(tmpl string, filepath string, ctx fileContext) error {
	if _, err := os.Stat(filepath); !os.IsNotExist(err) {
		return err
	}
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("Failed to create file: %s", err)
	}
	path, err := FSString(*liveFlag, tmpl)
	if err != nil {
		return fmt.Errorf("Failed to read template: %s", err)
	}
	t, err := template.New(tmpl).Funcs(
		template.FuncMap{
			"cast":             cast,
			"joinInterfaces":   joinInterfaces,
			"excludeFunctions": excludeFunctions,
			"onlyScalars":      onlyScalars,
		},
	).Parse(path)
	if err != nil {
		return fmt.Errorf("Failed to parse template: %s", err)
	}
	if err := t.Execute(file, ctx); err != nil {
		return fmt.Errorf("Failed to execute template: %s", err)
	}
	file.Close()
	return nil
}

var swiftCast = map[string]string{
	"ID": "String",
}

func cast(in string) string {
	if swift, ok := swiftCast[in]; ok {
		return swift
	}
	return in
}

func (s *Schema) ScalarKind() []Type {
	var exclude = map[string]bool{
		"Boolean": true,
		"Float":   true,
		"Int":     true,
		"String":  true,
	}
	types := s.ForKind(ScalarKind)

	var out []Type
	for _, t := range types {
		if _, ok := exclude[t.Name]; !ok {
			out = append(out, t)
		}
	}
	return out
}

func (s *Schema) InputKind() []Type {
	return s.ForKind(InputObjectKind)
}

func (s *Schema) PayloadKind() []Type {
	return s.ForKind(PayloadObjectKind)
}

func (s *Schema) InterfaceKind() []Type {
	return s.ForKind(InterfaceKind)
}

func (s *Schema) UnionKind() []Type {
	return s.ForKind(UnionKind)
}

func (s *Schema) EnumKind() []Type {
	return s.ForKind(EnumKind)
}

func (s *Schema) ObjectKind() []Type {
	types := s.ForKind(ObjectKind)

	var out []Type
	for _, t := range types {
		if !t.IsMeta() {
			out = append(out, t)
		}
	}
	return out
}

func (s *Schema) ForKind(in TypeKind) (out []Type) {
	for _, t := range s.Types {
		if t.Kind == in {
			out = append(out, t)
		}
	}
	return
}

// Fields

func (in *Field) SwiftField() string {
	return in.SwiftFieldWithPrefix("")
}

func (in *Field) SwiftFieldWithPrefix(prefix string) string {
	var out string
	if in.IsFunction() {
		args := joinArgs(in.Args, prefix)
		if prefix == "" {
			out = fmt.Sprintf("func %s(%s) -> %s?", in.Name, args, scalar(in.Type.Name))
		} else {
			out = fmt.Sprintf("func %s(%s) -> %s.%s?", in.Name, args, prefix, scalar(in.Type.Name))
		}
	} else {
		if len(prefix) > 0 {
			out = fmt.Sprintf("let %s: %s.%s?", in.Name, prefix, scalar(in.Type.Name))
		} else {
			out = fmt.Sprintf("let %s: %s?", in.Name, scalar(in.Type.Name))
		}
	}
	return out
}

func (in *Field) SwiftProtocolField() string {
	return in.SwiftProtocolFieldWithPrefix("")
}

func (in *Field) SwiftProtocolFieldWithPrefix(prefix string) string {
	var out string
	if len(prefix) > 0 {
		out = fmt.Sprintf("var %s: %s.%s? { get }", in.Name, prefix, scalar(in.Type.Name))
	} else {
		out = fmt.Sprintf("let %s: %s? { get }", in.Name, scalar(in.Type.Name))
	}
	return out
}

func (in *Field) SwiftCase() string {
	return fmt.Sprintf("case %s(%s)", strings.ToLower(in.Name), in.Name)
}

func joinArgs(in []Field, prefix string) string {
	var out []string
	for _, i := range in {
		var stmt string
		if prefix == "" {
			stmt = fmt.Sprintf("%s: %s?", i.Name, scalar(i.Type.Name))
		} else {
			stmt = fmt.Sprintf("%s: %s.%s?", i.Name, prefix, scalar(i.Type.Name))
		}
		out = append(out, stmt)
	}
	return strings.Join(out, ", ")
}

func joinInterfaces(in []Field) string {
	var out []string
	for _, i := range in {
		stmt := fmt.Sprintf("%s", i.Name)
		out = append(out, stmt)
	}
	out = append(out, "Codable")
	return strings.Join(out, ", ")
}

var graphQLToSwift = map[string]string{
	"Boolean": "Bool",
}

func scalar(in *string) string {
	var typename string
	if in != nil {
		typename = *in
	}
	if out, ok := graphQLToSwift[typename]; ok {
		return out
	}
	return typename
}

func excludeFunctions(in []Field) (out []Field) {
	for _, field := range in {
		if !field.IsFunction() {
			out = append(out, field)
		}
	}
	return
}

func onlyScalars(in []Field) (out []Field) {
	for _, field := range in {
		if field.IsScalar() {
			out = append(out, field)
		}
	}
	return
}
