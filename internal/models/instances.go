// (C) Copyright 2021 Hewlett Packard Enterprise Development LP
package models

// CreateInstanceBody
type CreateInstanceBody struct {
	// Cloud ID
	ZoneId            int32                                 `json:"zoneId"`
	Instance          *CreateInstanceBodyInstance           `json:"instance"`
	Volumes           []CreateInstanceBodyVolumes           `json:"volumes"`
	NetworkInterfaces []CreateInstanceBodyNetworkInterfaces `json:"networkInterfaces"`
	Config            *CreateInstanceBodyConfig             `json:"config"`
}

// CreateInstanceBodyInstance
type CreateInstanceBodyInstance struct {
	Name         string                                  `json:"name"`
	Site         *CreateInstanceBodyInstanceSite         `json:"site"`
	InstanceType *CreateInstanceBodyInstanceInstanceType `json:"instanceType"`
	Layout       *CreateInstanceBodyInstanceLayout       `json:"layout"`
	Plan         *CreateInstanceBodyInstancePlan         `json:"plan"`
}

// CreateInstanceBodyConfig
type CreateInstanceBodyConfig struct {
	// Virtual Image ID(Required when VMware InstanceType is used)
	Template       int32  `json:"template,omitempty"`
	ResourcePoolId string `json:"resourcePoolId"`
	// To specify agent install (on/off)
	NoAgent              string `json:"noAgent,omitempty"`
	SmbiosAssetTag       string `json:"smbiosAssetTag,omitempty"`
	HostId               string `json:"hostId,omitempty"`
	VmwareDomainName     string `json:"vmwareDomainName,omitempty"`
	VmwareCustomSpec     string `json:"vmwareCustomSpec,omitempty"`
	NestedVirtualization string `json:"nestedVirtualization,omitempty"`
	CreateUser           bool   `json:"createUser,omitempty"`
}

// CreateInstanceBodyInstanceInstanceType
type CreateInstanceBodyInstanceInstanceType struct {
	// Instance type code
	Code string `json:"code"`
}

// CreateInstanceBodyInstanceLayout
type CreateInstanceBodyInstanceLayout struct {
	// The layout id for the instance type that you want to provision.
	Id int32 `json:"id"`
}

// CreateInstanceBodyInstancePlan
type CreateInstanceBodyInstancePlan struct {
	// Service Plan ID
	Id int32 `json:"id"`
}

// CreateInstanceBodyInstanceSite
type CreateInstanceBodyInstanceSite struct {
	// Group ID
	Id int32 `json:"id"`
}

// CreateInstanceBodyNetwork
type CreateInstanceBodyNetwork struct {
	Id int32 `json:"id"`
}

// CreateInstanceBodyNetworkInterfaces
type CreateInstanceBodyNetworkInterfaces struct {
	Network                *CreateInstanceBodyNetwork `json:"network"`
	NetworkInterfaceTypeId int32                      `json:"networkInterfaceTypeId"`
}

// CreateInstanceBodyVolumes
type CreateInstanceBodyVolumes struct {
	Id         int32 `json:"id,omitempty"`
	RootVolume bool  `json:"rootVolume,omitempty"`
	// Name/type of the LV being created
	Name        string `json:"name"`
	Size        int32  `json:"size,omitempty"`
	StorageType int32  `json:"storageType,omitempty"`
	// The ID of the specific datastore. Auto selection can be specified as auto or autoCluster (for clusters).
	DatastoreId int32 `json:"datastoreId"`
}

// GetInstanceResponse
type GetInstanceResponse struct {
	Instance *GetInstanceResponseInstance `json:"instance,omitempty"`
	Success  bool                         `json:"success,omitempty"`
}

// GetInstanceResponseInstance
type GetInstanceResponseInstance struct {
	Id                  int32                                       `json:"id,omitempty"`
	Uuid                string                                      `json:"uuid,omitempty"`
	AccountId           int32                                       `json:"accountId,omitempty"`
	Tenant              *GetInstanceResponseInstanceTenant          `json:"tenant,omitempty"`
	InstanceType        *GetInstanceResponseInstanceInstanceType    `json:"instanceType,omitempty"`
	Group               *GetInstanceResponseInstanceGroup           `json:"group,omitempty"`
	Cloud               *GetInstanceResponseInstanceCloud           `json:"cloud,omitempty"`
	Containers          []int32                                     `json:"containers,omitempty"`
	Servers             []int32                                     `json:"servers,omitempty"`
	ConnectionInfo      []GetInstanceResponseInstanceConnectionInfo `json:"connectionInfo,omitempty"`
	Layout              *GetInstanceResponseInstanceLayout          `json:"layout,omitempty"`
	Plan                *GetInstanceResponseInstancePlan            `json:"plan,omitempty"`
	Name                string                                      `json:"name,omitempty"`
	Description         string                                      `json:"description,omitempty"`
	Config              *GetInstanceResponseInstanceConfig          `json:"config,omitempty"`
	Volumes             []GetInstanceResponseInstanceVolumes        `json:"volumes,omitempty"`
	Controllers         string                                      `json:"controllers,omitempty"`
	Interfaces          []GetInstanceResponseInstanceInterfaces     `json:"interfaces,omitempty"`
	CustomOptions       *interface{}                                `json:"customOptions,omitempty"`
	InstanceVersion     string                                      `json:"instanceVersion,omitempty"`
	Labels              []string                                    `json:"labels,omitempty"`
	Tags                []GetInstanceResponseInstanceTags           `json:"tags,omitempty"`
	Evars               []GetInstanceResponseInstanceEvars          `json:"evars,omitempty"`
	MaxMemory           int32                                       `json:"maxMemory,omitempty"`
	MaxStorage          int32                                       `json:"maxStorage,omitempty"`
	MaxCores            int32                                       `json:"maxCores,omitempty"`
	HourlyCost          int32                                       `json:"hourlyCost,omitempty"`
	HourlyPrice         int32                                       `json:"hourlyPrice,omitempty"`
	DateCreated         string                                      `json:"dateCreated,omitempty"`
	LastUpdated         string                                      `json:"lastUpdated,omitempty"`
	HostName            string                                      `json:"hostName,omitempty"`
	FirewallEnabled     bool                                        `json:"firewallEnabled,omitempty"`
	NetworkLevel        string                                      `json:"networkLevel,omitempty"`
	AutoScale           bool                                        `json:"autoScale,omitempty"`
	Locked              bool                                        `json:"locked,omitempty"`
	Status              string                                      `json:"status,omitempty"`
	StatusDate          string                                      `json:"statusDate,omitempty"`
	ExpireCount         int32                                       `json:"expireCount,omitempty"`
	ExpireWarningSent   bool                                        `json:"expireWarningSent,omitempty"`
	ShutdownCount       int32                                       `json:"shutdownCount,omitempty"`
	ShutdownWarningSent bool                                        `json:"shutdownWarningSent,omitempty"`
	CreatedBy           *GetInstanceResponseInstanceCreatedBy       `json:"createdBy,omitempty"`
	Owner               *GetInstanceResponseInstanceCreatedBy       `json:"owner,omitempty"`
}

//GetInstanceResponseInstanceCloud
type GetInstanceResponseInstanceCloud struct {
	Id   int32  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// GetInstanceResponseInstanceConfig
type GetInstanceResponseInstanceConfig struct {
	Template             int32                                            `json:"template,omitempty"`
	ResourcePoolId       string                                           `json:"resourcePoolId,omitempty"`
	NestedVirtualization string                                           `json:"nestedVirtualization,omitempty"`
	CreateUser           bool                                             `json:"createUser,omitempty"`
	NoAgent              string                                           `json:"noAgent,omitempty"`
	CreateBackup         bool                                             `json:"createBackup,omitempty"`
	MemoryDisplay        string                                           `json:"memoryDisplay,omitempty"`
	LayoutSize           int32                                            `json:"layoutSize,omitempty"`
	Backup               *GetInstanceResponseInstanceConfigBackup         `json:"backup,omitempty"`
	RemovalOptions       *GetInstanceResponseInstanceConfigRemovalOptions `json:"removalOptions,omitempty"`
}

// GetInstanceResponseInstanceConfigBackup
type GetInstanceResponseInstanceConfigBackup struct {
	ProviderBackupType int32  `json:"providerBackupType,omitempty"`
	JobAction          string `json:"jobAction,omitempty"`
	JobName            string `json:"jobName,omitempty"`
	Name               string `json:"name,omitempty"`
}

// GetInstanceResponseInstanceConfigRemovalOptions
type GetInstanceResponseInstanceConfigRemovalOptions struct {
	Force           bool  `json:"force,omitempty"`
	KeepBackups     bool  `json:"keepBackups,omitempty"`
	ReleaseEIPs     bool  `json:"releaseEIPs,omitempty"`
	RemoveVolumes   bool  `json:"removeVolumes,omitempty"`
	RemoveResources bool  `json:"removeResources,omitempty"`
	UserId          int32 `json:"userId,omitempty"`
}

// GetInstanceResponseInstanceConnectionInfo
type GetInstanceResponseInstanceConnectionInfo struct {
	Ip string `json:"ip,omitempty"`
}

// GetInstanceResponseInstanceCreatedBy
type GetInstanceResponseInstanceCreatedBy struct {
	Id       int32  `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
}

// GetInstanceResponseInstanceEvars
type GetInstanceResponseInstanceEvars struct {
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
	Export bool   `json:"export,omitempty"`
	Masked bool   `json:"masked,omitempty"`
}

// GetInstanceResponseInstanceGroup
type GetInstanceResponseInstanceGroup struct {
	Id   int32  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// GetInstanceResponseInstanceInstanceType
type GetInstanceResponseInstanceInstanceType struct {
	Id       int32  `json:"id,omitempty"`
	Code     string `json:"code,omitempty"`
	Category string `json:"category,omitempty"`
	Name     string `json:"name,omitempty"`
}

// GetInstanceResponseInstanceInterfaces
type GetInstanceResponseInstanceInterfaces struct {
	Id      int32                               `json:"id,omitempty"`
	Row     int32                               `json:"row,omitempty"`
	Network *GetInstanceResponseInstanceNetwork `json:"network,omitempty"`
}

// GetInstanceResponseInstanceLayout
type GetInstanceResponseInstanceLayout struct {
	Id                int32  `json:"id,omitempty"`
	Name              string `json:"name,omitempty"`
	ProvisionTypeCode string `json:"provisionTypeCode,omitempty"`
}

// GetInstanceResponseInstanceNetwork
type GetInstanceResponseInstanceNetwork struct {
	Id                     int32                                   `json:"id,omitempty"`
	Subnet                 string                                  `json:"subnet,omitempty"`
	Group                  string                                  `json:"group,omitempty"`
	DhcpServer             bool                                    `json:"dhcpServer,omitempty"`
	Name                   string                                  `json:"name,omitempty"`
	Pool                   *GetInstanceResponseInstanceNetworkPool `json:"pool,omitempty"`
	IpAddress              string                                  `json:"ipAddress,omitempty"`
	IpMode                 string                                  `json:"ipMode,omitempty"`
	NetworkInterfaceTypeId int32                                   `json:"networkInterfaceTypeId,omitempty"`
}

// GetInstanceResponseInstanceNetworkPool
type GetInstanceResponseInstanceNetworkPool struct {
	Id int32 `json:"id,omitempty"`
}

// GetInstanceResponseInstancePlan
type GetInstanceResponseInstancePlan struct {
	Id   int32  `json:"id,omitempty"`
	Code string `json:"code,omitempty"`
	Name string `json:"name,omitempty"`
}

// GetInstanceResponseInstanceTags
type GetInstanceResponseInstanceTags struct {
	Id    int32  `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// GetInstanceResponseInstanceTenant
type GetInstanceResponseInstanceTenant struct {
	Id   int32  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// GetInstanceResponseInstanceVolumes
type GetInstanceResponseInstanceVolumes struct {
	Size              int32  `json:"size,omitempty"`
	Name              string `json:"name,omitempty"`
	RootVolume        bool   `json:"rootVolume,omitempty"`
	Id                string `json:"id,omitempty"`
	DatastoreId       int32  `json:"datastoreId,omitempty"`
	MaxStorage        int32  `json:"maxStorage,omitempty"`
	DeviceDisplayName string `json:"deviceDisplayName,omitempty"`
}

// ResizeInstanceBody
type ResizeInstanceBody struct {
	Instance *ResizeInstanceBodyInstance `json:"instance,omitempty"`
}

// ResizeInstanceBodyInstance
type ResizeInstanceBodyInstance struct {
	// Instance ID
	Id   int32                           `json:"id,omitempty"`
	Plan *ResizeInstanceBodyInstancePlan `json:"plan,omitempty"`
	// Can be used to grow just the logical volume of the instance instead of choosing a plan
	Volumes               []ResizeInstanceBodyInstanceVolumes `json:"volumes,omitempty"`
	DeleteOriginalVolumes bool                                `json:"deleteOriginalVolumes,omitempty"`
}

// ResizeInstanceBodyInstancePlan
type ResizeInstanceBodyInstancePlan struct {
	// Service Plan ID
	Id int32 `json:"id,omitempty"`
}

// ResizeInstanceBodyInstanceVolumes
type ResizeInstanceBodyInstanceVolumes struct {
	Id          int32  `json:"id,omitempty"`
	RootVolume  bool   `json:"rootVolume,omitempty"`
	Name        string `json:"name,omitempty"`
	Size        int32  `json:"size,omitempty"`
	StorageType int32  `json:"storageType,omitempty"`
	DatastoreId int32  `json:"datastoreId,omitempty"`
}

// SnapshotBody
type SnapshotBody struct {
	Snapshot *SnapshotBodySnapshot `json:"snapshot,omitempty"`
}

// SnapshotBodySnapshot
type SnapshotBodySnapshot struct {
	// Optional name for the snapshot being created
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// UpdateInstanceBody
type UpdateInstanceBody struct {
	Instance *UpdateInstanceBodyInstance `json:"instance,omitempty"`
}

// UpdateInstanceBodyInstance
type UpdateInstanceBodyInstance struct {
	// Unique name scoped to your account for the instance
	Name string `json:"name,omitempty"`
	// Optional description field
	Description string `json:"description,omitempty"`
	// Add or update value of Metadata tags, Array of objects having a name and value
	AddTags *interface{} `json:"addTags,omitempty"`
	// Remove Metadata tags, Array of objects having a name and an optional value. If value is passed, it must match to be removed
	RemoveTags *interface{} `json:"removeTags,omitempty"`
}
