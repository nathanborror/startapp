package gen

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/nathanborror/startapp/def"
	"github.com/nathanborror/startapp/graphql"
)

const (
	templatesFolder = "templates"
	clientsFolder   = "clients"
	projectPrefix   = "Project"
)

var ignoredFiles = map[string]bool{
	"":          true,
	".DS_Store": true,
	"templates": true,
}

type Project struct {
	Name       string // The name of the project
	Domain     string // The domain the project is hosted on e.g. example.com
	Definition def.Definition
	Clients    []Client
	Templates  TemplateFiles
	dest       string
	err        error
}

type Client struct {
	Kind         ClientKind
	Name         string
	BundleDomain string // The bundle identifier domain e.g. com.example
	TeamID       string // Team identifier, generally used for iOS clients
	HasBackend   bool
	HasTests     bool
	Templates    TemplateFiles
}

type ClientKind string

const (
	IOSClientKind ClientKind = "ios"
	WebClientKind            = "www"
)

type TemplateFiles map[string]string // template name : file-path

// NewProject returns a new Project.
func NewProject(name string, dest string, domain string) *Project {
	return &Project{
		Name:      strings.ToLower(name),
		Domain:    domain,
		Templates: make(TemplateFiles),
		dest:      dest,
	}
}

// AddIOSClient appends a new iOS client to the Project.
func (p *Project) AddIOSClient(name string, teamID string, hasBackend bool, hasTests bool) {
	fmt.Printf("Adding iOS client: %s\n", name)
	client := Client{
		Kind:         IOSClientKind,
		Name:         name,
		BundleDomain: reverseDomain(p.Domain),
		TeamID:       teamID,
		HasBackend:   hasBackend,
		HasTests:     hasTests,
		Templates:    make(TemplateFiles),
	}
	p.Clients = append(p.Clients, client)
}

// ReadGraphQLSchema reads a GraphQL schema and adds a new Definition to the
// Project that will be used when rendering project templates.
func (p *Project) ReadGraphQLSchema(filename string) {
	fmt.Printf("Using GraphQL schema: %s\n", filename)
	if p.err != nil {
		return
	}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		p.err = err
		return
	}
	schema := graphql.New()
	if err := schema.Parse(data); err != nil {
		p.err = err
		return
	}
	p.Definition = def.New(schema)
}

// Write renders all the Project files and writes them out to their
// repsective directories.
func (p *Project) Write() {
	if p.err != nil {
		return
	}
	fmt.Println("Writing Templates...")
	p.ReadTemplates(_escData)
	fmt.Println("Writing template files...")
	p.WriteTemplateFiles()
	fmt.Println("Writing Go API scaffolding...")
	p.WriteGoScaffoldingForAPI()
	fmt.Println("Writing Swift GraphQL scaffolding...")
	p.WriteSwiftScaffoldingForGraphQL()
}

// Err returns the first encountered error.
func (p *Project) Err() error {
	return p.err
}

// ReadTemplates reads all the statically generated templates into the Project.
func (p *Project) ReadTemplates(templates map[string]*_escFile) {
	for name, file := range templates {
		filename := strings.TrimPrefix(file.local, templatesFolder+"/")

		// Ignore files
		if ignoredFiles[filename] {
			continue
		}

		// Rename files prefixed with 'Project' to the actual project name
		filename = strings.Replace(filename, projectPrefix, strings.Title(p.Name), -1)

		// Skip directories
		if file.isDir {
			continue
		}

		// Check for client templates
		if strings.HasPrefix(filename, clientsFolder) {
			for i, client := range p.Clients {
				prefix := fmt.Sprintf("%s/%s", clientsFolder, string(client.Kind))
				if strings.HasPrefix(filename, prefix) {
					p.Clients[i].Templates[name] = filename
					continue
				}
			}
		} else {
			p.Templates[name] = filename
		}
	}
}

// WriteTemplateFiles writes all rendered templates to the file system.
func (p *Project) WriteTemplateFiles() {
	p.writeTemplates(p.Templates)
	for _, client := range p.Clients {
		p.writeTemplates(client.Templates)
	}
}

// Copy copies a given file to the project root.
func (p *Project) Copy(filename string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		p.err = err
		return
	}

	ext := filepath.Ext(filename)
	name := filename[0 : len(filename)-len(ext)]

	file := NewFile(name, ext)
	file.WriteBytes(data)
	file.Write(p.dest, p.Name)
	file.PanicOnErr()
}

// WriteGoScaffoldingForAPI writes all the api scaffolding.
func (p *Project) WriteGoScaffoldingForAPI() {
	root := filepath.Join(p.dest, p.Name)
	dir := "api"

	for _, obj := range p.Definition.Objects {
		file := NewFile(strings.ToLower(obj.Name), "go")
		file.WriteAPIResolver("github.com/nathanborror/"+root, obj)
		file.GoFormat()
		file.Write(root, dir)
		file.PanicOnErr()
	}
	if len(p.Definition.Scalars) > 0 {
		file := NewFile("scalars", "go")
		file.WriteScalars(p.Definition.Scalars)
		file.GoFormat()
		file.Write(root, dir)
		file.PanicOnErr()
	}
	if len(p.Definition.Queries) > 0 {
		file := NewFile("queries", "go")
		file.WriteAPIQueries("github.com/nathanborror/" + root)
		file.GoFormat()
		file.Write(root, dir)
		file.PanicOnErr()
	}
	if len(p.Definition.Mutations) > 0 {
		file := NewFile("mutations", "go")
		file.WriteAPIMutations("github.com/nathanborror/" + root)
		file.GoFormat()
		file.Write(root, dir)
		file.PanicOnErr()
	}
}

// WriteSwiftScaffoldingForGraphQL writes empty '.graphql' files to the client's
// iOS GraphQL directory. Use these stubs to write your GraphQL queries for your
// clients to invoke.
func (p *Project) WriteSwiftScaffoldingForGraphQL() {
	var iosClient Client

	// Find iOS Client
	for _, client := range p.Clients {
		if client.Kind == IOSClientKind {
			iosClient = client
			break
		}
	}

	root := filepath.Join(p.dest, p.Name)
	dir := iosClient.GraphQLDir(p.Name)

	// Write mutation files
	for _, fn := range p.Definition.Mutations {
		file := NewFile(strings.ToLower(fn.Name), "graphql")
		file.Write(root, dir)
		file.PanicOnErr()
	}

	// Write query files
	for _, fn := range p.Definition.Queries {
		if fn.Return.IsInterface {
			for _, tp := range fn.Return.PossibleTypes {
				file := NewFile("node."+strings.ToLower(tp.Name), "graphql")
				file.Write(root, dir)
				file.PanicOnErr()
			}
		} else {
			file := NewFile(strings.ToLower(fn.Name), "graphql")
			file.Write(root, dir)
			file.PanicOnErr()
		}
	}
}

// IOSClient returns an iOS Client.
func (p *Project) IOSClient() Client {
	for _, c := range p.Clients {
		if c.Kind == IOSClientKind {
			return c
		}
	}
	return Client{}
}

// GraphQLDir returns the GraphQL folder that's located inside the client dir.
func (c *Client) GraphQLDir(projectName string) string {
	if c.Kind == IOSClientKind {
		return fmt.Sprintf("clients/ios/Sources/%sKit/GraphQL", strings.Title(projectName))
	}
	return ""
}

func (p *Project) writeTemplates(templates TemplateFiles) {
	for name, filename := range templates {
		p.writeTemplate(name, filename)
	}
}

func (p *Project) writeTemplate(name string, filename string) {

	writename := filepath.Join(p.dest, p.Name, filename)
	writepath := filepath.Dir(writename)

	// Check if file already exists
	if _, err := os.Stat(writename); !os.IsNotExist(err) {
		fmt.Printf("\tFile Already Exists: %s\n", filename)
		return
	}

	// Create directories if file path doesn't exist
	if _, err := os.Stat(writepath); os.IsNotExist(err) {
		if err := os.MkdirAll(writepath, 0755); err != nil {
			p.err = err
			return
		}
	}

	// Create file
	file, err := os.Create(writename)
	if err != nil {
		p.err = err
		return
	}

	// Read cached template data
	data, err := FSString(false, name)
	if err != nil {
		p.err = err
		return
	}

	// Create and execute template from data and project context
	tmpl, err := template.New(name).Funcs(
		template.FuncMap{
			"titlecase":                  strings.Title,
			"uppercase":                  strings.ToUpper,
			"lowercase":                  strings.ToLower,
			"joinArgsForSwift":           def.JoinArgsForSwift,
			"joinInterfacesForSwift":     def.JoinInterfacesForSwift,
			"joinArgsForGraphQL":         def.JoinArgsForGraphQL,
			"joinArgsForGraphQLVars":     def.JoinArgsForGraphQLVars,
			"graphQLScalarToSwiftScalar": def.GraphQLScalarToSwiftScalar,
			"swiftScalar":                def.ToSwiftScalar,
			"excludeSwiftScalars":        def.ExcludeSwiftScalars,
		},
	).Parse(data)
	if err != nil {
		p.err = err
		return
	}
	if err := tmpl.Execute(file, p); err != nil {
		p.err = err
		return
	}

	// Close file
	file.Close()
	fmt.Println("\tCreated File: ", writename)
}

func reverseDomain(in string) string {
	s := strings.Split(in, ".")
	reverse(s)
	return strings.Join(s, ".")
}

func reverse(in []string) {
	last := len(in) - 1
	for i := 0; i < len(in)/2; i++ {
		in[i], in[last-i] = in[last-i], in[i]
	}
}
