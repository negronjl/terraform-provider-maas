package api

import (
	"github.com/roblox/terraform-provider-maas/pkg/api/params"
	"github.com/roblox/terraform-provider-maas/pkg/maas/entity"
	"github.com/roblox/terraform-provider-maas/pkg/maas/entity/subnet"
)

// Subnet represents the MaaS Subnet endpoint
type Subnet interface {
	Delete(id int) error
	Get(id int) (*entity.Subnet, error)
	GetIPAddresses(id int, WithUsername, WithSummary bool) ([]subnet.IPAddress, error)
	GetReservedIPRanges(id int) ([]subnet.ReservedIPRange, error)
	GetStatistics(id int, IncludeRanges, IncludeSuggestions bool) (*subnet.Statistics, error)
	GetUnreservedIPRanges(id int) ([]subnet.IPRange, error)
	Put(id int, params *params.Subnet) (*entity.Subnet, error)
}
