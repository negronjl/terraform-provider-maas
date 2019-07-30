package gmaw_test

import (
	"log"
	"os"
	"testing"

	"github.com/juju/gomaasapi"
	. "github.com/roblox/terraform-provider-maas/pkg/gmaw"
	"github.com/roblox/terraform-provider-maas/pkg/maas"
)

const apiURL string = "http://localhost:5240/MAAS"

var client *gomaasapi.MAASObject

func TestMain(m *testing.M) {
	var err error
	client, err = GetClient(apiURL, "some:secret:key", "2.0")
	if err != nil {
		log.Fatal(err)
	}
	res := m.Run()
	os.Exit(res)
}

func TestNewMachine(t *testing.T) {
	NewMachine(client)
}

func TestMachine_Get(t *testing.T) {
	tests := []testCase{
		{URL: "machines/42", StatusCode: 200, Response: "Machines!"}, // TODO Make a sample file
		{URL: "machines/43", StatusCode: 404, Response: "Not Found"},
	}

	machine := NewMachine(client)
	runTestCases(t, tests, func(tc testCase) ([]byte, error) {
		return machine.Get(tc.URL[9:])
	})
}

func TestMachine_Commission(t *testing.T) {
	tests := []testCase{
		{URL: "machines/42?op=commission", StatusCode: 200, Response: "Machines!"}, // TODO Make a sample file
		{URL: "machines/43?op=commission", StatusCode: 404, Response: "Not Found"},
	}

	machine := NewMachine(client)
	runTestCases(t, tests, func(tc testCase) ([]byte, error) {
		return machine.Commission(tc.URL[9:12], maas.MachineCommissionParams{})
	})
}

func TestMachine_Deploy(t *testing.T) {
	tests := []testCase{
		{URL: "machines/42?op=deploy", StatusCode: 200, Response: "Machines!"}, // TODO Make a sample file
		{URL: "machines/43?op=deploy", StatusCode: 403, Response: "The user does not have permission to deploy this machine."},
		{URL: "machines/44?op=deploy", StatusCode: 404, Response: "Not Found"},
		{URL: "machines/45?op=deploy", StatusCode: 503, Response: "MAAS attempted to allocate an IP address, and there were no IP addresses available on the relevant cluster interface."},
	}

	machine := NewMachine(client)
	runTestCases(t, tests, func(tc testCase) ([]byte, error) {
		return machine.Deploy(tc.URL[9:12], maas.MachineDeployParams{})
	})
}

func TestMachine_Lock(t *testing.T) {
	tests := []testCase{
		{URL: "machines/42?op=lock", StatusCode: 200, Response: "Machines!"}, // TODO Make a sample file
		{URL: "machines/43?op=lock", StatusCode: 403, Response: "The user does not have permission to lock the machine."},
		{URL: "machines/44?op=lock", StatusCode: 404, Response: "Not Found"},
	}

	machine := NewMachine(client)
	runTestCases(t, tests, func(tc testCase) ([]byte, error) {
		return machine.Lock(tc.URL[9:12], "some-comment")
	})
}
