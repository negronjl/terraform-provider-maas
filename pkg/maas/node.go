package maas

// NodeStatus correlates to a Node Status code returned from the MAAS API.
type NodeStatus int

// The statuses are defined in src/maasserver/enum.py as of 06-2018, which can
// be found at https://github.com/maas/maas/blob/master/src/maasserver/enum.py
// The definitions are in `class NODE_STATUS`, which starts on L48 right now.
// The statuses map to sequential ints, hence the iota - if a status is added or
// changed, presumably they will maintain the sequential numbering. This way we only
// have to maintain the correct order instead of having to ensure each value is
// mapped correctly.
// NodeStatusDefault, defined at the bottom, is an exception - it has the same
// value as NodeStatusNew.
// There are some smoke tests in node_test.go to verify some of the more
// relevant statuses such as "commissioning" and "deployed".
const (
	NodeStatusNew NodeStatus = iota
	NodeStatusCommissioning
	NodeStatusFailedCommissioning
	NodeStatusMissing
	NodeStatusReady
	NodeStatusReserved
	NodeStatusDeployed
	NodeStatusRetired
	NodeStatusBroken
	NodeStatusDeploying
	NodeStatusAllocated
	NodeStatusFailedDeployment
	NodeStatusReleasing
	NodeStatusFailedReleasing
	NodeStatusDiskErasing
	NodeStatusFailedDiskErasing
	NodeStatusRescueMode
	NodeStatusEnteringRescureMode
	NodeStatusFailedEnteringRescueMode
	NodeStatusExitingRescueMode
	NodeStatusFailedExitingRescueMode
	NodeStatusTesting
	NodeStatusFailedTesting

	NodeStatusDefault NodeStatus = 0
)
