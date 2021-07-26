package gmaw_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jarcoal/httpmock"

	"github.com/roblox/terraform-provider-maas/pkg/api"
	"github.com/roblox/terraform-provider-maas/pkg/api/params"
	. "github.com/roblox/terraform-provider-maas/pkg/gmaw"
	"github.com/roblox/terraform-provider-maas/pkg/maas/entity"
	"github.com/roblox/terraform-provider-maas/test/helper"
)

func TestNewNetworkInterfaces(t *testing.T) {
	NewNetworkInterfaces(client)
}

func TestNetworkInterfaces(t *testing.T) {
	// Ensure the type implements the interface
	var _ api.NetworkInterfaces = (*NetworkInterfaces)(nil)

	// Create a new interfaces client to be used in the tests
	interfacesClient := NewNetworkInterfaces(client)

	// Register HTTPMock responders
	sid200, sid404 := "1234", "2345"
	url200 := fmt.Sprintf("/MAAS/api/2.0/nodes/%s/interfaces/", sid200)
	url404 := fmt.Sprintf("/MAAS/api/2.0/nodes/%s/interfaces/", sid404)

	var wants []entity.NetworkInterface
	if err := helper.TestdataFromJSON("maas/interfaces.json", &wants); err != nil {
		t.Fatal(err)
	}
	want := new(entity.NetworkInterface)
	if err := helper.TestdataFromJSON("maas/interface.json", want); err != nil {
		t.Fatal(err)
	}
	httpmock.RegisterResponder("GET", url200,
		httpmock.NewJsonResponderOrPanic(http.StatusOK, wants))
	httpmock.RegisterResponder("POST", url200,
		httpmock.NewJsonResponderOrPanic(http.StatusOK, want))
	httpmock.RegisterResponder("GET", url404,
		httpmock.NewStringResponder(http.StatusNotFound, "Not Found"))
	httpmock.RegisterResponder("POST", url404,
		httpmock.NewStringResponder(http.StatusNotFound, "Not Found"))

	t.Run("Get", func(t *testing.T) {
		t.Run("200", func(t *testing.T) {
			t.Parallel()
			got, err := interfacesClient.Get(sid200)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(wants, got, cmpopts.EquateEmpty()); diff != "" {
				t.Fatalf("json.Decode() mismatch (-want +got):\n%s", diff)
			}
		})
		t.Run("404", func(t *testing.T) {
			t.Parallel()
			got, err := interfacesClient.Get(sid404)
			if diff := cmp.Diff(([]entity.NetworkInterface{}), got, cmpopts.EquateEmpty()); diff != "" {
				t.Fatalf("json.Decode() mismatch (-want +got):\n%s", diff)
			}
			if err.Error() != "ServerError: 404 (Not Found)" {
				t.Fatal(err)
			}
		})
	})

	t.Run("CreateBond", func(t *testing.T) {
		t.Run("200", func(t *testing.T) {
			t.Parallel()
			got, err := interfacesClient.CreateBond(sid200, &params.NetworkInterfaceBond{})
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Fatalf("json.Decode() mismatch (-want +got):\n%s", diff)
			}
		})
		t.Run("404", func(t *testing.T) {
			t.Parallel()
			got, err := interfacesClient.CreateBond(sid404, &params.NetworkInterfaceBond{})
			if diff := cmp.Diff((&entity.NetworkInterface{}), got, cmpopts.EquateEmpty()); diff != "" {
				t.Fatalf("json.Decode() mismatch (-want +got):\n%s", diff)
			}
			if err.Error() != "ServerError: 404 (Not Found)" {
				t.Fatal(err)
			}
		})
	})

	t.Run("CreateBridge", func(t *testing.T) {
		t.Run("200", func(t *testing.T) {
			t.Parallel()
			got, err := interfacesClient.CreateBridge(sid200, &params.NetworkInterfaceBridge{})
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Fatalf("json.Decode() mismatch (-want +got):\n%s", diff)
			}
		})
		t.Run("404", func(t *testing.T) {
			t.Parallel()
			got, err := interfacesClient.CreateBridge(sid404, &params.NetworkInterfaceBridge{})
			if diff := cmp.Diff((&entity.NetworkInterface{}), got, cmpopts.EquateEmpty()); diff != "" {
				t.Fatalf("json.Decode() mismatch (-want +got):\n%s", diff)
			}
			if err.Error() != "ServerError: 404 (Not Found)" {
				t.Fatal(err)
			}
		})
	})

	t.Run("CreatePhysical", func(t *testing.T) {
		t.Run("200", func(t *testing.T) {
			t.Parallel()
			got, err := interfacesClient.CreatePhysical(sid200, &params.NetworkInterfacePhysical{})
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Fatalf("json.Decode() mismatch (-want +got):\n%s", diff)
			}
		})
		t.Run("404", func(t *testing.T) {
			t.Parallel()
			got, err := interfacesClient.CreatePhysical(sid404, &params.NetworkInterfacePhysical{})
			if diff := cmp.Diff((&entity.NetworkInterface{}), got, cmpopts.EquateEmpty()); diff != "" {
				t.Fatalf("json.Decode() mismatch (-want +got):\n%s", diff)
			}
			if err.Error() != "ServerError: 404 (Not Found)" {
				t.Fatal(err)
			}
		})
	})

	t.Run("CreateVLAN", func(t *testing.T) {
		t.Run("200", func(t *testing.T) {
			t.Parallel()
			got, err := interfacesClient.CreateVLAN(sid200, &params.NetworkInterfaceVLAN{})
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Fatalf("json.Decode() mismatch (-want +got):\n%s", diff)
			}
		})
		t.Run("404", func(t *testing.T) {
			t.Parallel()
			got, err := interfacesClient.CreateVLAN(sid404, &params.NetworkInterfaceVLAN{})
			if diff := cmp.Diff((&entity.NetworkInterface{}), got, cmpopts.EquateEmpty()); diff != "" {
				t.Fatalf("json.Decode() mismatch (-want +got):\n%s", diff)
			}
			if err.Error() != "ServerError: 404 (Not Found)" {
				t.Fatal(err)
			}
		})
	})
}
