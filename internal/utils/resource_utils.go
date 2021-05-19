package utils

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func ListOfMap() *schema.Schema {
	return &schema.Schema{
		Type: schema.TypeList,
		Elem: &schema.Schema{
			Type: schema.TypeMap,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}
