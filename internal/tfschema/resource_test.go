package tfschema_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/roblox/terraform-provider-maas/internal/provider"
	. "github.com/roblox/terraform-provider-maas/internal/tfschema"
)

type simpleStruct struct {
	ID       int    `computed:"true"`
	Name     string `required:"true" forcenew:"true"`
	Tries    uint8  `default:"3" optional:"true"`
	Password string `sensitive:"true"`
	DoStuff  bool   `default:"true" description:"Whether or not to do stuff"`
}

func (s simpleStruct) GetMetadata() interface{} {
	return struct{}{}
}

var simpleSchema = map[string]*schema.Schema{
	"id": &schema.Schema{
		Type:     schema.TypeInt,
		Computed: true,
	},
	"name": &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
	"tries": &schema.Schema{
		Type:     schema.TypeInt,
		Default:  3,
		Optional: true,
	},
	"password": &schema.Schema{
		Type:      schema.TypeString,
		Sensitive: true,
	},
	"dostuff": &schema.Schema{
		Type:        schema.TypeBool,
		Default:     true,
		Description: "Whether or not to do stuff",
	},
}

func TestNewSchemas(t *testing.T) {
	tests := []struct {
		name string
		got  provider.Endpoint
		want map[string]*schema.Schema
	}{
		{"simple", simpleStruct{}, simpleSchema},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NewResource(tc.got)
			if err != nil {
				t.Fatal(err)
			}
			diff := cmp.Diff(tc.want, got.TFResource())
			if diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
