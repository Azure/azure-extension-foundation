# Azure extension foundation
This repository contains the foundation source code for Azure Virtual Machine extension developers.
This source code is meant to be used by developpers publishing Virtual Machine extensions and the source code is open sourced under MIT License for reference. You can read the User Guide below.

* [Learn more: Azure Virtual Machine Extensions](https://azure.microsoft.com/en-us/documentation/articles/virtual-machines-extensions-features/)

# Usage
### Status reporting, sequence tracking and settings manipulation

```go
package main

import (
	"azure-extension-foundation/sequence"
	"azure-extension-foundation/settings"
	"fmt"
	"os"
)

// extension specific PublicSettings
type PublicSettings struct {
	Script   string   `json:"script"`
	FileURLs []string `json:"fileUris"`
}

// extension specific ProtectedSettings
type ProtectedSettings struct {
	SecretString       string   `json:"secretString"`
	SecretScript       string   `json:"secretScript"`
	FileURLs           []string `json:"fileUris"`
	StorageAccountName string   `json:"storageAccountName"`
	StorageAccountKey  string   `json:"storageAccountKey"`
}

func main() {
	extensionMrseq, environmentMrseq, err := sequence.GetMostRecentSequenceNumber()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	shouldRun := sequence.ShouldBeProcessed(extensionMrseq, environmentMrseq)
	if !shouldRun {
		fmt.Printf("environment mrseq has already been processed by extension (environment mrseq : %v, extension mrseq : %v)\n", environmentMrseq, extensionMrseq)
		os.Exit(-1)
	}

	err = sequence.SetExtensionMostRecentSequenceNumber(environmentMrseq)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	err = status.ReportTransitioning(environmentMrseq, "install", "installation in progress")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	var publicSettings PublicSettings
	var protectedSettings ProtectedSettings
	err = settings.GetExtensionSettings(environmentMrseq, &publicSettings, &protectedSettings)
	if err != nil {
		status.ReportError(environmentMrseq, "install", err.Error())
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	err = status.ReportSuccess(environmentMrseq, "install", "installation in complete")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
}
```
### Simple http client
``` go
func main() {
	client := httputil.NewSecureHttpClient(httputil.NoRetry)
	status, response, err := client.Get("http://www.microsoft.com/", [header])
	if err != nil {
		fmt.Println("error issuing get call")
		os.Exit(-1)
	}
}
```

### MSI
``` go
// struct definition; snippet from msi/msi.go
type Msi struct {
	AccessToken  string `json:"access_token"`
	ClientID     string `json:"client_id"`
	ExpiresIn    string `json:"expires_in"`
	ExpiresOn    string `json:"expires_on"`
	ExtExpiresIn string `json:"ext_expires_in"`
	NotBefore    string `json:"not_before"`
	Resource     string `json:"resource"`
	TokenType    string `json:"token_type"`
}
```

``` go
func main(){
	secureHttpClient := httputil.NewSecureHttpClient(httputil.NoRetry)
	msiProvider := NewMsiProvider(secureHttpClient)
	msi, err := msiProvider.GetMsi()
	if err != nil {
		fmt.Println("error getting msi")
		os.Exit(-1)
	}
}
```

# Contributing

This project welcomes contributions and suggestions.  Most contributions require you to agree to a
Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us
the rights to use your contribution. For details, visit https://cla.microsoft.com.

When you submit a pull request, a CLA-bot will automatically determine whether you need to provide
a CLA and decorate the PR appropriately (e.g., label, comment). Simply follow the instructions
provided by the bot. You will only need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or
contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.
