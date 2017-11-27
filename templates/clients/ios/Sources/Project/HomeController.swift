{{ $name := .IOSClient.Name }}
import UIKit
import {{$name}}Kit

class HomeController: AppViewController {

    override func viewDidLoad() {
        super.viewDidLoad()

        title = "{{$name}}"
    }

    override func stateDidChange(_ state: State) {
        super.stateDidChange(state)
    }
}