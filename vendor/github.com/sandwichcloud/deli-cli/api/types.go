package api

import (
	"time"

	"net"

	"github.com/satori/go.uuid"
)

type pageLinks struct {
	HREF string `json:"href"`
	REL  string `json:"rel"`
}

type AuthDiscover struct {
	Default *string           `json:"default"`
	Github  *GithubAuthDriver `json:"github"`
}

type GithubAuthDriver struct {
	//Github auth driver has no options
}

type Project struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type ProjectList struct {
	Links    []pageLinks `json:"projects_links"`
	Projects []Project   `json:"projects"`
}

type ProjectQuota struct {
	VCPU     int `json:"vcpu"`
	Ram      int `json:"ram"`
	Disk     int `json:"disk"`
	UsedVCPU int `json:"used_vcpu"`
	UsedRam  int `json:"used_ram"`
	UsedDisk int `json:"used_disk"`
}

type Region struct {
	Name           string    `json:"name"`
	Datacenter     string    `json:"datacenter"`
	ImageDatastore string    `json:"image_datastore"`
	ImageFolder    string    `json:"image_folder"`
	Schedulable    bool      `json:"schedulable"`
	State          string    `json:"state"`
	ErrorMessage   string    `json:"error_message"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type RegionList struct {
	Links   []pageLinks `json:"regions_links"`
	Regions []Region    `json:"regions"`
}

type Zone struct {
	Name                 string    `json:"name"`
	RegionName           string    `json:"region_name"`
	VMCluster            string    `json:"vm_cluster"`
	VMDatastore          string    `json:"vm_datastore"`
	VMFolder             string    `json:"vm_folder"`
	CoreProvisionPercent int       `json:"core_provision_percent"`
	RamProvisionPercent  int       `json:"ram_provision_percent"`
	Schedulable          bool      `json:"schedulable"`
	State                string    `json:"state"`
	ErrorMessage         string    `json:"error_message"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type ZoneList struct {
	Links []pageLinks `json:"zones_links"`
	Zones []Zone      `json:"zones"`
}

type Volume struct {
	Name         string    `json:"name"`
	ZoneName     string    `json:"zone_name"`
	Size         int       `json:"size"`
	AttachedTo   string    `json:"attached_to"`
	Task         string    `json:"task"`
	State        string    `json:"state"`
	ErrorMessage string    `json:"error_message"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type VolumeList struct {
	Links   []pageLinks `json:"volumes_links"`
	Volumes []Volume    `json:"volumes"`
}

type Image struct {
	Name         string    `json:"name"`
	ProjectName  string    `json:"project_name"`
	FileName     string    `json:"file_name"`
	RegionName   string    `json:"region_name"`
	Visibility   string    `json:"visibility"`
	State        string    `json:"state"`
	ErrorMessage string    `json:"error_message"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ImageList struct {
	Links  []pageLinks `json:"images_links"`
	Images []Image     `json:"images"`
}

type TokenInfo struct {
	Username           string   `json:"username,omitempty"`
	Driver             string   `json:"driver,omitempty"`
	ServiceAccountName string   `json:"service_account_name,omitempty"`
	ProjectName        string   `json:"project_name,omitempty"`
	GlobalRoles        []string `json:"global_roles"`
	ProjectRoles       []string `json:"project_roles,omitempty"`
}

type Network struct {
	Name         string    `json:"name"`
	PortGroup    string    `json:"port_group"`
	Cidr         string    `json:"cidr"`
	Gateway      net.IP    `json:"gateway"`
	DNSServers   []net.IP  `json:"dns_servers"`
	PoolStart    net.IP    `json:"pool_start"`
	PoolEnd      net.IP    `json:"pool_end"`
	RegionName   string    `json:"region_name"`
	State        string    `json:"state"`
	ErrorMessage string    `json:"error_message"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type NetworkList struct {
	Links    []pageLinks `json:"networks_links"`
	Networks []Network   `json:"networks"`
}

type Keypair struct {
	Name      string `json:"name" structs:"name"`
	PublicKey string `json:"public_key" structs:"public_key"`
}

type KeypairList struct {
	Links    []pageLinks `json:"keypairs_links"`
	KeyPairs []Keypair   `json:"keypairs"`
}

type NetworkPort struct {
	ID          uuid.UUID `json:"id"`
	NetworkName string    `json:"network_name"`
	IPAddress   net.IP    `json:"ip_address"`
	State       string    `json:"state"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type NetworkPortList struct {
	Links        []pageLinks   `json:"network-ports_links"`
	NetworkPorts []NetworkPort `json:"network-ports"`
}

type Flavor struct {
	Name      string    `json:"name"`
	VCPUS     int       `json:"vcpus"`
	Ram       int       `json:"ram"`
	Disk      int       `json:"disk"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FlavorList struct {
	Links   []pageLinks `json:"flavors_links"`
	Flavors []Flavor    `json:"flavors"`
}

type Instance struct {
	Name               string            `json:"name"`
	ImageName          string            `json:"image_name"`
	NetworkPortID      uuid.UUID         `json:"network_port_id"`
	RegionName         string            `json:"region_name"`
	ZoneName           string            `json:"zone_name"`
	ServiceAccountName string            `json:"service_account_name"`
	Tags               map[string]string `json:"tags"`
	UserData           string            `json:"user_data"`
	KeypairNames       []string          `json:"keypair_names"`
	FlavorName         string            `json:"flavor_name"`
	VCPUS              int               `json:"vcpus"`
	Ram                int               `json:"ram"`
	Disk               int               `json:"disk"`
	State              string            `json:"state"`
	PowerState         string            `json:"power_state"`
	Task               string            `json:"task"`
	ErrorMessage       string            `json:"error_message"`
	CreatedAt          time.Time         `json:"created_at"`
	UpdatedAt          time.Time         `json:"updated_at"`
}

type InstanceInitialVolume struct {
	Size       int  `json:"size"`
	AutoDelete bool `json:"auto_delete"`
}

type InstanceList struct {
	Links     []pageLinks `json:"instances_links"`
	Instances []Instance  `json:"instances"`
}

type Permissions struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

type PermissionList struct {
	Links       []pageLinks   `json:"permissions_links"`
	Permissions []Permissions `json:"permissions"`
}

type Role struct {
	Name        string    `json:"name"`
	Permissions []string  `json:"permissions"`
	State       string    `json:"state"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type RoleList struct {
	Links []pageLinks `json:"roles_links"`
	Roles []Role      `json:"roles"`
}

type ServiceAccount struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Keys      []string  `json:"keys"`
	State     string    `json:"state"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ServiceAccountList struct {
	Links           []pageLinks      `json:"service-accounts_links"`
	ServiceAccounts []ServiceAccount `json:"service-accounts"`
}

type PolicyBinding struct {
	Role    string   `json:"role"`
	Members []string `json:"members"`
}

type Policy struct {
	Bindings        []PolicyBinding `json:"bindings"`
	ResourceVersion string          `json:"resource_version"`
}
