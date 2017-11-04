/// This is a generated file, do not edit

import Foundation

extension Date {

    public static let rfc3339: DateFormatter = {
        let formatter = DateFormatter()
        formatter.calendar = Calendar(identifier: .iso8601)
        formatter.locale = Locale(identifier: "en_US_POSIX")
        formatter.timeZone = TimeZone(secondsFromGMT: 0)
        formatter.dateFormat = "yyyy-MM-dd'T'HH:mm:ss'Z"
        return formatter
    }()

    public var rfc3339: String {
        return Date.rfc3339.string(from: self)
    }

    public init?(rfc3339String str: String?) {
        guard let str = str else {
            return nil
        }
        guard let date = Date.rfc3339.date(from: str) else {
            return nil
        }
        self.init(timeInterval: 0, since: date)
    }
}
