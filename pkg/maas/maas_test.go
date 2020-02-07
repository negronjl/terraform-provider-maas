package maas_test

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"

	. "github.com/roblox/terraform-provider-maas/pkg/maas"
)

type emptyStruct struct{}

type simpleStruct struct {
	Name       string
	ID         int
	unexported float64
	Unsigned   uint8
}

type structWithTags struct {
	PascalCase  string `json:"snake_case"`
	ManyTags    string `json:"many_tags" sql:"-" computed:"true"`
	NotIncluded string `json:"-"`
}

type structWithArrays struct {
	Names []string
	IDs   []string
}

var (
	empty  = emptyStruct{}
	simple = simpleStruct{
		Name:       "Brian",
		ID:         31,      // nolint: gomnd
		unexported: 638.427, // nolint: gomnd
		Unsigned:   42,      // nolint: gomnd
	}
	simpleWithEmpties = simpleStruct{
		Name:       "Frank",
		unexported: 12, // nolint: gomnd
	}
	tags = structWithTags{
		PascalCase:  "always",
		ManyTags:    "sure",
		NotIncluded: "never",
	}
	arrays = structWithArrays{
		Names: []string{"Robin", "Little John", "Mervyn"},
		IDs:   []string{"Hood", "?", "Sheriff of Rottingham"},
	}
	arraysSortOf = structWithArrays{
		Names: []string{"None"},
		IDs:   []string{},
	}
)

var (
	emptyVals             = url.Values{}
	simpleVals url.Values = map[string][]string{
		"name":       []string{"Brian"},
		"id":         []string{"31"},
		"unexported": []string{"638.427"},
		"unsigned":   []string{"42"},
	}
	simpleWEVals url.Values = map[string][]string{
		"name":       []string{"Frank"},
		"id":         []string{"0"},
		"unexported": []string{"12"},
		"unsigned":   []string{"0"},
	}
	tagsVals url.Values = map[string][]string{
		"snake_case": []string{"always"},
		"many_tags":  []string{"sure"},
	}
	arraysVals url.Values = map[string][]string{
		"names": []string{"Robin", "Little John", "Mervyn"},
		"ids":   []string{"Hood", "?", "Sheriff of Rottingham"},
	}
	arraysSOVals url.Values = map[string][]string{
		"names": []string{"None"},
	}
)

func TestToQSP(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		want  url.Values
	}{
		{name: "empty", input: empty, want: emptyVals},
		{name: "simple", input: simple, want: simpleVals},
		{name: "simple with empties", input: simpleWithEmpties, want: simpleWEVals},
		{name: "json", input: tags, want: tagsVals},
		{name: "arrays", input: arrays, want: arraysVals},
		{name: "tricky arrays", input: arraysSortOf, want: arraysSOVals},
	}

	for _, testCase := range tests {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			got := ToQSP(tc.input)
			diff := cmp.Diff(tc.want, got)
			if diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func ExampleToQSP() {
	type whatev struct {
		Foo int `json:"bar"`
		Baz string
		bat int
	}
	data := whatev{
		Foo: 42, // nolint: gomnd
		Baz: "hello",
		bat: 10, // nolint: gomnd
	}
	res := ToQSP(data)
	fmt.Println(res.Encode()) // Prints bar=42&bat=10&baz=hello
}
