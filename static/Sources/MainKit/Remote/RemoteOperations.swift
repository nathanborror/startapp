import Foundation

// Fragments
{{range .Schema.ObjectKind}}{{if not .IsEdge}}{{if not .IsConnection}}
fileprivate let {{.Name}}Fragment = """
fragment {{.Name}}Fragment on {{.Name}} { {{range .Fields|excludeFunctions|onlyScalars}}
    {{.Name}}{{end}}
}
"""
{{end}}{{end}}{{end}}
