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
	"github.com/roblox/terraform-provider-maas/pkg/maas/entity/subnet"
	"github.com/roblox/terraform-provider-maas/test/helper"
)

func TestNewSubnet(t *testing.T) {
	NewSubnet(client)
}

func TestSubnet(t *testing.T) {
	// Ensure the type implements the interface
	var _ api.Subnet = (*Subnet)(nil)

	// Create a new subnet client to be used in the tests
	subnetClient := NewSubnet(client)

	t.Run("Delete", func(t *testing.T) {
		t.Run("204", func(t *testing.T) {
			t.Parallel()
			httpmock.RegisterResponder("DELETE", "/MAAS/api/2.0/subnets/1/",
				httpmock.NewStringResponder(http.StatusNoContent, ""))
			if err := subnetClient.Delete(1); err != nil {
				t.Fatal(err)
			}
		})
		t.Run("404", func(t *testing.T) {
			t.Parallel()
			httpmock.RegisterResponder("DELETE", "/MAAS/api/2.0/subnets/0/",
				httpmock.NewStringResponder(http.StatusNotFound, "Not Found"))
			if err := subnetClient.Delete(0); err.Error() != "ServerError: 404 (Not Found)" {
				t.Fatal(err)
			}
		})
	})

	t.Run("Get", func(t *testing.T) {
		t.Parallel()
		want := new(entity.Subnet)
		if err := helper.TestdataFromJSON("maas/subnet.json", want); err != nil {
			t.Fatal(err)
		}
		httpmock.RegisterResponder("GET", "/MAAS/api/2.0/subnets/123/",
			httpmock.NewJsonResponderOrPanic(http.StatusOK, want))
		got, err := subnetClient.Get(123)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(want, got, cmpopts.EquateEmpty()); diff != "" {
			t.Fatalf("json.Decode() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("GetIPAddresses", func(t *testing.T) {
		t.Run("200", func(t *testing.T) {
			t.Parallel()
			var want []subnet.IPAddress
			if err := helper.TestdataFromJSON("maas/subnets/ipaddresses.json", &want); err != nil {
				t.Fatal(err)
			}
			httpmock.RegisterResponder("GET", "/MAAS/api/2.0/subnets/2/",
				httpmock.NewJsonResponderOrPanic(http.StatusOK, &want))
			got, err := subnetClient.GetIPAddresses(2, true, true)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Fatalf("json.Decode() mismatch (-want +got):\n%s", diff)
			}
		})
		t.Run("404", func(t *testing.T) {
			t.Parallel()
			httpmock.RegisterResponder("GET", "/MAAS/api/2.0/subnets/3/",
				httpmock.NewStringResponder(http.StatusNotFound, "Not Found"))
			res, err := subnetClient.GetIPAddresses(3, true, true)
			if res != nil {
				t.Fatal("Expected result to be nil")
			}
			if err.Error() != "ServerError: 404 (Not Found)" {
				t.Fatal(err)
			}
		})
	})

	t.Run("GetReservedIPRanges", func(t *testing.T) {
		t.Run("200", func(t *testing.T) {
			t.Parallel()
			var want []subnet.ReservedIPRange
			if err := helper.TestdataFromJSON("maas/subnets/reservedipranges.json", &want); err != nil {
				t.Fatal(err)
			}
			httpmock.RegisterResponder("GET", "/MAAS/api/2.0/subnets/4/",
				httpmock.NewJsonResponderOrPanic(http.StatusOK, &want))
			got, err := subnetClient.GetReservedIPRanges(4)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Fatalf("json.Decode() mismatch (-want +got):\n%s", diff)
			}
		})
		t.Run("404", func(t *testing.T) {
			t.Parallel()
			httpmock.RegisterResponder("GET", "/MAAS/api/2.0/subnets/5/",
				httpmock.NewStringResponder(http.StatusNotFound, "Not Found"))
			res, err := subnetClient.GetReservedIPRanges(5)
			if res != nil {
				t.Fatal("Expected result to be nil")
			}
			if err.Error() != "ServerError: 404 (Not Found)" {
				t.Fatal(err)
			}
		})
	})

	t.Run("GetStatistics", func(t *testing.T) {
		t.Run("200", func(t *testing.T) {
			t.Parallel()
			want := new(subnet.Statistics)
			if err := helper.TestdataFromJSON("maas/subnets/statistics.json", want); err != nil {
				t.Fatal(err)
			}
			httpmock.RegisterResponder("GET", "/MAAS/api/2.0/subnets/6/",
				httpmock.NewJsonResponderOrPanic(http.StatusOK, want))
			got, err := subnetClient.GetStatistics(6, false, false)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Fatalf("json.Decode() mismatch (-want +got):\n%s", diff)
			}
		})
		t.Run("404", func(t *testing.T) {
			t.Parallel()
			httpmock.RegisterResponder("GET", "/MAAS/api/2.0/subnets/7/",
				httpmock.NewStringResponder(http.StatusNotFound, "Not Found"))
			got, err := subnetClient.GetStatistics(7, false, false)
			if diff := cmp.Diff((&subnet.Statistics{}), got, cmpopts.EquateEmpty()); diff != "" {
				t.Fatalf("json.Decode() mismatch (-want +got):\n%s", diff)
			}
			if err.Error() != "ServerError: 404 (Not Found)" {
				t.Fatal(err)
			}
		})
	})

	t.Run("GetUnreservedIPRanges", func(t *testing.T) {
		t.Run("200", func(t *testing.T) {
			t.Parallel()
			var want []subnet.IPRange
			if err := helper.TestdataFromJSON("maas/subnets/ipranges.json", &want); err != nil {
				t.Fatal(err)
			}
			httpmock.RegisterResponder("GET", "/MAAS/api/2.0/subnets/8/",
				httpmock.NewJsonResponderOrPanic(http.StatusOK, &want))
			got, err := subnetClient.GetUnreservedIPRanges(8)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Fatalf("json.Decode() mismatch (-want +got):\n%s", diff)
			}
		})
		t.Run("404", func(t *testing.T) {
			t.Parallel()
			httpmock.RegisterResponder("GET", "/MAAS/api/2.0/subnets/9/",
				httpmock.NewStringResponder(http.StatusNotFound, "Not Found"))
			res, err := subnetClient.GetUnreservedIPRanges(9)
			if res != nil {
				t.Fatal("Expected result to be nil")
			}
			if err.Error() != "ServerError: 404 (Not Found)" {
				t.Fatal(err)
			}
		})
	})

	t.Run("Put", func(t *testing.T) {
		t.Run("200", func(t *testing.T) {
			t.Parallel()
			want := new(entity.Subnet)
			if err := helper.TestdataFromJSON("maas/subnet.json", want); err != nil {
				t.Fatal(err)
			}
			httpmock.RegisterResponder("PUT", "/MAAS/api/2.0/subnets/10/",
				httpmock.NewJsonResponderOrPanic(http.StatusOK, want))
			res, err := subnetClient.Put(10, &params.Subnet{})
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(want, res, cmpopts.EquateEmpty()); diff != "" {
				t.Fatalf("json.Decode() mismatch (-want +got):\n%s", diff)
			}
		})
		t.Run("404", func(t *testing.T) {
			t.Parallel()
			httpmock.RegisterResponder("PUT", "/MAAS/api/2.0/subnets/11/",
				httpmock.NewStringResponder(http.StatusNotFound, "Not Found"))
			got, err := subnetClient.Put(11, &params.Subnet{})
			if diff := cmp.Diff((&entity.Subnet{}), got, cmpopts.EquateEmpty()); diff != "" {
				t.Fatalf("json.Decode() mismatch (-want +got):\n%s", diff)
			}
			if err.Error() != "ServerError: 404 (Not Found)" {
				t.Fatal(err)
			}
		})
	})
}
