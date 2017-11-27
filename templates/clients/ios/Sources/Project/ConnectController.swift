{{ $name := .IOSClient.Name }}
import UIKit
import {{$name}}Kit

final class ConnectController: AppTableViewController {

    lazy private(set) var emailCell: FormInputCell = {
        let cell = FormInputCell()
        cell.labelView.attributedText = "Email".attributedString(.label)
        cell.fieldView.attributedPlaceholder = "Your Email Address".attributedString(.title, color: .empty)
        cell.fieldView.delegate = self
        cell.fieldView.keyboardType = .emailAddress
        cell.fieldView.autocapitalizationType = .none
        cell.fieldView.returnKeyType = .next
        cell.selectionStyle = .none
        return cell
    }()

    lazy private(set) var passwordCell: FormInputCell = {
        let cell = FormInputCell()
        cell.labelView.attributedText = "Password".attributedString(.label)
        cell.fieldView.attributedPlaceholder = "Your Password".attributedString(.title, color: .empty)
        cell.fieldView.delegate = self
        cell.fieldView.isSecureTextEntry = true
        cell.fieldView.autocapitalizationType = .none
        cell.fieldView.returnKeyType = .done
        cell.selectionStyle = .none
        return cell
    }()

    override func viewDidLoad() {
        super.viewDidLoad()
        
        title = "Connect"

        // Offer Registration option
        let registerButton = UIBarButtonItem(title: "Register", style: .plain, target: self, action: #selector(handleRegister))
        navigationItem.rightBarButtonItem = registerButton
    }

    override func stateDidChange(_ state: State) {
        super.stateDidChange(state)

        // Interests
        let authError = state.authorization.error

        errorView.error = authError
        tableView.reloadData()
    }

    func handleSubmit() {
        guard let email = emailCell.value else {
            emailCell.fieldView.becomeFirstResponder()
            return
        }
        guard let password = passwordCell.value else {
            passwordCell.fieldView.becomeFirstResponder()
            return
        }
        {{$name}}.connect(email: email, password: password)
    }

    @objc func handleRegister() {
        let controller = RegisterController(style: .grouped)
        show(controller, sender: self)
    }

    override func tableView(_ tableView: UITableView, numberOfRowsInSection section: Int) -> Int {
        return 2
    }

    override func tableView(_ tableView: UITableView, cellForRowAt indexPath: IndexPath) -> UITableViewCell {
        switch indexPath.row {
        case 0:  return emailCell
        default: return passwordCell
        }
    }

    override func tableView(_ tableView: UITableView, viewForFooterInSection section: Int) -> UIView? {
        return errorView
    }

    override func tableView(_ tableView: UITableView, didSelectRowAt indexPath: IndexPath) {
        switch indexPath.row {
        case 0:  emailCell.fieldView.becomeFirstResponder()
        default: passwordCell.fieldView.becomeFirstResponder()
        }
    }
}

extension ConnectController: UITextFieldDelegate {

    func textFieldShouldReturn(_ textField: UITextField) -> Bool {
        switch textField {
        case emailCell.fieldView:
            passwordCell.fieldView.becomeFirstResponder()
        default:
            handleSubmit()
        }
        return false
    }
}
