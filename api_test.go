package apitest

import (
	"testing"
)

func Test1(t *testing.T) {
	tcs, err := Load("testdata.csv", ',', true)
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range tcs {
		tc.T = t
		tc.RunCheckStatusCode()
	}

}
