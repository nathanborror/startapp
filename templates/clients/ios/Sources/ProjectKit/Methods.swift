{{ $name := .IOSClient.Name }}
import Foundation

extension {{$name}} {

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