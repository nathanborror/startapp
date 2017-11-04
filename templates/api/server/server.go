package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/nathanborror/{{.Name}}/pkg/auth"
	"github.com/nathanborror/{{.Name}}/state"

	graphql "github.com/neelance/graphql-go"
)

type Handler struct {
	State  state.Stater
	Schema *graphql.Schema
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var args struct {
		Query         string                 `json:"query"`
		OperationName string                 `json:"operationName"`
		Variables     map[string]interface{} `json:"variables"`
	}

	// Read request body
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Ensure Content-Type header is application/json
	ctype := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ctype, "application/json") {
		http.Error(w, "RequestBad: Content-Type must be 'application/json'", 400)
		return
	}

	// Extract arguments
	if err := json.Unmarshal(data, &args); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Parse Session from Request
	session := auth.FromRequest(r)

	// Determine account ID
	accountID, _ := h.State.ReadAuthToken(session.Token)
	session.ID = accountID

	// Add Session to context
	ctx := session.Context(r.Context())

	// Execute
	resp := h.Schema.Exec(ctx, args.Query, args.OperationName, args.Variables)

	out, err := json.Marshal(resp)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(out)
}
