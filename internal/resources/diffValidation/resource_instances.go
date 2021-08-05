//  (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package diffvalidation

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/utils"
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
	notAllowed := []string{"network", "plan_id", "scale"}

	for _, param := range notAllowed {
		if i.diff.HasChange(param) {
			oldParam, _ := i.diff.GetChange(param)
			if i.isEmpty(oldParam) {
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

	return nil
}

func (i *Instance) isEmpty(param interface{}) bool {
	switch paramVal := param.(type) {
	case []interface{}:
		return len(paramVal) == 0
	case string:
		return paramVal == ""
	case int:
		return paramVal == 0
	default:
		return false
	}
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
	if i.isEmpty(oldVol) {
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

	// If the current power state is On, then only Off, suspend and restart will allowed
	if oldPower.(string) == utils.PowerOn {
		if newPower.(string) == utils.PowerOff || newPower.(string) == utils.Suspend || newPower.(string) == utils.Restart {
			return nil
		}
	} else {
		// If the current state is "" due to new instance creation or the current state is Off, suspend and restart,
		// then the new state allowed is only PowerOn
		if newPower.(string) == utils.PowerOn {
			return nil
		}
	}

	return fmt.Errorf("power operation not allowed from %s state to %s state", oldPower.(string), newPower.(string))
}
