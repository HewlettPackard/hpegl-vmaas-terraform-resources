//  (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package diffValidation

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Instance(diff *schema.ResourceDiff) error {
	notAllowed := []string{"network", "plan_id", "scale"}

	for _, param := range notAllowed {
		if diff.HasChange(param) {
			return fmt.Errorf("modifying %q is not allowed. "+
				"Please fix the configuration and re try", param)
		}
	}

	if err := instanceVolumeDiffValidate(diff); err != nil {
		return err
	}

	return nil
}

func instanceVolumeDiffValidate(diff *schema.ResourceDiff) error {
	if !diff.HasChange("volume") {
		return nil
	}
	oldVol, newVol := diff.GetChange("volume")

	newVolMap := make(map[string]bool)

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
				return fmt.Errorf("renaming the root/primary volume '%s' is not allowed. "+
					"Please fix your configuration and retry", tVol["name"].(string))
			}
		}
	}

	return nil
}
