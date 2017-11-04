package gen

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/nathanborror/startapp/def"
)

// File represents a file object that can written to a file system.
type File struct {
	buf  bytes.Buffer
	name string
	ext  string
	err  error
}

// NewFile returns a new File object.
func NewFile(name string, ext string) File {
	return File{name: name, ext: ext}
}

// Write takes the current buffer and writes a file to the given directory.
func (f *File) Write(elem ...string) {
	if f.err != nil {
		return
	}
	if len(elem) == 0 {
		f.err = fmt.Errorf("Missing directory")
		return
	}
	dir := filepath.Join(elem...)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			f.err = err
			return
		}
	}
	filePath := filepath.Join(dir, f.filename())
	if err := ioutil.WriteFile(filePath, f.buf.Bytes(), 0644); err != nil {
		f.err = err
		return
	}
}

// WriteBytes writes the given bytes to the file buffer.
func (f *File) WriteBytes(in []byte) {
	if _, err := f.buf.Write(in); err != nil {
		f.err = err
		return
	}
}

// GoFormat formats the file buffer using Go Format.
func (f *File) GoFormat() {
	if f.err != nil {
		return
	}
	if f.ext != "go" {
		f.err = fmt.Errorf("File must be a .go file")
		return
	}
	src, err := format.Source(f.buf.Bytes())
	if err != nil {
		fmt.Printf("*** Could not gofmt %s: %v\n", f.filename(), err)
		f.err = err
		return
	}
	f.buf.Reset()
	f.buf.Write(src)
	return
}

// PanicOnErr panics if an error exists.
func (f *File) PanicOnErr() {
	if f.err != nil {
		log.Fatal(f.err)
	}
}

// Err returns the first encoutnered error.
func (f *File) Err() error {
	return f.err
}

func (f *File) filename() string {
	return fmt.Sprintf("%s.%s", f.name, f.ext)
}

func (f *File) printf(format string, args ...interface{}) {
	fmt.Fprintf(&f.buf, format+"\n", args...)
}

// Go

func (f *File) WriteAPIResolver(projectPath string, t def.TypeDef) {
	f.printf("package api")
	f.printf("import ( \"%s/%s\"", projectPath, "state")
	f.printf("graphql \"github.com/neelance/graphql-go\" )")

	f.printf("type %sResolver struct {", strings.ToLower(t.Name))
	f.printf("%s *state.%s", strings.ToLower(t.Name), t.Name)
	f.printf("}")

	for _, field := range t.Fields {
		if field.Type.IsScalar {
			f.printf("func (r *%sResolver) %s() %s {", strings.ToLower(t.Name), strings.Title(field.Name), def.ToGoScalar(field.Type.Name))
			f.printf("return r.%s.%s", strings.ToLower(t.Name), strings.Title(field.Name))
			f.printf("}\n")
			continue
		}
		if field.Type.IsEnum {
			f.printf("// Implement Enum: %s\n", field.Name)
			continue
		}
		if field.Type.IsInterface {
			f.printf("// Implement Interface: %s\n", field.Name)
			continue
		}
		f.printf("func (r *%sResolver) %s() (*state.%s, error) {", strings.ToLower(t.Name), strings.Title(field.Name), def.ToGoScalar(field.Type.Name))
		f.printf("return nil, fmt.Errorf(\"Not Implemented\")")
		f.printf("}\n")
	}
}

func (f *File) WriteAPIQueries(projectPath string) {
	f.printf("package api")
	f.printf("import \"%s/%s\"", projectPath, "state")
}

func (f *File) WriteAPIMutations(projectPath string) {
	f.printf("package api")
	f.printf("import \"%s/%s\"", projectPath, "state")
}

func (f *File) WriteScalars(scalars []def.TypeDef) {
	f.printf("package api")
	for _, s := range scalars {
		f.printf("// %s Scalar\n", s.Name)
		f.printf("type %s struct {\n\t// TODO: Implement\n}\n", s.Name)

		f.printf("func (%s) ImplementsGraphQLType(name string) bool {", s.Name)
		f.printf("return name == \"%s\"", s.Name)
		f.printf("}\n")

		f.printf("func (s *%s) UnmarshalGraphQL(input interface{}) error {", s.Name)
		f.printf("return fmt.Errorf(\"Not Implemented\")")
		f.printf("}\n")
	}
}
