godog-restful (BDD HTTP steps implementation for godog)
============

godog-restful is a golang module for Rest API testing. Contains reusable steps for [godog](https://github.com/DATA-DOG/godog) BDD (behaviour-driven development) tool. It can be used for testing REST APIs and interacting with JSON data over HTTP. 

Quick Start
-----------

First you must have Go installed on your machine, as instructed in [Installing Go](http://golang.org/doc/install.html).

### Installation

Download godog-restful using the `go get` tool:

    go get github.com/niminmptechonly/godog-restful/godoghttp
    
### Usage

In your godog test application spec file, add the following import statement to import the godog-restful specs:

    import . "github.com/niminmptechonly/godog-restful/godoghttp"
    
And make a call to the FeatureCOntext function in godog-restful like:

    func FeatureContext(s *godog.Suite) {
        HTTPFeatureContext(s)
        //s.Step(`^Test spec$`, iTestSpec)
    }
    
Set the BASE_URL to required hostname URL:
    
    export BASE_URL=http://localhost:8888
    
    
and you can use the steps in your feature files as:

    Feature: Test http
    
      Background:
        Given I disable security check

      Scenario: Test success response
        Given I set header "Content-Type" with value "application/json"
        Given I attach the file "/godog-restful/example.txt" as "example"
        Given I set headers with values:
          | Content-Type    | application/json |
          | Accept          | application/json |
        When I send a "POST" request to "/endpoint/test/" with body:
        """
          {
            "test" : [{
              "key" : "value",
              "key" : "value"
            }]
          }
        """
        When I send a "POST" request to "/endpoint/test/"
        When I send a "POST" request to "/endpoint/test/" with values:
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


    
    
godog-http steps
---------------
```
 Given:
    - I disable security check
    - I set header "{}" with value "{}"
    - I attach the file "{}" as "{}"
    - I set headers with values:
    
 When:
    - I send a {} request to "{}" with body:
    - I send a {} request to "{}" with values:
    - I send a {} request to "{}"

 Then:
    - the response code should be {}
    - the response should contain json:
    - the response should contain text:
    - the response header "{}" should be "{}"
    - the response body path "{}" should be "{}" (refer GJSON package for possible paths)
    - print response

```

Reference
---------

The response body path comparisons are based on the golang [GJSON] (https://github.com/tidwall/gjson) package
Examples  :
```
{
  "name": {"first": "Tom", "last": "Anderson"},
  "age":37,
  "children": ["Sara","Alex","Jack"],
  "fav.movie": "Deer Hunter",
  "friends": [
    {"first": "Dale", "last": "Murphy", "age": 44},
    {"first": "Roger", "last": "Craig", "age": 68},
    {"first": "Jane", "last": "Murphy", "age": 47}
  ]
}
```
```
"name.last"          >> "Anderson"
"age"                >> 37
"children"           >> ["Sara","Alex","Jack"]
"children.#"         >> 3
"children.1"         >> "Alex"
"child*.2"           >> "Jack"
"c?ildren.0"         >> "Sara"
"fav\.movie"         >> "Deer Hunter"
"friends.#.first"    >> ["Dale","Roger","Jane"]
"friends.1.last"     >> "Craig"
```
