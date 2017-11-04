import Foundation

public enum {{.Name|titlecase}}Error: Error {

    case authenticationRequestBad(String)
    case authenticationRequired(String)
    case authenticationUnauthorized(String)

    case requestBadType(String)
    case requestBad(String)
    case requestNotUnderstood(String)
    case requestMethodNotAllowed(String)

    case notFound(String)
    case configurationFailure(String)
    case programmerFailure(String)
    case serviceUnavailable(String)
}
