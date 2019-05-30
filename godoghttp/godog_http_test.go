package godoghttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTests() apiFeature {
	api := apiFeature{}
	api.requestHeader = make(http.Header)
	api.responseHeader = make(http.Header)
	api.client = &http.Client{}
	return api
}

func TestSetHeaderWithValue(t *testing.T) {
	api := setupTests()
	err := api.setHeaderWithValue("header", "value")
	assert.Nil(t, err)
}

func TestResetResponse(t *testing.T) {
	api := setupTests()
	api.resetResponse()
	assert.Empty(t, api.responseBody, "response should be emtpy")
}

func TestSetHeadersWithValues(t *testing.T) {
	api := setupTests()
	dt := createDataTable()

	err := api.setHeadersWithValues(dt)
	assert.Nil(t, err)
}

func TestSetHeadersWithValuesInvalidDataTable(t *testing.T) {
	api := setupTests()
	dt := createDataTableWithOneCell()

	err := api.setHeadersWithValues(dt)
	assert.NotNil(t, err)
}

func TestResponseCodeShouldBeEqual(t *testing.T) {
	api := setupTests()
	api.responseCode = 200
	err := api.responseCodeShouldBe(200)
	assert.Nil(t, err)
}

func TestResponseCodeShouldBeNotEqual(t *testing.T) {
	api := setupTests()
	api.responseCode = 500
	err := api.responseCodeShouldBe(200)
	assert.NotNil(t, err)
	assert.Equal(t, "expected response code to be: 200, but actual is: 500, response message: ", err.Error())
}

func TestResponseShouldContainJSON(t *testing.T) {
	api := setupTests()
	testInput := `{"firstName":"John","lastName":"Dow"}`
	rawIn := json.RawMessage(testInput)
	actual, err := rawIn.MarshalJSON()
	if err != nil {
		fmt.Print("Error marshalling json")
	}
	api.responseBody = actual

	docString := new(gherkin.DocString)
	docString.Content = `{"firstName":"John","lastName":"Dow"}`

	err = api.responseShouldContainJSON(docString)
	assert.Nil(t, err)
}

func TestResponseShouldNotContainJSON(t *testing.T) {
	api := setupTests()
	testInput := `{"firstName":"John","lastName":"Dow"}`
	rawIn := json.RawMessage(testInput)
	actual, err := rawIn.MarshalJSON()
	if err != nil {
		fmt.Print("Error marshalling json")
	}
	api.responseBody = actual

	docString := new(gherkin.DocString)
	docString.Content = `{"firstName":"Ben","lastName":"Dow"}`

	err = api.responseShouldContainJSON(docString)
	assert.NotNil(t, err)
}

func TestResponseShouldContainJSONInvalidJSONInput(t *testing.T) {
	api := setupTests()
	testInput := `{"firstName":"John","lastName":"Dow"}`
	rawIn := json.RawMessage(testInput)
	actual, err := rawIn.MarshalJSON()
	if err != nil {
		fmt.Print("Error marshalling json")
	}
	api.responseBody = actual

	docString := new(gherkin.DocString)
	docString.Content = ""

	err = api.responseShouldContainJSON(docString)
	assert.NotNil(t, err)
}

func TestResponseShouldContainJSONInvalidJSONResponse(t *testing.T) {
	api := setupTests()
	testInput := ``
	rawIn := json.RawMessage(testInput)
	actual, err := rawIn.MarshalJSON()
	if err != nil {
		fmt.Print("Error marshalling json")
	}
	api.responseBody = actual

	docString := new(gherkin.DocString)
	docString.Content = `{"firstName":"John","lastName":"Dow"}`

	err = api.responseShouldContainJSON(docString)
	assert.NotNil(t, err)
}

func TestResponseShouldContainText(t *testing.T) {
	api := setupTests()

	testInput := `{"firstName":"John","lastName":"Dow"}`
	actual, err := json.Marshal(testInput)
	if err != nil {
		fmt.Print("Error marshalling json")
	}
	api.responseBody = actual

	docString := new(gherkin.DocString)
	docString.Content = "John"

	err = api.responseShouldContainText(docString)
	assert.Nil(t, err)
}

func TestResponseShouldNotContainText(t *testing.T) {
	api := setupTests()

	testInput := `{"firstName":"John","lastName":"Dow"}`
	actual, err := json.Marshal(testInput)
	if err != nil {
		fmt.Print("Error marshalling json")
	}
	api.responseBody = actual

	docString := new(gherkin.DocString)
	docString.Content = "Ben"

	err = api.responseShouldContainText(docString)
	assert.NotNil(t, err)
}

func TestResponseHeaderShouldBePresent(t *testing.T) {
	api := setupTests()
	api.responseHeader.Add("expected_header", "value")

	err := api.responseHeaderShouldBe("expected_header", "value")
	assert.Nil(t, err)
}

func TestResponseHeaderShouldNotBePresent(t *testing.T) {
	api := setupTests()
	api.responseHeader.Add("actual_header", "value")

	err := api.responseHeaderShouldBe("expected_header", "value")
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Errorf("header : %s not found in the response", "expected_header"), err)
}

func TestResponseHeaderShouldBePresentValueNotEqual(t *testing.T) {
	api := setupTests()
	api.responseHeader.Add("expected_header", "actual_value")

	err := api.responseHeaderShouldBe("expected_header", "expected_value")
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Errorf("expected header (%s) value: %s, not equal to actual response header value: %s",
		"expected_header", "expected_value", "actual_value"), err)
}

func TestResponseBodyPathShouldBePresent(t *testing.T) {
	api := setupTests()

	testInput := `{"name":{"firstName":"John","lastName":"Dow"}}`
	rawIn := json.RawMessage(testInput)
	actual, err := rawIn.MarshalJSON()
	if err != nil {
		fmt.Print("Error marshalling json")
	}

	api.responseBody = actual
	err = api.responseBodyPathShouldBe("name.firstName", "John")

	assert.Nil(t, err)
}

func TestResponseBodyPathShouldNotBePresent(t *testing.T) {
	api := setupTests()

	testInput := `{"name":{"firstName":"John","lastName":"Dow"}}`
	rawIn := json.RawMessage(testInput)
	actual, err := rawIn.MarshalJSON()
	if err != nil {
		fmt.Print("Error marshalling json")
	}

	api.responseBody = actual
	err = api.responseBodyPathShouldBe("name.middleName", "John")

	assert.NotNil(t, err)
	assert.Equal(t, fmt.Errorf("response body path : %s not found in the response", "name.middleName"), err)
}

func TestResponseBodyPathShouldBePresentValueNotEqual(t *testing.T) {
	api := setupTests()

	testInput := `{"name":{"firstName":"John","lastName":"Dow"}}`
	rawIn := json.RawMessage(testInput)
	actual, err := rawIn.MarshalJSON()
	if err != nil {
		fmt.Print("Error marshalling json")
	}

	api.responseBody = actual
	err = api.responseBodyPathShouldBe("name.firstName", "Ben")

	assert.NotNil(t, err)
	assert.Equal(t, fmt.Errorf("expected response value : %s not equal to the actual value : %s for the response "+
		"body path : %s", "Ben", "John", "name.firstName"), err)
}

func TestPrintResponse(t *testing.T) {
	api := setupTests()
	api.responseBody = []byte{}
	err := api.printResponse()
	assert.Nil(t, err)
}

func TestAttachTheFileAs(t *testing.T) {
	api := setupTests()

	filePath := "godog_http.go"
	fileName := "file"
	err := api.attachTheFileAs(filePath, fileName)
	assert.Nil(t, err)
}

func TestAttachTheFileAsFileNotExist(t *testing.T) {
	api := setupTests()

	filePath := "unknown.go"
	fileName := "file"
	err := api.attachTheFileAs(filePath, fileName)
	assert.NotNil(t, err)
}

func TestSendHTTPRequestOk(t *testing.T) {
	api := setupTests()
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.URL.String(), "/test")
		rw.Write([]byte(`OK`))
	}))

	api.client = server.Client()
	api.config = &EnvConfig{}
	api.config.BaseURL = server.URL
	err := api.sendHTTPRequest("POST", "/test", nil)

	assert.Nil(t, err)
	assert.Equal(t, 200, api.responseCode)
}

func TestSendHTTPRequestServerError(t *testing.T) {
	api := setupTests()
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.URL.String(), "/test")
		rw.WriteHeader(500)
		rw.Write([]byte(`OK`))
	}))

	api.client = server.Client()
	api.config = &EnvConfig{}
	api.config.BaseURL = server.URL
	err := api.sendHTTPRequest("POST", "/test", nil)

	assert.Nil(t, err)
	assert.Equal(t, 500, api.responseCode)
}

func TestSendHTTPRequestBaseURLNotSet(t *testing.T) {
	api := setupTests()
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.URL.String(), "/test")
		rw.Write([]byte(`OK`))
	}))

	api.client = server.Client()
	api.config = &EnvConfig{}
	err := api.sendHTTPRequest("POST", "/test", nil)

	assert.NotNil(t, err)
	assert.Equal(t, fmt.Errorf("base URL not set for endpoint : %s", "/test"), err)
}

func TestSendHTTPRequestErrorResponse(t *testing.T) {
	api := setupTests()
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.URL.String(), "/test")
		rw.Write([]byte(`OK`))
	}))

	api.client = server.Client()
	api.config = &EnvConfig{}
	api.config.BaseURL = "http://localhost:0000"
	err := api.sendHTTPRequest("POST", "/test", nil)

	assert.NotNil(t, err)
}

func TestSendHTTPRequestInvalidHTTPMethod(t *testing.T) {
	api := setupTests()
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.URL.String(), "/test")
		rw.Write(nil)
	}))

	api.client = server.Client()
	api.config = &EnvConfig{}
	api.config.BaseURL = server.URL
	err := api.sendHTTPRequest("POST", "/test", nil)

	assert.Nil(t, err)
	assert.Equal(t, 200, api.responseCode)
}

func TestSendARequestToWithBody(t *testing.T) {
	api := setupTests()

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.URL.String(), "/test")
		rw.Write([]byte(`OK`))
	}))

	api.client = server.Client()
	api.config = &EnvConfig{}
	api.config.BaseURL = server.URL

	docString := new(gherkin.DocString)
	docString.Content = `{"firstName":"John","lastName":"Dow"}`

	err := api.sendARequestToWithBody("POST", "/test", docString)

	assert.Nil(t, err)
}

func TestSendARequestToWithValues(t *testing.T) {
	api := setupTests()

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.URL.String(), "/test")
		rw.Write([]byte(`OK`))
	}))

	api.client = server.Client()
	api.config = &EnvConfig{}
	api.config.BaseURL = server.URL

	dataTable := createDataTable()

	err := api.sendARequestToWithValues("POST", "/test", dataTable)

	assert.Nil(t, err)
}

func TestSendARequestToWithValuesInvalidInput(t *testing.T) {
	api := setupTests()

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.URL.String(), "/test")
		rw.Write([]byte(`OK`))
	}))

	api.client = server.Client()
	api.config = &EnvConfig{}
	api.config.BaseURL = server.URL

	dataTable := createDataTableWithOneCell()

	err := api.sendARequestToWithValues("POST", "/test", dataTable)

	assert.NotNil(t, err)
}

func TestSendARequestTo(t *testing.T) {
	api := setupTests()

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.URL.String(), "/test")
		rw.Write([]byte(`OK`))
	}))

	api.client = server.Client()
	api.config = &EnvConfig{}
	api.config.BaseURL = server.URL
	api.requestBody = new(bytes.Buffer)

	err := api.sendARequestTo("POST", "/test")

	assert.Nil(t, err)
}

func TestHTTPFeatureContext(t *testing.T) {
	assert.NotPanics(t, func() {
		HTTPFeatureContext(&godog.Suite{})
	})
}

func TestDisableSecurityCheck(t *testing.T) {
	api := setupTests()
	api.disableSecurityCheck()
}

func createDataTable() *gherkin.DataTable {
	dt := new(gherkin.DataTable)
	dt.Type = "DataTable"

	tc := make([]*gherkin.TableCell, 2)
	tc[0] = new(gherkin.TableCell)
	tc[0].Value = "test"
	tc[1] = new(gherkin.TableCell)
	tc[1].Value = "value"

	tr := make([]*gherkin.TableRow, 1)
	tr[0] = new(gherkin.TableRow)
	tr[0].Cells = tc

	dt.Rows = tr

	return dt
}

func createDataTableWithOneCell() *gherkin.DataTable {
	dt := new(gherkin.DataTable)
	dt.Type = "DataTable"

	tc := make([]*gherkin.TableCell, 1)
	tc[0] = new(gherkin.TableCell)
	tc[0].Value = "test"

	tr := make([]*gherkin.TableRow, 1)
	tr[0] = new(gherkin.TableRow)
	tr[0].Cells = tc

	dt.Rows = tr

	return dt
}
