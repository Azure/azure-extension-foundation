// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package metadata

import (
	"encoding/json"
	"github.com/Azure/azure-extension-foundation/httputil"
	"io/ioutil"
	"testing"
)

var dummyMetadataJson = `{
        "compute": {
            "location": "some-location",
            "name": "some-computer",
            "offer": "some-offer",
            "osType": "Linux",
            "placementGroupId": "",
            "platformFaultDomain": "0",
            "platformUpdateDomain": "0",
            "publisher": "some-publisher",
            "resourceGroupName": "some-resourceGroup",
            "sku": "some-sku",
            "subscriptionId": "aaaa0000-aa00-aa00-aa00-aaaaaa000000",
            "tags": "",
            "version": "1.1.10",
            "vmId": "bbbb0000-bb00-bb00-bb00-bbbbbb000000",
            "vmSize": "Standard_D1"
        },
        "network": {
            "interface": [
                {
                    "ipv4": {
                        "ipAddress": [
                            {
                                "privateIpAddress": "0.0.0.0",
                                "publicIpAddress": "10.0.1.0"
                            }
                        ],
                        "subnet": [
                            {
                                "address": "10.0.0.0",
                                "prefix": "24"
                            }
                        ]
                    },
                    "ipv6": {
                        "ipAddress": []
                    },
                    "macAddress": "000AAABBB11",
                    "randomKey" : "randomValue"
                }
            ]
        }
    }`

func TestGetMetadataObjectFromJson(t *testing.T) {
	prov := NewMetadataProvider(&httputil.MockHttpClient{Getfunc: func(url string, headers map[string]string) (responseCode int, body []byte, err error) {
		return 200, []byte(dummyMetadataJson), nil
	}})

	metadata, err := prov.GetMetadata()
	if err != nil {
		t.Fatal(err.Error())
	}
	hostname := metadata.Compute.Name
	if hostname != "some-computer" {
		t.Fatalf("Hostname does not match. Expected: \"some-computer\", Actual: \"%s\"", hostname)
	}
	ipAddress := metadata.GetIpV4PublicAddress()
	if ipAddress != "10.0.1.0" {
		t.Fatalf("Ip address does not match. Expected: \"10.0.1.0\", Actual: \"%s\"", ipAddress)
	}
}

func TestEmptyHostnameFromMetadata(t *testing.T) {
	testJson := `{
        "compute": {
            "location": "some-location",
            "vmId": "bbbb0000-bb00-bb00-bb00-bbbbbb000000",
            "vmSize": "Standard_D1"
        },
        "network": {
        }
    }`

	prov := NewMetadataProvider(&httputil.MockHttpClient{Getfunc: func(url string, headers map[string]string) (responseCode int, body []byte, err error) {
		return 200, []byte(testJson), nil
	}})
	metadata, err := prov.GetMetadata()
	if err != nil {
		t.Fatal(err.Error())
	}
	hostname := metadata.Compute.Name
	if hostname != "" {
		t.Fatalf("Hostname \"%s\" returned when no hostname was expected", hostname)
	}
}

func TestRealMetadata(t *testing.T) {
	t.Skip() // for testing on Azure VM only
	client := httputil.NewSecureHttpClient(httputil.NoRetry)
	prov := provider{client}
	metadata, err := prov.GetMetadata()
	if err != nil {
		t.Fatal(err.Error())
	}
	jsonBytes, err := json.Marshal(metadata)
	if err != nil {
		t.Fatal(err.Error())
	}
	ioutil.WriteFile("testoutput/metadata.json", jsonBytes, 0666)
}
