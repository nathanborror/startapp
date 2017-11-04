package main

//go:generate esc -o gen/templates.go -pkg gen -prefix "/" templates

import (
	"github.com/nathanborror/startapp/cmd"
)

func main() {
	cmd.Execute()
}
