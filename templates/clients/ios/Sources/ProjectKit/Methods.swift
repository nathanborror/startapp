import Foundation

extension {{.Name|titlecase}} {

    public static func ping() {
        remote.ping { result in
            result.onSuccess {
                print($0)
            }
            result.onFailure {
                print($0)
            }
        }
    }

    // Example

    public static func viewer() {
        remote.viewer(token: nil) { result in
            result.onSuccess {
                print($0)
            }
            result.onFailure {
                print($0)
            }
        }
    }
}