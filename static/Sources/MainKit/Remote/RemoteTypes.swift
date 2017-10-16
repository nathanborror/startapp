/// This is a generated file, do not edit

import Foundation

// Interfaces
{{range .Schema.InterfaceKind}}
protocol {{.Name}} { {{range .Fields}}
    {{.SwiftProtocolFieldWithPrefix "Remote"}}{{end}}
}{{end}}

extension Remote { // Scalars
    {{range .Schema.ScalarKind}}
    typealias {{.Name}} = {{.Name|cast}}
    {{end}}
}

extension Remote { // Unions
    {{range .Schema.UnionKind}}
    enum {{.Name}} { {{range .PossibleTypes}}
        {{.SwiftCase}}{{end}}
    }
    {{end}}
}

extension Remote { // Objects
    {{range .Schema.ObjectKind}}
    struct {{.Name}}: {{.Interfaces|joinInterfaces}} { {{range .Fields|excludeFunctions}}
        {{.SwiftField}}{{end}}
    }
    {{end}}
}

extension Remote { // Inputs
    {{range .Schema.InputKind}}
    struct {{.Name}}: Codable { {{range .InputFields}}
        {{.SwiftField}}{{end}}
    }
    {{end}}
}

extension Remote { // Payloads
    {{range .Schema.PayloadKind}}
    struct {{.Name}}: Codable { {{range .Fields}}
        {{.SwiftField}}{{end}}
    }
    {{end}}
}
