package openstackapi

import (
	"fmt"

	"github.com/fassisrosa/beats/openstackbeat/openstackapi/mitaka"
	"github.com/fassisrosa/beats/openstackbeat/openstackapi/models"
)

type OpenStackAPI interface {
	GetAllInfo (mainUrl string) ([]models.VirtualMachine, error)
}


func NewOpenStackAPI (version string) (OpenStackAPI, error) {
	if version == "mitaka" {
		return mitaka.MitakaOpenStackAPI{}, nil
	}
	return nil, fmt.Errorf("Unsupported OpenStack version '%s'",version)
}
