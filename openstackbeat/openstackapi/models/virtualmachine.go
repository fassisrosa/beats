package models

import (
	"fmt"
	"time"
)

type VirtualMachine struct {
	Id string
	Name string
	Ip string
	Hypervisor string
	User string
	Tenant string
	Flavor string
	MemorySize int64
	DiskSize int64
	NumberCpus int64
	Status string
	CreatedAt time.Time
}

// TODO: Ip

func (vm VirtualMachine) String() string {
	return fmt.Sprintf("VM[id: %s, name: %s, status: %s, user: %s, tenant: %s, flavor: %s, mem: %d, disk: %d, vcpus: %d, hypervisor: %s]", vm.Id, vm.Name, vm.Status, vm.User, vm.Tenant, vm.Flavor, vm.MemorySize, vm.DiskSize, vm.NumberCpus, vm.Hypervisor)
}
