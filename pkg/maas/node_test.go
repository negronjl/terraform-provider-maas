package maas_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	. "github.com/roblox/terraform-provider-maas/pkg/maas"
)

func TestNodeStatus(t *testing.T) {
	tests := []struct {
		name string
		got  NodeStatus
		want NodeStatus
	}{
		{name: "new", got: NodeStatusNew, want: 0},
		{name: "default", got: NodeStatusDefault, want: 0},
		{name: "commissioning", got: NodeStatusCommissioning, want: 1}, // nolint: gomnd
		{name: "ready", got: NodeStatusReady, want: 4},                 // nolint: gomnd
		{name: "deployed", got: NodeStatusDeployed, want: 6},           // nolint: gomnd
		{name: "deploying", got: NodeStatusDeploying, want: 9},         // nolint: gomnd
		{name: "allocated", got: NodeStatusAllocated, want: 10},        // nolint: gomnd
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
