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

// Arguments

type connectionArgs struct {
	First  *int32
	Before *string
}

// Encoders

func encodeID(in string) graphql.ID {
	return graphql.ID(in)
}

func encodeTime(in time.Time) string {
	return in.UTC().Format(time.RFC3339)
}

func encodeCursor(in string) string {
	return base64.StdEncoding.EncodeToString([]byte(in))
}

func encodePageInfo(in state.PageInfo) *pageInfoResolver {
	return &pageInfoResolver{
		startID:     in.StartID,
		endID:       in.EndID,
		hasNext:     in.HasNext,
		hasPrevious: in.HasPrevious,
	}
}

// Decoders

func decodeID(in graphql.ID) string {
	return string(in)
}

func decodeCursor(in string) string {
	str, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return ""
	}
	return string(str)
}

// Conversions

func int32toIntPtr(in *int32) *int {
	if in == nil {
		return nil
	}
	out := int(*in)
	return &out
}
