# startapp

Creates a base Xcode project environment based on the way I like to work.
Adjust the static template folder to your favorite environment and run 
`go generate`. Any file or folder with the name `Main` will get replaced with 
the app name.

Template context:

```
{
    Name       string
    HasKit     bool
    HasTests   bool
    HasUITests bool
}
```

Dependencies:

- https://github.com/mjibson/esc
- https://github.com/yonaskolb/XcodeGen
