package api

import (
	"github.com/roblox/terraform-provider-maas/pkg/api/params"
	"github.com/roblox/terraform-provider-maas/pkg/maas/entity"
)

type Interfaces interface {
	Get(systemID string) ([]entity.Interface, error)
	CreateBond(systemID string, params *params.InterfaceBond) (*entity.Interface, error)
	CreateBridge(systemID string, params *params.InterfaceBridge) (*entity.Interface, error)
	CreatePhysical(systemID string, params *params.InterfacePhysical) (*entity.Interface, error)
	CreateVLAN(systemID string, params *params.InterfaceVLAN) (*entity.Interface, error)
}
