package provider

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform/helper/schema"
)

// ToSchemaMap transforms a struct{} into a map for a TF schema.Resource.
// The first return value of this function is the Schema property of the
// resource, and the function can handle recursion (ie nested structs).
func ToSchemaMap(in interface{}) (map[string]*schema.Schema, error) {
	st := reflect.TypeOf(in)
	sch := make(map[string]*schema.Schema)

	for i := 0; i < st.NumField(); i++ {
		f := st.Field(i)
		res, err := toSchema(f, in)
		if err != nil {
			return nil, err
		}
		sch[f.Name] = res
	}
	return sch, nil
}

// toSchema introspects a struct field.
// The second parameter is the struct itself
func toSchema(f reflect.StructField, in interface{}) (*schema.Schema, error) {
	s := schema.Schema{}

	// Figure out the schema.ValueType from the Go type
	switch f.Type.Kind() {
	case reflect.String:
		s.Type = schema.TypeString
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		s.Type = schema.TypeInt
	case reflect.Bool:
		s.Type = schema.TypeBool
	case reflect.Float32, reflect.Float64:
		s.Type = schema.TypeFloat
	}

	// Add any attributes defined in struct tags
	if val, ok := f.Tag.Lookup("required"); ok {
		if val == "true" {
			s.Required = true
		} else {
			s.Required = false
		}
	}
	if val, ok := f.Tag.Lookup("optional"); ok {
		if val == "true" {
			s.Optional = true
		} else {
			s.Optional = false
		}
	}
	if val, ok := f.Tag.Lookup("forcenew"); ok {
		if val == "true" {
			s.ForceNew = true
		} else {
			s.ForceNew = false
		}
	}
	if val, ok := f.Tag.Lookup("computed"); ok {
		if val == "true" {
			s.Computed = true
		} else {
			s.Computed = false
		}
	}
	if _, ok := f.Tag.Lookup("statefunc"); ok {
		sv := reflect.ValueOf(in)
		fn := sv.MethodByName(fmt.Sprintf("%sStateFunc", f.Name))
		if fn != (reflect.Value{}) {
			fi := fn.Interface()
			s.StateFunc = fi.(schema.SchemaStateFunc)
		} else {
			return &s, fmt.Errorf("Cannot find method %sStateFunc", f.Name)
		}
	}

	return &s, nil
}
