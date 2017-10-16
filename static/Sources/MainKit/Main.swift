import Foundation

public class {{.Name}} {

    enum Endpoint: String {
        case localhost  = "http://localhost:8080/graphql"
        case production = "http://YOUR_DOMAIN/graphql"
        case testing    = "mock://"
    }

}