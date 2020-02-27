package gmaw_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"

	"github.com/juju/errors"

	"github.com/google/go-cmp/cmp"
	. "github.com/roblox/terraform-provider-maas/pkg/gmaw"
)

// The test cases for the endpoint types use this type in unit tests
type testCase struct {
	URL        string
	Verb       string
	StatusCode int
	Response   string
}

// The endpoint types' methods use this function to verify functionality in unit tests
func runTestCases(t *testing.T, tests []testCase, f func(testCase) ([]byte, error)) {
	defer httpmock.Reset()
	for _, testCase := range tests {
		tc := testCase
		t.Run(tc.URL, func(t *testing.T) {
			httpmock.RegisterResponder(tc.Verb, fmt.Sprintf("%s/api/2.0/%s", apiURL, tc.URL),
				httpmock.NewStringResponder(tc.StatusCode, tc.Response))

			if err := verifyResponse(tc, f); err != nil {
				t.Error(err)
			}
		})
	}
}

func verifyResponse(tc testCase, f func(testCase) ([]byte, error)) error {
	res, err := f(tc)

	if tc.StatusCode == http.StatusOK {
		if err != nil {
			return fmt.Errorf("Unexpected server error: %s", err)
		}

		// Verify the response from the server is unchanged
		if diff := cmp.Diff(tc.Response, string(res)); diff != "" {
			return fmt.Errorf(diff)
		}
	} else {
		// Verify the error is unchanged
		expected := fmt.Sprintf("ServerError: %d (%s)", tc.StatusCode, tc.Response)
		if diff := cmp.Diff(expected, err.Error()); diff != "" {
			return fmt.Errorf(diff)
		}

		// And that the response is empty
		// TODO: Is this right?
		if string(res) != "" {
			return fmt.Errorf("Unexpected response: %s", res)
		}
	}
	return nil
}

func TestGetClient(t *testing.T) {
	url := "http://1.2.3.4:5240/MAAS"
	version := "2.0"

	// Raise an error, just for fun
	if _, err := GetClient(url, "invalid:api_key", version); !errors.IsNotValid(err) {
		t.Fatalf("Expected error to be NotValid with an invalid API key")
	}

	// Okay, get a client this time
	res, err := GetClient(url, "secr3t:key:s3cret", version)
	if err != nil {
		t.Fatalf("Received error from gomaasapi: %s", err)
	}

	url = fmt.Sprintf("%s/api/%s/", url, version)
	if diff := cmp.Diff(res.URL().String(), url); diff != "" {
		t.Fatal(diff)
	}
}
