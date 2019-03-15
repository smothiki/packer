package common

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
)

// This step configures the VM to enable the VRDP server
// on the guest machine.
//
// Uses:
//   driver Driver
//   ui packer.Ui
//   vmName string
//
// Produces:
// vrdp_port unit - The port that VRDP is configured to listen on.
type StepConfigureVRDP struct {
	VRDPBindAddress string
	VRDPPortMin     uint
	VRDPPortMax     uint
}

func (s *StepConfigureVRDP) Run(_ context.Context, state multistep.StateBag) multistep.StepAction {
	driver := state.Get("driver").(Driver)
	ui := state.Get("ui").(packer.Ui)
	vmName := state.Get("vmName").(string)

	log.Printf("Looking for available port between %d and %d on %s", s.VRDPPortMin, s.VRDPPortMax, s.VRDPBindAddress)

	command := []string{
		"modifyvm", vmName,
		"--vrdeaddress", fmt.Sprintf("%s", s.VRDPBindAddress),
		"--vrdeauthtype", "null",
		"--vrde", "on",
		"--vrdeport",
		fmt.Sprintf("%d", vrdpPort),
	}
	if err := driver.VBoxManage(command...); err != nil {
		err := fmt.Errorf("Error enabling VRDP: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	state.Put("vrdpIp", s.VRDPBindAddress)
	state.Put("vrdpPort", vrdpPort)

	return multistep.ActionContinue
}

func (s *StepConfigureVRDP) Cleanup(state multistep.StateBag) {}
