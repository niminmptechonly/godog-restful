package godoghttp

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/caarlos0/env"
	"github.com/tidwall/gjson"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// EnvConfig : env config for godog-http feature spec
type EnvConfig struct {
	BaseURL string `env:"BASE_URL" envDefault:"http://localhost:8888"`
}

type apiFeature struct {
	config         *EnvConfig
	client         *http.Client
	responseCode   int
	responseBody   []byte
	responseHeader http.Header
	requestHeader  http.Header
	requestCookie  http.Cookie
	requestBody    *bytes.Buffer
}

func (api *apiFeature) resetResponse() {
	api.responseBody = []byte{}
	api.responseCode = 0
}

func (api *apiFeature) setHeaderWithValue(header, value string) error {
	api.requestHeader.Set(header, value)
	return nil
}

func (api *apiFeature) setHeadersWithValues(headersTable *gherkin.DataTable) (err error) {
	if len(headersTable.Rows[0].Cells) != 2 {
		err = fmt.Errorf("expected two columns for event table row, got: %d", len(headersTable.Rows[0].Cells))
		return
	}

	for _, row := range headersTable.Rows {
		api.requestHeader.Add(row.Cells[0].Value, row.Cells[1].Value)
	}

	return
}

func (api *apiFeature) attachTheFileAs(filePath, fileName string) (err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}

	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(fileName, filepath.Base(filePath))
	if err != nil {
		return
	}
	_, err = io.Copy(part, file)

	err = writer.Close()
	if err != nil {
		return
	}

	api.requestBody = body
	api.requestHeader.Set("Content-Type", writer.FormDataContentType())

	return
}

func (api *apiFeature) sendARequestToWithBody(method, endpoint string, requestJSON *gherkin.DocString) (err error) {
	data := strings.NewReader(requestJSON.Content)
	return api.sendHTTPRequest(method, endpoint, data)
}

func (api *apiFeature) sendARequestToWithValues(method, endpoint string, valuesTable *gherkin.DataTable) (err error) {
	variablesMap := make(map[string]string)
	if len(valuesTable.Rows[0].Cells) != 2 {
		err = fmt.Errorf("expected two columns for event table row, got: %d", len(valuesTable.Rows[0].Cells))
		return
	}
	for _, row := range valuesTable.Rows {
		variablesMap[row.Cells[0].Value] = row.Cells[1].Value
	}
	jsonString, err := json.Marshal(variablesMap)
	if err != nil {
		return
	}

	return api.sendHTTPRequest(method, endpoint, strings.NewReader(string(jsonString)))
}

func (api *apiFeature) sendARequestTo(method, endpoint string) (err error) {
	return api.sendHTTPRequest(method, endpoint, api.requestBody)
}

func (api *apiFeature) sendHTTPRequest(method, endpoint string, body io.Reader) (err error) {

	if len(strings.TrimSpace(api.config.BaseURL)) == 0 {
		err = fmt.Errorf("base URL not set for endpoint : %s", endpoint)
		return
	}

	endpoint = api.config.BaseURL + endpoint
	req, err := http.NewRequest(method, endpoint, body)

	if err != nil {
		return
	}

	if api.requestHeader != nil {
		req.Header = api.requestHeader
	}

	response, err := api.client.Do(req)
	if err != nil {
		return
	}

	if response == nil {
		err = fmt.Errorf("unable to get response from endpoint: %s", endpoint)
		return
	}

	defer response.Body.Close()

	api.responseCode = response.StatusCode
	api.responseHeader = response.Header

	api.responseBody, err = ioutil.ReadAll(response.Body)
	if err != nil {
		err = fmt.Errorf("error error reading response body.  error:  %s", err)
	}

	return
}

func (api *apiFeature) responseCodeShouldBe(responseCode int) error {
	if responseCode != api.responseCode {
		if api.responseCode >= 400 {
			return fmt.Errorf("expected response code to be: %d, but actual is: %d, response message: %s",
				responseCode, api.responseCode, string(api.responseBody))
		}
		return fmt.Errorf("expected response code to be: %d, but actual is: %d", responseCode, api.responseCode)
	}
	return nil
}

func (api *apiFeature) responseShouldContainJSON(responseJSON *gherkin.DocString) (err error) {
	var expected, actual []byte
	var data interface{}
	var actualData interface{}

	if err = json.Unmarshal([]byte(responseJSON.Content), &data); err != nil {
		err = fmt.Errorf("error unmarshalling expected data: %s", err)
	}
	if expected, err = json.Marshal(data); err != nil {
		err = fmt.Errorf("error marshalling expected data: %s", err)
	}

	actual = api.responseBody
	if err := json.Unmarshal(actual, &actualData); err != nil {
		err = fmt.Errorf("error unmarshalling actual data: %s", err)
	}
	if actual, err = json.Marshal(actualData); err != nil {
		err = fmt.Errorf("error marshalling actual data: %s", err)
	}
	if string(actual) != string(expected) {
		err = fmt.Errorf("expected json %s, does not match actual: %s", string(expected), string(actual))
	}
	return
}

func (api *apiFeature) responseShouldContainText(responseText *gherkin.DocString) (err error) {

	var actual []byte
	actual = api.responseBody

	if !strings.Contains(string(actual), responseText.Content) {
		err = fmt.Errorf("expected text %s, not found in actual response: %s", string(responseText.Content), string(actual))
	}

	return
}

func (api *apiFeature) responseHeaderShouldBe(header, value string) (err error) {
	headerValue := api.responseHeader.Get(header)
	if headerValue == "" && len(headerValue) == 0 {
		err = fmt.Errorf("header : %s not found in the response", header)
		return
	}
	if !strings.EqualFold(headerValue, value) {
		err = fmt.Errorf("expected header (%s) value: %s, not equal to actual response header value: %s",
			header, value, headerValue)
		return
	}

	return
}

func (api *apiFeature) responseBodyPathShouldBe(jsonPath, value string) (err error) {
	actual := gjson.Get(string(api.responseBody), jsonPath)
	if !actual.Exists() {
		err = fmt.Errorf("response body path : %s not found in the response", jsonPath)
		return
	}
	if !strings.EqualFold(actual.String(), value) {
		err = fmt.Errorf("expected response value : %s not equal to the actual value : %s for the response "+
			"body path : %s", value, actual.String(), jsonPath)
	}

	return
}

func (api *apiFeature) disableSecurityCheck() error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	api.client.Transport = tr
	return nil
}

func (api *apiFeature) printResponse() error {
	fmt.Println(string(api.responseBody))
	return nil
}

// HTTPFeatureContext : FeatureContext for godog-http steps
func HTTPFeatureContext(s *godog.Suite) {
	api := &apiFeature{}

	s.BeforeScenario(func(interface{}) {
		api.resetResponse()

		var envConfig EnvConfig
		env.Parse(&envConfig)
		api.config = &envConfig

		api.requestHeader = make(http.Header)
		api.requestBody = new(bytes.Buffer)
		api.client = &http.Client{}
	})

	s.Step(`^I disable security check$`, api.disableSecurityCheck)
	s.Step(`^I set header "([^"]*)" with value "([^"]*)"$`, api.setHeaderWithValue)
	s.Step(`^I set headers with values:$`, api.setHeadersWithValues)
	s.Step(`^I attach the file "([^"]*)" as "([^"]*)"$`, api.attachTheFileAs)
	s.Step(`^I send a "(GET|POST|PUT|DELETE)" request to "([^"]*)" with body:$`, api.sendARequestToWithBody)
	s.Step(`^I send a "(GET|POST|PUT|DELETE)" request to "([^"]*)" with values:$`, api.sendARequestToWithValues)
	s.Step(`^I send a "(GET|POST|PUT|DELETE)" request to "([^"]*)"$`, api.sendARequestTo)
	s.Step(`^the response code should be (\d+)$`, api.responseCodeShouldBe)
	s.Step(`^the response should contain json:$`, api.responseShouldContainJSON)
	s.Step(`^the response should contain text:$`, api.responseShouldContainText)
	s.Step(`^the response header "([^"]*)" should be "([^"]*)"$`, api.responseHeaderShouldBe)
	s.Step(`^the response body path "([^"]*)" should be "([^"]*)"$`, api.responseBodyPathShouldBe)
	s.Step(`^print response$`, api.printResponse)
}
