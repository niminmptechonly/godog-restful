Feature: Test http

  Background:
    Given I disable security check

  Scenario: Test success response
    Given I set header "Content-Type" with value "application/json"
    Given I attach the file "/godog-http/example.txt" as "example"
    Given I set headers with values:
      | Content-Type    | application/json |
      | Accept          | application/json |
    When I send a "POST" request to "/endpoint/test" with body:
    """
      {
	"hello" : "11"
}
    """
    When I send a "POST" request to "/endpoint/test"
    When I send a "POST" request to "/endpoint/test" with values:
      | test    | 1 |
      | test2   | 2 |
      | test3   | 3 |
    Then the response header "Content-Type" should be "application/json"
    Then the response code should be 200
    Then the response should contain json:
    """
    {
       "test":[
          {
             "key":"value",
             "key":"value"
          }
       ]
    }
    """
    Then the response should contain text:
    """
    test
    """
    Then the response body path "test.0.key" should be "value"
    And print response

