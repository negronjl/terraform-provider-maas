package api

import (
	"github.com/roblox/terraform-provider-maas/pkg/api/params"
	"github.com/roblox/terraform-provider-maas/pkg/maas/entity"
)

// VLANs represents the MaaS Vlans endpoint
type VLANs interface {
	Get(fabricID int) ([]entity.VLAN, error)
	Post(fabricID int, params *params.VLAN) (*entity.VLAN, error)
}
