import XCTest
@testable import {{.Name|titlecase}}Kit

class Tests: XCTestCase {
    
    var session: MockRemoteSession!

    override func setUp() {
        super.setUp()
        session = MockRemoteSession()
        {{.Name|titlecase}}.replace(service: Service(.localhost, session: session))
    }
    
    override func tearDown() {
        super.tearDown()
        session = nil
    }
    
    func testServicePing() {
        session.nextData = "{}".data(using: .utf8)
        {{.Name|titlecase}}.ping()
    }
}