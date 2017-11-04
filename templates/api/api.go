package api

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/nathanborror/{{.Name}}/api/server"
	"github.com/nathanborror/{{.Name}}/state"
	graphql "github.com/neelance/graphql-go"
)

type Backends struct {
	State state.Stater
}

type rootResolver struct {
	*Backends
}

func Configure(schemaFilename string, backends Backends) {
	schema := Schema(schemaFilename, backends)
	registerHandlers(schema, backends)
}

func Schema(schemaFilename string, backends Backends) *graphql.Schema {
	root := rootResolver{&backends}
	return graphql.MustParseSchema(readFileContents(schemaFilename), &root)
}

func registerHandlers(schema *graphql.Schema, backends Backends) {
	http.Handle("/graphql", &server.Handler{Schema: schema, State: backends.State})
}

func readFileContents(filename string) string {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Failed to open schema file '%s': %v\n", filename, err)
		os.Exit(1)
	}
	return string(file)
}

// Encoders

func encodeID(i string) graphql.ID {
	return graphql.ID(i)
}

func encodeTime(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

// Decoders

func decodeID(i graphql.ID) string {
	return string(i)
}
