package mitaka

import (
	"fmt"
	"strconv"

	"github.com/fassisrosa/beats/restClient"
	"github.com/fassisrosa/beats/openstackbeat/openstackapi/models"
)

type MitakaOpenStackAPI struct { }

func (api MitakaOpenStackAPI) GetAllInfo(mainUrl string) (virtualMachinesInfo []models.VirtualMachine, err error) {

	client, err := restClient.NewRestClient()
	if err != nil {
		return 
	}

	accessInfo, err := getAuthenticationInfo(client, mainUrl)
	if err != nil {
		return
	}
	flavorIdToName, err := getFlavorsInfo(client, accessInfo)
	if err != nil {
		return
	}
	userIdToName, err := getUsersInfo(client, accessInfo)
	if err != nil {
		return
	}
	tenantIdToName, err := getTenantsInfo(client, accessInfo)
	if err != nil {
		return
	}
	virtualMachinesInfo, err = getServersInfo(client, accessInfo, userIdToName, tenantIdToName, flavorIdToName)
	if err != nil {
		return 
	}
/*
	hostIdToName, err = getHostsInfo(client, accessInfo)
	if err != nil {
		return 
	}
*/


	return
}

func getAuthenticationInfo(client restClient.RestClient, mainUrl string) (accessInfo map[string]string, err error) {
	// set main url
	accessInfo = map[string]string{ "main": mainUrl }

	jsonStr := []byte(`{"auth": {"tenantName":"admin", "passwordCredentials": {"username": "admin", "password": "admin"}}}`)
	bodyJson, err := client.GetObject("POST", accessInfo["main"]+"/v2.0/tokens", jsonStr)
	if err != nil {
		return 
	}
	accessJson := bodyJson.GetObject("access")
	serviceCatalogsJson := accessJson.GetArray("serviceCatalog")
	for _, oneServiceJsonIfc := range serviceCatalogsJson {
		oneServiceJson := oneServiceJsonIfc.(restClient.JsonObject)
		oneServiceName := oneServiceJson.GetString("name")
		oneServiceEndpoints := oneServiceJson.GetArray("endpoints")
		for _, oneEndpointIfc := range oneServiceEndpoints {
			oneEndpoint := oneEndpointIfc.(restClient.JsonObject)
			accessInfo[oneServiceName] = oneEndpoint.GetString("publicURL")
		}
	}

	tokenJson := accessJson.GetObject("token")
	client.AddHeader("X-Auth-Token", tokenJson.GetString("id"))

	return accessInfo, nil
}

func getFlavorsInfo(client restClient.RestClient, accessInfo map[string]string) (flavorIdToInfo map[string]map[string]string, err error) {
	if accessInfo["nova"] == "" {
		err = fmt.Errorf("Could not find 'nova' endpoint")
		return
	}
	bodyJson, err := client.GetObject("GET", accessInfo["nova"] + "/flavors/detail?all_tenants=1", nil)
	if err != nil {
		return 
	}
	flavorIdToInfo = map[string]map[string]string{}
	for _, oneFlavorJsonIfc := range bodyJson.GetArray("flavors") {
		oneFlavorJson := oneFlavorJsonIfc.(restClient.JsonObject)
		flavorId := oneFlavorJson.GetString("id")
		flavorIdToInfo[flavorId] = map[string]string{}
		flavorIdToInfo[flavorId]["name"] = oneFlavorJson.GetString("name")
		flavorIdToInfo[flavorId]["ram"] = fmt.Sprintf("%d",oneFlavorJson.GetInteger("ram"))
		flavorIdToInfo[flavorId]["disk"] = fmt.Sprintf("%d",oneFlavorJson.GetInteger("disk"))
		flavorIdToInfo[flavorId]["vcpus"] = fmt.Sprintf("%d",oneFlavorJson.GetInteger("vcpus"))
	}
	return
}

func getUsersInfo(client restClient.RestClient, accessInfo map[string]string) (userIdToName map[string]string, err error) {
	if accessInfo["main"] == "" {
		err = fmt.Errorf("Could not find 'main' endpoint")
		return
	}
	bodyJson, err := client.GetObject("GET", accessInfo["main"] + "/v3/users?all_tenants=1", nil)
	if err != nil {
		return
	}
	userIdToName = map[string]string{}
	usersJson := bodyJson.GetArray("users")
	for _, oneUserJsonIfc := range usersJson {
		oneUserJson := oneUserJsonIfc.(restClient.JsonObject)
		userId := oneUserJson.GetString("id")
		userName := oneUserJson.GetString("name")
		userIdToName[userId] = userName
	}
	return
}

func getTenantsInfo(client restClient.RestClient, accessInfo map[string]string) (tenantIdToName map[string]string, err error) {
	if accessInfo["main"] == "" {
		err = fmt.Errorf("Could not find 'main' endpoint")
		return
	}
	bodyJson, err := client.GetObject("GET", accessInfo["main"] + "/v2.0/tenants", nil)
	if err != nil {
		return
	}
	tenantIdToName = map[string]string{}
	tenantsJson := bodyJson.GetArray("tenants")
	for _, oneTenantJsonIfc := range tenantsJson {
		oneTenantJson := oneTenantJsonIfc.(restClient.JsonObject)
		tenantId := oneTenantJson.GetString("id")
		tenantName := oneTenantJson.GetString("name")
		tenantIdToName[tenantId] = tenantName
	}
	return
}

func getServersInfo(client restClient.RestClient, accessInfo map[string]string, userIdToName map[string]string, tenantIdToName map[string]string, flavorIdToInfo map[string]map[string]string) (virtualMachinesInfo []models.VirtualMachine, err error) {
	if accessInfo["nova"] == "" {
		err = fmt.Errorf("Could not find 'nova' endpoint")
		return
	}
	bodyJson, err := client.GetObject("GET", accessInfo["nova"]+"/servers/detail?all_tenants=1", nil)
	if err != nil {
		return
	}
	serversJson := bodyJson.GetArray("servers")
	virtualMachinesInfo =  make([]models.VirtualMachine, len(serversJson), len(serversJson))
	for idx, oneServerJsonIfc := range serversJson {
		oneServerJson := oneServerJsonIfc.(restClient.JsonObject)
		oneVirtualMachineInfo := models.VirtualMachine{}
		oneVirtualMachineInfo.Id = oneServerJson.GetString("id")
		oneVirtualMachineInfo.Name = oneServerJson.GetString("name")
		oneVirtualMachineInfo.Status = oneServerJson.GetString("status")
		oneVirtualMachineInfo.Hypervisor = oneServerJson.GetString("OS-EXT-SRV-ATTR:hypervisor_hostname")
		userId := oneServerJson.GetString("user_id")
		tenantId := oneServerJson.GetString("tenant_id")
		oneVirtualMachineInfo.User = userIdToName[userId]
		oneVirtualMachineInfo.Tenant = tenantIdToName[tenantId]

		flavorJson := oneServerJson.GetObject("flavor")
		flavorId := flavorJson.GetString("id")
		oneVirtualMachineInfo.Flavor = flavorIdToInfo[flavorId]["name"]
		oneVirtualMachineInfo.MemorySize, _ = strconv.ParseInt(flavorIdToInfo[flavorId]["ram"], 10, 64)
		oneVirtualMachineInfo.DiskSize, _ = strconv.ParseInt(flavorIdToInfo[flavorId]["disk"], 10, 64)
		oneVirtualMachineInfo.NumberCpus, _ = strconv.ParseInt(flavorIdToInfo[flavorId]["vcpus"], 10, 64)

		virtualMachinesInfo[idx] = oneVirtualMachineInfo
	}
	return
}

