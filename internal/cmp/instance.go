package cmp

import (
	cmp_client "github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/client"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/models"
)

type Instance struct {
	client cmp_client.APIClient
}

func (i *Instance) CreateInstance(instanceBody models.CreateInstanceBody) error {

	return nil
}
