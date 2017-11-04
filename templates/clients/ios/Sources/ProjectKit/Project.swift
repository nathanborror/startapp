import Foundation

public let {{.Name|titlecase}}DidChangeNotification = Notification.Name("{{.Name}}DidChangeNotification")

public struct {{.Name|titlecase}} {
    
    public static var state: State { return manager.state }
    public static let bundleIdentifier = "{{.IOSClient.BundleID}}"

    fileprivate(set) static var current: Service = Service(.localhost)
    
    static var remote: Remote { return current.remote }
    static var manager: Manager<State> { return current.manager }
    
    static func replace(service: Service) {
        current = service
    }
}

public class Service {

    enum Endpoint: String {
        case localhost  = "http://localhost:8080/graphql"
        case production = "http://{{.Domain}}/graphql"
    }

    internal var manager: Manager<State>
    internal var remote: Remote!

    init(_ endpoint: Endpoint, session: RemoteSession? = nil) {
        let initialState = Service.resumeState(from: endpoint.rawValue)

        self.manager = Manager(state: initialState)
        self.manager.subscribe(self, action: Service.managerStateChanged)

        if let session = session {
            self.remote = Remote(session: session, endpoint: endpoint.rawValue)
            return
        }
        let sessionConfig = URLSessionConfiguration.default
        let session = URLSession(configuration: sessionConfig, delegate: self.remote, delegateQueue: nil)
        self.remote = Remote(session: session, endpoint: endpoint.rawValue)
    }

    private func managerStateChanged(state: State) {
        let nc = NotificationCenter.default
        if Thread.isMainThread {
            nc.post(name: {{.Name|titlecase}}DidChangeNotification, object: nil)
            return
        }
        DispatchQueue.main.async {
            nc.post(name: {{.Name|titlecase}}DidChangeNotification, object: nil)
        }
    }

    private static func resumeState(from bucket: String) -> State {
        return State()
    }
}
