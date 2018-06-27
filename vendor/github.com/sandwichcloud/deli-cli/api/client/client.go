package client

import (
	"fmt"
	"net/http"

	"context"

	"net"

	"github.com/sandwichcloud/deli-cli/api"
	"github.com/sandwichcloud/deli-cli/api/client/auth"
	"github.com/sandwichcloud/deli-cli/api/client/flavor"
	"github.com/sandwichcloud/deli-cli/api/client/image"
	"github.com/sandwichcloud/deli-cli/api/client/instance"
	"github.com/sandwichcloud/deli-cli/api/client/keypair"
	"github.com/sandwichcloud/deli-cli/api/client/network"
	"github.com/sandwichcloud/deli-cli/api/client/permission"
	"github.com/sandwichcloud/deli-cli/api/client/policy"
	"github.com/sandwichcloud/deli-cli/api/client/project"
	"github.com/sandwichcloud/deli-cli/api/client/region"
	"github.com/sandwichcloud/deli-cli/api/client/role"
	"github.com/sandwichcloud/deli-cli/api/client/serviceAccount"
	"github.com/sandwichcloud/deli-cli/api/client/volume"
	"github.com/sandwichcloud/deli-cli/api/client/zone"
	"golang.org/x/oauth2"
)

type SandwichClient struct {
	APIServer *string
	token     *oauth2.Token
}

type ClientInterface interface {
	createOAuthClient() *http.Client
	Auth() AuthClientInterface
	Project() ProjectClientInterface
	Region() RegionClientInterface
	Zone() ZoneClientInterface
	Image(projectName string) ImageClientInterface
	Network() NetworkClientInterface
	NetworkPort(projectName string) NetworkPortClientInterface
	Keypair(projectName string) KeypairClientInterface
	Flavor() FlavorClientInterface
	Volume(projectName string) VolumeClientInterface
	Instance(projectName string) InstanceClientInterface
	Permission() PermissionClientInterface
	SystemRole() RoleClientInterface
	ProjectRole(projectName string) RoleClientInterface
	SystemServiceAccount() ServiceAccountClientInterface
	ProjectServiceAccount(projectName string) ServiceAccountClientInterface
	SystemPolicy() PolicyClientInterface
	ProjectPolicy(projectName string) PolicyClientInterface
	SetToken(token *oauth2.Token)
}

type AuthClientInterface interface {
	Login(username, password string) (*oauth2.Token, error)
}

type ProjectClientInterface interface {
	Create(name string) (*api.Project, error)
	Get(name string) (*api.Project, error)
	Delete(name string) error
	List(limit int, marker string) (*api.ProjectList, error)
	GetQuota(projectName string) (*api.ProjectQuota, error)
	SetQuota(projectName string, vcpu, ram, disk int) error
}

type RegionClientInterface interface {
	Create(name, datacenter, imageDatastore, imageFolder string) (*api.Region, error)
	Get(name string) (*api.Region, error)
	Delete(name string) error
	List(limit int, marker string) (*api.RegionList, error)
	ActionSchedule(name string, schedulable bool) error
}

type ZoneClientInterface interface {
	Create(name, regionName, vmCluster, vmDatastore, vmFolder string, coreProvisionPercent, ramProvisionPercent int) (*api.Zone, error)
	Get(name string) (*api.Zone, error)
	Delete(name string) error
	List(regionName string, limit int, marker string) (*api.ZoneList, error)
	ActionSchedule(name string, schedulable bool) error
}

type VolumeClientInterface interface {
	Create(name, zoneName string, size int) (*api.Volume, error)
	Get(name string) (*api.Volume, error)
	Delete(name string) error
	List(limit int, marker string) (*api.VolumeList, error)
	ActionAttach(name, instanceName string) error
	ActionDetach(name string) error
	ActionGrow(name string, newSize int) error
	ActionClone(name, newName string) (*api.Volume, error)
}

type ImageClientInterface interface {
	Create(name, regionName, fileName string) (*api.Image, error)
	Get(name string) (*api.Image, error)
	Delete(name string) error
	List(limit int, marker string) (*api.ImageList, error)
}

type KeypairClientInterface interface {
	Create(name, publicKey string) (*api.Keypair, error)
	Get(name string) (*api.Keypair, error)
	Delete(name string) error
	List(limit int, marker string) (*api.KeypairList, error)
}

type NetworkPortClientInterface interface {
	Get(id string) (*api.NetworkPort, error)
	List(limit int, marker string) (*api.NetworkPortList, error)
	Delete(id string) error
}

type FlavorClientInterface interface {
	Create(name string, vcpus, ram, disk int) (*api.Flavor, error)
	Get(name string) (*api.Flavor, error)
	Delete(name string) error
	List(limit int, marker string) (*api.FlavorList, error)
}

type InstanceClientInterface interface {
	Create(name, imageName, regionName, zoneName, networkName, serviceAccountName, flavorName string, disk int, keypairNames []string, initialVolumes []api.InstanceInitialVolume, tags map[string]string, userData string) (*api.Instance, error)
	Get(name string) (*api.Instance, error)
	Delete(name string) error
	List(imageName string, limit int, marker string) (*api.InstanceList, error)
	ActionStop(name string, hard bool, timeout int) error
	ActionStart(name string) error
	ActionRestart(name string, hard bool, timeout int) error
	ActionImage(instanceName string, imageName string) (*api.Image, error)
}

type NetworkClientInterface interface {
	Create(name, regionName, portGroup, cidr string, gateway, poolStart, poolEnd net.IP, dnsServers []net.IP) (*api.Network, error)
	Get(name string) (*api.Network, error)
	Delete(name string) error
	List(region_name string, limit int, marker string) (*api.NetworkList, error)
}

type PermissionClientInterface interface {
	Get(name string) (*api.Permissions, error)
	List(limit int, marker string) (*api.PermissionList, error)
}

type RoleClientInterface interface {
	Create(name string, permissions []string) (*api.Role, error)
	Get(name string) (*api.Role, error)
	List(limit int, marker string) (*api.RoleList, error)
	Update(name string, permissions []string) error
	Delete(name string) error
}

type ServiceAccountClientInterface interface {
	Create(name string) (*api.ServiceAccount, error)
	Get(name string) (*api.ServiceAccount, error)
	Delete(name string) error
	List(limit int, marker string) (*api.ServiceAccountList, error)
	CreateKey(serviceAccountName, keyName string) (*oauth2.Token, error)
	DeleteKey(serviceAccountName, keyName string) error
}

type PolicyClientInterface interface {
	Get() (*api.Policy, error)
	Set(policy api.Policy) error
}

func (client *SandwichClient) createOAuthClient() *http.Client {
	ctx := context.Background()
	return oauth2.NewClient(ctx, oauth2.StaticTokenSource(client.token))
}

func (client *SandwichClient) Auth() AuthClientInterface {
	authClient := &auth.AuthClient{APIServer: client.APIServer}

	if client.token != nil {
		authClient.HttpClient = client.createOAuthClient()
	}

	return authClient
}

func (client *SandwichClient) Project() ProjectClientInterface {
	return &project.ProjectClient{APIServer: client.APIServer, HttpClient: client.createOAuthClient()}
}

func (client *SandwichClient) Region() RegionClientInterface {
	return &region.RegionClient{APIServer: client.APIServer, HttpClient: client.createOAuthClient()}
}

func (client *SandwichClient) Zone() ZoneClientInterface {
	return &zone.ZoneClient{APIServer: client.APIServer, HttpClient: client.createOAuthClient()}
}

func (client *SandwichClient) Volume(projectName string) VolumeClientInterface {
	return &volume.VolumeClient{APIServer: client.APIServer, HttpClient: client.createOAuthClient(), ProjectName: projectName}
}

func (client *SandwichClient) Image(projectName string) ImageClientInterface {
	return &image.ImageClient{APIServer: client.APIServer, HttpClient: client.createOAuthClient(), ProjectName: projectName}
}

func (client *SandwichClient) Network() NetworkClientInterface {
	return &network.NetworkClient{APIServer: client.APIServer, HttpClient: client.createOAuthClient()}
}

func (client *SandwichClient) NetworkPort(projectName string) NetworkPortClientInterface {
	return &network.NetworkPortClient{APIServer: client.APIServer, HttpClient: client.createOAuthClient(), ProjectName: projectName}
}

func (client *SandwichClient) Keypair(projectName string) KeypairClientInterface {
	return &keypair.KeypairClient{APIServer: client.APIServer, HttpClient: client.createOAuthClient(), ProjectName: projectName}
}

func (client *SandwichClient) Flavor() FlavorClientInterface {
	return &flavor.FlavorClient{APIServer: client.APIServer, HttpClient: client.createOAuthClient()}
}

func (client *SandwichClient) Instance(projectName string) InstanceClientInterface {
	return &instance.InstanceClient{APIServer: client.APIServer, HttpClient: client.createOAuthClient(), ProjectName: projectName}
}

func (client *SandwichClient) Permission() PermissionClientInterface {
	return &permission.PermissionClient{APIServer: client.APIServer, HttpClient: client.createOAuthClient()}
}

func (client *SandwichClient) SystemRole() RoleClientInterface {
	return &role.RoleClient{APIServer: client.APIServer, HttpClient: client.createOAuthClient(), Type: "system/roles"}
}

func (client *SandwichClient) ProjectRole(projectName string) RoleClientInterface {
	return &role.RoleClient{APIServer: client.APIServer, HttpClient: client.createOAuthClient(), Type: fmt.Sprintf("projects/%s/roles", projectName)}
}

func (client *SandwichClient) SystemServiceAccount() ServiceAccountClientInterface {
	return &serviceAccount.ServiceAccountClient{APIServer: client.APIServer, HttpClient: client.createOAuthClient(), Type: "system/service-accounts"}
}

func (client *SandwichClient) ProjectServiceAccount(projectName string) ServiceAccountClientInterface {
	return &serviceAccount.ServiceAccountClient{APIServer: client.APIServer, HttpClient: client.createOAuthClient(), Type: fmt.Sprintf("projects/%s/service-accounts", projectName)}
}

func (client *SandwichClient) SystemPolicy() PolicyClientInterface {
	return &policy.PolicyClient{APIServer: client.APIServer, HttpClient: client.createOAuthClient(), Type: "system/policy"}
}

func (client *SandwichClient) ProjectPolicy(projectName string) PolicyClientInterface {
	return &policy.PolicyClient{APIServer: client.APIServer, HttpClient: client.createOAuthClient(), Type: fmt.Sprintf("projects/%s/policy", projectName)}
}

func (client *SandwichClient) SetToken(token *oauth2.Token) {
	client.token = token
}
