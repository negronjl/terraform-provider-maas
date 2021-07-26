package node_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	. "github.com/roblox/terraform-provider-maas/pkg/maas/entity/node"
)

func TestStatus(t *testing.T) {
	tests := []struct {
		name string
		got  Status
		want Status
	}{
		{name: "new", got: StatusNew, want: 0},
		{name: "default", got: StatusDefault, want: 0},
		{name: "commissioning", got: StatusCommissioning, want: 1}, // nolint: gomnd
		{name: "ready", got: StatusReady, want: 4},                 // nolint: gomnd
		{name: "deployed", got: StatusDeployed, want: 6},           // nolint: gomnd
		{name: "deploying", got: StatusDeploying, want: 9},         // nolint: gomnd
		{name: "allocated", got: StatusAllocated, want: 10},        // nolint: gomnd
	}

	for _, testCase := range tests {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			diff := cmp.Diff(tc.want, tc.got)
			if diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
