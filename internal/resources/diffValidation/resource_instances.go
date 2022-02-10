//  (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package diffvalidation

import (
	"fmt"
	"strings"

	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Instance struct {
	diff *schema.ResourceDiff
}

func NewInstanceValidate(diff *schema.ResourceDiff) *Instance {
	return &Instance{
		diff: diff,
	}
}

func (i *Instance) DiffValidate() error {
	notAllowed := []string{"plan_id", "scale"}

	for _, param := range notAllowed {
		if i.diff.HasChange(param) {
			oldParam, _ := i.diff.GetChange(param)
			if utils.IsEmpty(oldParam) {
				continue
			}

			return fmt.Errorf("modifying %q is not allowed. "+
				"Please fix the configuration and re try", param)
		}
	}

	err := i.instanceValidatePowerTransition()
	if err != nil {
		return err
	}

	if err := i.instanceVolumeDiffValidate(); err != nil {
		return err
	}

	if err := i.instanceTemplateValidate(); err != nil {
		return err
	}

	if err := i.instanceValidateSnapshotDeletion(); err != nil {
		return err
	}

	return nil
}

func (i *Instance) instanceTemplateValidate() error {
	configSet := i.diff.Get("config").(*schema.Set)
	if configSet == nil {
		return nil
	}

	config := configSet.List()
	// for clone config can be nil
	if len(config) > 0 {
		c0 := config[0].(map[string]interface{})
		templateID := c0["template_id"].(int)
		if strings.ToLower(i.diff.Get("instance_type_code").(string)) == "vmware" {
			if templateID == 0 {
				return fmt.Errorf("template_id is required for 'vmware' instance type code")
			}
		}
	}

	return nil
}

func (i *Instance) instanceVolumeDiffValidate() error {
	if !i.diff.HasChange("volume") {
		return nil
	}

	oldVol, newVol := i.diff.GetChange("volume")
	newVolMap := make(map[string]bool)

	// Validate is the volume names are unique
	if err := i.instanceValidateVolumeNameIsUnique(newVol.([]interface{})); err != nil {
		return err
	}

	// If create operation validation should be skipped
	if utils.IsEmpty(oldVol) {
		return nil
	}

	// Validate if the primary volume is being modified
	for _, vol := range newVol.([]interface{}) {
		tVol := vol.(map[string]interface{})
		newVolMap[tVol["name"].(string)] = tVol["root"].(bool)
	}

	for _, vol := range oldVol.([]interface{}) {
		tVol := vol.(map[string]interface{})
		isRoot, ok := newVolMap[tVol["name"].(string)]

		if tVol["root"].(bool) {
			if !ok {
				return fmt.Errorf("deleting root/primary volume '%s' is not allowed. "+
					"Please fix your configuration and retry", tVol["name"].(string))
			} else if !isRoot {
				return fmt.Errorf("interchanging the root/primary volume '%s' is not allowed. "+
					"Please fix your configuration and retry", tVol["name"].(string))
			}
		}
	}

	return nil
}

func (i *Instance) instanceValidateVolumeNameIsUnique(vol []interface{}) error {
	volumes := make(map[string]bool)
	for _, v := range vol {
		tVol := v.(map[string]interface{})
		if _, ok := volumes[tVol["name"].(string)]; !ok {
			volumes[tVol["name"].(string)] = true

			continue
		}

		return fmt.Errorf("volume names should be unique")
	}

	return nil
}

func (i *Instance) instanceValidatePowerTransition() error {
	if !i.diff.HasChange("power") {
		return nil
	}

	oldPower, newPower := i.diff.GetChange("power")
	oldPowerStr := oldPower.(string)
	newPowerStr := newPower.(string)

	// if oldPower is empty it can be either create operation as well as
	// update operation
	if oldPower == "" {
		// if status is empty it indicates this is a create operation
		if _, status := i.diff.GetChange("status"); status == "" {
			if newPowerStr != utils.PowerOn {
				return fmt.Errorf("while creating instance only %s is supported", utils.PowerOn)
			}

			return nil
		}
		// In case of update operation consider oldPower = powerOn
		oldPowerStr = utils.PowerOn
	}
	// If the current power state is On, then only Off, suspend and restart will allowed
	// If the current state is "" due to new instance creation
	if oldPowerStr == utils.PowerOn {
		if newPowerStr == utils.PowerOff || newPowerStr == utils.Suspend || newPowerStr == utils.Restart {
			return nil
		}
	} else if newPowerStr == utils.PowerOn {
		// The current state is Off, suspend and restart,
		// then the new state allowed is only PowerOn
		if oldPower == utils.PowerOff || oldPower == utils.Suspend || oldPower == utils.Restart {
			return nil
		}
	}

	return fmt.Errorf("power operation not allowed from %s state to %s state", oldPowerStr, newPowerStr)
}

func (i *Instance) instanceValidateSnapshotDeletion() error {
	if !i.diff.HasChange("snapshot") {
		return nil
	}
	// if there is a change in snapshot, then get the new value
	_, newSnapshotValue := i.diff.GetChange("snapshot")
	newMap := utils.GetlistMap(newSnapshotValue)
	// if an attempt to delete the snapshot, then return error
	if len(newMap) <= 0 {
		return fmt.Errorf("deleting snapshot is not supported currently")
	}

	return nil
}
