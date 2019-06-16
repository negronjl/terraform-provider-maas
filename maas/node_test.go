package maas_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	. "github.com/roblox/terraform-provider-maas/maas"
)

func TestNodeStatus(t *testing.T) {
	tests := []struct {
		name string
		got  NodeStatus
		want NodeStatus
	}{
		{name: "new", got: NodeStatusNew, want: 0},
		{name: "default", got: NodeStatusDefault, want: 0},
		{name: "commissioning", got: NodeStatusCommissioning, want: 1},
		{name: "ready", got: NodeStatusReady, want: 4},
		{name: "deployed", got: NodeStatusDeployed, want: 6},
		{name: "deploying", got: NodeStatusDeploying, want: 9},
		{name: "allocated", got: NodeStatusAllocated, want: 10},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			diff := cmp.Diff(tc.want, tc.got)
			if diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
