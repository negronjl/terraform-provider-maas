package tfschema

import (
	"reflect"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

// Resource wraps a Terraform resource's Schema.
// It provides functionality for parsing a struct{} into
// a map[string]*Schema for use in Terraform resources and providers.
type Resource map[string]*Schema

// NewResource transforms a struct{} into a map for a TF schema.Resource.
// The first return value of this function is the Schema property of the
// resource, and the function can handle recursion (ie nested structs).
// TODO: Make this type handle recursion so it can handle a struct within a struct
func NewResource(typ Endpoint) (Resource, error) {
	st := reflect.TypeOf(typ)
	resource := make(Resource)

	for i := 0; i < st.NumField(); i++ {
		f := st.Field(i)
		sch := NewSchema(&f)
		if err := sch.Parse(&f); err != nil {
			return nil, err
		}
		resource[strings.ToLower(f.Name)] = sch
	}

	mdType := typ.GetMetadata()
	mdResource := make(Resource)
	st = reflect.TypeOf(mdType)
	for i := 0; i < st.NumField(); i++ {
		f := st.Field(i)
		sch := NewSchema(&f)
		if err := sch.Parse(&f); err != nil {
			return nil, err
		}
		// The todo above will eliminate this
		mdResource[strings.ToLower(f.Name)] = sch
	}
	// This will work with recursion support
	// resource["metadata"] = mdResource

	return resource, nil
}

// TFResource returns a type that is compatible with Terraform's schema.Resource
func (s Resource) TFResource() map[string]*schema.Schema {
	res := make(map[string]*schema.Schema)

	for idx, val := range s {
		res[idx] = val.TFSchema()
	}
	return res
}
