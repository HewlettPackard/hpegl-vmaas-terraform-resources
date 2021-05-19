package cmp

import "github.com/hpe-hcss/vmaas-terraform-resources/internal/models"

type instance interface {
	CreateInstance(instanceBody models.CreateInstanceBody) error
}
