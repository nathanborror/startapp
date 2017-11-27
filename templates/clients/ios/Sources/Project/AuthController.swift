{{ $name := .IOSClient.Name }}
import UIKit

final class AuthController: UINavigationController {

    convenience init() {
        self.init(rootViewController: ConnectController(style: .grouped))
    }
}
