{{ $name := .IOSClient.Name }}
import XCTest
@testable import {{$name}}Kit

class Tests: XCTestCase {
    
    var session: MockRemoteSession!

    override func setUp() {
        super.setUp()
        session = MockRemoteSession()
        {{$name}}.replace(service: Service(.localhost, session: session))
    }
    
    override func tearDown() {
        super.tearDown()
        session = nil
    }
}