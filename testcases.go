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

		gotStatusCode   int
		gotResponseBody string
		// ValidateResponseFunc // TODO: later

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

func (tc *TestCase) RunCheckStatusCode() error {
	if tc.SkipFlag {
		tc.Log("Skipped test case:", tc.Name)
		return nil
	}

	err := tc.executeHttp()
	if err != nil {
		return fmt.Errorf("RunCheckStatusCode() executeHttp() error: %v", err)
	}

	if tc.WantStatusCode != tc.gotStatusCode {
		return fmt.Errorf("RunCheckStatusCode() Status Code: wanted %d but got %d", tc.WantStatusCode, tc.gotStatusCode)
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
