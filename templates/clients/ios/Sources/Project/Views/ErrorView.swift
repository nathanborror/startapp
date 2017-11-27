import UIKit

class ErrorView: UIView {

    var error: Error? {
        didSet { handleErrorChange() }
    }

    var inset = UIEdgeInsets(top: 10, left: 16, bottom: 10, right: 16)

    private(set) lazy var labelView: UILabel = {
        let view = UILabel()
        view.numberOfLines = 0
        self.addSubview(view)
        return view
    }()
    
    override func sizeThatFits(_ size: CGSize) -> CGSize {
        let labelFit = size.insetBy(inset).infiniteHeight()
        return labelView.sizeThatFits(labelFit).outsetBy(inset)
    }
    
    override func layoutSubviews() {
        super.layoutSubviews()
        labelView.frame = bounds.insetBy(inset)
    }

    func handleErrorChange() {
        guard let error = error else {
            labelView.attributedText = nil
            return
        }
        labelView.attributedText = error.localizedDescription.attributedString(.caption, color: .important)
    }
}
