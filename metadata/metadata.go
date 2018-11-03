package metadata

import (
	"encoding/json"
	"fmt"
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
	httpClient httpClient
}

func NewMetadataProvider(client httpClient) provider {
	return provider{httpClient: client}
}

type httpClient interface {
	Get(url string, headers map[string]string) (responseCode int, body []byte, err error)
}

type MetadataCompute struct {
	Loocation             string      `json:"location"`
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

func (provider *provider)GetMetadata() (Metadata, error) {
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
