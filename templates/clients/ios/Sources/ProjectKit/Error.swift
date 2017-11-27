{{ $name := .IOSClient.Name }}
import Foundation

public struct {{$name}}Error: Error, Codable {
    public struct Location: Codable {
        public let line: Int
        public let column: Int
    }
    public let message: String
    public let locations: [Location]?

    public var localizedDescription: String {
        return message
    }
}

extension {{$name}}Error {

    init(remote: RemoteError) {
        self.message = remote.message
        self.locations = remote.locations?.map { Location(line: $0.line, column: $0.column) }
    }
}
