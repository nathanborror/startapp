{{ $name := .IOSClient.Name }}
import Foundation

extension State {

    init() {
        self.account = Account()
        self.authorization = Authorization()
        self.error = nil
    }
}

extension Authorization {

    init() {
        self.token = nil
        self.stage = .disconnected
        self.error = nil
    }
}

extension Account {

    init() {
        self.id = ""
        self.name = ""
        self.email = ""
        self.created = Date()
        self.modified = Date()
        self.error = nil
    }

    mutating func apply(remote: Remote.Account?) {
        guard let remote = remote else { return }
        self.id = remote.id ?? self.id
        self.name = remote.name ?? self.name
        self.email = remote.email ?? self.email
        self.created = Date(rfc3339String: remote.created) ?? self.created
        self.modified = Date(rfc3339String: remote.modified) ?? self.modified
    }
}
