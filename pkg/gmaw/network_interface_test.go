package gmaw_test

// 403 on the tag ones
// 400 on the SetDefaultGateway

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/roblox/terraform-provider-maas/pkg/api/params"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jarcoal/httpmock"

	"github.com/roblox/terraform-provider-maas/pkg/api"
	. "github.com/roblox/terraform-provider-maas/pkg/gmaw"
	"github.com/roblox/terraform-provider-maas/pkg/maas/entity"
	"github.com/roblox/terraform-provider-maas/test/helper"
)

func TestNewNetworkInterface(t *testing.T) {
	NewNetworkInterface(client)
}

func TestNetworkInterface(t *testing.T) {
	// Ensure the type implements the interface
	var _ api.NetworkInterface = (*NetworkInterface)(nil)

	// Create a new interfaces client to be used in the tests
	interfaceClient := NewNetworkInterface(client)

	// Load test data
	want := new(entity.NetworkInterface)
	if err := helper.TestdataFromJSON("maas/interface.json", want); err != nil {
		t.Fatal(err)
	}

	// Register HTTPMock responders
	sid := "45tywk"
	ifc200, ifc400, ifc403, ifc404 := 123, 234, 345, 456
	url200 := fmt.Sprintf("/MAAS/api/2.0/nodes/%s/interfaces/%d/", sid, ifc200)
	url400 := fmt.Sprintf("/MAAS/api/2.0/nodes/%s/interfaces/%d/", sid, ifc400)
	url403 := fmt.Sprintf("/MAAS/api/2.0/nodes/%s/interfaces/%d/", sid, ifc403)
	url404 := fmt.Sprintf("/MAAS/api/2.0/nodes/%s/interfaces/%d/", sid, ifc404)
	httpmock.RegisterResponder("GET", url200,
		httpmock.NewJsonResponderOrPanic(http.StatusOK, want))
	httpmock.RegisterResponder("POST", url200,
		httpmock.NewJsonResponderOrPanic(http.StatusOK, want))
	httpmock.RegisterResponder("PUT", url200,
		httpmock.NewJsonResponderOrPanic(http.StatusOK, want))
	httpmock.RegisterResponder("DELETE", url200,
		httpmock.NewStringResponder(http.StatusNoContent, ""))
	// The API documentation doesn't specify the exact messages for the next two responders.
	httpmock.RegisterResponder("POST", url400,
		httpmock.NewStringResponder(http.StatusBadRequest, "Bad Request"))
	httpmock.RegisterResponder("POST", url403,
		httpmock.NewStringResponder(http.StatusForbidden, "Forbidden"))
	httpmock.RegisterResponder("GET", url404,
		httpmock.NewStringResponder(http.StatusNotFound, "Not Found"))
	httpmock.RegisterResponder("POST", url404,
		httpmock.NewStringResponder(http.StatusNotFound, "Not Found"))
	httpmock.RegisterResponder("PUT", url404,
		httpmock.NewStringResponder(http.StatusNotFound, "Not Found"))
	httpmock.RegisterResponder("DELETE", url404,
		httpmock.NewStringResponder(http.StatusNotFound, "Not Found"))

	t.Run("Delete", func(t *testing.T) {
		t.Run("204", func(t *testing.T) {
			t.Parallel()
			err := interfaceClient.Delete(sid, ifc200)
			if err != nil {
				t.Fatal(err)
			}
		})
		t.Run("404", func(t *testing.T) {
			t.Parallel()
			err := interfaceClient.Delete(sid, ifc404)
			if err.Error() != "ServerError: 404 (Not Found)" {
				t.Fatal(err)
			}
		})
	})

	t.Run("Get", func(t *testing.T) {
		t.Run("200", func(t *testing.T) {
			t.Parallel()
			got, err := interfaceClient.Get(sid, ifc200)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Fatalf("json.Decode() mismatch (-want +got):\n%s", diff)
			}
		})
		t.Run("404", func(t *testing.T) {
			t.Parallel()
			got, err := interfaceClient.Get(sid, ifc404)
			if diff := cmp.Diff((&entity.NetworkInterface{}), got, cmpopts.EquateEmpty()); diff != "" {
				t.Fatalf("json.Decode() mismatch (-want +got):\n%s", diff)
			}
			if err.Error() != "ServerError: 404 (Not Found)" {
				t.Fatal(err)
			}
		})
	})

	t.Run("AddTag", func(t *testing.T) {
		t.Run("200", func(t *testing.T) {
			t.Parallel()
			got, err := interfaceClient.AddTag(sid, ifc200, "some_tag")
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Fatalf("json.Decode() mismatch (-want +got):\n%s", diff)
			}
		})
		t.Run("403", func(t *testing.T) {
			t.Parallel()
			got, err := interfaceClient.AddTag(sid, ifc403, "some_tag")
			if diff := cmp.Diff((&entity.NetworkInterface{}), got, cmpopts.EquateEmpty()); diff != "" {
				t.Fatalf("json.Decode() mismatch (-want +got):\n%s", diff)
			}
			if err.Error() != "ServerError: 403 (Forbidden)" {
				t.Fatal(err)
			}
		})
		t.Run("404", func(t *testing.T) {
			t.Parallel()
			got, err := interfaceClient.AddTag(sid, ifc404, "some_tag")
			if diff := cmp.Diff((&entity.NetworkInterface{}), got, cmpopts.EquateEmpty()); diff != "" {
				t.Fatalf("json.Decode() mismatch (-want +got):\n%s", diff)
			}
			if err.Error() != "ServerError: 404 (Not Found)" {
				t.Fatal(err)
			}
		})
	})

	t.Run("Disconnect", func(t *testing.T) {
		t.Run("200", func(t *testing.T) {
			t.Parallel()
			got, err := interfaceClient.Disconnect(sid, ifc200)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Fatalf("json.Decode() mismatch (-want +got):\n%s", diff)
			}
		})
		t.Run("404", func(t *testing.T) {
			t.Parallel()
			got, err := interfaceClient.Disconnect(sid, ifc404)
			if diff := cmp.Diff((&entity.NetworkInterface{}), got, cmpopts.EquateEmpty()); diff != "" {
				t.Fatalf("json.Decode() mismatch (-want +got):\n%s", diff)
			}
			if err.Error() != "ServerError: 404 (Not Found)" {
				t.Fatal(err)
			}
		})
	})

	t.Run("LinkSubnet", func(t *testing.T) {
		t.Run("200", func(t *testing.T) {
			t.Parallel()
			got, err := interfaceClient.LinkSubnet(sid, ifc200, &params.NetworkInterfaceLink{})
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Fatalf("json.Decode() mismatch (-want +got):\n%s", diff)
			}
		})
		t.Run("404", func(t *testing.T) {
			t.Parallel()
			got, err := interfaceClient.LinkSubnet(sid, ifc404, &params.NetworkInterfaceLink{})
			if diff := cmp.Diff((&entity.NetworkInterface{}), got, cmpopts.EquateEmpty()); diff != "" {
				t.Fatalf("json.Decode() mismatch (-want +got):\n%s", diff)
			}
			if err.Error() != "ServerError: 404 (Not Found)" {
				t.Fatal(err)
			}
		})
	})

	t.Run("RemoveTag", func(t *testing.T) {
		t.Run("200", func(t *testing.T) {
			t.Parallel()
			got, err := interfaceClient.RemoveTag(sid, ifc200, "some_tag")
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Fatalf("json.Decode() mismatch (-want +got):\n%s", diff)
			}
		})
		t.Run("403", func(t *testing.T) {
			t.Parallel()
			got, err := interfaceClient.RemoveTag(sid, ifc403, "some_tag")
			if diff := cmp.Diff((&entity.NetworkInterface{}), got, cmpopts.EquateEmpty()); diff != "" {
				t.Fatalf("json.Decode() mismatch (-want +got):\n%s", diff)
			}
			if err.Error() != "ServerError: 403 (Forbidden)" {
				t.Fatal(err)
			}
		})
		t.Run("404", func(t *testing.T) {
			t.Parallel()
			got, err := interfaceClient.RemoveTag(sid, ifc404, "some_tag")
			if diff := cmp.Diff((&entity.NetworkInterface{}), got, cmpopts.EquateEmpty()); diff != "" {
				t.Fatalf("json.Decode() mismatch (-want +got):\n%s", diff)
			}
			if err.Error() != "ServerError: 404 (Not Found)" {
				t.Fatal(err)
			}
		})
	})

	t.Run("SetDefaultGateway", func(t *testing.T) {
		t.Run("200", func(t *testing.T) {
			t.Parallel()
			got, err := interfaceClient.SetDefaultGateway(sid, ifc200, 12)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Fatalf("json.Decode() mismatch (-want +got):\n%s", diff)
			}
		})
		t.Run("400", func(t *testing.T) {
			t.Parallel()
			got, err := interfaceClient.SetDefaultGateway(sid, ifc400, 13)
			if diff := cmp.Diff((&entity.NetworkInterface{}), got, cmpopts.EquateEmpty()); diff != "" {
				t.Fatalf("json.Decode() mismatch (-want +got):\n%s", diff)
			}
			if err.Error() != "ServerError: 400 (Bad Request)" {
				t.Fatal(err)
			}
		})
		t.Run("404", func(t *testing.T) {
			t.Parallel()
			got, err := interfaceClient.SetDefaultGateway(sid, ifc404, 14)
			if diff := cmp.Diff((&entity.NetworkInterface{}), got, cmpopts.EquateEmpty()); diff != "" {
				t.Fatalf("json.Decode() mismatch (-want +got):\n%s", diff)
			}
			if err.Error() != "ServerError: 404 (Not Found)" {
				t.Fatal(err)
			}
		})
	})

	t.Run("UnlinkSubnet", func(t *testing.T) {
		t.Run("200", func(t *testing.T) {
			t.Parallel()
			got, err := interfaceClient.UnlinkSubnet(sid, ifc200, 2)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Fatalf("json.Decode() mismatch (-want +got):\n%s", diff)
			}
		})
		t.Run("404", func(t *testing.T) {
			t.Parallel()
			got, err := interfaceClient.UnlinkSubnet(sid, ifc404, 3)
			if diff := cmp.Diff((&entity.NetworkInterface{}), got, cmpopts.EquateEmpty()); diff != "" {
				t.Fatalf("json.Decode() mismatch (-want +got):\n%s", diff)
			}
			if err.Error() != "ServerError: 404 (Not Found)" {
				t.Fatal(err)
			}
		})
	})

	t.Run("Put", func(t *testing.T) {
		t.Run("200", func(t *testing.T) {
			t.Parallel()
			got, err := interfaceClient.Put(sid, ifc200, &params.NetworkInterfacePhysical{})
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Fatalf("json.Decode() mismatch (-want +got):\n%s", diff)
			}
		})
		t.Run("404", func(t *testing.T) {
			t.Parallel()
			got, err := interfaceClient.Put(sid, ifc404, &params.NetworkInterfacePhysical{})
			if diff := cmp.Diff((&entity.NetworkInterface{}), got, cmpopts.EquateEmpty()); diff != "" {
				t.Fatalf("json.Decode() mismatch (-want +got):\n%s", diff)
			}
			if err.Error() != "ServerError: 404 (Not Found)" {
				t.Fatal(err)
			}
		})
	})
}
