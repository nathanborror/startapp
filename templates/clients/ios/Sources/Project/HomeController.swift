import UIKit
import {{.Name|titlecase}}Kit

class HomeController: UIViewController {

    override func viewDidLoad() {
        super.viewDidLoad()
        title = "{{.Name|titlecase}}"
    }
}
