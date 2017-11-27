import UIKit

enum StringAttributes {
    case title
    case regular
    case caption
    case label
    
    var fontSize: CGFloat {
        switch self {
        case .title:   return 17
        case .regular: return 15
        case .caption: return 11
        case .label:   return 12
        }
    }
    
    var lineHeight: CGFloat {
        switch self {
        case .title:   return 20
        case .regular: return 20
        case .caption: return 14
        case .label:   return 20
        }
    }
    
    var fontWeight: UIFont.Weight {
        switch self {
        case .title: return .semibold
        default:     return .regular
        }
    }
    
    var font: UIFont {
        return .systemFont(ofSize: self.fontSize, weight: self.fontWeight)
    }
    
    var color: UIColor {
        switch self {
        case .caption: return .regularLight
        case .label:   return .regularLight
        default:       return .regular
        }
    }
    
    var paragraphStyle: NSParagraphStyle {
        let style = NSMutableParagraphStyle()
        style.maximumLineHeight = self.lineHeight
        style.minimumLineHeight = self.lineHeight
        style.lineSpacing = (self.lineHeight - self.fontSize) / 2
        return style
    }
    
    var attributes: [NSAttributedStringKey: Any] {
        return [
            .font: self.font,
            .foregroundColor: self.color,
            .paragraphStyle: self.paragraphStyle,
        ]
    }
}


extension String {
    
    func attributedString(_ attributes: StringAttributes, color: UIColor? = nil) -> NSAttributedString {
        var attrs = attributes.attributes
        if let color = color {
            attrs[.foregroundColor] = color
        }
        return NSAttributedString(string: self, attributes: attrs)
    }
}

extension UIColor {
    
    static let tint         = #colorLiteral(red: 0, green: 0, blue: 0, alpha: 1)
    static let important    = #colorLiteral(red: 1, green: 0.4436733723, blue: 0.464625299, alpha: 1)
    static let empty        = #colorLiteral(red: 0.6593592763, green: 0.6987351775, blue: 0.7414059043, alpha: 1)
    static let emptyLight   = #colorLiteral(red: 0.9362974763, green: 0.9362974763, blue: 0.9362974763, alpha: 1)
    static let regular      = #colorLiteral(red: 0, green: 0, blue: 0, alpha: 1)
    static let regularLight = #colorLiteral(red: 0.6593592763, green: 0.6987351775, blue: 0.7414059043, alpha: 1)
}
