{{ $name := .IOSClient.Name }}
import XCTest
@testable import {{$name}}Kit

class Tests: XCTestCase {
    
    override func setUp() {
        super.setUp()
    }
    
    override func tearDown() {
        super.tearDown()
    }
    
    func testExample() {
    }
}
