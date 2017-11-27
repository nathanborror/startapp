# startapp

Creates a base project environment based on the way I like to work.
Adjust the `/templates` folder to your needs and run `go generate`. Any file 
or folder with the name `Project` will be replaced with the app name.

Dependencies:

- https://github.com/mjibson/esc
- https://github.com/yonaskolb/XcodeGen

## Tasks

- [ ] Support recursive connection types in generated Swift code
- [ ] Ignore ID scalar in favor of neelance/go-graphql's implementation
- [ ] Fix mutation arguments in Swift
- [x] Fix camelCase on Swift mutation strings
- [ ] Swift connection edges aren't generating `edges: [Edge]` correctly