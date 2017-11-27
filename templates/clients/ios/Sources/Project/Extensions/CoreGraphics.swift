import UIKit

extension UIEdgeInsets {

    /// Returns the sum left and right insets.
    var totalHorizontal: CGFloat {
        return left + right
    }

    /// Returns the sum top and bottom insets.
    var totalVertical: CGFloat {
        return top + bottom
    }

    /// Returns a point based on the inset's top and left values.
    var origin: CGPoint {
        return CGPoint(x: left, y: top)
    }
}

extension CGSize {

    /// Returns a size inset by the given insets.
    func insetBy(_ inset: UIEdgeInsets) -> CGSize {
        return CGSize(width: width - inset.totalHorizontal, height: height - inset.totalVertical)
    }

    /// Returns a size with its height set to infinity.
    func infiniteHeight() -> CGSize {
        return CGSize(width: width, height: .greatestFiniteMagnitude)
    }

    /// Returns a size that adds the given insets.
    func outsetBy(_ inset: UIEdgeInsets) -> CGSize {
        return CGSize(width: width + inset.totalHorizontal, height: height + inset.totalVertical)
    }
}

extension CGRect {

    /// Returns a rect with the given insets subtracted from it's size while maintaining it's center point.
    func insetBy(_ inset: UIEdgeInsets) -> CGRect {
        return CGRect(x: minX + inset.left, y: minY + inset.top, width: width - inset.totalHorizontal, height: height - inset.totalVertical)
    }
}

extension CGPoint {

    /// Returns a CGPoint offset by the given CGRect's maxY position.
    func offsetY(_ rect: CGRect) -> CGPoint {
        return CGPoint(x: x, y: y + rect.maxY)
    }

    /// Returns a CGPoint offset by the given CGRect's maxX position.
    func offsetX(_ rect: CGRect) -> CGPoint {
        return CGPoint(x: x + rect.maxX, y: y)
    }
}
