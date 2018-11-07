// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package metadata

import (
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-extension-foundation/httputil"
)

const metadataUrl = "http://169.254.169.254/metadata/instance?api-version=2017-08-01"

type Metadata struct {
	Compute MetadataCompute `json:"compute"`
	Network MetadataNetwork `json:"network"`
}

type MetadataNetwork struct {
	Intrfc []map[string]interface{} `json:"interface"`
}

type provider struct {
	httpClient httputil.HttpClient
}

func NewMetadataProvider(client httputil.HttpClient) provider {
	return provider{httpClient: client}
}

type MetadataCompute struct {
	Location              string      `json:"location"`
	Name                  string      `json:"name"`
	Offer                 string      `json:"offer"`
	OsType                string      `json:"osType"`
	PlacementGroupId      string      `json:"placementGroupId"`
	PlatformFaultDomain   string      `json:"platformFaultDomain"`
	PlatformUpdateDomatin string      `json:"platformUpdateDomain"`
	Publisher             string      `json:"publisher"`
	ResourceGroupName     string      `json:"resourceGroupName"`
	Sku                   string      `json:"sku"`
	SubscriptionId        string      `json:"subscriptionId"`
	Tags                  interface{} `json:"tags"`
	Version               string      `json:"version"`
	VmId                  string      `json:"vmId"`
	VmSize                string      `json:"vmSize"`
}

func GetMetadataFromJsonString(jsonString *string) (Metadata, error) {
	retval := Metadata{}
	data := []byte(*jsonString)
	err := json.Unmarshal(data, &retval)
	return retval, err
}

func (metadata *Metadata) GetIpV4PublicAddress() string {
	defaultIp := "0.0.0.0"
	interface0Bytes, err := json.Marshal(metadata.Network.Intrfc[0]["ipv4"])
	if err != nil {
		return defaultIp
	}
	var interface0ipv4 map[string][]map[string]string
	err = json.Unmarshal(interface0Bytes, &interface0ipv4)
	if err != nil {
		return defaultIp
	}
	retval := ""
	if len(interface0ipv4["ipAddress"]) > 0 {
		retval = interface0ipv4["ipAddress"][0]["publicIpAddress"]
	}

	if retval == "" {
		return defaultIp
	}
	return retval
}

func (metadata *Metadata) GetAzureResourceId() string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/virtualMachines/%s",
		metadata.Compute.SubscriptionId, metadata.Compute.ResourceGroupName, metadata.Compute.Name)
}

func (provider *provider) GetMetadata() (Metadata, error) {
	retval := Metadata{}
	responseCode, responseBody, err := provider.httpClient.Get(metadataUrl, map[string]string{"Metadata": "true"})
	if err != nil {
		return retval, err
	}
	responseString := string(responseBody[:])
	if responseCode != 200 {
		return retval, fmt.Errorf("Get request for metadata returned return code %v.\nResponse Body: %s", responseCode, responseString)
	}
	err = json.Unmarshal(responseBody[:], &retval)
	return retval, err
}
