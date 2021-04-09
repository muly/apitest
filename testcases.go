package apitest

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/gocarina/gocsv"
)

type (
	TestCase struct {
		Name             string `csv:"name"`
		RequestBody      string `csv:"request_body"`
		HttpVerb         string `csv:"http_verb"`
		Uri              string `csv:"uri"`
		WantStatusCode   int    `csv:"want_status_code"`
		WantResponseBody string `csv:"want_response_body"`
		SkipFlag         bool   `csv:"skip_flag"`
		//

		*testing.T
	}
)

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

func (tc TestCase) RunCheckStatusCode() (responseBody []byte) {
	if tc.SkipFlag {
		tc.Log("Skipped test case:")
		return
	}

	// prepare request
	body := strings.NewReader(tc.RequestBody)
	req, err := http.NewRequest(tc.HttpVerb, tc.Uri, body)
	if err != nil {
		tc.Error(err)
	}

	c := http.Client{}

	resp, err := c.Do(req)
	if err != nil {
		tc.Error(err)
	}

	if tc.WantStatusCode != resp.StatusCode {
		tc.Error(tc.Name, ": Status Code: wanted ", tc.WantStatusCode, " but got ", resp.StatusCode)
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		tc.Error(err)
	}

	return b
}
