package gmaw_test

import (
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

func TestNewVLANs(t *testing.T) {
	NewVLANs(client)
}

func TestVLANs(t *testing.T) {
	// Ensure the type implements the interface
	var _ api.VLANs = (*VLANs)(nil)

	// Create a new vlans client to be used in the tests
	vlansClient := NewVLANs(client)

	t.Run("Get", func(t *testing.T) {
		var vlans []entity.VLAN
		if err := helper.TestdataFromJSON("maas/vlans.json", &vlans); err != nil {
			t.Fatal(err)
		}
		t.Run("200", func(t *testing.T) {
			t.Parallel()
			httpmock.RegisterResponder("GET", "/MAAS/api/2.0/fabrics/123/vlans/",
				httpmock.NewJsonResponderOrPanic(http.StatusOK, vlans))
			res, err := vlansClient.Get(123)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(vlans, res, cmpopts.EquateEmpty()); diff != "" {
				t.Fatalf("json.Decode(VLANs) mismatch (-want +got):\n%s", diff)
			}
		})
		t.Run("404", func(t *testing.T) {
			t.Parallel()
			httpmock.RegisterResponder("GET", "/MAAS/api/2.0/fabrics/234/vlans/",
				httpmock.NewStringResponder(http.StatusNotFound, "Not Found"))

			got, err := vlansClient.Get(234)
			if diff := cmp.Diff(([]entity.VLAN{}), got, cmpopts.EquateEmpty()); diff != "" {
				t.Fatalf("json.Decode() mismatch (-want +got):\n%s", diff)
			}
			if err.Error() != "ServerError: 404 (Not Found)" {
				t.Fatal(err)
			}
		})
	})
	t.Run("Post", func(t *testing.T) {
		vlan := new(entity.VLAN)
		if err := helper.TestdataFromJSON("maas/vlan.json", vlan); err != nil {
			t.Fatal(err)
		}
		t.Run("200", func(t *testing.T) {
			t.Parallel()
			httpmock.RegisterResponder("POST", "/MAAS/api/2.0/fabrics/123/vlans/",
				httpmock.NewJsonResponderOrPanic(http.StatusOK, vlan))

			p := new(params.VLAN)
			got, err := vlansClient.Post(123, p)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(vlan, got, cmpopts.EquateEmpty()); diff != "" {
				t.Fatalf("json.Decode(VLANs) mismatch (-want +got):\n%s", diff)
			}
		})
		t.Run("404", func(t *testing.T) {
			t.Parallel()
			httpmock.RegisterResponder("POST", "/MAAS/api/2.0/fabrics/234/vlans/",
				httpmock.NewStringResponder(http.StatusNotFound, "Not Found"))

			p := new(params.VLAN)
			got, err := vlansClient.Post(234, p)
			if diff := cmp.Diff((&entity.VLAN{}), got, cmpopts.EquateEmpty()); diff != "" {
				t.Fatalf("json.Decode() mismatch (-want +got):\n%s", diff)
			}
			if err.Error() != "ServerError: 404 (Not Found)" {
				t.Fatal(err)
			}
		})
	})
}
