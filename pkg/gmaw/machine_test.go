package gmaw_test

import (
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/juju/gomaasapi"
	. "github.com/roblox/terraform-provider-maas/pkg/gmaw"
	"github.com/roblox/terraform-provider-maas/pkg/maas"
)

const apiURL string = "http://localhost:5240/MAAS"

var client *gomaasapi.MAASObject

func TestMain(m *testing.M) {
	httpmock.Activate()
	var err error
	client, err = GetClient(apiURL, "some:secret:key", "2.0")
	if err != nil {
		log.Fatal(err)
	}
	res := m.Run()
	httpmock.DeactivateAndReset()
	os.Exit(res)
}

func TestNewMachine(t *testing.T) {
	NewMachine(client)
}

func TestMachine_Get(t *testing.T) {
	tests := []testCase{
		{URL: "machines/42/", Verb: "GET", StatusCode: http.StatusOK, Response: "Machines!"}, // TODO Make a sample file
		{URL: "machines/43/", Verb: "GET", StatusCode: http.StatusNotFound, Response: "Not Found"},
	}

	machine := NewMachine(client)
	runTestCases(t, tests, func(tc testCase) ([]byte, error) {
		return machine.Get(tc.URL[9:])
	})
}

func TestMachine_Commission(t *testing.T) {
	tests := []testCase{
		{URL: "machines/42/?op=commission", Verb: "POST",
			StatusCode: http.StatusOK, Response: "Machines!"}, // TODO Make a sample file
		{URL: "machines/43/?op=commission", Verb: "POST", StatusCode: http.StatusNotFound, Response: "Not Found"},
	}

	machine := NewMachine(client)
	runTestCases(t, tests, func(tc testCase) ([]byte, error) {
		return machine.Commission(tc.URL[9:11], maas.MachineCommissionParams{})
	})
}

func TestMachine_Deploy(t *testing.T) {
	tests := []testCase{
		{URL: "machines/42/?op=deploy", Verb: "POST",
			StatusCode: http.StatusOK, Response: "Machines!"}, // TODO Make a sample file
		{URL: "machines/43/?op=deploy", Verb: "POST", StatusCode: http.StatusForbidden,
			Response: "The user does not have permission to deploy this machine."},
		{URL: "machines/44/?op=deploy", Verb: "POST", StatusCode: http.StatusNotFound, Response: "Not Found"},
		{URL: "machines/45/?op=deploy", Verb: "POST", StatusCode: http.StatusServiceUnavailable,
			Response: "MAAS attempted to allocate an IP address, and there were no IP addresses available on the relevant cluster interface."}, // nolint: lll
	}

	machine := NewMachine(client)
	runTestCases(t, tests, func(tc testCase) ([]byte, error) {
		return machine.Deploy(tc.URL[9:11], &maas.MachineDeployParams{})
	})
}

func TestMachine_Lock(t *testing.T) {
	tests := []testCase{
		{URL: "machines/42/?op=lock", Verb: "POST",
			StatusCode: http.StatusOK, Response: "Machines!"}, // TODO Make a sample file
		{URL: "machines/43/?op=lock", Verb: "POST", StatusCode: http.StatusForbidden,
			Response: "The user does not have permission to lock the machine."},
		{URL: "machines/44/?op=lock", Verb: "POST", StatusCode: http.StatusNotFound, Response: "Not Found"},
	}

	machine := NewMachine(client)
	runTestCases(t, tests, func(tc testCase) ([]byte, error) {
		return machine.Lock(tc.URL[9:11], "some-comment")
	})
}
