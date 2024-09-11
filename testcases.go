package apitest

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/gocarina/gocsv"
)

type (
	TestCase struct {
		Name             string `csv:"name"`
		RequestBody      string `csv:"request_body"`
		RequestMethod    string `csv:"request_method"`
		RequestHeaders   map[string]string
		RequestPath      string `csv:"request_path"`
		WantStatusCode   int    `csv:"want_status_code"`
		WantResponseBody string `csv:"want_response_body"`
		SkipFlag         bool   `csv:"skip_flag"`
		SkipBodyCheck    bool   `csv:"skip_body_check"`

		gotStatusCode   int
		gotResponseBody string

		ValidateResponseFunc func(string) error

		*testing.T
	}
)

var baseURL = ""

func Init(url string) {
	baseURL = url
}

func Load(filePath string, delim rune, hasHeader bool) ([]TestCase, error) {
	testcases := []TestCase{}

	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	err = gocsv.UnmarshalFile(f, &testcases)
	if err != nil {
		return nil, err
	}

	return testcases, nil
}

func (tc *TestCase) RunCheck() error {
	if tc.SkipFlag {
		tc.Log("Skipped test case:", tc.Name)
		return nil
	}

	if err := tc.executeHttp(); err != nil {
		return fmt.Errorf("RunCheck() error: %v", err)
	}

	if tc.WantStatusCode != tc.gotStatusCode {
		return fmt.Errorf("RunCheck() Status Code mismatch: wanted %d but got %d", tc.WantStatusCode, tc.gotStatusCode)
	}

	if tc.SkipBodyCheck {
		return nil
	}

	// if testcase has the want response body, use it to validate response body
	if tc.WantResponseBody != "" && tc.WantResponseBody != tc.gotResponseBody {
		return fmt.Errorf("RunCheck() Response body mismatch: wanted %s but got %s", tc.WantResponseBody, tc.gotResponseBody)
	}

	// if validation function is provided, use it to validate response body
	if tc.ValidateResponseFunc != nil {
		if err := tc.ValidateResponseFunc(tc.gotResponseBody); err != nil {
			return fmt.Errorf("RunCheck() Response body mismatch: wanted %s but got %s", tc.WantResponseBody, tc.gotResponseBody)
		}
	}

	return nil
}

func (tc *TestCase) executeHttp() error {
	urlString := fmt.Sprintf("%s%s", baseURL, tc.RequestPath)
	if baseURL == "" {
		urlString = tc.RequestPath // i.e. the complete path is provided in the test case

	}

	req, err := http.NewRequest(tc.RequestMethod, urlString, bytes.NewBufferString(tc.RequestBody))
	if err != nil {
		return fmt.Errorf("executeHttp() http.NewRequest() error: %v", err)
	}

	for k, v := range tc.RequestHeaders {
		req.Header.Set(k, v)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("executeHttp() client.Do() error: %v", err)
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("executeHttp() io.ReadAll() error: %v", err)
	}

	tc.gotResponseBody = string(b)
	tc.gotStatusCode = resp.StatusCode

	return nil
}
