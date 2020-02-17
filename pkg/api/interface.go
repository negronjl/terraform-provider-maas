package api

import (
	"github.com/roblox/terraform-provider-maas/pkg/api/params"
	"github.com/roblox/terraform-provider-maas/pkg/maas/entity"
)

type Interface interface {
	Delete(systemID string, id int) error
	Get(systemID string, id int) (*entity.Interface, error)
	AddTag(systemID string, id int, tag string) (*entity.Interface, error)
	Disconnect(systemID string, id int) (*entity.Interface, error)
	LinkSubnet(systemID string, id int, params *params.InterfaceLink) (*entity.Interface, error)
	RemoveTag(systemID string, id int, tag string) (*entity.Interface, error)
	SetDefaultGateway(systemID string, id, linkID int) (*entity.Interface, error)
	UnlinkSubnet(systemID string, id, linkID int) (*entity.Interface, error)
	Put(systemID string, id int, params interface{}) (*entity.Interface, error)
}
