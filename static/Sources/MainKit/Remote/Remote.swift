/// This is a generated file, do not edit

import Foundation

struct Remote {
}

// Operations

struct RemoteQuery: Codable {
    let query: String
    let variables: [String: String]?
}

struct RemoteMutation<T: Codable>: Codable {
    let query: String
    let variables: RemoteInput<T>?

    init(query: String, input: T? = nil) {
        self.query = query
        guard let input = input else {
            self.variables = nil
            return
        }
        self.variables = RemoteInput(input: input)
    }
}

struct RemoteInput<T: Codable>: Codable {
    let input: T
}

protocol RemoteResponse: Codable {
    var errors: [RemoteError]? { get }
}

enum RemoteResult<Value> {
    case success(Value)
    case progress(Float)
    case failure(RemoteError)

    func onSuccess(_ handler: (Value) -> ()) {
        if case .success(let v) = self { handler(v) }
    }
    func onProgress(_ handler: (Float) -> ()) {
        if case .progress(let p) = self { handler(p) }
    }
    func onFailure(_ handler: (RemoteError) -> ()) {
        if case .failure(let e) = self { handler(e) }
    }
}

// Errors

struct RemoteError: Error, Codable {
    struct Location: Codable {
        let line: Int
        let column: Int
    }
    let message: String
    let locations: [Location]?

    var localizedDescription: String {
        return message
    }
}

extension RemoteError {

    init(description: String) {
        self.message = description
        self.locations = nil
    }

    init(error: DecodingError) {
        self.message = error.localizedDescription
        self.locations = nil
    }
}
