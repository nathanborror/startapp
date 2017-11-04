import Foundation

public struct State {
    public var authorization: Authorization
    public var account: Account
    public var error: Error?
}

public struct Authorization {
    public var token: String?
    public var stage: Stage
    public var error: Error?

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
    public var error: Error?
}
