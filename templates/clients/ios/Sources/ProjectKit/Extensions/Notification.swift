import Foundation

extension NotificationCenter {
    
    @discardableResult
    public static func subscribe(to name: Notification.Name, using: @escaping (Notification) -> Void) -> NSObjectProtocol {
        return self.default.addObserver(forName: name, object: nil, queue: nil, using: using)
    }
    
    public static func unsubscribe(_ observer: Any, from name: Notification.Name) {
        self.default.removeObserver(observer, name: name, object: nil)
    }
    
    public static func unsubscribe(_ observer: Any) {
        self.default.removeObserver(observer)
    }
}
