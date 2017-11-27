import UIKit

class FormInputCell: UITableViewCell {
    
    static let reuseIdentifier = "FormInputCell"
    
    var labelInset = UIEdgeInsets(top: 5, left: 16, bottom: 0, right: 16)
    var fieldInset = UIEdgeInsets(top: 0, left: 16, bottom: 10, right: 16)
    var value: String? {
        return fieldView.text
    }
    
    private(set) lazy var labelView: UILabel = {
        let view = UILabel()
        view.backgroundColor = UIColor(red: 0, green: 1, blue: 0, alpha: 0.1)
        self.contentView.addSubview(view)
        return view
    }()
    
    private(set) lazy var fieldView: UITextField = {
        let view = UITextField()
        view.clearButtonMode = .always
        view.backgroundColor = UIColor(red: 1, green: 0, blue: 0, alpha: 0.1)
        self.contentView.addSubview(view)
        return view
    }()
    
    override func sizeThatFits(_ size: CGSize) -> CGSize {
        let labelFit = size.insetBy(labelInset).infiniteHeight()
        let labelSize = labelView.sizeThatFits(labelFit).outsetBy(labelInset)
        
        let fieldFit = size.insetBy(fieldInset).infiniteHeight()
        let fieldSize = fieldView.sizeThatFits(fieldFit).outsetBy(fieldInset)
        
        return CGSize(width: max(labelSize.width, fieldSize.width), height: labelSize.height + fieldSize.height)
    }
    
    override func layoutSubviews() {
        super.layoutSubviews()
        
        let labelFit = bounds.size.insetBy(labelInset)
        let labelSize = labelView.sizeThatFits(labelFit)
        labelView.frame = CGRect(origin: labelInset.origin, size: labelSize)
        
        let fieldFit = bounds.size.insetBy(fieldInset)
        var fieldSize = fieldView.sizeThatFits(fieldFit)
        fieldSize.width = bounds.width - fieldInset.totalHorizontal
        fieldView.frame = CGRect(origin: fieldInset.origin.offsetY(labelView.frame), size: fieldSize)
    }
    
    override func prepareForReuse() {
        super.prepareForReuse()
        fieldView.text = nil
    }
}

