package gmaw_test

import (
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jarcoal/httpmock"

	"github.com/roblox/terraform-provider-maas/pkg/api"
	. "github.com/roblox/terraform-provider-maas/pkg/gmaw"
)

func TestNewMAASServer(t *testing.T) {
	NewMAASServer(client)
}

func TestMAASServer(t *testing.T) {
	// Ensure the type implements the interface
	var _ api.MAASServer = (*MAASServer)(nil)

	// Create a new client to be used in the tests
	maasClient := NewMAASServer(client)

	t.Run("Get", func(t *testing.T) {
		t.Parallel()
		want := "the_value"
		httpmock.RegisterResponder("GET", "/MAAS/api/2.0/maas/",
			httpmock.NewStringResponder(http.StatusOK, want))
		got, err := maasClient.Get("the_key")
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(want, got, cmpopts.EquateEmpty()); diff != "" {
			t.Fatalf("Returned value mismatch (-want +got):\n%s", diff)
		}
	})
	t.Run("Post", func(t *testing.T) {
		t.Parallel()
		httpmock.RegisterResponder("POST", "/MAAS/api/2.0/maas/",
			httpmock.NewStringResponder(http.StatusOK, "OK"))

		err := maasClient.Post("key", "value")
		if err != nil {
			t.Fatal(err)
		}
	})
}
