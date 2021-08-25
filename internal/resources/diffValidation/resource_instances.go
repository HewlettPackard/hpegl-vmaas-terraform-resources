//  (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package diffvalidation

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/utils"
)

var errPrimaryNetworkUpdation = fmt.Errorf("primary network updation/deletion is not supported")

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

	if err := i.validateIsPrimaryNetworkChanged(); err != nil {
		return err
	}

	return nil
}

func (i *Instance) validateIsPrimaryNetworkChanged() error {
	if i.diff.HasChange("network") {
		o, n := i.diff.GetChange("network")
		oldMap := utils.GetlistMap(o)
		newMap := utils.GetlistMap(n)

		// skip upon create operation
		if len(oldMap) == 0 {
			return nil
		}
		// check whether primary network has been changed?
		for i := range oldMap {
			if oldMap[i]["is_primary"].(bool) {
				if len(newMap) < i {
					return errPrimaryNetworkUpdation
				}
				if !reflect.DeepEqual(newMap[i], oldMap[i]) {
					return errPrimaryNetworkUpdation
				}

				break
			}
		}
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
	// If the current power state is On, then only Off, suspend and restart will allowed
	// If the current state is "" due to new instance creation
	if oldPowerStr == utils.PowerOn || oldPowerStr == "" {
		if newPowerStr == utils.PowerOff || newPowerStr == utils.Suspend || newPowerStr == utils.Restart {
			return nil
		}
	} else if newPowerStr == utils.PowerOn {
		// The current state is Off, suspend and restart,
		// then the new state allowed is only PowerOn
		return nil
	}

	return fmt.Errorf("power operation not allowed from %s state to %s state", oldPowerStr, newPowerStr)
}
