/*package tfschema provides helpers for working with Terraform's helper/schema.Schema.*/
package tfschema

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
)

// Schema wraps a Terraform Schema with additional functionality for parsing struct fields.
// Use this type to convert a struct field into a Terraform schema, using struct
// tags to handle the various schema attributes (eg required, computed).
type Schema struct {
	schema.Schema
	BoolTags   map[string]*bool
	StringTags map[string]*string
	FuncTags   map[string]interface{}
}

func NewSchema(f *reflect.StructField) *Schema {
	s := &Schema{
		Schema: schema.Schema{},
	}
	s.BoolTags = map[string]*bool{
		"required":  &s.Required,
		"optional":  &s.Optional,
		"forcenew":  &s.ForceNew,
		"computed":  &s.Computed,
		"sensitive": &s.Sensitive,
	}
	s.StringTags = map[string]*string{
		"description": &s.Description,
		"deprecated":  &s.Deprecated,
		"removed":     &s.Removed,
	}
	s.FuncTags = map[string]interface{}{
		"DefaultFunc":      &s.DefaultFunc,
		"DiffSuppressFunc": &s.DiffSuppressFunc,
		"ValidateFunc":     &s.ValidateFunc,
		"StateFunc":        &s.StateFunc,
	}
	return s
}

// Parse introspects the struct field to populate the helper.Schema.
// It will find the Terraform ValueType of the field, search for any
// struct tags, and set the default value, if there is one.
func (s *Schema) Parse(f *reflect.StructField) (err error) {
	if err = s.SetValueType(f.Type.Kind()); err != nil {
		return
	}

	if err = s.ParseTags(f.Tag); err != nil {
		return
	}

	if val, ok := f.Tag.Lookup("default"); ok {
		err = s.SetDefault(val)
	}
	return
}

// TFSchema returns the underlying Terraform schema
func (s *Schema) TFSchema() *schema.Schema {
	return &s.Schema
}

// SetValueType sets the schema's Terraform ValueType
func (s *Schema) SetValueType(t reflect.Kind) error {
	switch t {
	case reflect.String:
		s.Type = schema.TypeString
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		s.Type = schema.TypeInt
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		s.Type = schema.TypeInt
	case reflect.Bool:
		s.Type = schema.TypeBool
	case reflect.Float32, reflect.Float64:
		s.Type = schema.TypeFloat
	default:
		return fmt.Errorf("no ValueType equivalent of %s", t)
	}
	return nil
}

// SetDefault sets the schema's Default value.
func (s *Schema) SetDefault(val string) (err error) {
	switch s.Type {
	case schema.TypeString:
		s.Default = val
	case schema.TypeInt:
		s.Default, err = strconv.Atoi(val)
	case schema.TypeBool:
		s.Default = (val == "true")
	case schema.TypeFloat:
		s.Default, err = strconv.ParseFloat(val, 10)
	default:
		err = fmt.Errorf("default only valid for str, bool, num (received %v)", s.Type)
	}
	return err
}

// ParseTags sets values in the schema based on the values of struct tags.
func (s *Schema) ParseTags(tags reflect.StructTag) error {
	for tag, prop := range s.BoolTags {
		if val, ok := tags.Lookup(tag); ok {
			*prop = val == "true"
		}
	}
	for tag, prop := range s.StringTags {
		if val, ok := tags.Lookup(tag); ok {
			*prop = val
		}
	}
	return nil
}

/* func (s *Schema) AddHelpers() error {
	for tag := range s.FuncTags {
		if _, ok := tags.Lookup(strings.ToLower(tag)); !ok {
			continue
		}
		funcName := fmt.Sprintf("%s%s", f.Name, tag)
		sv := reflect.ValueOf(in)
		fn := sv.MethodByName(funcName)
		if fn.IsValid() {
			fv := reflect.ValueOf(funcTags[tag])
			fv.Set(reflect.Indirect(fn))
		} else {
			return fmt.Errorf("Cannot find method %s", funcName)
		}
	}
	return nil
} */
