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

		if err := tc.Execute(); err != nil{
			tc.Errorf(`testcase "%s" failed to Execute: %s`, tc.Name, err)
		}
		if err := tc.Check(); err != nil{
			tc.Errorf(`testcase "%s" failed with check: %s`, tc.Name, err)
		}
	}
}
