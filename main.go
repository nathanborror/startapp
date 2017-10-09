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
	liveFlag    = flag.Bool("liveassets", false, "Serve Static Assets from Disk relative to CWD")
)

type fileContext struct {
	Name       string
	HasKit     bool
	HasTests   bool
	HasUITests bool
	TeamID     string
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
	t, err := template.New(tmpl).Parse(path)
	if err != nil {
		return fmt.Errorf("Failed to parse template: %s", err)
	}
	if err := t.Execute(file, ctx); err != nil {
		return fmt.Errorf("Failed to execute template: %s", err)
	}
	file.Close()
	return nil
}
