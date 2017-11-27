{{ $name := .IOSClient.Name }}
import Foundation

public struct State {
    public var authorization: Authorization
    public var account: Account
    public var error: {{$name}}Error?
}

public struct Authorization {
    public var token: String?
    public var stage: Stage
    public var error: {{$name}}Error?

    public enum Stage {
        case connected
        case connecting
        case registering
        case disconnected
    }
}

public struct Account {
    public var id: String
    public var name: String
    public var email: String
    public var created: Date
    public var modified: Date
    public var error: {{$name}}Error?
}
