# Azure extension foundation
This repository contains the foundation source code for Azure Virtual Machine extension developers.
This source code is meant to be used by Microsoft Azure employees publishing Virtual Machine extensions and the source code is open sourced under Apache 2.0 License for reference. You can read the User Guide below.

* [Learn more: Azure Virtual Machine Extensions](https://azure.microsoft.com/en-us/documentation/articles/virtual-machines-extensions-features/)

# Usage

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

	sequence.SetExtensionMostRecentSequenceNumber(environmentMrseq)

	var publicSettings PublicSettings
	var protectedSettings ProtectedSettings
	settings.GetExtensionSettings(environmentMrseq, &publicSettings, &protectedSettings)
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
