{{ $name := .IOSClient.Name }}
import UIKit
import {{$name}}Kit

@UIApplicationMain
class AppDelegate: UIResponder, UIApplicationDelegate {

    var window: UIWindow?

    func application(_ application: UIApplication, didFinishLaunchingWithOptions launchOptions: [UIApplicationLaunchOptionsKey : Any]? = nil) -> Bool {
        window = UIWindow(frame: UIScreen.main.bounds)
        window?.backgroundColor = .white
        window?.rootViewController = AppController()
        window?.makeKeyAndVisible()
        return true
    }
}

final class AppController: UINavigationController {

    convenience init() {
        self.init(rootViewController: HomeController())

        // Listen for State changes
        NotificationCenter.subscribe(to: {{$name}}DidChangeNotification) { [weak self] notif in
            self?.handleStateChange({{$name}}.state)
        }
    }

    override func viewDidAppear(_ animated: Bool) {
        super.viewDidAppear(animated)
        handleStateChange({{$name}}.state)
    }

    func handleStateChange(_ state: State) {

        // Interests
        let authStage = state.authorization.stage

        // Require account authorization
        guard authStage == .connected else {
            if presentedViewController is AuthController {
                return // Already showing auth controller
            }
            visibleViewController?.present(AuthController(), animated: false, completion: nil)
            return
        }

        // Dismiss auth when connected
        if presentedViewController is AuthController {
            presentedViewController?.dismiss(animated: true, completion: nil)
        }
    }
}

class AppViewController: UIViewController {
    let errorView = ErrorView()

    deinit {
        NotificationCenter.unsubscribe(self, from: {{$name}}DidChangeNotification)
    }

    override func viewDidLoad() {
        super.viewDidLoad()
        NotificationCenter.subscribe(to: {{$name}}DidChangeNotification) { [weak self] notif in
            self?.stateDidChange({{$name}}.state)
        }
    }

    func stateDidChange(_ state: State) {}
}

class AppTableViewController: UITableViewController {
    let errorView = ErrorView()

    deinit {
        NotificationCenter.unsubscribe(self, from: {{$name}}DidChangeNotification)
    }

    override func viewDidLoad() {
        super.viewDidLoad()
        NotificationCenter.subscribe(to: {{$name}}DidChangeNotification) { [weak self] notif in
            self?.stateDidChange({{$name}}.state)
        }
    }

    func stateDidChange(_ state: State) {}
}
